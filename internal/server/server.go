package server

import (
	"money-tracker/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	*fiber.App
	Db        *gorm.DB
	Validator *validator.Validate
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "money-tracker",
			AppName:      "money-tracker",
		}),
		Db:        database.New(),
		Validator: validator.New(),
	}

	return server
}
