package entity

import (
	"time"

	"gorm.io/gorm"
)

type Subcategory struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Slug       string         `json:"slug"`
	CategoryID int            `json:"category_id"`
	UserID     int            `json:"user_id"`
	CreatedAt  *time.Time     `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
}
