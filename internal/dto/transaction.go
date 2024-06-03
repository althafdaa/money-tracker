package dto

import (
	"money-tracker/internal/database/entity"
	"time"

	"gorm.io/gorm"
)

type GetAllQueryParams struct {
	Offset        *int
	Page          int                    `query:"page" validate:"required,numeric,min=1"`
	Limit         int                    `query:"limit" validate:"required,numeric,min=1,max=20"`
	Search        string                 `query:"search" validate:"omitempty"`
	CategoryID    []int                  `query:"category_id" validate:"omitempty"`
	SubcategoryID []int                  `query:"subcategory_id" validate:"omitempty"`
	Type          entity.TransactionType `query:"type" validate:"omitempty,oneof=income expense"`
	StartedAt     string                 `query:"started_at" validate:"omitempty,datetime=2006-01-02"`
	EndedAt       string                 `query:"ended_at" validate:"omitempty,datetime=2006-01-02"`
}

type CreateUpdateTransaction struct {
	Amount        int
	UserID        int
	CategoryID    int
	TransactionAt time.Time
	SubcategoryID *int
	Description   *string
}

type TransactionResponse struct {
	ID              int                    `json:"id"`
	Amount          int                    `json:"amount"`
	UserID          int                    `json:"-"`
	TransactionType entity.TransactionType `json:"transaction_type"`
	TransactionAt   time.Time              `json:"transaction_at"`
	Category        entity.Category        `json:"category"`
	Subcategory     *entity.Subcategory    `json:"subcategory"`
	Description     *string                `json:"description"`
	CreatedAt       *time.Time             `json:"created_at"`
	UpdatedAt       *time.Time             `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `json:"-"`
}

type TransactionsResponse []TransactionResponse

type Total struct {
	Total        int `json:"total"`
	TotalIncome  int `json:"total_income"`
	TotalExpense int `json:"total_expense"`
}

type TransactionsWithTotalResponse struct {
	Transactions TransactionsResponse `json:"transactions"`
	Total        Total                `json:"total"`
}
