package habitsRepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	xpmodel "kai-back/internal/modules/xp/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) dbFromTx(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}

	return r.db
}

func (r *Repository) FindUserHabits(ctx context.Context, userID uuid.UUID) ([]habitsmodel.UserHabit, error) {
	var habits []habitsmodel.UserHabit

	err := r.db.WithContext(ctx).
		Preload("HabitCatalog").
		Preload("HabitRecords", func(db *gorm.DB) *gorm.DB {
			return db.Order("fecha DESC")
		}).
		Where("usuario_id = ? AND activo = true", userID).
		Order("fecha_inicio DESC").
		Find(&habits).
		Error
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *Repository) CountDailyCompleted(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (int, error) {
	var completed sql.NullInt64

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.WithContext(ctx).
		Table("registros_habito AS rh").
		Select("COUNT(DISTINCT rh.habito_usuario_id)").
		Joins("JOIN habitos_usuario AS hu ON hu.id = rh.habito_usuario_id").
		Where("rh.usuario_id = ? AND hu.activo = true AND rh.fecha = CURRENT_DATE AND rh.completado = true", userID).
		Scan(&completed).
		Error
	if err != nil {
		return 0, err
	}

	return int(completed.Int64), nil
}

func (r *Repository) FindCategories(
	ctx context.Context,
) ([]habitsdto.HabitCategoryResponse, error) {

	var categories []habitsdto.HabitCategoryResponse

	err := r.db.
		WithContext(ctx).
		Model(&xpmodel.XPCategory{}).
		Select(`
			categorias_xp.id,
			categorias_xp.nombre as name,
			categorias_xp.descripcion as description
		`).
		Group(`
			categorias_xp.id,
			categorias_xp.nombre,
			categorias_xp.descripcion
		`).
		Order("categorias_xp.nombre ASC").
		Scan(&categories).
		Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *Repository) FindCatalogByCategoryExcludingUser(
	ctx context.Context,
	userID uuid.UUID,
	categoryID uuid.UUID,
) ([]habitsmodel.HabitCatalog, error) {

	var habits []habitsmodel.HabitCatalog

	subQuery := r.db.
		Table("habitos_usuario").
		Select("habito_catalogo_id").
		Where("usuario_id = ?", userID).
		Where("activo = true")

	err := r.db.
		WithContext(ctx).
		Where("activo = true").
		Where("categoria_xp_id = ?", categoryID).
		Where("id NOT IN (?)", subQuery).
		Order("nombre ASC").
		Find(&habits).
		Error

	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *Repository) FindCatalogByID(
	ctx context.Context,
	catalogID uuid.UUID,
) (*habitsmodel.HabitCatalog, error) {

	var habit habitsmodel.HabitCatalog

	err := r.db.
		WithContext(ctx).
		Where("activo = true").
		First(&habit, "id = ?", catalogID).
		Error

	if err != nil {
		return nil, err
	}

	return &habit, nil
}

func (r *Repository) UserHasActiveHabit(
	ctx context.Context,
	userID uuid.UUID,
	catalogID uuid.UUID,
) (bool, error) {

	var count int64

	err := r.db.
		WithContext(ctx).
		Model(&habitsmodel.UserHabit{}).
		Where("usuario_id = ?", userID).
		Where("habito_catalogo_id = ?", catalogID).
		Where("activo = true").
		Count(&count).
		Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repository) FindUserHabitByCatalogID(
	ctx context.Context,
	userID uuid.UUID,
	habitCatalogID uuid.UUID,
) (*habitsmodel.UserHabit, error) {

	var habit habitsmodel.UserHabit

	err := r.db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("habito_catalogo_id = ?", habitCatalogID).
		First(&habit).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &habit, nil
}

func (r *Repository) ReactivateHabit(
	ctx context.Context,
	habitID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Model(&habitsmodel.UserHabit{}).
		Where("id = ?", habitID).
		Update("activo", true).
		Error
}

func (r *Repository) ReactivateHabitWithTodayRecord(
	ctx context.Context,
	userHabit *habitsmodel.UserHabit,
) error {

	return r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {

			if err := tx.
				Model(&habitsmodel.UserHabit{}).
				Where("id = ?", userHabit.ID).
				Update("activo", true).
				Error; err != nil {
				return err
			}

			var count int64

			if err := tx.
				Model(&habitsmodel.HabitRecord{}).
				Where("habito_usuario_id = ?", userHabit.ID).
				Where("fecha = CURRENT_DATE").
				Count(&count).
				Error; err != nil {
				return err
			}

			if count == 0 {

				record := habitsmodel.HabitRecord{
					UsuarioID:   userHabit.UserID,
					UserHabitID: userHabit.ID,
					Fecha:       time.Now(),
					Completado:  false,
					XPGanada:    0,
				}

				if err := tx.Create(&record).Error; err != nil {
					return err
				}
			}

			return nil
		})
}

func (r *Repository) SelectHabit(
	ctx context.Context,
	userHabit *habitsmodel.UserHabit,
	initialRecord *habitsmodel.HabitRecord,
) error {

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(userHabit).Error; err != nil {
			return err
		}

		initialRecord.UserHabitID = userHabit.ID

		if err := tx.Create(initialRecord).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *Repository) FindHabitDetailByID(
	ctx context.Context,
	userID uuid.UUID,
	habitUserID uuid.UUID,
) (*habitsmodel.UserHabit, error) {

	var habit habitsmodel.UserHabit

	err := r.db.
		WithContext(ctx).
		Model(&habitsmodel.UserHabit{}).
		Preload("HabitCatalog").
		Preload("HabitRecords", func(db *gorm.DB) *gorm.DB {
			return db.Order("fecha DESC")
		}).
		Where(
			"id = ? AND usuario_id = ? AND activo = true",
			habitUserID,
			userID,
		).
		First(&habit).
		Error

	if err != nil {
		return nil, err
	}

	return &habit, nil
}

func (r *Repository) FindUserHabitByID(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	habitID uuid.UUID,
) (*habitsmodel.UserHabit, error) {

	var habit habitsmodel.UserHabit

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Preload("HabitCatalog").
		Where("id = ?", habitID).
		Where("usuario_id = ?", userID).
		First(&habit).
		Error

	if err != nil {
		return nil, err
	}

	return &habit, nil
}

func (r *Repository) DeactivateHabit(
	ctx context.Context,
	habitID uuid.UUID,
) error {

	return r.db.
		WithContext(ctx).
		Model(&habitsmodel.UserHabit{}).
		Where("id = ?", habitID).
		Update("activo", false).
		Error
}

func (r *Repository) FindActiveUserHabits(
	ctx context.Context,
	userID uuid.UUID,
) ([]habitsmodel.UserHabit, error) {

	var habits []habitsmodel.UserHabit

	err := r.db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("activo = true").
		Find(&habits).
		Error

	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *Repository) FindTodayRecords(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) ([]habitsmodel.HabitRecord, error) {

	var records []habitsmodel.HabitRecord

	err := r.db.
		WithContext(ctx).
		Where("usuario_id = ?", userID).
		Where("fecha = CURRENT_DATE").
		Find(&records).
		Error

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository) FindTodayRecord(
	ctx context.Context,
	tx *gorm.DB,
	userHabitID uuid.UUID,
	date time.Time,
) (*habitsmodel.HabitRecord, error) {

	var record habitsmodel.HabitRecord

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Where("habito_usuario_id = ?", userHabitID).
		Where("fecha = CURRENT_DATE").
		Find(&record).
		Error

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *Repository) CreateHabitRecords(
	ctx context.Context,
	tx *gorm.DB,
	records []habitsmodel.HabitRecord,
) error {
	db := r.db

	if tx != nil {
		db = tx
	}

	if len(records) == 0 {
		return nil
	}

	return db.
		WithContext(ctx).
		Create(&records).
		Error
}

func (r *Repository) UpdateHabitRecord(
	ctx context.Context,
	tx *gorm.DB,
	record *habitsmodel.HabitRecord,
) error {
	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Model(&habitsmodel.HabitRecord{}).
		Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"completado":       record.Completado,
			"valor_registrado": record.ValorRegistrado,
			"xp_ganada":        record.XPGanada,
		}).
		Error
}

func (r *Repository) FindStreakByHabit(
	ctx context.Context,
	tx *gorm.DB,
	userHabitID uuid.UUID,
) (*habitsmodel.Streak, error) {

	var streak habitsmodel.Streak

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Where("habito_usuario_id = ?", userHabitID).
		First(&streak).
		Error

	if err != nil {
		return nil, err
	}

	return &streak, nil
}

func (r *Repository) CreateStreak(
	ctx context.Context,
	tx *gorm.DB,
	streak *habitsmodel.Streak,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Create(streak).
		Error
}

func (r *Repository) UpdateStreak(
	ctx context.Context,
	tx *gorm.DB,
	streak *habitsmodel.Streak,
) error {

	db := r.db

	if tx != nil {
		db = tx
	}

	return db.
		WithContext(ctx).
		Save(streak).
		Error
}

func (r *Repository) FindTodayHabitRecords(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
	date time.Time,
) ([]habitsmodel.HabitRecord, error) {
	var records []habitsmodel.HabitRecord

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Joins("JOIN habitos_usuario AS hu ON hu.id = registros_habito.habito_usuario_id").
		Where("registros_habito.usuario_id = ?", userID).
		Where("hu.activo = ?", true).
		Where("registros_habito.fecha = CURRENT_DATE").
		Order("registros_habito.fecha DESC").
		Find(&records).
		Error

	if err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Repository) FindLastCompletedHabitDate(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) (*time.Time, error) {

	var lastDate sql.NullTime

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Model(&habitsmodel.HabitRecord{}).
		Select("MAX(fecha)").
		Where("usuario_id = ?", userID).
		Where("completado = true").
		Scan(&lastDate).
		Error

	if err != nil {
		return nil, err
	}

	return &lastDate.Time, nil
}

func (r *Repository) FindCompletedHabitDates(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) ([]time.Time, error) {

	var dates []time.Time

	db := r.db

	if tx != nil {
		db = tx
	}

	err := db.
		WithContext(ctx).
		Model(&habitsmodel.HabitRecord{}).
		Distinct("fecha").
		Where("usuario_id = ?", userID).
		Where("completado = true").
		Order("fecha DESC").
		Pluck("fecha", &dates).
		Error

	if err != nil {
		return nil, err
	}

	return dates, nil
}
