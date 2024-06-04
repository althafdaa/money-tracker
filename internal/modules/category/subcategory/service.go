package subcategory

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"

	"gorm.io/gorm"
)

type SubcategoryService interface {
	CreateSubcategory(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	DeleteSubcategoryByID(id int) *domain.Error
	UpdateSubcategoryByID(id int, userID int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	GetOneSubcategoryByID(id int) (*entity.Subcategory, *domain.Error)
}
type subcategoryService struct {
	subcategoryRepository SubcategoryRepository
	utils                 utils.Utils
}

// GetOneSubcategoryByID implements SubcategoryService.
func (s *subcategoryService) GetOneSubcategoryByID(id int) (*entity.Subcategory, *domain.Error) {
	res, resErr := s.subcategoryRepository.GetOneByID(id)

	if resErr != nil {
		return nil, resErr
	}

	return res, nil
}

// CreateSubcategory implements SubcategoryService.
func (s *subcategoryService) CreateSubcategory(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	res, resErr := s.subcategoryRepository.CreateOne(&dto.SubcategoryBody{
		Name:       body.Name,
		CategoryID: body.CategoryID,
		UserID:     body.UserID,
	})

	if resErr != nil {
		if errors.Is(resErr.Err, gorm.ErrDuplicatedKey) {
			return nil, &domain.Error{
				Code: 400,
				Err:  errors.New("SUBCATEGORY_ALREADY_EXISTS"),
			}
		}
		return nil, resErr
	}

	return res, nil
}

// DeleteSubcategoryById implements SubcategoryService.
func (s *subcategoryService) DeleteSubcategoryByID(id int) *domain.Error {
	err := s.subcategoryRepository.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateSubcategory implements SubcategoryService.
func (s *subcategoryService) UpdateSubcategoryByID(id int, userID int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	current, currentErr := s.subcategoryRepository.GetOneByID(id)

	if currentErr != nil {
		if errors.Is(currentErr.Err, gorm.ErrRecordNotFound) {
			return nil, &domain.Error{
				Code: 404,
				Err:  errors.New("SUBCATEGORY_NOT_FOUND"),
			}
		}
		return nil, currentErr
	}

	if current.UserID != userID {
		return nil, &domain.Error{
			Code: 403,
			Err:  errors.New("FORBIDDEN"),
		}
	}

	res, resErr := s.subcategoryRepository.UpdateOne(id, &dto.SubcategoryBody{
		Name: body.Name,
	})

	if resErr != nil {
		return nil, resErr
	}

	return res, nil

}

func NewSubcategoryService(subcategoryRepository SubcategoryRepository, utils utils.Utils) SubcategoryService {
	return &subcategoryService{
		subcategoryRepository,
		utils,
	}
}
