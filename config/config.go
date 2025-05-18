package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	SessionKey  string
}

var AppConfig Config

func LoadConfig() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:3225@localhost/webapp?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = "secret-key-123"
	}

	AppConfig = Config{
		DatabaseURL: dbURL,
		ServerPort:  port,
		SessionKey:  sessionKey,
	}

	log.Println("Конфигурация загружена успешно")
} 