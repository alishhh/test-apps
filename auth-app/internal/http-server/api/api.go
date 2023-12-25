package api

import (
	"encoding/json"
	"net/http"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(accessToken string, refreshToken string, expiresIn int) Response {
	return Response{
		Status:       StatusOK,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

type Request struct {
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
}

type Response struct {
	Status       string `json:"status,omitempty"`
	Error        string `json:"error,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int    `json:"expiresIn,omitempty"`
}

func ResponseOK(w http.ResponseWriter, r *http.Request, accessToken string, refreshToken string, expiresIn int) {
	req := OK(accessToken, refreshToken, expiresIn)
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
