package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	Port        string
	DBURL       string
	JWTSecret   string
	ReadTimeout time.Duration
}

func Load() *Config {
	godotenv.Load()
	cfg := &Config{
		Env:         getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		DBURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
		ReadTimeout: 5 * time.Second,
	}

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return val
}
