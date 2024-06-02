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
	ID              int             `json:"id"`
	Amount          int             `json:"amount"`
	CategoryID      int             `json:"category_id"`
	TransactionType TransactionType `json:"transaction_type"`
	TransactionAt   time.Time       `json:"transaction_at"`
	SubcategoryID   *int            `json:"subcategory_id"`
	Description     *string         `json:"description"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	UserID          int             `json:"-"`
	DeletedAt       gorm.DeletedAt  `json:"-"`
}

type TransactionResponse struct {
	ID              int             `json:"id"`
	Amount          int             `json:"amount"`
	UserID          int             `json:"-"`
	TransactionType TransactionType `json:"transaction_type"`
	TransactionAt   time.Time       `json:"transaction_at"`
	Category        Category        `json:"category"`
	Subcategory     *Subcategory    `json:"subcategory"`
	Description     *string         `json:"description"`
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"-"`
}

type TransactionRaw struct {
	ID     int `json:"id"`
	Amount int `json:"amount"`
	UserID int `json:"user_id"`

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
	DeletedAt       gorm.DeletedAt  `json:"-"`
}
