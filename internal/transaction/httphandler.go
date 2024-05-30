package transaction

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/dto"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService TransactionService
	validator          *validator.Validate
}

func (t *TransactionHandler) CreateOneIncomeTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	type RequestBody struct {
		Amount        int64  `json:"amount" validate:"required, numeric"`
		Description   string `json:"description"`
		CategoryID    int64  `json:"category_id" validate:"required, numeric"`
		SubcategoryID int64  `json:"subcategory_id" validate:"numeric"`
		TransactionAt string `json:"transaction_at" validate:"required, datetime"`
	}

	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": 500})
	}

	if err := t.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionTime, err := time.Parse("2006-01-02", req.TransactionAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}
	res, transactionErr := t.transactionService.CreateOneIncomeTransaction(&entity.Transaction{
		Amount:          req.Amount,
		UserID:          user.UserID,
		CategoryID:      req.CategoryID,
		SubcategoryID:   &req.SubcategoryID,
		TransactionType: entity.Income,
		Description:     &req.Description,
		TransactionAt:   transactionTime,
	})

	if transactionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func (t *TransactionHandler) CreateOneExpenseTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	type RequestBody struct {
		Amount        int64  `json:"amount" validate:"required, numeric"`
		Description   string `json:"description"`
		CategoryID    int64  `json:"category_id" validate:"required, numeric"`
		SubcategoryID int64  `json:"subcategory_id" validate:"numeric"`
		TransactionAt string `json:"transaction_at" validate:"required, datetime"`
	}

	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": 500})
	}

	if err := t.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionTime, err := time.Parse("2006-01-02", req.TransactionAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}
	res, transactionErr := t.transactionService.CreateOneExpenseTransaction(&entity.Transaction{
		Amount:          req.Amount * -1,
		UserID:          user.UserID,
		CategoryID:      req.CategoryID,
		SubcategoryID:   &req.SubcategoryID,
		TransactionType: entity.Income,
		Description:     &req.Description,
		TransactionAt:   transactionTime,
	})

	if transactionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func NewTransactionHandler(
	transactionService TransactionService,
	validator *validator.Validate,
) *TransactionHandler {
	return &TransactionHandler{
		transactionService,
		validator,
	}
}
