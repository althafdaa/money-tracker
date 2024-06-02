package auth

import (
	"context"
	"encoding/json"
	"io"
	"money-tracker/internal/config"
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
	GenerateAndUpdateNewToken(user *entity.User, refreshTokenID int) (*dto.AuthResponse, *domain.Error)
	tokenGenerator(user *entity.User) (*dto.NewTokenDto, *domain.Error)
}

type authService struct {
	refreshTokenService refreshtoken.RefreshTokenService
	config              *config.Config
}

// GenerateAndUpdateNewToken implements AuthService.
func (a *authService) GenerateAndUpdateNewToken(user *entity.User, refreshTokenID int) (*dto.AuthResponse, *domain.Error) {
	token, err := a.tokenGenerator(user)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	data, dataErr := a.refreshTokenService.UpdateRefreshTokenByRefreshTokenID(refreshTokenID, &entity.RefreshToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiredAt:    &token.RefreshTokenExpiresIn,
		UpdatedAt:    &now,
		UserID:       user.ID,
	})

	if dataErr != nil {
		return nil, dataErr
	}

	return &dto.AuthResponse{
		AccessToken:           data.AccessToken,
		TokenType:             "Bearer",
		RefreshToken:          data.RefreshToken,
		AccessTokenExpiresIn:  int(data.ExpiredAt.Unix()),
		RefreshTokenExpiresIn: int(data.ExpiredAt.Unix()),
	}, nil

}

// tokenGenerator implements AuthService.
func (a *authService) tokenGenerator(user *entity.User) (*dto.NewTokenDto, *domain.Error) {
	secret := os.Getenv("JWT_SECRET")
	expiredTime := time.Now().Add(time.Minute * 60)
	claims := &dto.ATClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, signError := token.SignedString([]byte(secret))

	if signError != nil {
		return nil, &domain.Error{
			Err:  signError,
			Code: fiber.StatusInternalServerError,
		}

	}

	expiredRefreshTime := time.Now().Add(time.Hour * 24 * 7)
	refreshToken := utils.GenerateRandomCode(32)

	return &dto.NewTokenDto{
		AccessToken:           accessToken,
		TokenType:             "Bearer",
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  expiredTime,
		RefreshTokenExpiresIn: expiredRefreshTime,
	}, nil
}

// GenerateToken implements AuthService.
func (a *authService) GenerateNewToken(user *entity.User) (*dto.AuthResponse, *domain.Error) {
	newToken, tokenERr := a.tokenGenerator(user)

	if tokenERr != nil {
		return nil, tokenERr
	}

	_, err := a.refreshTokenService.GenerateRefreshToken(&entity.RefreshToken{
		UserID:       user.ID,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiredAt:    &newToken.RefreshTokenExpiresIn,
	})

	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:           newToken.AccessToken,
		TokenType:             newToken.TokenType,
		RefreshToken:          newToken.RefreshToken,
		AccessTokenExpiresIn:  int(newToken.AccessTokenExpiresIn.Unix()),
		RefreshTokenExpiresIn: int(newToken.RefreshTokenExpiresIn.Unix()),
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
	token, err := a.config.GoogleOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: fiber.StatusInternalServerError,
		}
	}
	return token, nil
}

func NewAuthService(
	refreshTokenService refreshtoken.RefreshTokenService,
	config *config.Config,
) AuthService {
	return &authService{
		refreshTokenService,
		config,
	}
}
