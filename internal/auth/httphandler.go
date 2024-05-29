package auth

import (
	"errors"
	"money-tracker/internal/dto"
	"money-tracker/internal/user"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	googleConfig *oauth2.Config
	authService  AuthService
	validator    *validator.Validate
	userService  user.UserService
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

func (a *AuthHandler) AuthGoogle(c *fiber.Ctx) error {
	code := new(dto.AuthGoogleBody)
	err := c.BodyParser(code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}
	err = a.validator.Struct(code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.StatusBadRequest,
		})
	}

	googleUserData, userErr := a.authService.ParseTokenToUser(code.Code)

	if userErr != nil {
		return c.Status(userErr.Code).JSON(fiber.Map{
			"error": userErr.Err.Error(),
			"code":  userErr.Code,
		})
	}
	checkedUser, existErr := a.userService.CheckEmail(googleUserData.Email)

	if existErr != nil {
		return c.Status(existErr.Code).JSON(fiber.Map{
			"error": existErr.Err.Error(),
			"code":  existErr.Code,
		})
	}

	if checkedUser != nil {
		token, err := a.authService.GenerateNewToken(checkedUser)
		if err != nil {
			return c.Status(err.Code).JSON(fiber.Map{
				"error": err.Err.Error(),
				"code":  err.Code,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"code": fiber.StatusOK,
			"data": token,
		})
	} else {
		checkedUser, err := a.userService.CreateUserFromGoogle(googleUserData)
		if err != nil {
			return c.Status(err.Code).JSON(fiber.Map{
				"error": err.Err.Error(),
				"code":  err.Code,
			})
		}

		token, err := a.authService.GenerateNewToken(checkedUser)

		if err != nil {
			return c.Status(err.Code).JSON(fiber.Map{
				"error": err.Err.Error(),
				"code":  err.Code,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"code": fiber.StatusCreated,
			"data": token,
		})
	}

}

func NewAuthHandler(
	googleConfig *oauth2.Config,
	authService AuthService,
	validator *validator.Validate,
	userService user.UserService,
) *AuthHandler {
	return &AuthHandler{googleConfig, authService, validator, userService}
}
