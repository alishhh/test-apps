package remove

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alishhh/url-shortener-app/internal/http-server/handler/api"
	sl "github.com/alishhh/url-shortener-app/internal/logger"
	"github.com/alishhh/url-shortener-app/internal/storage"
	"github.com/go-chi/chi"
)

type URLRemover interface {
	RemoveURL(ctx context.Context, alias string) error
}

func Remove(ctx context.Context, log *slog.Logger, urlDeleter URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			http.Error(w, "alias is empty", http.StatusBadRequest)
			api.ResponseError(w, r, errors.New("alias is empty").Error())
			return
		}

		err := urlDeleter.RemoveURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			http.Error(w, err.Error(), http.StatusNotFound)
			api.ResponseError(w, r, err.Error())
			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			api.ResponseError(w, r, err.Error())
			return
		}

		api.ResponseOK(w, r, alias, nil)
	}
}
