package subcategory

import "gorm.io/gorm"

type SubcategoryRepository interface{}
type subcategoryRepository struct {
	db *gorm.DB
}

func NewSubcategoryRepository(db *gorm.DB) SubcategoryRepository {
	return &subcategoryRepository{
		db,
	}
}
