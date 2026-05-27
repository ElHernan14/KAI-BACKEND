package habits

import (
	"context"
	"log"
	"net/http"
	"time"

	habitsdto "kai-back/internal/modules/habits/dto"
	habitsmodel "kai-back/internal/modules/habits/models"
	errorHandler "kai-back/internal/shared/errors"

	"github.com/google/uuid"
)

type ServicePort interface {
	GetUserHabits(ctx context.Context, userID uuid.UUID) (*habitsdto.HabitsViewResponse, error)
}

type Service struct {
	repository RepositoryPort
}

func NewService(repository RepositoryPort) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetUserHabits(ctx context.Context, userID uuid.UUID) (*habitsdto.HabitsViewResponse, error) {
	userHabits, err := s.repository.FindUserHabits(ctx, userID)
	if err != nil {
		log.Printf("Error fetching user habits: %v", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "error al obtener hábitos del usuario")
	}

	completedToday, err := s.repository.CountDailyCompleted(ctx, userID)
	if err != nil {
		log.Printf("Error counting daily completed habits: %v", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "error al contar hábitos completados hoy")
	}

	habitsResponse := make([]habitsdto.UserHabitResponse, 0, len(userHabits))
	for _, habit := range userHabits {
		habitsResponse = append(habitsResponse, buildUserHabitResponse(habit))
	}

	total := len(habitsResponse)
	if completedToday > total {
		completedToday = total
	}

	return &habitsdto.HabitsViewResponse{
		DailyProgress: habitsdto.DailyProgress{
			Total:     total,
			Completed: completedToday,
			Pending:   total - completedToday,
		},
		Habits: habitsResponse,
	}, nil
}

func buildUserHabitResponse(habit habitsmodel.UserHabit) habitsdto.UserHabitResponse {
	records := make([]habitsdto.HabitRecordResponse, 0, len(habit.HabitRecords))
	completedToday := false

	for _, record := range habit.HabitRecords {
		if isToday(record.Fecha) && record.Completado {
			completedToday = true
		}

		records = append(records, habitsdto.HabitRecordResponse{
			ID:              record.ID.String(),
			Date:            record.Fecha,
			Completed:       record.Completado,
			RegisteredValue: record.ValorRegistrado,
			XPEarned:        record.XPGanada,
			CreatedAt:       record.CreatedAt,
		})
	}

	return habitsdto.UserHabitResponse{
		ID:             habit.ID.String(),
		HabitCatalogID: habit.HabitCatalogID.String(),
		Name:           habit.HabitCatalog.Name,
		Description:    habit.HabitCatalog.Description,
		Category:       habit.HabitCatalog.Category,
		CareType:       habit.HabitCatalog.CareType,
		Difficulty:     habit.HabitCatalog.Difficulty,
		BaseXP:         habit.HabitCatalog.BaseXP,
		Custom:         habit.Custom,
		Active:         habit.Active,
		StartDate:      habit.StartDate,
		HabitImage:     habit.HabitCatalog.HabitImage,
		CompletedToday: completedToday,
		Records:        records,
	}
}

func isToday(value time.Time) bool {
	now := time.Now()
	year, month, day := now.Date()
	valueYear, valueMonth, valueDay := value.Date()

	return year == valueYear && month == valueMonth && day == valueDay
}
