package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfig() *Config {
	err := godotenv.Load("configs/server.env")
	if err != nil {
		log.Println("No configs/server.env file found or error loading it")
	}

	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "monitoring"),
		DBUser:     getEnv("DB_USER", "monitoring_user"),
		DBPassword: getEnv("DB_PASSWORD", "StrongPassword"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
