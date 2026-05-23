package users

import (
	"context"
	kaidto "kai-back/internal/modules/kai/dto"
	usersdto "kai-back/internal/modules/users/dto"
	userRepository "kai-back/internal/modules/users/repository"
	xpdto "kai-back/internal/modules/xp/dto"
	authshared "kai-back/internal/shared/auth"
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

	//Construyo XP response
	xpResponse := make([]xpdto.UserXPSummary, 0)

	for _, xp := range user.XPCategory {

		categoryName := ""

		if xp.Category != nil {
			categoryName = xp.Category.Name
		}

		xpResponse = append(xpResponse, xpdto.UserXPSummary{
			CategoryID:   xp.CategoryID.String(),
			CategoryName: categoryName,
			Value:        xp.Value,
		})
	}

	//Construyo Attributes Kai response
	attributesResponse := make([]kaidto.KaiAttributeSummary, 0)

	for _, attribute := range user.KaiAttributes {

		attributeName := ""

		if attribute.AttributeType != nil {
			attributeName = attribute.AttributeType.Name
		}

		attributesResponse = append(attributesResponse, kaidto.KaiAttributeSummary{
			AttributeID:   attribute.AttributeTypeID.String(),
			AttributeName: attributeName,
			Value:         attribute.Value,
		})
	}

	//Construyo response final
	response := &usersdto.MeResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		ProfileBase:  user.BaseProfile,
		KaiStage:     user.KaiStage,
		GlobalStreak: user.GlobalStreak,
		InactiveDays: user.InactiveDays,

		XP:         xpResponse,
		Attributes: attributesResponse,
	}

	//Agrego atributo dominante de Kai state
	var dominantAttribute *kaidto.DominantAttributeResponse = &kaidto.DominantAttributeResponse{
		ID:   uuid.Nil,
		Name: "",
	}
	//Agrego referencia de currentMode
	currentMode := ""
	if user.KaiState.CurrentMode != nil || user.KaiState.DominantAttribute != nil {
		if user.KaiState.CurrentMode != nil {
			currentMode = *user.KaiState.CurrentMode
		}
		if user.KaiState.DominantAttribute != nil {
			dominantAttribute = &kaidto.DominantAttributeResponse{
				ID:   user.KaiState.DominantAttribute.ID,
				Name: user.KaiState.DominantAttribute.Name,
			}
		}
	}

	//Agrego estado kai
	response.KaiState = &kaidto.KaiStateSummary{
		CurrentState:      user.KaiState.CurrentState,
		CurrentStage:      user.KaiState.CurrentStage,
		CurrentMode:       currentMode,
		Energy:            user.KaiState.Energy,
		BondLevel:         user.KaiState.BondLevel,
		RecoveryMode:      user.KaiState.RecoveryMode,
		DominantAttribute: dominantAttribute,
		LastEvolution:     user.KaiState.LastEvolution,
	}

	//Agrego user_configuracion
	response.Configuration = &usersdto.UserConfigResponse{
		NotificationsEnabled:   user.Configuration.NotificationsEnabled,
		SoundsEnabled:          user.Configuration.SoundsEnabled,
		ShowStreaks:            user.Configuration.ShowStreaks,
		DiscreteMode:           user.Configuration.DiscreteMode,
		KaiIntensity:           user.Configuration.KaiIntensity,
		LockWithPIN:            user.Configuration.LockWithPIN,
		AllowEmotionalMessages: user.Configuration.AllowEmotionalMessages,
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

func (s *Service) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	req usersdto.ChangePasswordRequest,
) error {

	user, err := s.repository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		log.Println("error fetching user:", err)
		return errorHandler.NewAppError(
			http.StatusNotFound,
			"usuario no encontrado",
		)
	}

	if err := authshared.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return errorHandler.NewAppError(http.StatusUnauthorized, "Contraseña actual incorrecta")
	}

	if req.CurrentPassword == req.NewPassword {
		return errorHandler.NewAppError(
			http.StatusBadRequest,
			"la nueva contraseña no puede ser igual a la actual",
		)
	}

	newPasswordHash, err := authshared.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("error hashing password user:", err)
		return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo modificar contraseña.")
	}

	return s.repository.UpdatePassword(
		ctx,
		userID,
		newPasswordHash,
	)
}
