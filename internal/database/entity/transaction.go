package entity

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID              int64           `json:"id"`
	Amount          int64           `json:"amount"`
	UserID          int64           `json:"user_id"`
	CategoryID      int64           `json:"category_id"`
	TransactionType TransactionType `json:"transaction_type"`
	TransactionAt   time.Time       `json:"transaction_at"`
	SubcategoryID   *int64          `json:"subcategory_id"`
	Description     *string         `json:"description"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}
