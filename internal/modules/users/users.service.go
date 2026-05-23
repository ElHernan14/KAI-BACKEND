package users

import (
	"context"
	kaidto "kai-back/internal/modules/kai/dto"
	usersdto "kai-back/internal/modules/users/dto"
	userRepository "kai-back/internal/modules/users/repository"
	errorHandler "kai-back/internal/shared/errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repository userRepository.UsersRepository
}

func NewService(repositoryPort userRepository.UsersRepository) *Service {
	return &Service{
		repository: repositoryPort,
	}
}

func (s *Service) GetMe(
	ctx context.Context,
	userID uuid.UUID,
) (*usersdto.MeResponse, error) {

	user, err := s.repository.FindMeByID(ctx, userID)
	if err != nil {
		log.Println("error fetching user:", err)
		return nil, errorHandler.NewAppError(http.StatusNotFound, "No se encontró el usuario")
	}

	response := &usersdto.MeResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		ProfileBase:  user.BaseProfile,
		KaiStage:     user.KaiStage,
		GlobalStreak: user.GlobalStreak,
		InactiveDays: user.InactiveDays,

		KaiState: &kaidto.KaiStateSummary{
			CurrentState: user.KaiState.CurrentState,
			CurrentStage: user.KaiState.CurrentStage,
			Energy:       user.KaiState.Energy,
			BondLevel:    user.KaiState.BondLevel,
			RecoveryMode: user.KaiState.RecoveryMode,
		},

		Configuration: &usersdto.UserConfigResponse{
			NotificationsEnabled:   user.Configuration.NotificationsEnabled,
			SoundsEnabled:          user.Configuration.SoundsEnabled,
			ShowStreaks:            user.Configuration.ShowStreaks,
			DiscreteMode:           user.Configuration.DiscreteMode,
			KaiIntensity:           user.Configuration.KaiIntensity,
			LockWithPIN:            user.Configuration.LockWithPIN,
			AllowEmotionalMessages: user.Configuration.AllowEmotionalMessages,
		},
	}

	if user.Configuration.ReminderTime != nil {
		formatted := user.Configuration.ReminderTime.Format("15:04")
		response.Configuration.ReminderTime = &formatted
	}

	return response, nil
}

func (s *Service) UpdateMe(
	ctx context.Context,
	userID uuid.UUID,
	req usersdto.UpdateMeRequest,
) (*usersdto.MeResponse, error) {

	user, err := s.repository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		log.Println("error fetching user:", err)
		return nil, errorHandler.NewAppError(
			http.StatusNotFound,
			"usuario no encontrado",
		)
	}

	user.Name = strings.TrimSpace(req.Name)

	if req.ProfileBase != nil {
		profile := strings.TrimSpace(*req.ProfileBase)
		user.BaseProfile = &profile
	}

	if err := s.repository.Update(ctx, user); err != nil {
		log.Println("error updating user:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo actualizar usuario.")
	}

	return s.GetMe(ctx, userID)
}
