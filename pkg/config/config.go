package config

import (
	"fmt"
	"os"
	"strconv"
)

func LoadConfig() *Config {
	return &Config{
		SERVER_PORT:      getEnvOrDefault("SERVER_PORT", ":8080"),
		DB_USER:          getEnvOrDefault("DB_USER", "postgres"),
		DB_PASSWORD:      getEnvOrDefault("DB_PASSWORD", "password"),
		DB_HOST:          getEnvOrDefault("DB_HOST", "localhost"),
		DB_PORT:          getEnvOrDefault("DB_PORT", "5432"),
		DB_NAME:          getEnvOrDefault("DB_NAME", "shorten-url-services"),
		DB_SSLMODE:       getEnvOrDefault("DB_SSLMODE", "disable"),
		QUEUE_NAME:       getEnvOrDefault("FACTORIAL_CAL_SERVICES_QUEUE_NAME", "default-queue"),
		SWAGGER_HOST:     getEnvOrDefault("SWAGGER_HOST", "localhost:8080"),
		REDIS_HOST:       getEnvOrDefault("REDIS_HOST", "localhost"),
		REDIS_PORT:       getEnvOrDefault("REDIS_PORT", "6379"),
		REDIS_PASSWORD:   getEnvOrDefault("REDIS_PASSWORD", ""),
		SECRET_KEY:       getEnvOrDefault("SECRET_KEY", "secret"),
		SHORT_URL_LENGTH: getEnvOrDefault("SHORT_URL_LENGTH", "8"),
	}
}

type Config struct {
	SERVER_PORT      string `mapstructure:"SERVER_PORT"`
	DB_USER          string `mapstructure:"DB_USER"`
	DB_PASSWORD      string `mapstructure:"DB_PASSWORD"`
	DB_HOST          string `mapstructure:"DB_HOST"`
	DB_PORT          string `mapstructure:"DB_PORT"`
	DB_NAME          string `mapstructure:"DB_NAME"`
	DB_SSLMODE       string `mapstructure:"DB_SSLMODE"`
	QUEUE_NAME       string `mapstructure:"QUEUE_NAME"`
	SWAGGER_HOST     string `mapstructure:"SWAGGER_HOST"`
	REDIS_HOST       string `mapstructure:"REDIS_HOST"`
	REDIS_PORT       string `mapstructure:"REDIS_PORT"`
	REDIS_PASSWORD   string `mapstructure:"REDIS_PASSWORD"`
	SECRET_KEY       string `mapstructure:"SECRET_KEY"`
	SHORT_URL_LENGTH string `mapstructure:"SHORT_URL_LENGTH"`
}

func (c *Config) DSN() string {
	// postgres://postgres:secret@localhost:5432/mydb?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DB_USER, c.DB_PASSWORD, c.DB_HOST, c.DB_PORT, c.DB_NAME, c.DB_SSLMODE)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.REDIS_HOST, c.REDIS_PORT)
}

// Validate validates the configuration and returns an error if required fields are missing or invalid
func (c *Config) Validate() error {
	// Required database fields
	if c.DB_HOST == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DB_PORT == "" {
		return fmt.Errorf("DB_PORT is required")
	}
	if c.DB_NAME == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DB_USER == "" {
		return fmt.Errorf("DB_USER is required")
	}

	// Validate DB_PORT is numeric
	if _, err := strconv.Atoi(c.DB_PORT); err != nil {
		return fmt.Errorf("DB_PORT must be numeric: %w", err)
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
