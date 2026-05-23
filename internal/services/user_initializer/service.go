package userinitializer

import (
	"context"
	kaimodel "kai-back/internal/modules/kai/models"
	kairepository "kai-back/internal/modules/kai/repository"
	usersmodel "kai-back/internal/modules/users/models"
	usersrepository "kai-back/internal/modules/users/repository"
	xpmodel "kai-back/internal/modules/xp/models"
	xprepository "kai-back/internal/modules/xp/repository"
	errorHandler "kai-back/internal/shared/errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	db *gorm.DB

	userRepository usersrepository.UsersRepository

	kaiRepository kairepository.KaiRepository

	xpRepository xprepository.XpRepository
}

func New(
	db *gorm.DB,
	userRepository usersrepository.UsersRepository,
	kaiRepository kairepository.KaiRepository,
	xpRepository xprepository.XpRepository,
) Service {
	return &service{
		db:             db,
		userRepository: userRepository,
		kaiRepository:  kaiRepository,
		xpRepository:   xpRepository,
	}
}

func (s *service) InitializeNewUser(
	ctx context.Context,
	user *usersmodel.User,
) error {

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. usuario
		if err := s.userRepository.CreateUser(ctx, tx, user); err != nil {
			log.Println("error creating user:", err)
			return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo crear nuevo usuario.")
		}

		// 2. config
		config := buildInitialConfiguration(user.ID)

		if err := s.userRepository.CreateConfiguration(ctx, tx, config); err != nil {
			log.Println("error initializing configuration user:", err)
			return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo inicializar configuracion_usuario.")
		}

		// 3. kai state
		kaiState := buildInitialKaiState(user.ID)

		if err := s.kaiRepository.CreateKaiState(ctx, tx, kaiState); err != nil {
			log.Println("error initializing user kai state:", err)
			return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo inicializar estado kai de usuario.")
		}

		// 4. categorias xp
		if err := s.initializeUserXPCategories(ctx, tx, user.ID); err != nil {
			log.Println("error initializing user xp:", err)
			return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo inicializar xp de usuario.")
		}

		// 5. atributos kai
		if err := s.initializeKaiAttributes(
			ctx,
			tx,
			user.ID,
		); err != nil {
			log.Println("error initializing attribute types of user:", err)
			return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo inicializar atributos kai de usuario.")
		}

		return nil
	})
}

func (s *service) initializeKaiAttributes(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) error {
	attributeTypes, err := s.kaiRepository.FindAllAttributeTypes(ctx)
	if err != nil {
		log.Println("error getting attribute types:", err)
		return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo obtener atributos kai.")
	}

	attributes := make([]kaimodel.KaiAttribute, 0, len(attributeTypes))

	for _, attrType := range attributeTypes {

		attributes = append(attributes, kaimodel.KaiAttribute{
			UserID:          userID,
			AttributeTypeID: attrType.ID,
			Value:           0,
		})
	}

	return s.kaiRepository.CreateAttributes(
		ctx,
		tx,
		attributes,
	)
}

func (s *service) initializeUserXPCategories(
	ctx context.Context,
	tx *gorm.DB,
	userID uuid.UUID,
) error {
	categories, err := s.xpRepository.FindAllCategories(ctx)
	if err != nil {
		log.Println("error getting categories:", err)
		return errorHandler.NewAppError(http.StatusInternalServerError, "No se pudo obtener categorías.")
	}

	userXP := make([]xpmodel.UserXP, 0, len(categories))

	for _, category := range categories {

		userXP = append(userXP, xpmodel.UserXP{
			UserID:     userID,
			CategoryID: category.ID,
			Value:      0,
		})
	}

	return s.xpRepository.CreateUserXP(
		ctx,
		tx,
		userXP,
	)
}
