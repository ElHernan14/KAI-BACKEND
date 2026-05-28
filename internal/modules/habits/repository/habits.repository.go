package habitsRepository

import (
	"context"
	"database/sql"

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

func (r *Repository) CountDailyCompleted(ctx context.Context, userID uuid.UUID) (int, error) {
	var completed sql.NullInt64

	err := r.db.WithContext(ctx).
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
