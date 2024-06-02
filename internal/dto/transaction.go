package dto

import (
	"money-tracker/internal/database/entity"
	"time"
)

type GetAllQueryParams struct {
	Page          int                    `query:"page" validate:"required,numeric,min=1"`
	Limit         int                    `query:"limit" validate:"required,numeric,min=1,max=20"`
	Search        string                 `query:"search" validate:"omitempty"`
	CategoryID    []int                  `query:"category_id" validate:"omitempty"`
	SubcategoryID []int                  `query:"subcategory_id" validate:"omitempty"`
	Type          entity.TransactionType `query:"type" validate:"omitempty,oneof=income expense"`
	StartedAt     string                 `query:"started_at" validate:"omitempty,datetime=2006-01-02"`
	EndedAt       string                 `query:"ended_at" validate:"omitempty,datetime=2006-01-02"`
}

type FindAllTransactionFilter struct {
	Offset        int
	Limit         int
	Search        string
	CategoryID    []int
	SubcategoryID []int
	Type          entity.TransactionType
	StartedAt     string
	EndedAt       string
}

type CreateUpdateTransactionDto struct {
	Amount        int
	UserID        int
	CategoryID    int
	TransactionAt time.Time
	SubcategoryID *int
	Description   *string
}
