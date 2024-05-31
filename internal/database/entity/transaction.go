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

type TransactionRaw struct {
	ID     int64 `json:"id"`
	Amount int64 `json:"amount"`
	UserID int64 `json:"user_id"`

	CategoryID   int          `json:"category_id"`
	CategorySlug string       `json:"category_slug"`
	CategoryName string       `json:"category_name"`
	CategoryType CategoryType `json:"category_type"`

	SubcategoryName *string `json:"subcategory_name"`
	SubcategoryID   *int    `json:"subcategory_id"`
	SubcategorySlug *string `json:"subcategory_slug"`

	TransactionType TransactionType `json:"transaction_type"`
	TransactionAt   time.Time       `json:"transaction_at"`
	Description     *string         `json:"description"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt
}

type PgTransaction struct {
	ID              int64
	Amount          int64
	UserID          int64
	Category        PgCategory
	TransactionType TransactionType
	TransactionAt   time.Time
	Subcategory     *PgSubcategory
	Description     *string
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	DeletedAt       gorm.DeletedAt
}
