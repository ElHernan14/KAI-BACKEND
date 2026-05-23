package auth

import (
	"context"
	"fmt"
	"log"
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
		log.Println("error getting user by email:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, fmt.Sprintf("No se pudo obtener usuario con email: %s", email))
	}
	if existingUser != nil {
		return nil, errorHandler.NewAppError(http.StatusConflict, "Email ya registrado.")
	}

	passwordHash, err := authshared.HashPassword(req.Password)
	if err != nil {
		log.Println("error hashing password user:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo registrar usuario.")
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
		log.Println("error getting user by email:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, fmt.Sprintf("No se pudo obtener usuario con email: %s", email))
	}
	if user == nil {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Credenciales invalidas")
	}

	if err := authshared.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Credenciales invalidas")
	}

	return s.buildAuthResponse(user)
}

func (s *Service) buildAuthResponse(user *usersmodel.User) (*authdto.AuthResponse, error) {
	token, err := authshared.GenerateToken(user.ID.String(), user.Email, s.jwtSecret, s.jwtTTL)
	if err != nil {
		log.Println("error generating user token:", err)
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Credenciales invalidas")
	}

	return &authdto.AuthResponse{
		Token: token,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
