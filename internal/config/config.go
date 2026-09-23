package config

import "os"

type Config struct{ HTTPAddr, DatabaseURL, BotToken string }

func Load() Config {
	return Config{HTTPAddr: get("HTTP_ADDR", ":8080"), DatabaseURL: get("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/splitmate?sslmode=disable"), BotToken: os.Getenv("TELEGRAM_BOT_TOKEN")}
}
func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
