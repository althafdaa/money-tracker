//go:build wireinject
// +build wireinject

package internal

import (
	"money-tracker/internal/auth"
	"money-tracker/internal/category"
	"money-tracker/internal/category/subcategory"
	"money-tracker/internal/middleware"
	refreshtoken "money-tracker/internal/refresh_token"
	"money-tracker/internal/router"
	"money-tracker/internal/transaction"
	"money-tracker/internal/user"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func InitializeServer(
	db *gorm.DB,
	validator *validator.Validate,
	googleConfig *oauth2.Config,
) *router.HTTP {
	wire.Build(
		router.NewHTTP,

		middleware.NewAuthMiddleware,

		subcategory.NewSubcategoryRepository,
		subcategory.NewSubcategoryService,

		category.NewCategoryRepository,
		category.NewCategoryService,
		category.NewCategoryHandler,

		auth.NewAuthHandler,
		auth.NewAuthService,

		user.NewUserRepository,
		user.NewUserService,

		refreshtoken.NewRefreshTokenRepository,
		refreshtoken.NewRefreshTokenService,

		transaction.NewTransactionRepository,
		transaction.NewTransactionService,
		transaction.NewTransactionHandler,
	)
	return nil
}
