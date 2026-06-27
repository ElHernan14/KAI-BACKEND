package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *Controller) {
	router.GET("/me", controller.Me)
	router.PUT("/me", controller.UpdateMe)
	router.PUT("/me/profile-photo", controller.UpdateProfilePhoto)
	router.PUT("/change-password", controller.ChangePassword)
	router.GET("/profile", controller.GetProfile)
	router.PUT("/profile", controller.PutUserProfile)
	router.PUT("/configuration", controller.PutUserConfigurationProfile)
}
