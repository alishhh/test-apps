package service

import (
	"log/slog"

	"magnum.kz/services/auth-app/internal/config"
	"magnum.kz/services/auth-app/internal/keycloak"
)

type Service struct {
	logger   *slog.Logger
	Keycloak *keycloak.Keycloak
}

func New(cfg *config.Config, log *slog.Logger, keycloak *keycloak.Keycloak) *Service {
	return &Service{
		logger:   log,
		Keycloak: keycloak,
	}
}
