package subcategory

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"time"

	"gorm.io/gorm"
)

type SubcategoryRepository interface {
	CreateOne(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	DeleteByID(id int) *domain.Error
	UpdateOne(id int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	GetOneByID(id int) (*entity.Subcategory, *domain.Error)
}
type subcategoryRepository struct {
	db *gorm.DB
}

// GetOneByID implements SubcategoryRepository.
func (s *subcategoryRepository) GetOneByID(id int) (*entity.Subcategory, *domain.Error) {
	var subcategory entity.Subcategory
	err := s.db.Table("subcategory").Where("id = ?", id).First(&subcategory).Error

	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	return &subcategory, nil
}

// CreateOne implements SubcategoryRepository.
func (s *subcategoryRepository) CreateOne(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	var subcategory entity.Subcategory
	err := s.db.Table("subcategory").Create(&entity.Subcategory{
		Name:       body.Name,
		CategoryID: body.CategoryID,
		UserID:     body.UserID,
	}).Scan(&subcategory).Error

	if err != nil {
		return nil, &domain.Error{
			Err:  err,
			Code: 500,
		}
	}

	return &subcategory, nil
}

// DeleteById implements SubcategoryRepository.
func (s *subcategoryRepository) DeleteByID(id int) *domain.Error {
	now := time.Now()
	err := s.db.Raw("update subcategory set deleted_at = ? where id = ?", &now, id).Error

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
	res := s.db.Raw("update subcategory set name = ? where id = ? returning *", body.Name, id).Scan(&subcategory)

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
