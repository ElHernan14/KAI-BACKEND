package main

import (
	"kai-back/internal/config"
	"kai-back/internal/database"
	"kai-back/routes"
	"log"
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
		log.Fatal("Error initializing database: ", err)
	}

	router := routes.SetupRouter(cfg, database.DB)

	log.Println("KAI backend API started")
	log.Println("Server running on port:", cfg.ServerPort)

	err = router.Run(":" + cfg.ServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
