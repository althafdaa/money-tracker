package category

import (
	"errors"
	"money-tracker/internal/category/subcategory"
	"money-tracker/internal/dto"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService    CategoryService
	subcategoryService subcategory.SubcategoryService
	validator          *validator.Validate
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {

	var body dto.CreateCategoryBody
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

	res, resErr := h.categoryService.CreateCategory(&body)

	if resErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": resErr.Err.Error(),
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
		return c.Status(createError.Code).JSON(fiber.Map{
			"error": createError.Err.Error(),
			"code":  createError.Code,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func (h *CategoryHandler) DeleteSubcategoryByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors.New("INVALID_REQUEST_PARAMS"),
			"code":  fiber.ErrBadRequest,
		})
	}

	newErr := h.subcategoryService.DeleteSubcategoryByID(id)

	if newErr != nil {
		return c.Status(newErr.Code).JSON(fiber.Map{
			"error": newErr.Err.Error(),
			"code":  fiber.ErrInternalServerError,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
	})
}

func (h *CategoryHandler) UpdateSubcategoryByID(c *fiber.Ctx) error {
	type RequestBody struct {
		Name string `json:"name" validate:"required,min=3,max=20"`
	}
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors.New("INVALID_REQUEST_PARAMS"),
			"code":  fiber.ErrBadRequest,
		})
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors.New("INVALID_REQUEST_BODY"),
			"code":  fiber.ErrBadRequest,
		})
	}

	err = h.validator.Struct(body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"code":  fiber.ErrBadRequest,
		})
	}

	res, resErr := h.subcategoryService.UpdateSubcategoryByID(id, &dto.SubcategoryBody{
		Name: body.Name,
	})

	if resErr != nil {
		return c.Status(resErr.Code).JSON(fiber.Map{
			"error": resErr.Err.Error(),
			"code":  resErr.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
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
