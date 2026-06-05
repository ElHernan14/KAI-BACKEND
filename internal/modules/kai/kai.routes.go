package kai

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, controller *KaiController) {
	router.GET("", controller.GetKaiDashboard)
}
