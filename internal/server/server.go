package server

import (
	"github.com/gofiber/fiber/v2"

	"money-tracker/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "money-tracker",
			AppName:      "money-tracker",
		}),

		db: database.New(),
	}

	return server
}
