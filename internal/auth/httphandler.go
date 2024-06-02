package auth

import (
	"errors"
	refreshtoken "money-tracker/internal/refresh_token"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService    AuthService
	validator      *validator.Validate
	refreshService refreshtoken.RefreshTokenService
}

func (a *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	url := a.authService.GenerateGoogleLoginUrl()
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
	if code == "" {
		err := errors.New("INVALID_CODE")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	data, err := a.authService.LoginWithGoogle(code)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})
}

func (a *AuthHandler) AuthGoogle(c *fiber.Ctx) error {
	type authRequestBody struct {
		Code string `json:"code" validate:"required"`
	}
	code := new(authRequestBody)
	if err := c.BodyParser(code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	if err := a.validator.Struct(code); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	data, err := a.authService.LoginWithGoogle(code.Code)

	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})
}

func (a *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")

	if refreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "REFRESH_TOKEN_INVALID",
			"code":  fiber.StatusBadRequest})
	}

	data, err := a.authService.RefreshToken(refreshToken)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": data,
	})
}

func (a *AuthHandler) Logout(c *fiber.Ctx) error {
	accToken := c.Get("Authorization")
	accToken = strings.TrimPrefix(accToken, "Bearer ")

	err := a.refreshService.LogoutRefreshToken(accToken)
	if err != nil {
		return c.Status(err.Code).JSON(fiber.Map{
			"error": err.Err.Error(),
			"code":  err.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": "LOGOUT_SUCCESS",
	})
}

func NewAuthHandler(
	authService AuthService,
	validator *validator.Validate,
	refreshTokenService refreshtoken.RefreshTokenService,
) *AuthHandler {
	return &AuthHandler{authService, validator, refreshTokenService}
}
