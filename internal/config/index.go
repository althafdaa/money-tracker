package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config interface {
	googleOauthConfig() *oauth2.Config
	exchange(code string) (*oauth2.Token, *domain.Error)
	parseGoogleUserData(token string) (*dto.GoogleUserData, *domain.Error)
	ParseAccessTokenToUserData(token string) (*dto.GoogleUserData, *domain.Error)
	AuthCodeURL() string
}

type config struct {
}

// AuthCodeURL implements Config.
func (a *config) AuthCodeURL() string {
	state := os.Getenv("GOOGLE_STATE")
	return a.googleOauthConfig().AuthCodeURL(state)
}

// ParseAccessTokenToUserData implements Config.
func (a *config) ParseAccessTokenToUserData(code string) (*dto.GoogleUserData, *domain.Error) {
	token, err := a.exchange(code)
	if err != nil {
		return nil, err
	}
	userData, err := a.parseGoogleUserData(token.AccessToken)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

// ParseGoogleUserData implements Config.
func (a *config) parseGoogleUserData(token string) (*dto.GoogleUserData, *domain.Error) {
	if token == "" {
		return nil, &domain.Error{
			Err:  errors.New("TOKEN_IS_EMPTY"),
			Code: fiber.StatusBadRequest}
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token)
	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: fiber.StatusInternalServerError,
		}
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: fiber.StatusInternalServerError,
		}
	}
	var googleUserData dto.GoogleUserData
	unmarshalErr := json.Unmarshal(userData, &googleUserData)
	if unmarshalErr != nil {
		return nil, &domain.Error{
			Err:  unmarshalErr,
			Code: fiber.StatusInternalServerError,
		}
	}

	return &googleUserData, nil
}

// Exchange implements Config.
func (a *config) exchange(code string) (*oauth2.Token, *domain.Error) {
	if code == "" {
		return nil, &domain.Error{
			Err:  errors.New("CODE_IS_EMPTY"),
			Code: fiber.StatusBadRequest,
		}
	}
	token, err := a.googleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: fiber.StatusInternalServerError,
		}
	}
	return token, nil
}

func (a *config) googleOauthConfig() *oauth2.Config {
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	client_url := os.Getenv("GOOGLE_REDIRECT_URL_DOMAIN")

	return &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: client_secret,
		Endpoint:     google.Endpoint,
		RedirectURL:  client_url + "/",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
	}
}

func NewConfigInit() Config {
	return &config{}
}
