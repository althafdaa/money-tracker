package dto

import "money-tracker/internal/database/entity"

type CreateCategoryBody struct {
	Name string              `json:"name" validate:"required,min=3,max=20"`
	Type entity.CategoryType `json:"type" validate:"required"`
}

type CreateCategoryRepoBody struct {
	Name string              `json:"name"`
	Slug string              `json:"slug"`
	Type entity.CategoryType `json:"type"`
}

type SubcategoryBody struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	CategoryID int    `json:"category_id"`
	UserID     int    `json:"user_id"`
}

type CategoryFilters struct {
	Type entity.CategoryType `json:"type"`
}

type CreateUpdateRequestBodyDto struct {
	Amount        int     `json:"amount" validate:"required,numeric"`
	CategoryID    int     `json:"category_id" validate:"required,numeric"`
	TransactionAt string  `json:"transaction_at" validate:"required,datetime=2006-01-02"`
	SubcategoryID *int    `json:"subcategory_id" validate:"omitempty,numeric"`
	Description   *string `json:"description" validate:"omitempty,max=255"`
}
