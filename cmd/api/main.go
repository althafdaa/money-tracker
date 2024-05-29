package main

import (
	"fmt"
	"money-tracker/internal/auth"
	"money-tracker/internal/config"
	"money-tracker/internal/router"
	"money-tracker/internal/server"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()
	cfg := config.NewConfig()
	googleCfg := cfg.GoogleOauthConfig()
	authService := auth.NewAuthService(googleCfg)
	authHandler := auth.NewAuthHandler(googleCfg, authService)
	routes := router.NewHTTP(authHandler)
	routes.RegisterFiberRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
