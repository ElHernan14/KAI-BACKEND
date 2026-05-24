package users

import (
	"context"
	usersdto "kai-back/internal/modules/users/dto"
	userRepository "kai-back/internal/modules/users/repository"
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
	// xpResponse := make([]xpdto.UserXPSummary, 0)

	// for _, xp := range user.XPCategory {

	// 	categoryName := ""

	// 	if xp.Category != nil {
	// 		categoryName = xp.Category.Name
	// 	}

	// 	xpResponse = append(xpResponse, xpdto.UserXPSummary{
	// 		CategoryID:   xp.CategoryID.String(),
	// 		CategoryName: categoryName,
	// 		Value:        xp.Value,
	// 	})
	// }

	// //Construyo Attributes Kai response
	// attributesResponse := make([]kaidto.KaiAttributeSummary, 0)

	// for _, attribute := range user.KaiAttributes {

	// 	attributeName := ""

	// 	if attribute.AttributeType != nil {
	// 		attributeName = attribute.AttributeType.Name
	// 	}

	// 	attributesResponse = append(attributesResponse, kaidto.KaiAttributeSummary{
	// 		AttributeID:   attribute.AttributeTypeID.String(),
	// 		AttributeName: attributeName,
	// 		Value:         attribute.Value,
	// 	})
	// }

	//Construyo response final
	response := &usersdto.MeResponse{}

	//Agrego usuario
	response.User = &usersdto.UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		ProfileBase:  user.BaseProfile,
		KaiStage:     user.KaiStage,
		GlobalStreak: user.GlobalStreak,
		InactiveDays: user.InactiveDays,
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

	//Construye los habitos del usuario
	// var habitsResponse []habitsdto.UserHabitResponse

	// for _, habit := range user.UserHabits {

	// 	var records []habitsdto.HabitRecordResponse

	// 	for _, record := range habit.HabitRecords {
	// 		records = append(records, habitsdto.HabitRecordResponse{
	// 			Fecha:           record.Fecha,
	// 			Completado:      record.Completado,
	// 			ValorRegistrado: record.ValorRegistrado,
	// 			XPGanada:        record.XPGanada,
	// 		})
	// 	}

	// 	habitsResponse = append(habitsResponse, habitsdto.UserHabitResponse{
	// 		ID:          habit.ID.String(),
	// 		Name:        habit.HabitCatalog.Name,
	// 		Description: habit.HabitCatalog.Description,

	// 		Category:   habit.HabitCatalog.Category,
	// 		CareType:   habit.HabitCatalog.CareType,
	// 		Difficulty: habit.HabitCatalog.Difficulty,

	// 		BaseXP: habit.HabitCatalog.BaseXP,

	// 		Active: habit.Active,

	// 		HabitImage: habit.HabitCatalog.HabitImage,

	// 		Records: records,
	// 	})
	// }

	return response, nil
}

func (s *Service) UpdateMe(
	ctx context.Context,
	userID uuid.UUID,
	req usersdto.UpdateMeRequest,
) error {

	user, err := s.repository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		log.Println("error fetching user:", err)
		return errorHandler.NewAppError(
			http.StatusNotFound,
			"usuario no encontrado",
		)
	}

	user.Name = strings.TrimSpace(req.Name)

	if req.ProfileBase != nil {
		profile := strings.TrimSpace(*req.ProfileBase)
		user.BaseProfile = &profile
	}

	return s.repository.Update(ctx, user)
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
