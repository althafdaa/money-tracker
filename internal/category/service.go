package category

import (
	"errors"
	"money-tracker/internal/category/subcategory"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type categoryAndSubcategory struct {
	Category    entity.Category     `json:"category"`
	Subcategory *entity.Subcategory `json:"subcategory"`
}
type CategoryService interface {
	CreateCategory(body *dto.CreateCategoryBody) (*entity.Category, *domain.Error)
	GetOneCategoryByID(id int) (*entity.Category, *domain.Error)
	GetOneCategoryAndSubcategoryByID(catID int, subcatID *int) (*categoryAndSubcategory, *domain.Error)
}
type categoryService struct {
	categoryRepo       CategoryRepository
	subcategoryService subcategory.SubcategoryService
}

// GetOneCategoryAndSubcategoryByID implements CategoryService.
func (c *categoryService) GetOneCategoryAndSubcategoryByID(catID int, subcatID *int) (*categoryAndSubcategory, *domain.Error) {
	println("getting category")
	cat, catErr := c.GetOneCategoryByID(catID)

	if catErr != nil {
		if errors.Is(catErr.Err, gorm.ErrRecordNotFound) {
			return nil, &domain.Error{
				Code: fiber.StatusNotFound,
				Err:  errors.New("CATEGORY_NOT_FOUND"),
			}
		}
		return nil, catErr
	}

	var sub *entity.Subcategory
	if subcatID != nil {
		var subErr *domain.Error
		sub, subErr = c.subcategoryService.GetOneSubcategoryByID(*subcatID)
		if subErr != nil {
			if errors.Is(subErr.Err, gorm.ErrRecordNotFound) {
				return nil, &domain.Error{
					Code: fiber.StatusNotFound,
					Err:  errors.New("SUBCATEGORY_NOT_FOUND"),
				}
			}
			return nil, subErr
		}
	}

	return &categoryAndSubcategory{
		Category:    *cat,
		Subcategory: sub,
	}, nil

}

// GetOneCategoryByID implements CategoryService.
func (c *categoryService) GetOneCategoryByID(id int) (*entity.Category, *domain.Error) {
	res, resErr := c.categoryRepo.GetOne(id)

	if resErr != nil {
		return nil, resErr
	}

	return res, nil
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

func NewCategoryService(categoryRepo CategoryRepository,
	subcategoryService subcategory.SubcategoryService,
) CategoryService {
	return &categoryService{
		categoryRepo,
		subcategoryService,
	}
}
