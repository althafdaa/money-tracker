package server

import (
	"money-tracker/internal/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	db *gorm.DB
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
