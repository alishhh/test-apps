package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel             string        `env:"LOG_LEVEL" env-required:"true"`
	AppName              string        `env:"APP_NAME" env-required:"true"`
	AppHost              string        `env:"APP_HOST" env-required:"true"`
	AppPort              string        `env:"APP_PORT" env-required:"true"`
	KeycloakHost         string        `env:"KEYCLOAK_HOST" env-required:"true"`
	KeycloakRealm        string        `env:"KEYCLOAK_REALM" env-required:"true"`
	KeycloakClientID     string        `env:"KEYCLOAK_CLIENT_ID" env-required:"true"`
	KeycloakClientSecret string        `env:"KEYCLOAK_CLIENT_SECRET" env-required:"true"`
	Timeout              time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout          time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
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
	if cfg.KeycloakHost == "" {
		log.Fatalf("KEYCLOAK_HOST is empty")
	}
	if cfg.KeycloakRealm == "" {
		log.Fatalf("KEYCLOAK_REALM is empty")
	}
	if cfg.KeycloakClientID == "" {
		log.Fatalf("KEYCLOAK_CLIENT_ID is empty")
	}
	if cfg.KeycloakClientSecret == "" {
		log.Fatalf("KEYCLOAK_CLIENT_SECRET is empty")
	}
	return &cfg
}
