package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
}

func (a *Config) GoogleOauthConfig() *oauth2.Config {
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	return &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     google.Endpoint,
		// RedirectURL:  "http://localhost:3000",
		RedirectURL: "http://localhost:8080/api/v1/auth/google/callback",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}
}

func NewConfigInit() *Config {
	return &Config{}
}
