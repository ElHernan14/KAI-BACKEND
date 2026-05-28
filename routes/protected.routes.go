package routes

import (
	habitsRoute "kai-back/internal/modules/habits/routes"
	"kai-back/internal/modules/home"
	"kai-back/internal/modules/users"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(api *gin.RouterGroup, container *AppContainer) {
	protected := api.Group("")
	protected.Use(container.AuthMiddleware)

	home.RegisterRoutes(protected.Group("/home"), container.HomeController)
	users.RegisterRoutes(protected.Group("/users"), container.UserController)
	habitsRoute.RegisterRoutes(protected.Group("/habits"), container.HabitsController)
}
