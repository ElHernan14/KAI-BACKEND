package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"

	"kai-back/internal/config"
	authdto "kai-back/internal/modules/auth/dto"
	authmodel "kai-back/internal/modules/auth/models"
	authrepository "kai-back/internal/modules/auth/repository"
	usersmodel "kai-back/internal/modules/users/models"
	userrepository "kai-back/internal/modules/users/repository"
	initializerUser "kai-back/internal/services/user_initializer"
	authshared "kai-back/internal/shared/auth"
	errorHandler "kai-back/internal/shared/errors"
	mailService "kai-back/internal/shared/mail"
	"kai-back/internal/shared/security"
	transaction "kai-back/internal/shared/transaction"
)

type Service struct {
	repository         authrepository.AuthRepository
	userRepository     userrepository.UsersRepository
	initializerUser    initializerUser.Service
	transactionManager transaction.TransactionManager
	jwtSecret          string
	jwtTTL             time.Duration
	googleClientID     string
	mailService        mailService.MailService
	config             config.Config
}

func NewService(
	repository authrepository.AuthRepository,
	userRepository userrepository.UsersRepository,
	initializerUser initializerUser.Service,
	transactionManager transaction.TransactionManager,
	jwtSecret string,
	jwtTTL time.Duration,
	googleClientID string,
	mailServ mailService.MailService,
	config config.Config,
) *Service {
	return &Service{
		repository:         repository,
		userRepository:     userRepository,
		initializerUser:    initializerUser,
		transactionManager: transactionManager,
		jwtSecret:          jwtSecret,
		jwtTTL:             jwtTTL,
		googleClientID:     googleClientID,
		mailService:        mailServ,
		config:             config,
	}
}

func (s *Service) Register(ctx context.Context, req authdto.RegisterRequest) (*authdto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	username := normalizeUsername(req.Username)

	existingUser, err := s.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		log.Println("error getting user by email:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, fmt.Sprintf("No se pudo obtener usuario con email: %s", email))
	}
	if existingUser != nil {
		return nil, errorHandler.NewAppError(http.StatusConflict, "Email ya registrado.")
	}

	existingUsername, err := s.userRepository.FindUserByUsername(ctx, username)
	if err != nil {
		log.Println("error getting user by username:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, fmt.Sprintf("No se pudo obtener usuario con username: %s", username))
	}
	if existingUsername != nil {
		return nil, errorHandler.NewAppError(http.StatusConflict, "Username ya registrado.")
	}

	passwordHash, err := authshared.HashPassword(req.Password)
	if err != nil {
		log.Println("error hashing password user:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo registrar usuario.")
	}

	user := &usersmodel.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		Username:     &username,
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
	identifier := normalizeIdentifier(req.Identifier)
	if identifier == "" {
		identifier = normalizeIdentifier(req.Email)
	}
	if identifier == "" {
		return nil, errorHandler.NewAppError(http.StatusBadRequest, "email o username es requerido")
	}

	user, err := s.userRepository.FindUserByEmailOrUsername(ctx, identifier)
	if err != nil {
		log.Println("error getting user by login identifier:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo obtener usuario")
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
		UserResponse: authdto.AuthUserResponse{
			Email:    user.Email,
			Name:     user.Name,
			Username: user.Username,
		},
		Token: token,
	}, nil
}

func (s *Service) RenewToken(ctx context.Context, userID string, email string) (*authdto.ValidateTokenResponse, error) {
	token, err := authshared.GenerateToken(userID, email, s.jwtSecret, s.jwtTTL)
	if err != nil {
		log.Println("error renewing user token:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo renovar token")
	}

	normalizedEmail := normalizeEmail(email)

	user, err := s.userRepository.FindUserByEmail(ctx, normalizedEmail)
	if err != nil {
		log.Println("error getting user by email:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, fmt.Sprintf("No se pudo obtener usuario con email: %s", email))
	}
	if user == nil {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Credenciales invalidas")
	}

	return &authdto.ValidateTokenResponse{
		Valid: true,
		UserResponse: authdto.AuthUserResponse{
			Email:    user.Email,
			Name:     user.Name,
			Username: user.Username,
		},
		Token: token,
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func normalizeIdentifier(identifier string) string {
	return strings.ToLower(strings.TrimSpace(identifier))
}

func (s *Service) GoogleLogin(ctx context.Context, req authdto.GoogleLoginRequest) (*authdto.AuthResponse, error) {
	payload, err := idtoken.Validate(
		ctx,
		req.IDToken,
		s.googleClientID,
	)

	log.Println("payload", payload)

	if err != nil {
		log.Println("invalid google token:", err)
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Token Google inválido")
	}

	email, ok := payload.Claims["email"].(string)

	if !ok {
		return nil, errorHandler.NewAppError(http.StatusUnauthorized, "Email no encontrado")
	}

	name, _ := payload.Claims["name"].(string)

	email = normalizeEmail(email)

	user, err := s.userRepository.FindUserByEmail(
		ctx,
		email,
	)

	if err != nil {
		log.Println("error finding user:", err)
		return nil, errorHandler.NewAppError(http.StatusInternalServerError, "Error obteniendo usuario")
	}

	if user == nil {

		user = &usersmodel.User{
			Name:         name,
			Email:        email,
			KaiStage:     "cachorro",
			GlobalStreak: 0,
			InactiveDays: 0,
		}

		if err := s.initializerUser.InitializeNewUser(
			ctx,
			user,
		); err != nil {

			return nil, err
		}
	}

	return s.buildAuthResponse(user)
}

func (s *Service) ForgotPassword(
	ctx context.Context,
	req authdto.ForgotPasswordRequest,
) error {

	user, err := s.repository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		// No revelamos si existe o no.
		return nil
	}

	code, err := security.GenerateResetCode()
	if err != nil {
		return err
	}

	codeHash := security.HashResetCode(code)

	resetCode := &authmodel.PasswordResetCode{
		UserID:    user.ID,
		CodeHash:  codeHash,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if err := s.repository.CreatePasswordResetCode(ctx, resetCode); err != nil {
		return err
	}

	return s.mailService.SendPasswordResetEmail(
		user.Email,
		code,
	)
}

func (s *Service) ResetPassword(
	ctx context.Context,
	req authdto.ResetPasswordRequest,
) error {

	user, err := s.repository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return errorHandler.NewAppError(
			http.StatusBadRequest,
			"codigo invalido o expirado",
		)
	}

	codeHash := security.HashResetCode(req.Code)

	return s.transactionManager.WithTransaction(ctx, func(tx *gorm.DB) error {

		resetCode, err := s.repository.FindValidPasswordResetCode(
			ctx,
			user.ID,
			codeHash,
		)
		if err != nil {
			return errorHandler.NewAppError(
				http.StatusBadRequest,
				"codigo invalido o expirado",
			)
		}

		passwordHash, err := authshared.HashPassword(req.Password)
		if err != nil {
			return err
		}

		if err := s.repository.UpdateUserPassword(
			ctx,
			tx,
			user.ID,
			passwordHash,
		); err != nil {
			return err
		}

		return s.repository.MarkPasswordResetCodeUsed(
			ctx,
			tx,
			resetCode.ID,
		)
	})
}
