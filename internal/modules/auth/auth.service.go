package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	authdto "kai-back/internal/modules/auth/dto"
	kaimodel "kai-back/internal/modules/kai/models"
	usersmodel "kai-back/internal/modules/users/models"
	authshared "kai-back/internal/shared/auth"
	errorHandler "kai-back/internal/shared/errors"
)

type RepositoryPort interface {
	FindUserByEmail(ctx context.Context, email string) (*usersmodel.User, error)
	RegisterUserWithInitialState(
		ctx context.Context,
		user *usersmodel.User,
		config *usersmodel.UserConfiguration,
		kaiState *kaimodel.KaiState,
	) error
}

type Service struct {
	repository RepositoryPort
	jwtSecret  string
	jwtTTL     time.Duration
}

func NewService(repository RepositoryPort, jwtSecret string, jwtTTL time.Duration) *Service {
	return &Service{
		repository: repository,
		jwtSecret:  jwtSecret,
		jwtTTL:     jwtTTL,
	}
}

func (s *Service) Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	existingUser, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errorHandler.NewAppError(http.StatusConflict, "email ya registrado")
	}

	passwordHash, err := authshared.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &usersmodel.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: passwordHash,
		KaiStage:     "cachorro",
		GlobalStreak: 0,
		InactiveDays: 0,
	}

	config := &usersmodel.UserConfiguration{
		NotificationsEnabled:   true,
		SoundsEnabled:          true,
		ShowStreaks:            true,
		DiscreteMode:           false,
		KaiIntensity:           "normal",
		LockWithPIN:            false,
		AllowEmotionalMessages: true,
	}

	kaiState := &kaimodel.KaiState{
		CurrentState:        "con_ganas_de_empezar",
		CurrentStage:        "cachorro",
		Energy:              100,
		DaysWithoutActivity: 0,
		RecoveryMode:        false,
		BondLevel:           1,
	}

	if err := s.repository.RegisterUserWithInitialState(ctx, user, config, kaiState); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *Service) Login(ctx context.Context, req authdto.LoginRequest) (*authdto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "credenciales invalidas")
	}

	if err := authshared.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "credenciales invalidas")
	}

	return s.buildAuthResponse(user)
}

func (s *Service) buildAuthResponse(user *usersmodel.User) (*authdto.AuthResponse, error) {
	token, err := authshared.GenerateToken(user.ID.String(), user.Email, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return nil, err
	}

	return &authdto.AuthResponse{
		Token: token,
		User: authdto.AuthUserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
