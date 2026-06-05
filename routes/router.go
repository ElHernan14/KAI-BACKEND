package routes

import (
	"kai-back/internal/config"
	"kai-back/internal/middleware"
	"kai-back/internal/shared/response"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	router := gin.New()
	container := NewAppContainer(db, cfg)
	_ = router.SetTrustedProxies(nil)

	router.Use(
		middleware.RequestIDMiddleware(),
		middleware.LoggerMiddleware(),
		middleware.ErrorMiddleware(),
	)

	router.GET("/health", healthHandler)
	if err := os.MkdirAll("uploads/profiles", 0755); err == nil {
		router.Static("/uploads/profiles", "uploads/profiles")
	}

	api := router.Group("/api/v1")

	RegisterPublicRoutes(api, container)
	RegisterProtectedRoutes(api, container)

	return router
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, response.Success(gin.H{
		"status": "ok",
	}))
}
