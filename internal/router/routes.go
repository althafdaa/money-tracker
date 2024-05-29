package router

import (
	"money-tracker/internal/auth"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth *auth.AuthHandler
}

func (s *HTTP) RegisterFiberRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Post("/google", s.auth.AuthGoogle)
}

func NewHTTP(auth *auth.AuthHandler) *HTTP {
	return &HTTP{
		auth,
	}
}
