package category

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"
)

type CategoryService interface {
	CreateSlug(name string) (*entity.Category, *domain.Error)
}
type categoryService struct {
	categoryRepo CategoryRepository
}

// CreateSlug implements CategoryService.
func (c *categoryService) CreateSlug(name string) (*entity.Category, *domain.Error) {
	slug, err := utils.Slugify(name)
	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	res, resErr := c.categoryRepo.CreateOne(&dto.CategoryBody{Name: name, Slug: slug})

	if resErr != nil {
		return nil, resErr
	}

	return res, nil
}

func NewCategoryService(categoryRepo CategoryRepository) CategoryService {
	return &categoryService{categoryRepo}
}
