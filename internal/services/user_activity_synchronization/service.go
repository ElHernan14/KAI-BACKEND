package useractivitysynchronization

import (
	"context"
	habitsrepo "kai-back/internal/modules/habits/repository"
	kai "kai-back/internal/modules/kai/repository"
	messagesmodel "kai-back/internal/modules/messages/models"
	messagesrepo "kai-back/internal/modules/messages/repository"
	userrepo "kai-back/internal/modules/users/repository"
	transaction "kai-back/internal/shared/transaction"
	"kai-back/internal/shared/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	transactionManager transaction.TransactionManager

	userRepository    userrepo.UsersRepository
	habitRepository   habitsrepo.HabitsRepository
	kaiRepository     kai.KaiRepository
	messageRepository messagesrepo.MessageRepositoryPort
}

func New(
	transactionManager transaction.TransactionManager,
	userRepository userrepo.UsersRepository,
	habitRepository habitsrepo.HabitsRepository,
	kaiRepository kai.KaiRepository,
	messageRepository messagesrepo.MessageRepositoryPort,
) *Service {
	return &Service{
		transactionManager: transactionManager,
		userRepository:     userRepository,
		habitRepository:    habitRepository,
		kaiRepository:      kaiRepository,
		messageRepository:  messageRepository,
	}
}

func (s *Service) SyncUserActivityState(
	ctx context.Context,
	userID uuid.UUID,
	force bool,
) error {

	return s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {
			now := time.Now()
			user, err := s.userRepository.FindUserSyncState(ctx, tx, userID)
			if err != nil {
				return err
			}
			if !force &&
				utils.IsSameCalendarDay(user.ActivitySyncAt, now) &&
				!hasInvalidInactiveDays(user.InactiveDays) {
				return nil
			}

			kaiState, err := s.kaiRepository.FindKaiStateByUserID(ctx, tx, userID)
			if err != nil {
				return err
			}

			if shouldShowReturnMessage(kaiState.LastInteraction, now) {
				message, err := s.messageRepository.FindRandomMessageByType(
					ctx,
					tx,
					userID,
					messagesmodel.MessageTypeReturn,
				)
				if err != nil {
					return err
				}

				if message != nil {
					if err := s.messageRepository.CreateUserMessage(ctx, tx, &messagesmodel.UserMessage{
						UserID:    userID,
						MessageID: message.ID,
						Read:      false,
						ShownAt:   now,
					}); err != nil {
						return err
					}
					if err := s.kaiRepository.UpdateLastMessage(ctx, tx, userID, message.Message); err != nil {
						return err
					}
				}

				if err := s.kaiRepository.UpdateLastInteraction(ctx, tx, userID, now); err != nil {
					return err
				}
			}

			lastActivity, err := s.findLastCompletedHabitDate(
				ctx,
				tx,
				userID,
			)
			if err != nil {
				return err
			}
			if lastActivity == nil {
				if err := s.userRepository.UpdateActivityStats(ctx, tx, userID, 0, 0); err != nil {
					return err
				}
				if err := s.kaiRepository.UpdateTemporalState(
					ctx,
					tx,
					userID,
					kaiState.Energy,
					s.determineKaiState(0, 0),
					s.determineKaiMode(0),
					0,
				); err != nil {
					return err
				}
				return s.userRepository.UpdateActivitySyncAt(ctx, tx, userID, now)
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

			return s.userRepository.UpdateActivitySyncAt(ctx, tx, userID, now)
		},
	)
}
