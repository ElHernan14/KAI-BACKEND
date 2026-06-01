package habitsRoutes

import (
	controller "kai-back/internal/modules/habits/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, controller *controller.Controller) {
	router.GET("", controller.GetHabits)
	router.GET("/categories", controller.GetCategories)
	router.GET("/catalog", controller.GetCatalogByCategory)
	router.POST("/select", controller.SelectHabit)
	router.GET("/:habitUserId", controller.GetHabitDetail)
	router.DELETE("/:habitUserId/deactivate", controller.DeactivateHabit)
}
