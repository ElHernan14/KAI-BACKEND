package routes

import (
	"kai-back/internal/modules/habits"
	"kai-back/internal/modules/users"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(api *gin.RouterGroup, container *AppContainer) {
	protected := api.Group("")
	protected.Use(container.AuthMiddleware)

	users.RegisterRoutes(protected.Group("/users"))
	habits.RegisterRoutes(protected.Group("/habits"))
}
