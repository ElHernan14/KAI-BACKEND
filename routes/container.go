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
	//repositories
	authRepository := authrepository.NewRepository(db)
	userRepository := userrepository.NewRepository(db)
	xpRepository := xprepository.New(db)
	kaiRepository := kairepository.NewRepository(db)
	habitsRepository := habitsRepository.NewRepository(db)
	homeRepository := homerepository.NewRepository(db)
	messageRepo := messageRepository.NewMessageRepository(db)

	//services
	transaction := transactionGorm.NewGormTransactionManager(db)
	habitsDailyRecordsService := habitsServ.NewHabitsDailyRecordsService(habitsRepository)
	userActivitySyncService := userActivitySynchronizationService.New(transaction, userRepository, habitsRepository, kaiRepository)
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
		cfg.JWTSecret,
		time.Duration(cfg.JWTTTLHours)*time.Hour,
	)
	userService := usermodule.NewService(userRepository)
	habitsService := habitsServ.NewService(habitsRepository, habitsDailyRecordsService)
	homeService := homemodule.NewService(homeRepository, userRepository, userActivitySyncService, habitsDailyRecordsService)
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
