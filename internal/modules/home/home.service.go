package home

import (
	"context"
	"net/http"

	homedto "kai-back/internal/modules/home/dto"
	repository "kai-back/internal/modules/home/repository"
	errorHandler "kai-back/internal/shared/errors"

	"github.com/google/uuid"
)

type ServicePort interface {
	GetHome(ctx context.Context, userID uuid.UUID) (*homedto.HomeResponse, error)
}

type Service struct {
	repository repository.HomeRepository
}

func NewService(repository repository.HomeRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetHome(ctx context.Context, userID uuid.UUID) (*homedto.HomeResponse, error) {
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

	currentStreak, err := s.repository.FindCurrentStreak(ctx, userID)
	if err != nil {
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
		CurrentStreak: currentStreak,
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
