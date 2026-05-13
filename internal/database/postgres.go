package database

import (
	"database/sql"
	"fmt"
	"log"

	"kai-back/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg config.Config) (*sql.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("✅ PostgreSQL connected successfully")

	return db, nil
}
