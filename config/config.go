package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment   string
	Port          string
	DatabaseURL   string
	RedisHost     string
	RedisPort     string
	RedisPwd      string
	RedisDB       int
	JWTSecret     string
	JWTSauce      string
	JWTExpiration int
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RequestLimit  int
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Environment:   getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "3000"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPwd:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTSauce:      getEnv("JWT_SAUCE", ""),
		JWTExpiration: getEnvAsInt("JWT_EXPIRATION", 1),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USERNAME", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "aub-task"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		RequestLimit:  getEnvAsInt("RATE_LIMIT", 100),
	}
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
		)
	}
	if cfg.JWTSecret == "" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
