package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                int64          `json:"id"`
	Name              string         `json:"name"`
	Email             string         `json:"email"`
	Hash              string         `json:"-"`
	ProfilePictureUrl string         `json:"profile_picture_url"`
	CreatedAt         *time.Time     `json:"created_at"`
	UpdatedAt         *time.Time     `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at"`
}
