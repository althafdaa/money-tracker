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
	FindAll(userID int) (*[]entity.CategoryWithSubcategoryRaw, *domain.Error)
	GetOne(id int) (*entity.Category, *domain.Error)
	DeleteOne(id int) *domain.Error
}
type categoryRepository struct {
	db *gorm.DB
}

// GetOne implements CategoryRepository.
func (c *categoryRepository) GetOne(id int) (*entity.Category, *domain.Error) {
	var category entity.Category
	err := c.db.Table("category").Where("id = ?", id).First(&category).Error

	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &category, nil
}

// FindAll implements CategoryRepository.
func (c *categoryRepository) FindAll(userID int) (*[]entity.CategoryWithSubcategoryRaw, *domain.Error) {
	args := make([]interface{}, 0)
	var category []entity.CategoryWithSubcategoryRaw
	query := `
		select 
			c.*,
			s.id as subcategory_id,
			s."name" as subcategory_name,
			s.slug as subcategory_slug,
			s.created_at as subcategory_created_at,
			s.updated_at as subcategory_updated_at
		FROM
			category c
		LEFT JOIN
			subcategory s ON c.id = s.category_id
		WHERE
			s.user_id = ?
	`
	args = append(args, userID)

	err := c.db.Raw(query, args...).Scan(&category).Error

	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &category, nil
}

// CreateOne implements CategoryRepository.
func (c *categoryRepository) CreateOne(body *dto.CreateCategoryRepoBody) (*entity.Category, *domain.Error) {
	var category entity.Category
	err := c.db.Table("category").Create(&entity.Category{
		Name: body.Name,
		Slug: body.Slug,
		Type: body.Type,
	}).Scan(&category).Error

	if err != nil {
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
