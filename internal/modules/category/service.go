package category

import (
	"errors"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/modules/category/subcategory"
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
	GetAllCategories(userID int) (*[]dto.CategoryWithSubcategory, *domain.Error)
	GetOneCategoryByID(id int) (*entity.Category, *domain.Error)
	GetOneCategoryAndSubcategoryByID(catID int, subcatID *int) (*categoryAndSubcategory, *domain.Error)
}
type categoryService struct {
	categoryRepo       CategoryRepository
	subcategoryService subcategory.SubcategoryService
	utils              utils.Utils
}

// GetAllCategories implements CategoryService.
func (c *categoryService) GetAllCategories(userID int) (*[]dto.CategoryWithSubcategory, *domain.Error) {
	raw, rawErr := c.categoryRepo.FindAll(userID)
	if rawErr != nil {
		return nil, rawErr
	}

	var categories []dto.CategoryWithSubcategory

	for _, r := range *raw {
		var category *dto.CategoryWithSubcategory
		for i := range categories {
			if categories[i].ID == r.ID {
				category = &categories[i]
				break
			}
		}

		if category == nil {
			category = &dto.CategoryWithSubcategory{
				ID:        r.ID,
				Name:      r.Name,
				Slug:      r.Slug,
				Type:      r.Type,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
				DeletedAt: r.DeletedAt,
			}
			categories = append(categories, *category)
		}

		if *r.SubcategoryID != 0 {
			category.Subcategories = append(category.Subcategories, entity.Subcategory{
				ID:         *r.SubcategoryID,
				Name:       *r.SubcategoryName,
				CreatedAt:  r.SubcategoryCreatedAt,
				UpdatedAt:  r.SubcategoryUpdatedAt,
				UserID:     userID,
				CategoryID: r.ID,
			})

		}
	}

	return &categories, nil
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
	slug, err := c.utils.Slugify(body.Name)
	if err != nil {
		return nil, &domain.Error{
			Code: 500,
			Err:  err,
		}
	}

	res, resErr := c.categoryRepo.CreateOne(&dto.CreateCategoryRepoBody{Name: body.Name, Slug: slug, Type: body.Type})

	if resErr != nil {
		if errors.Is(resErr.Err, gorm.ErrDuplicatedKey) {
			return nil, &domain.Error{
				Code: 404,
				Err:  errors.New("CATEGORY_ALREADY_EXISTS"),
			}
		}
		return nil, resErr
	}

	return res, nil
}

func NewCategoryService(categoryRepo CategoryRepository,
	subcategoryService subcategory.SubcategoryService,
	utils utils.Utils,
) CategoryService {
	return &categoryService{
		categoryRepo,
		subcategoryService,
		utils,
	}
}
