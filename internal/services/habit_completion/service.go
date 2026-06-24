package habitcompletion

import (
	"context"

	habitsdto "kai-back/internal/modules/habits/dto"
	habit "kai-back/internal/modules/habits/repository"
	habitDailyRecordsServicePort "kai-back/internal/modules/habits/service"
	kai "kai-back/internal/modules/kai/repository"
	message "kai-back/internal/modules/messages/repository"
	xp "kai-back/internal/modules/xp/repository"
	userActivitySynchronizationService "kai-back/internal/services/user_activity_synchronization"
	transaction "kai-back/internal/shared/transaction"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	habitRepo                          habit.HabitsRepository
	xpRepo                             xp.XpRepository
	kaiRepo                            kai.KaiRepository
	messageRepo                        message.MessageRepositoryPort
	habitsDailyRecordsService          habitDailyRecordsServicePort.HabitsDailyRecordsServicePort
	UserActivitySynchronizationService userActivitySynchronizationService.ServicePort
	transactionManager                 transaction.TransactionManager
}

func NewService(
	habitRepo habit.HabitsRepository,
	xpRepo xp.XpRepository,
	kaiRepo kai.KaiRepository,
	messageRepo message.MessageRepositoryPort,
	userActivitySyncService userActivitySynchronizationService.ServicePort,
	habitsDailyRecordsService habitDailyRecordsServicePort.HabitsDailyRecordsServicePort,
	transactionManager transaction.TransactionManager,
) *Service {
	return &Service{
		habitRepo:                          habitRepo,
		xpRepo:                             xpRepo,
		kaiRepo:                            kaiRepo,
		messageRepo:                        messageRepo,
		UserActivitySynchronizationService: userActivitySyncService,
		habitsDailyRecordsService:          habitsDailyRecordsService,
		transactionManager:                 transactionManager,
	}
}

func (s *Service) CompleteHabit(
	ctx context.Context,
	userID uuid.UUID,
	userHabitID uuid.UUID,
	value *string,
) (*habitsdto.CompleteHabitResponse, error) {
	var result *habitsdto.CompleteHabitResponse

	err := s.habitsDailyRecordsService.EnsureTodayHabitRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.transactionManager.WithTransaction(
		ctx,
		func(tx *gorm.DB) error {
			habit, rec, err := s.validateHabit(ctx, tx, userID, userHabitID)
			if err != nil {
				return err
			}

			if _, err = s.completeRecord(ctx, tx, habit, value, rec); err != nil {
				return err
			}

			streak, err := s.updateStreak(ctx, tx, userID, habit)
			if err != nil {
				return err
			}

			xpGranted, err := s.grantXP(ctx, tx, userID, habit)
			if err != nil {
				return err
			}

			if err = s.grantKaiAttributes(ctx, tx, userID, habit, xpGranted); err != nil {
				return err
			}

			dominantAttribute, err := s.recalculateDominantAttribute(ctx, tx, userID)
			if err != nil {
				return err
			}

			evolution, err := s.updateKaiState(
				ctx,
				tx,
				userID,
				habit,
				streak,
				dominantAttribute,
			)
			if err != nil {
				return err
			}

			if err = s.generateMotivationalMessage(ctx, tx, userID, streak, dominantAttribute); err != nil {
				return err
			}

			currentStreak := 0
			if streak != nil {
				currentStreak = streak.CurrentDays
			}

			result = &habitsdto.CompleteHabitResponse{
				HabitoUsuarioID:     userHabitID,
				XPGanada:            xpGranted,
				RachaActual:         currentStreak,
				EnergiaActual:       evolution.State.Energy,
				NivelVinculo:        evolution.State.BondLevel,
				AtributoDominanteID: evolution.State.DominantAttributeID,
				Evoluciono:          evolution.Evolved,
				EtapaAnterior:       evolution.PreviousStage,
				EtapaActual:         evolution.State.CurrentStage,
				EventoEvolucion: habitsdto.EvolutionEventResponse{
					Activo:     evolution.EventActive,
					Etapa:      evolution.State.CurrentStage,
					IniciadoEn: evolution.State.LastEvolution,
					ExpiraEn:   evolution.EventExpiresAt,
				},
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	if err = s.UserActivitySynchronizationService.SyncUserActivityState(ctx, userID); err != nil {
		return nil, err
	}

	return result, nil
}
