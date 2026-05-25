package home

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type KaiSummaryRow struct {
	CurrentState string
	CurrentStage string
	Energy       int
	KaiImage     *string
	LastMessage  *string
	RecoveryMode bool
	BondLevel    int
}

type DailyHabitRow struct {
	ID        uuid.UUID
	Name      string
	Category  *string
	Completed bool
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindKaiSummary(ctx context.Context, userID uuid.UUID) (*KaiSummaryRow, error) {
	var row KaiSummaryRow

	err := r.db.WithContext(ctx).
		Table("estado_kai").
		Select(`
			estado_actual AS current_state,
			etapa_actual AS current_stage,
			energia AS energy,
			imagen_kai AS kai_image,
			ultimo_mensaje AS last_message,
			modo_recuperacion AS recovery_mode,
			nivel_vinculo AS bond_level
		`).
		Where("usuario_id = ?", userID).
		Take(&row).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *Repository) FindTotalXP(ctx context.Context, userID uuid.UUID) (int, error) {
	var total sql.NullInt64

	err := r.db.WithContext(ctx).
		Table("xp_usuario").
		Select("COALESCE(SUM(valor), 0)").
		Where("usuario_id = ?", userID).
		Scan(&total).
		Error
	if err != nil {
		return 0, err
	}

	return int(total.Int64), nil
}

func (r *Repository) FindCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	var current sql.NullInt64

	err := r.db.WithContext(ctx).
		Table("rachas").
		Select("COALESCE(MAX(dias_actuales), 0)").
		Where("usuario_id = ? AND activa = true", userID).
		Scan(&current).
		Error
	if err != nil {
		return 0, err
	}

	return int(current.Int64), nil
}

func (r *Repository) FindDailyHabits(ctx context.Context, userID uuid.UUID) ([]DailyHabitRow, error) {
	var habits []DailyHabitRow

	err := r.db.WithContext(ctx).
		Table("habitos_usuario AS hu").
		Select(`
			hu.id AS id,
			hc.nombre AS name,
			hc.categoria AS category,
			COALESCE(rh.completado, false) AS completed
		`).
		Joins("JOIN habitos_catalogo AS hc ON hc.id = hu.habito_catalogo_id").
		Joins(`
			LEFT JOIN registros_habito AS rh
			ON rh.habito_usuario_id = hu.id
			AND rh.fecha = CURRENT_DATE
		`).
		Where("hu.usuario_id = ? AND hu.activo = true", userID).
		Order("hc.nombre ASC").
		Scan(&habits).
		Error
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (r *Repository) FindFallbackMessage(ctx context.Context) (*string, error) {
	var message string

	err := r.db.WithContext(ctx).
		Table("mensajes_kai").
		Select("mensaje").
		Where("activo = true").
		Order("created_at DESC").
		Limit(1).
		Scan(&message).
		Error
	if err != nil {
		return nil, err
	}
	if message == "" {
		return nil, nil
	}

	return &message, nil
}
