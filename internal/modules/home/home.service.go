package home

import (
	"context"
	"net/http"
	"time"

	habitsDailyRecordsService "kai-back/internal/modules/habits/service"
	homedto "kai-back/internal/modules/home/dto"
	repository "kai-back/internal/modules/home/repository"
	kaievolution "kai-back/internal/modules/kai/evolution"
	kairepo "kai-back/internal/modules/kai/repository"
	messagesmodel "kai-back/internal/modules/messages/models"
	messagesrepo "kai-back/internal/modules/messages/repository"
	userrepo "kai-back/internal/modules/users/repository"
	userActivitySynchronizationService "kai-back/internal/services/user_activity_synchronization"
	errorHandler "kai-back/internal/shared/errors"

	"github.com/google/uuid"
)

type ServicePort interface {
	GetHome(ctx context.Context, userID uuid.UUID) (*homedto.HomeResponse, error)
}

type Service struct {
	repository                         repository.HomeRepository
	UserRepository                     userrepo.UsersRepository
	KaiRepository                      kairepo.KaiRepository
	MessageRepository                  messagesrepo.MessageRepositoryPort
	habitsDailyRecordsService          habitsDailyRecordsService.HabitsDailyRecordsServicePort
	UserActivitySynchronizationService userActivitySynchronizationService.ServicePort
}

func NewService(
	repository repository.HomeRepository,
	userRepo userrepo.UsersRepository,
	kaiRepo kairepo.KaiRepository,
	messageRepo messagesrepo.MessageRepositoryPort,
	userActivitySyncService userActivitySynchronizationService.ServicePort,
	habitsDailyRecordsService habitsDailyRecordsService.HabitsDailyRecordsServicePort,
) *Service {
	return &Service{
		repository:                         repository,
		UserRepository:                     userRepo,
		KaiRepository:                      kaiRepo,
		MessageRepository:                  messageRepo,
		UserActivitySynchronizationService: userActivitySyncService,
		habitsDailyRecordsService:          habitsDailyRecordsService,
	}
}

func (s *Service) GetHome(ctx context.Context, userID uuid.UUID) (*homedto.HomeResponse, error) {
	if err := s.habitsDailyRecordsService.EnsureTodayHabitRecords(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.UserActivitySynchronizationService.SyncUserActivityState(ctx, userID, false); err != nil {
		return nil, err
	}

	kai, err := s.repository.FindKaiSummary(ctx, userID)
	if err != nil {
		return nil, err
	}
	if kai == nil {
		return nil, errorHandler.NewAppError(http.StatusNotFound, "estado kai no encontrado")
	}

	totalXP, err := s.repository.FindTotalXP(ctx, userID)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		return nil, err
	}

	dailyHabits, err := s.repository.FindDailyHabits(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	evolutionActive, evolutionExpiresAt := kaievolution.EventWindow(kai.LastEvolution, now)
	message, err := s.resolveHomeMessage(ctx, userID, kai, evolutionActive, now)
	if err != nil {
		return nil, err
	}

	habitsResponse := make([]homedto.DailyHabitSummary, 0, len(dailyHabits))
	completed := 0
	for _, habit := range dailyHabits {
		if habit.Completed {
			completed++
		}
		habitsResponse = append(habitsResponse, homedto.DailyHabitSummary{
			ID:        habit.ID.String(),
			Name:      habit.Name,
			Category:  habit.Category,
			Completed: habit.Completed,
		})
	}

	total := len(habitsResponse)
	return &homedto.HomeResponse{
		Kai: homedto.KaiHomeSummary{
			CurrentState: kai.CurrentState,
			CurrentStage: kai.CurrentStage,
			Energy:       kai.Energy,
			BondLevel:    kai.BondLevel,
			RecoveryMode: kai.RecoveryMode,
			KaiImage:     kai.KaiImage,
		},
		Message: message,
		EvolutionEvent: homedto.EvolutionEventResponse{
			Active:    evolutionActive,
			Stage:     kai.CurrentStage,
			StartedAt: kai.LastEvolution,
			ExpiresAt: evolutionExpiresAt,
			Message:   evolutionMessage(evolutionActive, message),
		},
		TotalXP:       totalXP,
		CurrentStreak: user.GlobalStreak,
		HabitsToday:   habitsResponse,
		DailyProgress: homedto.DailyProgress{
			Total:     total,
			Completed: completed,
			Pending:   total - completed,
		},
	}, nil
}

func (s *Service) resolveHomeMessage(
	ctx context.Context,
	userID uuid.UUID,
	kai *repository.KaiSummaryRow,
	evolutionActive bool,
	now time.Time,
) (string, error) {
	if evolutionActive {
		messageType := evolutionMessageType(kai.CurrentStage)
		lastMessage, err := s.MessageRepository.FindLastUserMessage(ctx, nil, userID)
		if err != nil {
			return "", err
		}
		if lastMessage != nil && lastMessage.KaiMessage != nil &&
			lastMessage.KaiMessage.Type == messageType {
			return lastMessage.KaiMessage.Message, nil
		}

		return s.selectAndStoreMessage(ctx, userID, messageType, false, now, kai.LastMessage)
	}

	hasMessages, err := s.MessageRepository.HasUserMessages(ctx, userID)
	if err != nil {
		return "", err
	}
	if !hasMessages {
		return s.selectAndStoreMessage(
			ctx, userID, messagesmodel.MessageTypeWelcome, true, now, kai.LastMessage,
		)
	}

	if isFirstInteractionToday(kai.LastInteraction, now) {
		return s.selectAndStoreMessage(
			ctx, userID, messagesmodel.MessageTypeGreeting, true, now, kai.LastMessage,
		)
	}

	return s.resolveMessage(ctx, kai.LastMessage, false)
}

func (s *Service) selectAndStoreMessage(
	ctx context.Context,
	userID uuid.UUID,
	messageType string,
	updateInteraction bool,
	now time.Time,
	currentMessage *string,
) (string, error) {
	message, err := s.MessageRepository.FindRandomMessageByType(ctx, nil, userID, messageType)
	if err != nil {
		return "", err
	}

	if message != nil {
		if err := s.MessageRepository.CreateUserMessage(ctx, nil, &messagesmodel.UserMessage{
			UserID:    userID,
			MessageID: message.ID,
			Read:      false,
			ShownAt:   now,
		}); err != nil {
			return "", err
		}
		if err := s.KaiRepository.UpdateLastMessage(ctx, nil, userID, message.Message); err != nil {
			return "", err
		}
	}

	if updateInteraction {
		if err := s.KaiRepository.UpdateLastInteraction(ctx, nil, userID, now); err != nil {
			return "", err
		}
	}

	if message != nil {
		return message.Message, nil
	}
	return s.resolveMessage(ctx, currentMessage, false)
}

func isFirstInteractionToday(lastInteraction *time.Time, now time.Time) bool {
	if lastInteraction == nil {
		return true
	}

	year, month, day := lastInteraction.In(now.Location()).Date()
	nowYear, nowMonth, nowDay := now.Date()
	return year != nowYear || month != nowMonth || day != nowDay
}

func evolutionMessageType(stage string) string {
	if kaievolution.NormalizeStage(stage) == "adulto" {
		return messagesmodel.MessageTypeEvolutionAdult
	}
	return messagesmodel.MessageTypeEvolutionYoung
}

func evolutionMessage(active bool, message string) string {
	if !active {
		return ""
	}
	return message
}

func (s *Service) resolveMessage(
	ctx context.Context,
	kaiMessage *string,
	evolutionActive bool,
) (string, error) {
	if evolutionActive {
		message, err := s.repository.FindRandomEvolutionMessage(ctx)
		if err != nil {
			return "", err
		}
		if message != nil && *message != "" {
			return *message, nil
		}
	}

	if kaiMessage != nil && *kaiMessage != "" {
		return *kaiMessage, nil
	}

	message, err := s.repository.FindFallbackMessage(ctx)
	if err != nil {
		return "", err
	}
	if message != nil && *message != "" {
		return *message, nil
	}

	return "", nil
}
