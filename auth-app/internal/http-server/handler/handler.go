package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"magnum.kz/services/auth-app/internal/config"
	"magnum.kz/services/auth-app/internal/http-server/api"
	mwLogger "magnum.kz/services/auth-app/internal/http-server/middleware/logger"
	sl "magnum.kz/services/auth-app/internal/logger"
	"magnum.kz/services/auth-app/internal/service"
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

func (h *Handler) SetupRoutes(ctx context.Context) *chi.Mux {
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
	r.Post("/login", h.Login(ctx))
	r.Post("/refreshToken", h.RefreshToken(ctx))
	r.Post("/validateToken", h.ValidateToken(ctx))

	return r
}

func (h *Handler) Login(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.Request{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req); err != nil {
			h.logger.Error("failed to decode", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			api.ResponseError(w, r, err.Error())
			return
		}

		jwt, err := h.Service.Keycloak.Login(ctx, req.Username, req.Password)
		if err != nil {
			h.logger.Error("failed to login", sl.Err(err))
			http.Error(w, err.Error(), http.StatusForbidden)
			api.ResponseError(w, r, err.Error())
			return
		}

		resp := &api.Response{
			AccessToken:  jwt.AccessToken,
			RefreshToken: jwt.RefreshToken,
			ExpiresIn:    jwt.ExpiresIn,
		}
		api.ResponseOK(w, r, resp.AccessToken, resp.RefreshToken, resp.ExpiresIn)
	}
}

func (h *Handler) RefreshToken(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.Request{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req); err != nil {
			h.logger.Error("failed to decode", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			api.ResponseError(w, r, err.Error())
			return
		}

		jwt, err := h.Service.Keycloak.RefreshToken(ctx, req.RefreshToken)
		if err != nil {
			h.logger.Error("failed to refresh token", sl.Err(err))
			http.Error(w, err.Error(), http.StatusForbidden)
			api.ResponseError(w, r, err.Error())
			return
		}

		resp := &api.Response{
			AccessToken:  jwt.AccessToken,
			RefreshToken: jwt.RefreshToken,
			ExpiresIn:    jwt.ExpiresIn,
		}
		api.ResponseOK(w, r, resp.AccessToken, resp.RefreshToken, resp.ExpiresIn)
	}
}

func (h *Handler) ValidateToken(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.Request{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req); err != nil {
			h.logger.Error("failed to decode", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			api.ResponseError(w, r, err.Error())
			return
		}

		jwt, err := h.Service.Keycloak.ValidateToken(ctx, req.AccessToken)
		if err != nil {
			h.logger.Error("failed to refresh token", sl.Err(err))
			http.Error(w, err.Error(), http.StatusForbidden)
			api.ResponseError(w, r, err.Error())
			return
		}

		if !*jwt.Active {
			h.logger.Error("invalid or expired token", sl.Err(errors.New("invalid or expired token")))
			http.Error(w, errors.New("invalid or expired token").Error(), http.StatusUnauthorized)
			api.ResponseError(w, r, errors.New("invalid or expired token").Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
