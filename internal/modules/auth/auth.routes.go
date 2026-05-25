package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *Controller, authMiddleware gin.HandlerFunc) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
	router.GET("/validate", authMiddleware, controller.ValidateToken)
}
