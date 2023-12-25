package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alishhh/url-shortener-app/internal/config"
	authapp "github.com/alishhh/url-shortener-app/internal/http-client/auth-app"
	"github.com/alishhh/url-shortener-app/internal/http-server/handler"
	"github.com/alishhh/url-shortener-app/internal/logger"
	"github.com/alishhh/url-shortener-app/internal/service"
	"github.com/alishhh/url-shortener-app/internal/storage/sqlite"
)

type App struct {
	Config  *config.Config
	Logger  *slog.Logger
	Storage *sqlite.Storage
	Service *service.Service
	Handler *handler.Handler
}

func New() *App {
	config := config.New()
	logger := logger.New(config)
	storage := sqlite.New(config, logger)
	client := authapp.New(config, logger)
	service := service.New(config, logger, storage, client)
	handler := handler.New(config, logger, service)

	return &App{
		Config:  config,
		Logger:  logger,
		Storage: storage,
		Service: service,
		Handler: handler,
	}
}

func (app *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{
		Addr:         fmt.Sprintf("%s:%s", app.Config.AppHost, app.Config.AppPort),
		Handler:      app.Handler.SetupRoutes(ctx, app.Storage),
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
