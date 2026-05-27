package habits

import (
	"context"
	"database/sql"

	habitsmodel "kai-back/internal/modules/habits/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryPort interface {
	FindUserHabits(ctx context.Context, userID uuid.UUID) ([]habitsmodel.UserHabit, error)
	CountDailyCompleted(ctx context.Context, userID uuid.UUID) (int, error)
}

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
