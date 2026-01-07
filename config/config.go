package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Psp_api_key_dk string
	Psp_api_key_se string
	Psp_api_key_no string
	DatabaseURL    string
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

func LoadConfig() Config {
	loadEnvFile()
	cfg := Config{
		Psp_api_key_dk: mustGetenv("PSP_API_KEY_DK"),
		Psp_api_key_se: mustGetenv("PSP_API_KEY_SE"),
		Psp_api_key_no: mustGetenv("PSP_API_KEY_NO"),
		DatabaseURL:    mustGetenv("DATABASE_URL"),
	}
	return cfg
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}
