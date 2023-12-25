package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"magnum.kz/services/auth-app/internal/config"
	"magnum.kz/services/auth-app/internal/http-server/handler"
	"magnum.kz/services/auth-app/internal/keycloak"
	"magnum.kz/services/auth-app/internal/logger"
	"magnum.kz/services/auth-app/internal/service"
)

type App struct {
	Config   *config.Config
	Logger   *slog.Logger
	Keycloak *keycloak.Keycloak
	Service  *service.Service
	Handler  *handler.Handler
}

func New() *App {
	config := config.New()
	logger := logger.New(config)
	keycloak := keycloak.New(config, logger)
	service := service.New(config, logger, keycloak)
	handler := handler.New(config, logger, service)

	return &App{
		Config:   config,
		Logger:   logger,
		Keycloak: keycloak,
		Service:  service,
		Handler:  handler,
	}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", app.Config.AppHost, app.Config.AppPort),
		Handler:      app.Handler.SetupRoutes(ctx),
		ReadTimeout:  app.Config.Timeout,
		WriteTimeout: app.Config.Timeout,
		IdleTimeout:  app.Config.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			app.Logger.Error("failed to start server")
		}
	}()

	app.Logger.Info("server started")

	<-done
	app.Logger.Info("stopping server")

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Error("failed to stop server", logger.Err(err))
		return
	}
	app.Logger.Info("server stopped")
}
