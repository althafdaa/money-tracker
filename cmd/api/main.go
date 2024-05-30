package main

import (
	"fmt"
	"money-tracker/internal/auth"
	"money-tracker/internal/category"
	"money-tracker/internal/config"
	"money-tracker/internal/middleware"
	refreshtoken "money-tracker/internal/refresh_token"
	"money-tracker/internal/router"
	"money-tracker/internal/server"
	"money-tracker/internal/user"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()
	server.App.Use(cors.New())
	cfg := config.NewConfig()
	googleCfg := cfg.GoogleOauthConfig()

	authMiddleware := middleware.NewAuthMiddleware()

	categoryRepo := category.NewCategoryRepository(server.Db)
	categoryService := category.NewCategoryService(categoryRepo)
	categoryHandler := category.NewCategoryHandler(categoryService)

	refreshTokenRepo := refreshtoken.NewRefreshTokenRepository(server.Db)
	refreshTokenService := refreshtoken.NewRefreshTokenService(refreshTokenRepo)

	userRepo := user.NewUserRepository(server.Db)
	userService := user.NewUserService(userRepo)

	authService := auth.NewAuthService(googleCfg, refreshTokenService)
	authHandler := auth.NewAuthHandler(googleCfg, authService, server.Validator, userService, refreshTokenService)

	routes := router.NewHTTP(authHandler, categoryHandler, authMiddleware)
	routes.RegisterFiberRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
