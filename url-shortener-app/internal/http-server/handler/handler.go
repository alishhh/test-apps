package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/alishhh/url-shortener-app/internal/config"
	"github.com/alishhh/url-shortener-app/internal/http-server/handler/api"
	"github.com/alishhh/url-shortener-app/internal/http-server/handler/url/get"
	"github.com/alishhh/url-shortener-app/internal/http-server/handler/url/remove"
	"github.com/alishhh/url-shortener-app/internal/http-server/handler/url/save"
	mwLogger "github.com/alishhh/url-shortener-app/internal/http-server/middleware/logger"
	sl "github.com/alishhh/url-shortener-app/internal/logger"
	"github.com/alishhh/url-shortener-app/internal/service"
	"github.com/alishhh/url-shortener-app/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	logger  *slog.Logger
	Service *service.Service
}

func New(cfg *config.Config, log *slog.Logger, service *service.Service) *Handler {
	return &Handler{
		logger:  log,
		Service: service,
	}
}

func (h *Handler) SetupRoutes(ctx context.Context, storage *sqlite.Storage) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mwLogger.New(h.logger))
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	r.Route("/URL", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(h.ValidateToken)
			r.Post("/save", save.New(ctx, h.logger, storage))
			r.Delete("/remove/{alias}", remove.Remove(ctx, h.logger, storage))
		})

		r.Get("/", get.GetAll(ctx, h.logger, storage))
		r.Get("/{alias}", get.Get(ctx, h.logger, storage))
	})

	return r
}

func (h *Handler) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, "Bearer ")
		if authHeader == "" || reqToken == authHeader {
			h.logger.Error("missing authorization header", sl.Err(errors.New("missing authorization header")))
			http.Error(w, errors.New("missing authorization header").Error(), http.StatusForbidden)
			return
		}

		status, err := h.Service.Client.ValidateToken(reqToken)
		if err != nil {
			h.logger.Error("failed to validate token", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			api.ResponseError(w, r, err.Error())
			return
		}

		if status == 200 {
			w.WriteHeader(status)
			next.ServeHTTP(w, r)
		} else {
			h.logger.Error("unauthorized or invalid token", sl.Err(errors.New("unauthorized or invalid token")))
			http.Error(w, errors.New("unauthorized or invalid token").Error(), status)
			api.ResponseError(w, r, errors.New("unauthorized or invalid token").Error())
			return
		}
	})
}
