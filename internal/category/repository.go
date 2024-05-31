package category

import (
	"errors"
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
	err := c.db.Table("category").Create(&entity.Category{
		Name: body.Name,
		Slug: body.Slug,
		Type: body.Type,
	}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, &domain.Error{
				Code: 400,
				Err:  errors.New("CATEGORY_ALREADY_EXISTS"),
			}
		}
		return nil, &domain.Error{
			Err:  err,
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
