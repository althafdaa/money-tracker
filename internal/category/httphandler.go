package category

import (
	"errors"
	"money-tracker/internal/category/subcategory"
	"money-tracker/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService    CategoryService
	subcategoryService subcategory.SubcategoryService
	validator          *validator.Validate
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	type CreateCategoryBody struct {
		Name string `json:"name" validate:"required,min=3,max20"`
	}

	var body CreateCategoryBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors.New("INVALID_REQUEST_BODY"),
			"code":  fiber.ErrBadRequest,
		})
	}

	res, err := h.categoryService.CreateSlug(body.Name)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
			"code":  fiber.ErrInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func (h *CategoryHandler) CreateSubcategory(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	type CreateSubcategoryBody struct {
		Name       string `json:"name" validate:"required,min=3,max=20"`
		CategoryID int    `json:"category_id" validate:"required,numeric"`
	}

	var body CreateSubcategoryBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors.New("INVALID_REQUEST_BODY"),
			"code":  fiber.ErrBadRequest,
		})
	}

	err := h.validator.Struct(body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.ErrBadRequest,
		})
	}

	res, createError := h.subcategoryService.CreateSubcategory(&dto.SubcategoryBody{
		Name:       body.Name,
		CategoryID: body.CategoryID,
		UserID:     int(user.UserID),
	})

	if createError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": createError,
			"code":  fiber.ErrInternalServerError,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func NewCategoryHandler(
	categoryService CategoryService,
	subcategoryService subcategory.SubcategoryService,
	validator *validator.Validate,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService,
		subcategoryService,
		validator,
	}
}
