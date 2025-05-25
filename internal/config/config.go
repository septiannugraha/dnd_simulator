package config

import (
	"os"
)

type Config struct {
	Port        string
	MongoURI    string
	DatabaseName string
	JWTSecret   string
	OpenAIKey   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "dnd_simulator"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		OpenAIKey:   getEnv("OPENAI_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}