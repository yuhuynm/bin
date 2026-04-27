package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName   string
	Port      string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string
	SecretKey string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppName:   getEnv("APP_NAME", "go-bin"),
		Port:      getEnv("PORT", "8080"),
		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "admin"),
		DBPass:    getEnv("DB_PASSWORD", "adminpassword123"),
		DBName:    getEnv("DB_NAME", "go_bin"),
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
		SecretKey: getEnv("SECRET_KEY", "change-this-32-byte-secret-key!!"),
	}
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPass,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
