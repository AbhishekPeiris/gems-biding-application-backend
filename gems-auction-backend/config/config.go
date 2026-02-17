package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	MaxDBConns  int32
}

var AppConfig *Config

// LoadConfig loads environment variables into struct
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found. Using system environment variables")
	}

	maxConnsStr := getEnv("DB_MAX_CONNS", "10")
	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil {
		log.Fatal("Invalid DB_MAX_CONNS value")
	}

	AppConfig = &Config{
		Port:       getEnv("PORT", "8081"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "gems_auction"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecret"),
		MaxDBConns: int32(maxConns),
	}

	log.Println("✅ Configuration Loaded Successfully")
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
