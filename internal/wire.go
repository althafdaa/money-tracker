//go:build wireinject
// +build wireinject

package internal

import (
	"money-tracker/internal/config"
	"money-tracker/internal/middleware"
	"money-tracker/internal/modules/auth"
	"money-tracker/internal/modules/category"
	"money-tracker/internal/modules/category/subcategory"
	refreshtoken "money-tracker/internal/modules/refresh_token"
	"money-tracker/internal/modules/transaction"

	"money-tracker/internal/modules/user"

	"money-tracker/internal/utils"

	"money-tracker/internal/router"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeServer(
	db *gorm.DB,
	validator *validator.Validate,
) *router.HTTP {
	wire.Build(
		router.NewHTTP,

		config.NewConfigInit,
		utils.NewUtils,

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
