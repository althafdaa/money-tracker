package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ConfigGoogle struct {
	GoogleLoginConfig *oauth2.Config
}

func (c *ConfigGoogle) GoogleOauthConfig() *oauth2.Config {
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")

	return &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:3000",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}
}

func NewConfig() *ConfigGoogle {
	return &ConfigGoogle{}
}
