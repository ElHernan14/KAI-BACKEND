package home

import (
	"context"
	"errors"
	"time"

	habitsmodel "kai-back/internal/modules/habits/models"
	messagesmodel "kai-back/internal/modules/messages/models"
	xpmodel "kai-back/internal/modules/xp/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type KaiSummaryRow struct {
	CurrentState    string
	CurrentStage    string
	Energy          int
	KaiImage        *string
	LastMessage     *string
	RecoveryMode    bool
	BondLevel       int
	LastEvolution   *time.Time
	LastInteraction *time.Time
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
			nivel_vinculo AS bond_level,
			ultima_evolucion AS last_evolution,
			ultima_interaccion AS last_interaction
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
	var total int

	err := r.db.WithContext(ctx).
		Model(&xpmodel.UserXP{}).
		Select("COALESCE(SUM(valor), 0)").
		Where("usuario_id = ?", userID).
		Scan(&total).
		Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *Repository) FindCurrentStreak(ctx context.Context, userID uuid.UUID) (int, error) {
	var current int

	err := r.db.WithContext(ctx).
		Model(&habitsmodel.Streak{}).
		Select("COALESCE(MAX(dias_actuales), 0)").
		Where("usuario_id = ? AND activa = true", userID).
		Scan(&current).
		Error
	if err != nil {
		return 0, err
	}

	return current, nil
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
	var message messagesmodel.KaiMessage

	err := r.db.WithContext(ctx).
		Model(&messagesmodel.KaiMessage{}).
		Select("mensaje").
		Where("activo = true").
		Where("tipo NOT IN ?", []string{
			messagesmodel.MessageTypeWelcome,
			messagesmodel.MessageTypeGreeting,
			messagesmodel.MessageTypeReturn,
			messagesmodel.MessageTypeEvolutionYoung,
			messagesmodel.MessageTypeEvolutionAdult,
		}).
		Order("created_at DESC").
		Take(&message).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if message.Message == "" {
		return nil, nil
	}

	return &message.Message, nil
}

func (r *Repository) FindRandomEvolutionMessage(ctx context.Context) (*string, error) {
	var message messagesmodel.KaiMessage

	err := r.db.WithContext(ctx).
		Model(&messagesmodel.KaiMessage{}).
		Select("mensaje").
		Where("activo = true").
		Where("tipo IN ? OR contexto = ?", []string{
			messagesmodel.MessageTypeEvolutionYoung,
			messagesmodel.MessageTypeEvolutionAdult,
		}, "evolucion").
		Order("RANDOM()").
		Take(&message).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if message.Message == "" {
		return nil, nil
	}

	return &message.Message, nil
}
