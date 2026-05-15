package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
}
