package router

import (
	"money-tracker/internal/auth"
	"money-tracker/internal/category"
	"money-tracker/internal/middleware"
	"money-tracker/internal/transaction"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth           *auth.AuthHandler
	authMiddleware *middleware.AuthMiddleware
	transaction    *transaction.TransactionHandler
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
	auth.Get("/google/callback", s.auth.GoogleCallback)
	auth.Get("/refresh", s.auth.RefreshToken)
	auth.Post("/google/callback", s.auth.AuthGoogle)

	category := v1.Group("/category", s.authMiddleware.Init)
	category.Post("/", s.category.CreateCategory)
	category.Get("/", s.category.GetAllCategories)
	category.Post("/subcategory", s.category.CreateSubcategory)
	category.Delete("/subcategory/:id", s.category.DeleteSubcategoryByID)
	category.Put("/subcategory/:id", s.category.UpdateSubcategoryByID)

	transaction := v1.Group("/transaction", s.authMiddleware.Init)
	transaction.Post("/", s.transaction.CreateTransaction)
	transaction.Delete("/:id", s.transaction.DeleteTransactionByID)
	transaction.Put("/:id", s.transaction.UpdateTransactionByID)
	transaction.Get("/:id", s.transaction.GetOneTransactionByID)
	transaction.Get("/", s.transaction.GetAllTransactions)
	transaction.Get("/total", s.transaction.GetAllTransactionTotal)
}

func NewHTTP(
	auth *auth.AuthHandler,
	category *category.CategoryHandler,
	transaction *transaction.TransactionHandler,
	authMiddleware *middleware.AuthMiddleware,
) *HTTP {
	return &HTTP{
		auth,
		authMiddleware,
		transaction,
		category,
	}
}
