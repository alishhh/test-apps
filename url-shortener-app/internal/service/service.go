package service

import (
	"log/slog"

	"github.com/alishhh/url-shortener-app/internal/config"
	authapp "github.com/alishhh/url-shortener-app/internal/http-client/auth-app"
	"github.com/alishhh/url-shortener-app/internal/storage/sqlite"
)

type Service struct {
	logger  *slog.Logger
	Storage *sqlite.Storage
	Client  *authapp.Client
}

func New(cfg *config.Config, log *slog.Logger, storage *sqlite.Storage, client *authapp.Client) *Service {
	return &Service{
		logger:  log,
		Storage: storage,
		Client:  client,
	}
}
