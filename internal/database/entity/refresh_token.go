package entity

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID           int
	AccessToken  string
	RefreshToken string
	UserID       int
	ExpiredAt    *time.Time
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    gorm.DeletedAt
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
