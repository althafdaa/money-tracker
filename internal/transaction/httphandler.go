package transaction

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/dto"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionService TransactionService
	validator          *validator.Validate
}

type createUpdateRequestBody struct {
	Amount          int64                  `json:"amount" validate:"required, numeric"`
	CategoryID      int64                  `json:"category_id" validate:"required, numeric"`
	SubcategoryID   *int64                 `json:"subcategory_id" validate:"numeric"`
	TransactionAt   string                 `json:"transaction_at" validate:"required, datetime"`
	TransactionType entity.TransactionType `json:"transaction_type" validate:"required, oneof=income expense"`
	Description     *string                `json:"description"`
}

func (t *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	var req createUpdateRequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	if err := t.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionTime, err := time.Parse("2006-01-02", req.TransactionAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	res, transactionErr := t.transactionService.CreateOneTransaction(&entity.Transaction{
		Amount:          req.Amount,
		UserID:          user.UserID,
		CategoryID:      req.CategoryID,
		TransactionType: req.TransactionType,
		TransactionAt:   transactionTime,
		SubcategoryID:   req.SubcategoryID,
		Description:     req.Description,
	})

	if transactionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": res,
	})
}

func (t *TransactionHandler) DeleteTransactionByID(c *fiber.Ctx) error {
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionErr := t.transactionService.DeleteTransactionByID(transactionID)
	if transactionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": nil,
	})
}

func (t *TransactionHandler) UpdateTransactionByID(c *fiber.Ctx) error {
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	currentTransaction, transactionErr := t.transactionService.GetOneTransactionByID(transactionID)
	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	user := c.Locals("user").(*dto.ATClaims)

	if user.UserID != currentTransaction.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED", "code": fiber.StatusUnauthorized})
	}

	var req createUpdateRequestBody
	err = c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	if err := t.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionTime, err := time.Parse("2006-01-02", req.TransactionAt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	res, transactionErr := t.transactionService.UpdateTransactionByID(transactionID, &entity.Transaction{
		Amount:          req.Amount,
		UserID:          user.UserID,
		CategoryID:      req.CategoryID,
		TransactionType: req.TransactionType,
		TransactionAt:   transactionTime,
		SubcategoryID:   req.SubcategoryID,
		Description:     req.Description,
	})

	if transactionErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": res,
	})

}

func (t *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	res, transactionErr := t.transactionService.FindAllTransactions(int(user.UserID))
	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": res,
	})
}

func (t *TransactionHandler) GetOneTransactionByID(c *fiber.Ctx) error {
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	res, transactionErr := t.transactionService.GetOneTransactionByID(transactionID)
	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
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
