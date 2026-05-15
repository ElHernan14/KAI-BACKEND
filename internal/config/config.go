package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      int
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	ServerPort  string
	JWTSecret   string
	JWTTTLHours int
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system env")
	}

	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	jwtTTLHours, _ := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))

	return Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      port,
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "postgres"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "super_secret_key"),
		JWTTTLHours: jwtTTLHours,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
