package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	ServerPort     string
	JWTSecret      string
	JWTTTLHours    int
	GoogleClientID string
	FrontendURL    string
	SMTP           SMTPConfig
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system env")
	}

	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	jwtTTLHours, _ := strconv.Atoi(getEnv("JWT_TTL_HOURS", "1"))
	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         port,
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "postgres"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "super_secret_key"),
		JWTTTLHours:    jwtTTLHours,
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		FrontendURL:    getEnv("FRONTEND_URL", "smtp.gmail.com"),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     smtpPort,
			User:     getEnv("SMTP_USER", "hernanbonne98@gmail.com"),
			Password: getEnv("SMTP_PASSWORD", "app_password"),
			From:     getEnv("SMTP_FROM", "KAI <hernanbonne98@gmail.com>"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
