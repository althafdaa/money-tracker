package router

import (
	"money-tracker/internal/auth"
	"money-tracker/internal/category"
	"money-tracker/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth           *auth.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	category       *category.CategoryHandler
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
	auth.Post("/logout", s.authMiddleware.Init, s.auth.Logout)
	auth.Get("/google", s.auth.GoogleLogin)
	auth.Get("/refresh", s.auth.RefreshToken)
	auth.Post("/google/callback", s.auth.AuthGoogle)

	category := v1.Group("/category", s.authMiddleware.Init)
	category.Post("/", s.category.CreateCategory)
	category.Post("/subcategory", s.category.CreateSubcategory)
}

func NewHTTP(
	auth *auth.AuthHandler,
	category *category.CategoryHandler,
	authMiddleware *middleware.AuthMiddleware,
) *HTTP {
	return &HTTP{
		auth,
		authMiddleware,
		category,
	}
}
