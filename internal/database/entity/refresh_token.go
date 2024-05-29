package entity

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID           int64
	AccessToken  string
	RefreshToken string
	UserID       int64
	ExpiredAt    *time.Time
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt
}
