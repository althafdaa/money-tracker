package dto

import (
	"money-tracker/internal/database/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleUserData struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Verified_email bool   `json:"verified_email"`
	Name           string `json:"name"`
	Given_name     string `json:"given_name"`
	Family_name    string `json:"family_name"`
	Picture        string `json:"picture"`
	Locale         string `json:"locale"`
}

type ATClaims struct {
	UserID int
	Email  string
	jwt.RegisteredClaims
}

type ExpirationResponse struct {
	AccessTokenExpiresIn  int `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int `json:"refresh_token_expires_in"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpirationResponse
}

type SelfResponse struct {
	entity.User
	Token ExpirationResponse `json:"token"`
}

type NewTokenDto struct {
	AccessToken           string
	TokenType             string
	RefreshToken          string
	AccessTokenExpiresIn  time.Time
	RefreshTokenExpiresIn time.Time
}
