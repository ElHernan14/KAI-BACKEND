package main

import (
	"kai-back/internal/config"
	"kai-back/internal/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// load env config
	cfg := config.LoadConfig()

	// setup logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[KAI-API]")

	// init database
	err := database.InitGorm(cfg)
	if err != nil {
		log.Fatal("❌ Error initializing database: ", err)
	}

	// gin
	router := gin.Default()

	// test route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("🚀 KAI backend API started")
	log.Println("🚀 Server running on port:", cfg.ServerPort)

	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
