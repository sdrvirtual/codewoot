package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
		Host string
		URL string
	}

	Chatwoot struct {
		URL   string
	}

	Codechat struct {
		URL string
		GlobalToken string
	}

	Database struct {
		URL string
	}

	Authorization struct {
		Key string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}
	godotenv.Load()

	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.Host = getEnv("HOST", "0.0.0.0")
	cfg.Server.URL = getEnv("API_URL", "http://localhost:8080")

	cfg.Chatwoot.URL = os.Getenv("CHATWOOT_URL")

	cfg.Codechat.URL = os.Getenv("CODECHAT_URL")
	cfg.Codechat.GlobalToken = os.Getenv("CODECHAT_KEY")

	cfg.Database.URL = os.Getenv("DB_URL")
	cfg.Authorization.Key = os.Getenv("API_KEY")

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
