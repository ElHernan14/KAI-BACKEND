package habitcompletion

import (
	"context"
	habit "kai-back/internal/modules/habits/repository"
	habitDailyRecordsServicePort "kai-back/internal/modules/habits/service"
	kai "kai-back/internal/modules/kai/repository"
	message "kai-back/internal/modules/messages/repository"
	xp "kai-back/internal/modules/xp/repository"
	transaction "kai-back/internal/shared/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	habitRepo                 habit.HabitsRepository
	xpRepo                    xp.XpRepository
	kaiRepo                   kai.KaiRepository
	messageRepo               message.MessageRepositoryPort
	habitsDailyRecordsService habitDailyRecordsServicePort.HabitsDailyRecordsServicePort
	transactionManager        transaction.TransactionManager
}

func NewService(
	habitRepo habit.HabitsRepository,
	xpRepo xp.XpRepository,
	kaiRepo kai.KaiRepository,
	messageRepo message.MessageRepositoryPort,
	habitsDailyRecordsService habitDailyRecordsServicePort.HabitsDailyRecordsServicePort,
	transactionManager transaction.TransactionManager,
) *Service {
	return &Service{
		habitRepo:                 habitRepo,
		xpRepo:                    xpRepo,
		kaiRepo:                   kaiRepo,
		messageRepo:               messageRepo,
		habitsDailyRecordsService: habitsDailyRecordsService,
		transactionManager:        transactionManager,
	}
}

func (s *Service) CompleteHabit(
	ctx context.Context,
	userID uuid.UUID,
	userHabitID uuid.UUID,
	value *string,
) error {
	// Aseguramos que existan registros de hábitos para hoy antes de completar el hábito del usuario
	err := s.habitsDailyRecordsService.
		EnsureTodayHabitRecords(
			ctx,
			userID,
		)
	if err != nil {
		return err
	}
	// Iniciamos la transacción para completar el hábito del usuario y todas las operaciones relacionadas
	return s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {
			habit, rec, err := s.validateHabit(
				ctx,
				userID,
				userHabitID,
			)
			if err != nil {
				return err
			}

			record, err := s.completeRecord(
				ctx,
				tx,
				habit,
				value,
				rec,
			)
			if err != nil {
				return err
			}
			_ = record

			streak, err := s.updateStreak(
				ctx,
				tx,
				userID,
				habit,
			)
			if err != nil {
				return err
			}

			xpGranted, err := s.grantXP(
				ctx,
				tx,
				userID,
				habit,
			)
			if err != nil {
				return err
			}

			err = s.grantKaiAttributes(
				ctx,
				tx,
				userID,
				habit,
				xpGranted,
			)
			if err != nil {
				return err
			}

			dominantAttribute, err := s.recalculateDominantAttribute(
				ctx,
				tx,
				userID,
			)
			if err != nil {
				return err
			}

			dominantAttributeValue := uuid.Nil
			if dominantAttribute != nil {
				dominantAttributeValue = *dominantAttribute
			}

			err = s.updateKaiState(
				ctx,
				tx,
				userID,
				&dominantAttributeValue,
			)
			if err != nil {
				return err
			}

			err = s.generateMotivationalMessage(
				ctx,
				tx,
				userID,
				&dominantAttributeValue,
			)
			if err != nil {
				return err
			}

			_ = streak

			return nil
		},
	)
}
