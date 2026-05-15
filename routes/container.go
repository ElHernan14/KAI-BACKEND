package routes

import (
	"time"

	"kai-back/internal/config"
	"kai-back/internal/middleware"
	authmodule "kai-back/internal/modules/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppContainer struct {
	AuthController *authmodule.Controller
	AuthMiddleware gin.HandlerFunc
}

func NewAppContainer(db *gorm.DB, cfg config.Config) *AppContainer {
	//repositories
	authRepository := authmodule.NewRepository(db)

	//services
	authService := authmodule.NewService(
		authRepository,
		cfg.JWTSecret,
		time.Duration(cfg.JWTTTLHours)*time.Hour,
	)

	//controllers
	authController := authmodule.NewController(authService)

	return &AppContainer{
		AuthController: authController,
		AuthMiddleware: middleware.AuthMiddleware(cfg.JWTSecret),
	}
}
