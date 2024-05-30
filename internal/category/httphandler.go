package category

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService CategoryService
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	type CreateCategoryBody struct {
		Name string `json:"name"`
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

func NewCategoryHandler(
	categoryService CategoryService,
) *CategoryHandler {
	return &CategoryHandler{
		categoryService,
	}
}
