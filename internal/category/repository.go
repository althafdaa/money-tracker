package category

import "gorm.io/gorm"

type CategoryRepository interface{}
type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(
	db *gorm.DB,
) CategoryRepository {
	return &categoryRepository{
		db,
	}
}
