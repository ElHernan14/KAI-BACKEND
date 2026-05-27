package habits

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	router.GET("", controller.GetHabits)
}
