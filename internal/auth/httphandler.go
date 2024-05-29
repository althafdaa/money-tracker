package auth

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	googleConfig *oauth2.Config
	authService  AuthService
}

func (a *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := os.Getenv("GOOGLE_STATE")
	url := a.googleConfig.AuthCodeURL(state)
	c.Status(fiber.StatusFound)
	return c.Redirect(url)
}

func (a *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	stateQ := c.Query("state")
	stateEnv := os.Getenv("GOOGLE_STATE")
	if stateQ != stateEnv {
		err := errors.New("INVALID_STATE")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusUnauthorized,
		})
	}
	code := c.Query("code")
	token, err := a.authService.ExchangeToken(code)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	user, err := a.authService.ParseTokenToUser(token.AccessToken)

	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": user,
	})
}

func NewAuthHandler(googleConfig *oauth2.Config, authService AuthService) *AuthHandler {
	return &AuthHandler{googleConfig, authService}
}
