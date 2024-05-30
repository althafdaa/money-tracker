package category

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"
)

type CategoryService interface {
	CreateCategory(body *dto.CreateCategoryBody) (*entity.Category, *domain.Error)
}
type categoryService struct {
	categoryRepo CategoryRepository
}

// CreateSlug implements CategoryService.
func (c *categoryService) CreateCategory(body *dto.CreateCategoryBody) (*entity.Category, *domain.Error) {
	slug, err := utils.Slugify(body.Name)
	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	res, resErr := c.categoryRepo.CreateOne(&dto.CreateCategoryRepoBody{Name: body.Name, Slug: slug, Type: body.Type})

	if resErr != nil {
		return nil, resErr
	}

	return res, nil
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}
