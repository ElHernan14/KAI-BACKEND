package habitsService

import (
	"context"
	habitsmodel "kai-back/internal/modules/habits/models"
	habitRepo "kai-back/internal/modules/habits/repository"
	errorHandler "kai-back/internal/shared/errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type HabitsDailyRecordsServicePort interface {
	EnsureTodayHabitRecords(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type HabitsDailyRecordsService struct {
	habitRepo habitRepo.HabitsRepository
}

func NewHabitsDailyRecordsService(
	habitRepo habitRepo.HabitsRepository,
) *HabitsDailyRecordsService {
	return &HabitsDailyRecordsService{
		habitRepo: habitRepo,
	}
}

func (s *HabitsDailyRecordsService) EnsureTodayHabitRecords(
	ctx context.Context,
	userID uuid.UUID,
) error {

	activeHabits, err := s.habitRepo.
		FindActiveUserHabits(
			ctx,
			userID,
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al obtener hábitos activos del usuario",
		)
	}

	todayRecords, err := s.habitRepo.
		FindTodayRecords(
			ctx,
			userID,
			time.Now(),
		)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusInternalServerError,
			"error al obtener registros de hábitos del día",
		)
	}

	existingRecords := make(
		map[uuid.UUID]bool,
	)

	for _, record := range todayRecords {
		existingRecords[record.UserHabitID] = true
	}

	var recordsToCreate []habitsmodel.HabitRecord

	for _, habit := range activeHabits {

		if existingRecords[habit.ID] {
			continue
		}

		recordsToCreate = append(
			recordsToCreate,
			habitsmodel.HabitRecord{
				UsuarioID:   userID,
				UserHabitID: habit.ID,
				Fecha:       time.Now(),
				Completado:  false,
				XPGanada:    0,
			},
		)
	}

	return s.habitRepo.CreateHabitRecords(
		ctx,
		nil,
		recordsToCreate,
	)
}
