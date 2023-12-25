package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alishhh/url-shortener-app/internal/http-server/handler/api"
	sl "github.com/alishhh/url-shortener-app/internal/logger"
	"github.com/alishhh/url-shortener-app/internal/storage"
	"github.com/alishhh/url-shortener-app/internal/storage/sqlite"
	"github.com/go-chi/chi"
)

type URLGetter interface {
	GetURLs(ctx context.Context) ([]sqlite.Alias, error)
	GetURL(ctx context.Context, alias string) (string, error)
}

func Get(ctx context.Context, log *slog.Logger, urlGeter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			http.Error(w, "alias is empty", http.StatusBadRequest)
			api.ResponseError(w, r, errors.New("alias is empty").Error())
			return
		}

		res, err := urlGeter.GetURL(ctx, alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			http.Error(w, err.Error(), http.StatusNotFound)
			api.ResponseError(w, r, err.Error())
			return
		}

		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			api.ResponseError(w, r, err.Error())
			return
		}

		http.Redirect(w, r, res, http.StatusFound)
	}
}

func GetAll(ctx context.Context, log *slog.Logger, urlGeter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := urlGeter.GetURLs(ctx)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("urls not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			api.ResponseError(w, r, err.Error())
			return
		}

		if err != nil {
			log.Error("failed to get urls", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			api.ResponseError(w, r, err.Error())
			return
		}

		api.ResponseOK(w, r, "", res)
	}
}
