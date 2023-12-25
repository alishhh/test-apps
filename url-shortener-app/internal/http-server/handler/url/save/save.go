package save

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/alishhh/url-shortener-app/internal/http-server/handler/api"
	sl "github.com/alishhh/url-shortener-app/internal/logger"
	"github.com/alishhh/url-shortener-app/internal/storage"
)

type URLSaver interface {
	SaveURL(ctx context.Context, urlToSave string, alias string) error
}

func New(ctx context.Context, log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &api.Request{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(req); err != nil {
			log.Info("failed to decode", sl.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			api.ResponseError(w, r, err.Error())
			return
		}

		err := urlSaver.SaveURL(ctx, req.URL, req.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			http.Error(w, err.Error(), http.StatusConflict)
			api.ResponseError(w, r, err.Error())
			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			api.ResponseError(w, r, err.Error())
			return
		}

		api.ResponseOK(w, r, req.Alias, nil)
	}
}
