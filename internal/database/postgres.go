package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic(fmt.Sprintf("Error connecting to database: %s", err))
	}
	fmt.Println("db connected")
	db.Logger = logger.Default.LogMode((logger.Info))
	return db
}
