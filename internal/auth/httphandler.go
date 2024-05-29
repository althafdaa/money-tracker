package auth

import "github.com/gofiber/fiber/v2"

type AuthHandler struct {
}

func (a *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	return nil
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}
