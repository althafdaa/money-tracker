package main

import (
	"fmt"
	"money-tracker/internal"
	"money-tracker/internal/database/seeder"
	"money-tracker/internal/server"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {

	server := server.New()
	server.App.Use(cors.New())

	server.App.Use(requestid.New())
	server.App.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006, 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))

	seedErr := seeder.NewSeeder(server.Db).SeedEverything()
	if seedErr != nil {
		panic(fmt.Sprintf("FAILED_TO_SEED: %s", seedErr))
	}

	internal.
		InitializeServer(server.Db, server.Validator).
		RegisterFiberRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
