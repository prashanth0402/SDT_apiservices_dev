package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv string
	Port   string
	DBUser string
	DBPass string
	DBHost string
}

var (
	cfg  *AppConfig
	once sync.Once
)

// LoadConfig loads env variables once and caches them
func LoadConfig() *AppConfig {
	once.Do(func() {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "local" // default
		}

		// Load file only in local/dev
		switch env {
		case "local":
			if err := godotenv.Load(".env.development"); err != nil {
				log.Println("⚠️ No .env.development found, using system environment")
			}
		case "production":
			// Typically no file in production
			_ = godotenv.Load(".env.prod") // optional if you keep a prod file
		}

		cfg = &AppConfig{
			AppEnv: env,
			Port:   getEnv("PORT", "8080"),
			DBUser: getEnv("DB_USER", "root"),
			DBPass: getEnv("DB_PASS", ""),
			DBHost: getEnv("DB_HOST", "localhost"),
		}
	})

	return cfg
}

// Helper for default values
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
