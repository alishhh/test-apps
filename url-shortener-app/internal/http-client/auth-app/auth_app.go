package authapp

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/alishhh/url-shortener-app/internal/config"
)

type Client struct {
	Client *http.Client
	URL    string
}

func New(cfg *config.Config, logger *slog.Logger) *Client {
	return &Client{
		Client: &http.Client{},
		URL:    cfg.ClientURL,
	}
}

func (c *Client) ValidateToken(token string) (int, error) {
	mt := map[string]string{
		"accessToken": token,
	}
	req, _ := json.Marshal(mt)
	reqBody := bytes.NewBuffer(req)
	resp, err := c.Client.Post(c.URL+"/validateToken", "application/json", reqBody)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
