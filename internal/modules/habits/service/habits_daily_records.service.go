package habitsService

import (
	"context"
	habitsmodel "kai-back/internal/modules/habits/models"
	habitRepo "kai-back/internal/modules/habits/repository"
	userRepo "kai-back/internal/modules/users/repository"
	errorHandler "kai-back/internal/shared/errors"
	transaction "kai-back/internal/shared/transaction"
	"kai-back/internal/shared/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HabitsDailyRecordsServicePort interface {
	EnsureTodayHabitRecords(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type HabitsDailyRecordsService struct {
	habitRepo          habitRepo.HabitsRepository
	userRepository     userRepo.UsersRepository
	transactionManager transaction.TransactionManager
}

func NewHabitsDailyRecordsService(
	habitRepo habitRepo.HabitsRepository,
	userRepository userRepo.UsersRepository,
	transactionManager transaction.TransactionManager,
) *HabitsDailyRecordsService {
	return &HabitsDailyRecordsService{
		habitRepo:          habitRepo,
		userRepository:     userRepository,
		transactionManager: transactionManager,
	}
}

func (s *HabitsDailyRecordsService) EnsureTodayHabitRecords(
	ctx context.Context,
	userID uuid.UUID,
) error {
	return s.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) error {
		now := time.Now()
		user, err := s.userRepository.FindUserSyncState(ctx, tx, userID)
		if err != nil {
			return err
		}
		if utils.IsSameCalendarDay(user.HabitRecordsSyncedAt, now) {
			return nil
		}

		activeHabits, err := s.habitRepo.
			FindActiveUserHabits(
				ctx,
				tx,
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
				tx,
				userID,
				now,
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
					Fecha:       now,
					Completado:  false,
					XPGanada:    0,
				},
			)
		}

		if err := s.habitRepo.CreateHabitRecords(ctx, tx, recordsToCreate); err != nil {
			return err
		}

		return s.userRepository.UpdateHabitRecordsSyncedAt(ctx, tx, userID, now)
	})
}
