package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel    string        `env:"LOG_LEVEL" env-required:"true"`
	AppName     string        `env:"APP_NAME" env-required:"true"`
	AppHost     string        `env:"APP_HOST" env-required:"true"`
	AppPort     string        `env:"APP_PORT" env-required:"true"`
	DBPath      string        `env:"DB_PATH" env-required:"true"`
	ClientURL   string        `env:"CLIENT_URL" env-required:"true"`
	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("error reading envs: %s", err.Error())
	}
	if cfg.LogLevel == "" {
		log.Fatalf("LOG_LEVEL is empty")
	}
	if cfg.AppName == "" {
		log.Fatalf("APP_NAME is empty")
	}
	if cfg.AppHost == "" {
		log.Fatalf("APP_HOST is empty")
	}
	if cfg.AppPort == "" {
		log.Fatalf("APP_PORT is empty")
	}
	if cfg.ClientURL == "" {
		log.Fatalf("CLIENT_URL is empty")
	}
	return &cfg
}
