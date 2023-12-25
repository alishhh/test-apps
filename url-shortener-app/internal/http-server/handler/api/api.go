package api

import (
	"encoding/json"
	"net/http"

	"github.com/alishhh/url-shortener-app/internal/storage/sqlite"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(alias string, aliases []sqlite.Alias) Response {
	return Response{
		Status:  StatusOK,
		Alias:   alias,
		Aliases: aliases,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

type Request struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type RequestClient struct {
	RefreshToken string `json:"refreshToken"`
}

type Response struct {
	Status  string         `json:"status"`
	Error   string         `json:"error,omitempty"`
	Alias   string         `json:"alias,omitempty"`
	Aliases []sqlite.Alias `json:"aliases,omitempty"`
}

func ResponseOK(w http.ResponseWriter, r *http.Request, alias string, aliases []sqlite.Alias) {
	req := OK(alias, aliases)
	res, _ := json.Marshal(&req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ResponseError(w http.ResponseWriter, r *http.Request, msg string) {
	req := Error(msg)
	res, _ := json.Marshal(&req)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
