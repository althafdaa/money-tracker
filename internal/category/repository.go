package category

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"time"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateOne(body *dto.CreateCategoryRepoBody) (*entity.Category, *domain.Error)
	UpdateOne(id int, body *dto.CreateCategoryRepoBody) (*entity.Category, *domain.Error)
	DeleteOne(id int) *domain.Error
}
type categoryRepository struct {
	db *gorm.DB
}

// CreateOne implements CategoryRepository.
func (c *categoryRepository) CreateOne(body *dto.CreateCategoryRepoBody) (*entity.Category, *domain.Error) {
	var category entity.Category
	res := c.db.Exec("insert into category (name, slug, type) values (?, ?, ?) returning *", body.Name, body.Slug, body.Type).Scan(&category)

	if res.Error != nil {
		return nil, &domain.Error{
			Err:  res.Error,
			Code: 500,
		}
	}

	return &category, nil
}

// DeleteOne implements CategoryRepository.
func (c *categoryRepository) DeleteOne(id int) *domain.Error {
	now := time.Now()
	err := c.db.Raw("update category set deleted_at = ? where id = ?", &now, id).Error

	if err != nil {
		return &domain.Error{
			Code: 500,
			Err:  err,
		}
	}
	return nil
}

// UpdateOne implements CategoryRepository.
func (c *categoryRepository) UpdateOne(id int, body *dto.CreateCategoryRepoBody) (*entity.Category, *domain.Error) {
	var category entity.Category
	res := c.db.Exec("update category set name = ?, slug = ? where id = ?", body.Name, body.Slug, id).Scan(&category)

	if res.Error != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}
	return &category, nil
}

func NewCategoryRepository(
	db *gorm.DB,
) CategoryRepository {
	return &categoryRepository{
		db,
	}
}
