package config

import (
	"log"
	"os"
)

type Config struct {
	AppPort        string
	DatabaseURL string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		AppPort:        GetEnv("APP_PORT"),
		DatabaseURL: GetEnv("DATABASE_URL"),
		JWTSecret:   GetEnv("JWT_SECRET"),
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s environment variable not set.", key)
	}
	return value
}
