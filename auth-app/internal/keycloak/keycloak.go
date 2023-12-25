package keycloak

import (
	"context"
	"log/slog"

	"github.com/Nerzal/gocloak/v13"
	"magnum.kz/services/auth-app/internal/config"
)

type Keycloak struct {
	logger         *slog.Logger
	GocloackClient *gocloak.GoCloak
	Realm          string
	ClientID       string
	ClientSecret   string
}

func New(cfg *config.Config, log *slog.Logger) *Keycloak {
	client := gocloak.NewClient(cfg.KeycloakHost)
	return &Keycloak{
		logger:         log,
		Realm:          cfg.KeycloakRealm,
		ClientID:       cfg.KeycloakClientID,
		ClientSecret:   cfg.KeycloakClientSecret,
		GocloackClient: client,
	}
}

func (k *Keycloak) Login(ctx context.Context, username, password string) (*gocloak.JWT, error) {
	jwt, err := k.GocloackClient.Login(
		ctx,
		k.ClientID,
		k.ClientSecret,
		k.Realm,
		username,
		password,
	)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (k *Keycloak) RefreshToken(ctx context.Context, refreshToken string) (*gocloak.JWT, error) {
	jwt, err := k.GocloackClient.RefreshToken(
		ctx,
		refreshToken,
		k.ClientID,
		k.ClientSecret,
		k.Realm,
	)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (k *Keycloak) ValidateToken(ctx context.Context, accessToken string) (*gocloak.IntroSpectTokenResult, error) {
	jwt, err := k.GocloackClient.RetrospectToken(
		ctx,
		accessToken,
		k.ClientID,
		k.ClientSecret,
		k.Realm,
	)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}
