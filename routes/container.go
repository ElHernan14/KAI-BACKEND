package routes

import (
	"time"

	"kai-back/internal/config"
	"kai-back/internal/middleware"
	authmodule "kai-back/internal/modules/auth"
	authrepository "kai-back/internal/modules/auth/repository"
	homemodule "kai-back/internal/modules/home"
	homerepository "kai-back/internal/modules/home/repository"
	kairepository "kai-back/internal/modules/kai/repository"
	usermodule "kai-back/internal/modules/users"
	userrepository "kai-back/internal/modules/users/repository"
	xprepository "kai-back/internal/modules/xp/repository"
	initializeruserservice "kai-back/internal/services/user_initializer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppContainer struct {
	AuthController *authmodule.Controller
	HomeController *homemodule.Controller
	UserController *usermodule.Controller
	AuthMiddleware gin.HandlerFunc
}

func NewAppContainer(db *gorm.DB, cfg config.Config) *AppContainer {
	//repositories
	authRepository := authrepository.NewRepository(db)
	userRepository := userrepository.NewRepository(db)
	xpRepository := xprepository.New(db)
	kaiRepository := kairepository.NewRepository(db)

	//services
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
	homeRepository := homerepository.NewRepository(db)
	homeService := homemodule.NewService(homeRepository)

	//controllers
	authController := authmodule.NewController(authService)
	homeController := homemodule.NewController(homeService)
	userController := usermodule.NewController(userService)

	return &AppContainer{
		AuthController: authController,
		HomeController: homeController,
		UserController: userController,
		AuthMiddleware: middleware.AuthMiddleware(cfg.JWTSecret),
	}
}
