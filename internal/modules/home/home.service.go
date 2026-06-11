package home

import (
	"context"
	"net/http"

	habitsDailyRecordsService "kai-back/internal/modules/habits/service"
	homedto "kai-back/internal/modules/home/dto"
	repository "kai-back/internal/modules/home/repository"
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
	habitsDailyRecordsService          habitsDailyRecordsService.HabitsDailyRecordsServicePort
	UserActivitySynchronizationService userActivitySynchronizationService.ServicePort
}

func NewService(
	repository repository.HomeRepository,
	userRepo userrepo.UsersRepository,
	userActivitySyncService userActivitySynchronizationService.ServicePort,
	habitsDailyRecordsService habitsDailyRecordsService.HabitsDailyRecordsServicePort,
) *Service {
	return &Service{
		repository:                         repository,
		UserRepository:                     userRepo,
		UserActivitySynchronizationService: userActivitySyncService,
		habitsDailyRecordsService:          habitsDailyRecordsService,
	}
}

func (s *Service) GetHome(ctx context.Context, userID uuid.UUID) (*homedto.HomeResponse, error) {

	// Aseguramos que existan los registros diarios para hoy antes de obtener la información del home
	err := s.habitsDailyRecordsService.
		EnsureTodayHabitRecords(
			ctx,
			userID,
		)
	if err != nil {
		return nil, err
	}

	// Aseguramos Sincronizar toda la información temporal del usuario dependiente del paso del tiempo y de su actividad reciente.
	err = s.UserActivitySynchronizationService.SyncUserActivityState(
		ctx,
		userID,
	)
	if err != nil {
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

	// currentStreak, err := s.repository.FindCurrentStreak(ctx, userID)
	// if err != nil {
	// 	return nil, err
	// }

	user, err := s.UserRepository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		return nil, err
	}

	dailyHabits, err := s.repository.FindDailyHabits(ctx, userID)
	if err != nil {
		return nil, err
	}

	message, err := s.resolveMessage(ctx, kai.LastMessage)
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
		Message:       message,
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

func (s *Service) resolveMessage(ctx context.Context, kaiMessage *string) (string, error) {
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
