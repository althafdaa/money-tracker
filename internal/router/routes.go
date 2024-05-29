package router

import (
	"money-tracker/internal/auth"
	"money-tracker/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth           *auth.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

func (s *HTTP) RegisterFiberRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	v1 := app.Group("/api/v1")

	v1.Get("/restricted", s.authMiddleware.Init, func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	auth := v1.Group("/auth")
	auth.Post("/google", s.auth.AuthGoogle)
}

func NewHTTP(auth *auth.AuthHandler, authMiddleware *middleware.AuthMiddleware) *HTTP {
	return &HTTP{
		auth,
		authMiddleware,
	}
}
