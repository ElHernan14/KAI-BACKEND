package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	router.GET("/me", controller.Me)
	router.PUT("/me", controller.UpdateMe)
	router.PUT("/change-password", controller.ChangePassword)
}
