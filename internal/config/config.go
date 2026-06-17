package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseRrl string
	JWTSecret   string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseRrl: getRequiredEnv("DATABASE_URL"),
		JWTSecret:   getRequiredEnv("JWT_SECRET"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s is required", key)
	}

	return value
}
