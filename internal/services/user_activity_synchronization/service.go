package useractivitysynchronization

import (
	"context"
	habitsrepo "kai-back/internal/modules/habits/repository"
	kai "kai-back/internal/modules/kai/repository"
	userrepo "kai-back/internal/modules/users/repository"
	transaction "kai-back/internal/shared/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	transactionManager transaction.TransactionManager

	userRepository  userrepo.UsersRepository
	habitRepository habitsrepo.HabitsRepository
	kaiRepository   kai.KaiRepository
}

func New(
	transactionManager transaction.TransactionManager,
	userRepository userrepo.UsersRepository,
	habitRepository habitsrepo.HabitsRepository,
	kaiRepository kai.KaiRepository,
) *Service {
	return &Service{
		transactionManager: transactionManager,
		userRepository:     userRepository,
		habitRepository:    habitRepository,
		kaiRepository:      kaiRepository,
	}
}

func (s *Service) SyncUserActivityState(
	ctx context.Context,
	userID uuid.UUID,
) error {

	return s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {

			lastActivity, err := s.findLastCompletedHabitDate(
				ctx,
				tx,
				userID,
			)
			if err != nil {
				return err
			}

			inactiveDays := s.calculateInactiveDays(
				lastActivity,
			)

			globalStreak, err := s.calculateGlobalStreak(
				ctx,
				tx,
				userID,
			)
			if err != nil {
				return err
			}

			kaiState, err := s.kaiRepository.
				FindKaiStateByUserID(
					ctx,
					tx,
					userID,
				)
			if err != nil {
				return err
			}

			energy := s.calculateEnergy(
				kaiState.Energy,
				inactiveDays,
			)

			state := s.determineKaiState(
				globalStreak,
				inactiveDays,
			)

			mode := s.determineKaiMode(
				inactiveDays,
			)

			err = s.userRepository.
				UpdateActivityStats(
					ctx,
					tx,
					userID,
					globalStreak,
					inactiveDays,
				)
			if err != nil {
				return err
			}

			err = s.kaiRepository.
				UpdateTemporalState(
					ctx,
					tx,
					userID,
					energy,
					state,
					mode,
					inactiveDays,
				)
			if err != nil {
				return err
			}

			return nil
		},
	)
}
