package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
		Host string
	}

	Chatwoot struct {
		URL   string
		Token string
	}

	Codechat struct {
		URL string
		GlobalToken string
	}
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, loading from environment variables.")
	}

	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.Host = getEnv("HOST", "0.0.0.0")

	cfg.Chatwoot.URL = os.Getenv("CHATWOOT_URL")
	// cfg.Chatwoot.Token = os.Getenv("CHATWOOT_TOKEN")

	cfg.Codechat.URL = os.Getenv("CODECHAT_URL")
	cfg.Codechat.GlobalToken = os.Getenv("CODECHAT_KEY")

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
