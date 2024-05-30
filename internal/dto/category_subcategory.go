package dto

import "money-tracker/internal/database/entity"

type CreateCategoryBody struct {
	Name string              `json:"name" validate:"required,min=3,max20"`
	Type entity.CategoryType `json:"type" validate:"required,category_type"`
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
