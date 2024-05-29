package auth

import (
	"context"
	"encoding/json"
	"io"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type AuthService interface {
	ExchangeToken(code string) (*oauth2.Token, *domain.Error)
	ParseTokenToUser(token string) (*dto.GoogleUserData, *domain.Error)
}

type authService struct {
	googleConfig *oauth2.Config
}

// ParseTokenToUser implements AuthService.
func (a *authService) ParseTokenToUser(token string) (*dto.GoogleUserData, *domain.Error) {
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

// ExchangeToken implements AuthService.
func (a *authService) ExchangeToken(code string) (*oauth2.Token, *domain.Error) {
	token, err := a.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: fiber.StatusInternalServerError,
		}
	}
	return token, nil
}

func NewAuthService(googleConfig *oauth2.Config) AuthService {
	return &authService{
		googleConfig,
	}
}
