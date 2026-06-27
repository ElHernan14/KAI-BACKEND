package users

import (
	"context"
	habitDailyRecordsServicePort "kai-back/internal/modules/habits/service"
	usersdto "kai-back/internal/modules/users/dto"
	userRepository "kai-back/internal/modules/users/repository"
	xp "kai-back/internal/modules/xp/repository"
	userActivitySynchronizationService "kai-back/internal/services/user_activity_synchronization"
	authshared "kai-back/internal/shared/auth"
	errorHandler "kai-back/internal/shared/errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repository                         userRepository.UsersRepository
	xpRepo                             xp.XpRepository
	habitsDailyRecordsService          habitDailyRecordsServicePort.HabitsDailyRecordsServicePort
	UserActivitySynchronizationService userActivitySynchronizationService.ServicePort
}

func NewService(
	repositoryPort userRepository.UsersRepository,
	xpRepo xp.XpRepository,
	userActivitySyncService userActivitySynchronizationService.ServicePort,
	habitsDailyRecordsService habitDailyRecordsServicePort.HabitsDailyRecordsServicePort,
) *Service {
	return &Service{
		repository:                         repositoryPort,
		xpRepo:                             xpRepo,
		UserActivitySynchronizationService: userActivitySyncService,
		habitsDailyRecordsService:          habitsDailyRecordsService,
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

	//Construyo response final
	response := &usersdto.MeResponse{}

	//Agrego usuario
	response.User = &usersdto.UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		ProfileBase:  user.BaseProfile,
		ProfilePhoto: findProfilePhotoURL(userID),
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
		ReminderTime:           &user.Configuration.ReminderTime.String,
	}

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

func (s *Service) UpdateProfilePhoto(
	ctx context.Context,
	userID uuid.UUID,
) error {

	user, err := s.repository.FindUserByID(ctx, userID)
	if user == nil || err != nil {
		log.Println("error fetching user:", err)
		return errorHandler.NewAppError(
			http.StatusNotFound,
			"usuario no encontrado",
		)
	}

	return nil
}

func (s *Service) GetUserProfile(
	ctx context.Context,
	userID uuid.UUID,
) (*usersdto.UserProfileResponse, error) {

	if err := s.habitsDailyRecordsService.EnsureTodayHabitRecords(ctx, userID); err != nil {
		return nil, err
	}

	if err := s.UserActivitySynchronizationService.SyncUserActivityState(ctx, userID); err != nil {
		return nil, err
	}

	user, err := s.repository.FindUserProfileByID(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	totalXP, err := s.xpRepo.FindTotalUserXP(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	var userReminderTime *string
	if user.Configuration.ReminderTime.Valid {
		userReminderTime = &user.Configuration.ReminderTime.String
	}

	return &usersdto.UserProfileResponse{
		Usuario: usersdto.UserProfileDataResponse{
			ID:            user.ID.String(),
			Nombre:        user.Name,
			Email:         user.Email,
			Username:      user.Username,
			PerfilBase:    user.BaseProfile,
			EtapaKai:      user.KaiStage,
			FechaRegistro: user.RegisteredAt,
		},
		ConfiguracionUsuario: usersdto.UserConfigurationResponse{
			NotificacionesActivas:       user.Configuration.NotificationsEnabled,
			SonidosActivos:              user.Configuration.SoundsEnabled,
			MostrarRachas:               user.Configuration.ShowStreaks,
			ModoDiscreto:                user.Configuration.DiscreteMode,
			IntensidadKai:               user.Configuration.KaiIntensity,
			HorarioRecordatorio:         userReminderTime,
			BloquearConPIN:              user.Configuration.LockWithPIN,
			PermitirMensajesEmocionales: user.Configuration.AllowEmotionalMessages,
		},
		Resumen: usersdto.UserProfileSummary{
			XPTotal:      totalXP,
			RachaGlobal:  user.GlobalStreak,
			DiasInactivo: user.InactiveDays,
		},
	}, nil
}

func (s *Service) UpdateUserProfile(
	ctx context.Context,
	userID uuid.UUID,
	req usersdto.UpdateUserProfileRequest,
) error {

	updates := map[string]interface{}{}
	log.Printf("USUARIO: %v", req)

	if req.Nombre != nil {
		updates["nombre"] = *req.Nombre
	}

	if req.Username != nil {
		existingUsername, err := s.repository.FindUserByUsername(ctx, *req.Username)

		if err != nil {
			return err
		}

		if existingUsername != nil {
			return errorHandler.NewAppError(http.StatusConflict, "Username ya registrado.")
		}

		updates["username"] = *req.Username
	}

	if req.PerfilBase != nil {
		updates["perfil_base"] = *req.PerfilBase
	}

	return s.repository.UpdateUserProfile(
		ctx,
		nil,
		userID,
		updates,
	)
}

func (s *Service) UpdateUserConfiguration(
	ctx context.Context,
	userID uuid.UUID,
	req usersdto.UpdateUserConfigurationRequest,
) error {

	updates := map[string]interface{}{}

	if req.NotificacionesActivas != nil {
		updates["notificaciones_activas"] = *req.NotificacionesActivas
	}

	if req.SonidosActivos != nil {
		updates["sonidos_activos"] = *req.SonidosActivos
	}

	if req.MostrarRachas != nil {
		updates["mostrar_rachas"] = *req.MostrarRachas
	}

	if req.ModoDiscreto != nil {
		updates["modo_discreto"] = *req.ModoDiscreto
	}

	if req.IntensidadKai != nil {
		updates["intensidad_kai"] = *req.IntensidadKai
	}

	if req.HorarioRecordatorio != nil {
		if err := validateHorario(*req.HorarioRecordatorio); err != nil {
			return errorHandler.NewAppError(
				http.StatusBadRequest,
				"horario_recordatorio inválido, formato esperado HH:mm:ss",
			)
		}
		updates["horario_recordatorio"] = req.HorarioRecordatorio
	}

	if req.BloquearConPIN != nil {
		updates["bloquear_con_pin"] = *req.BloquearConPIN
	}

	if req.PermitirMensajesEmocionales != nil {
		updates["permitir_mensajes_emocionales"] = *req.PermitirMensajesEmocionales
	}

	return s.repository.UpdateUserConfiguration(
		ctx,
		nil,
		userID,
		updates,
	)
}

func validateHorario(horario string) error {
	_, err := time.Parse("15:04:05", horario)
	return err
}
