package routes

import (
	authmodule "kai-back/internal/modules/auth"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(api *gin.RouterGroup, container *AppContainer) {
	api.GET("/health", healthHandler)

	authmodule.RegisterRoutes(api.Group("/auth"), container.AuthController)
}
