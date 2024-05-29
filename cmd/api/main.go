package main

import (
	"fmt"
	"money-tracker/internal/config"
	"money-tracker/internal/server"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	config.GoogleConfig()

	server.RegisterFiberRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
