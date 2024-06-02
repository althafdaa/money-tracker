package entity

import (
	"time"

	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryIncome  CategoryType = "income"
	CategoryExpense CategoryType = "expense"
)

type Category struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Type      CategoryType   `json:"type"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (Category) TableName() string {
	return "category"
}

type CategoryWithSubcategoryRaw struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	Type      CategoryType   `json:"type"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`

	SubcategoryID        *int            `json:"subcategory_id"`
	SubcategoryName      *string         `json:"subcategory_name"`
	SubcategorySlug      *string         `json:"subcategory_slug"`
	SubcategoryCreatedAt *time.Time      `json:"subcategory_created_at"`
	SubcategoryUpdatedAt *time.Time      `json:"subcategory_updated_at"`
	SubcategoryDeletedAt *gorm.DeletedAt `json:"-"`
}
