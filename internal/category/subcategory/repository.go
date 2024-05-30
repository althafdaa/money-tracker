package subcategory

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"

	"gorm.io/gorm"
)

type SubcategoryRepository interface {
	CreateOne(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	DeleteById(id int) *domain.Error
	UpdateOne(id int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
}
type subcategoryRepository struct {
	db *gorm.DB
}

// CreateOne implements SubcategoryRepository.
func (s *subcategoryRepository) CreateOne(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	var subcategory entity.Subcategory
	res := s.db.Exec("insert into subcategory (name, slug, category_id, user_id) values (?, ?, ?, ?) returning *", body.Name, body.Slug, body.CategoryID, body.UserID).Scan(&subcategory)

	if res.Error != nil {
		return nil, &domain.Error{
			Err:  res.Error,
			Code: 500,
		}
	}

	return &subcategory, nil
}

// DeleteById implements SubcategoryRepository.
func (s *subcategoryRepository) DeleteById(id int) *domain.Error {
	err := s.db.Exec("delete from subcategory where id = ?", id).Error

	if err != nil {
		return &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	return nil
}

// UpdateOne implements SubcategoryRepository.
func (s *subcategoryRepository) UpdateOne(id int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	var subcategory entity.Subcategory
	res := s.db.Exec("update subcategory set name = ?, slug = ? where id = ?", body.Name, body.Slug, id).Scan(&subcategory)

	if res.Error != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  res.Error,
		}
	}
	return &subcategory, nil
}

func NewSubcategoryRepository(db *gorm.DB) SubcategoryRepository {
	return &subcategoryRepository{
		db,
	}
}
