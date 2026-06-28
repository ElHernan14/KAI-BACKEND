package routes

import (
	"time"

	"kai-back/internal/config"
	"kai-back/internal/middleware"
	authmodule "kai-back/internal/modules/auth"
	authrepository "kai-back/internal/modules/auth/repository"
	habitsController "kai-back/internal/modules/habits/controller"
	habitsRepository "kai-back/internal/modules/habits/repository"
	habitsServ "kai-back/internal/modules/habits/service"
	homemodule "kai-back/internal/modules/home"
	homerepository "kai-back/internal/modules/home/repository"
	kaiModule "kai-back/internal/modules/kai"
	kairepository "kai-back/internal/modules/kai/repository"
	messageRepository "kai-back/internal/modules/messages/repository"
	usermodule "kai-back/internal/modules/users"
	userrepository "kai-back/internal/modules/users/repository"
	xprepository "kai-back/internal/modules/xp/repository"
	habitCompletionService "kai-back/internal/services/habit_completion"
	userActivitySynchronizationService "kai-back/internal/services/user_activity_synchronization"
	initializeruserservice "kai-back/internal/services/user_initializer"
	mail "kai-back/internal/shared/mail"
	transactionGorm "kai-back/internal/shared/transaction"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppContainer struct {
	AuthController   *authmodule.Controller
	HabitsController *habitsController.Controller
	HomeController   *homemodule.Controller
	UserController   *usermodule.Controller
	KaiController    *kaiModule.KaiController
	AuthMiddleware   gin.HandlerFunc
}

func NewAppContainer(db *gorm.DB, cfg config.Config) *AppContainer {
	//config
	appConfig := config.LoadConfig()

	//repositories
	authRepository := authrepository.NewRepository(db)
	userRepository := userrepository.NewRepository(db)
	xpRepository := xprepository.New(db)
	kaiRepository := kairepository.NewRepository(db)
	habitsRepository := habitsRepository.NewRepository(db)
	homeRepository := homerepository.NewRepository(db)
	messageRepo := messageRepository.NewMessageRepository(db)

	//services
	mailService := mail.NewService(mail.SMTPConfig{
		Host:     appConfig.SMTP.Host,
		Port:     appConfig.SMTP.Port,
		User:     appConfig.SMTP.User,
		Password: appConfig.SMTP.Password,
		From:     appConfig.SMTP.From,
	})
	transaction := transactionGorm.NewGormTransactionManager(db)
	habitsDailyRecordsService := habitsServ.NewHabitsDailyRecordsService(
		habitsRepository,
		userRepository,
		transaction,
	)
	userActivitySyncService := userActivitySynchronizationService.New(
		transaction,
		userRepository,
		habitsRepository,
		kaiRepository,
		messageRepo,
	)
	habitCompletionService := habitCompletionService.NewService(
		habitsRepository,
		xpRepository,
		kaiRepository,
		messageRepo,
		userActivitySyncService,
		habitsDailyRecordsService,
		transaction,
	)
	initializerUserService := initializeruserservice.New(
		db,
		userRepository,
		kaiRepository,
		xpRepository,
	)
	authService := authmodule.NewService(
		authRepository,
		userRepository,
		initializerUserService,
		transaction,
		cfg.JWTSecret,
		time.Duration(cfg.JWTTTLHours)*time.Hour,
		cfg.GoogleClientID,
		mailService,
		appConfig,
	)
	userService := usermodule.NewService(userRepository, xpRepository, userActivitySyncService, habitsDailyRecordsService)
	habitsService := habitsServ.NewService(habitsRepository, habitsDailyRecordsService)
	homeService := homemodule.NewService(
		homeRepository,
		userRepository,
		kaiRepository,
		messageRepo,
		userActivitySyncService,
		habitsDailyRecordsService,
	)
	kaiService := kaiModule.NewKaiService(kaiRepository, messageRepo, habitsRepository, xpRepository, habitsDailyRecordsService, userActivitySyncService)

	//controllers
	authController := authmodule.NewController(authService)
	habitsController := habitsController.NewController(habitsService, habitCompletionService)
	homeController := homemodule.NewController(homeService)
	userController := usermodule.NewController(userService)
	kaiController := kaiModule.NewKaiController(kaiService)

	return &AppContainer{
		AuthController:   authController,
		HabitsController: habitsController,
		HomeController:   homeController,
		UserController:   userController,
		KaiController:    kaiController,
		AuthMiddleware:   middleware.AuthMiddleware(cfg.JWTSecret),
	}
}
