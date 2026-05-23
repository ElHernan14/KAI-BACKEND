package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	authdto "kai-back/internal/modules/auth/dto"
	authrepository "kai-back/internal/modules/auth/repository"
	usersmodel "kai-back/internal/modules/users/models"
	userrepository "kai-back/internal/modules/users/repository"
	initializerUser "kai-back/internal/services/user_initializer"
	authshared "kai-back/internal/shared/auth"
	errorHandler "kai-back/internal/shared/errors"
)

type Service struct {
	repository      authrepository.AuthRepository
	userRepository  userrepository.UsersRepository
	initializerUser initializerUser.Service
	jwtSecret       string
	jwtTTL          time.Duration
}

func NewService(
	repository authrepository.AuthRepository,
	userRepository userrepository.UsersRepository,
	initializerUser initializerUser.Service,
	jwtSecret string,
	jwtTTL time.Duration,
) *Service {
	return &Service{
		repository:      repository,
		userRepository:  userRepository,
		initializerUser: initializerUser,
		jwtSecret:       jwtSecret,
		jwtTTL:          jwtTTL,
	}
}

func (s *Service) Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	existingUser, err := s.userRepository.FindUserByEmail(ctx, email)
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

	if err := s.initializerUser.InitializeNewUser(ctx, user); err != nil {
		return nil, err
	}

	return s.buildAuthResponse(user)
}

func (s *Service) Login(ctx context.Context, req authdto.LoginRequest) (*authdto.AuthResponse, error) {
	email := normalizeEmail(req.Email)

	user, err := s.userRepository.FindUserByEmail(ctx, email)
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
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
