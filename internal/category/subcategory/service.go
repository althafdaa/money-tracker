package subcategory

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"
)

type SubcategoryService interface {
	CreateSubcategory(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
	DeleteSubcategoryByID(id int) *domain.Error
	UpdateSubcategoryByID(id int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error)
}
type subcategoryService struct {
	subcategoryRepository SubcategoryRepository
}

// CreateSubcategory implements SubcategoryService.
func (s *subcategoryService) CreateSubcategory(body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	slug, err := utils.Slugify(body.Name)
	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	res, resErr := s.subcategoryRepository.CreateOne(&dto.SubcategoryBody{
		Name:       body.Name,
		Slug:       slug,
		CategoryID: body.CategoryID,
		UserID:     body.UserID,
	})

	if resErr != nil {
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
func (s *subcategoryService) UpdateSubcategoryByID(id int, body *dto.SubcategoryBody) (*entity.Subcategory, *domain.Error) {
	updatedSlug, err := utils.Slugify(body.Name)

	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	res, resErr := s.subcategoryRepository.UpdateOne(id, &dto.SubcategoryBody{
		Name: body.Name,
		Slug: updatedSlug,
	})

	if resErr != nil {
		return nil, resErr
	}

	return res, nil

}

func NewSubcategoryService(subcategoryRepository SubcategoryRepository) SubcategoryService {
	return &subcategoryService{
		subcategoryRepository,
	}
}
