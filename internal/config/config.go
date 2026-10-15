package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv string

	HTTP HTTPConfig

	DB DatabaseConfig

	RabbitMQ RabbitMQConfig
}

type HTTPConfig struct {
	Addr string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RabbitMQConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    VHost    string
}

func Load() (Config, error) {

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	cfg := Config{
		AppEnv: getEnv("APP_ENV", "development"),
		HTTP: HTTPConfig{ Addr: getEnv("HTTP_ADDR", ":8080"),},
	}

	var err error

	cfg.DB.Host =  getEnv("DB_HOST", "localhost")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	if cfg.DB.Name == "" {
		return Config{}, errors.New("DB_NAME is required")
	}

	if cfg.DB.User == "" {
		return Config{}, errors.New("DB_USER is required")
	}

	if cfg.DB.Password == "" {
		return Config{}, errors.New("DB_PASSWORD is required")
	}

	cfg.DB.Port, err = getInt("DB_PORT", 5432)
	if err != nil {
		return Config{}, err
	}

	if cfg.DB.Port < 1 || cfg.DB.Port > 65535 {
		return Config{}, errors.New("DB_PORT must be between 1 and 65535")
	}

	cfg.DB.MaxOpenConns, err = getInt("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return Config{}, err
	}

	cfg.DB.MaxIdleConns, err = getInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return Config{}, err
	}

	if cfg.DB.MaxOpenConns < 1 {
		return Config{}, errors.New("DB_MAX_OPEN_CONNS must be positive")
	}

	if cfg.DB.MaxIdleConns < 0 ||
		cfg.DB.MaxIdleConns > cfg.DB.MaxOpenConns {
		return Config{}, errors.New(
			"DB_MAX_IDLE_CONNS must be between 0 and DB_MAX_OPEN_CONNS",
		)
	}

	lifetime := getEnv("DB_CONN_MAX_LIFETIME", "5m")
	cfg.DB.ConnMaxLifetime, err = time.ParseDuration(lifetime)
	if err != nil || cfg.DB.ConnMaxLifetime <= 0 {
		return Config{}, errors.New(
			"DB_CONN_MAX_LIFETIME must be a positive duration, e.g. 5m",
		)
	}

	
	cfg.RabbitMQ.Host = getEnv("RABBITMQ_HOST", "localhost")
	cfg.RabbitMQ.User = os.Getenv("RABBITMQ_USER")
	cfg.RabbitMQ.Password = os.Getenv("RABBITMQ_PASSWORD")
	cfg.RabbitMQ.VHost = getEnv("RABBITMQ_VHOST", "/")

	if cfg.RabbitMQ.User == "" {
		return Config{}, errors.New("RABBITMQ_USER is required")
	}

	if cfg.RabbitMQ.Password == "" {
		return Config{}, errors.New("RABBITMQ_PASSWORD is required")
	}

	cfg.RabbitMQ.Port, err = getInt("RABBITMQ_PORT", 5672)
	if err != nil {
		return Config{}, err
	}

	if cfg.RabbitMQ.Port < 1 || cfg.RabbitMQ.Port > 65535 {
		return Config{}, errors.New(
			"RABBITMQ_PORT must be between 1 and 65535",
		)
	}


	return cfg, nil

}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func getInt(key string, fallback int) (int, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}

	return parsed, nil
}