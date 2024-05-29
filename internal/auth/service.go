package auth

import (
	"context"
	"encoding/json"
	"io"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	refreshtoken "money-tracker/internal/refresh_token"
	"money-tracker/internal/utils"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthService interface {
	ExchangeToken(code string) (*oauth2.Token, *domain.Error)
	ParseTokenToUser(token string) (*dto.GoogleUserData, *domain.Error)
	GenerateNewToken(user *entity.User) (*dto.AuthResponse, *domain.Error)
}

type authService struct {
	googleConfig        *oauth2.Config
	refreshTokenService refreshtoken.RefreshTokenService
}

// GenerateToken implements AuthService.
func (a *authService) GenerateNewToken(user *entity.User) (*dto.AuthResponse, *domain.Error) {
	secret := os.Getenv("JWT_SECRET")
	expiredTime := time.Now().Add(time.Minute * 5)
	claims := &dto.ATClaims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, _ := token.SignedString(secret)

	expiredRefreshTime := time.Now().Add(time.Hour * 24 * 7)
	refreshToken := utils.GenerateRandomCode(32)

	_, err := a.refreshTokenService.GenerateRefreshToken(&entity.RefreshToken{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    &expiredRefreshTime,
	})

	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:           accessToken,
		TokenType:             "Bearer",
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  int(expiredTime.Unix()),
		RefreshTokenExpiresIn: int(expiredRefreshTime.Unix()),
	}, nil

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

func NewAuthService(googleConfig *oauth2.Config, refreshTokenService refreshtoken.RefreshTokenService) AuthService {
	return &authService{
		googleConfig,
		refreshTokenService,
	}
}
