package transaction

import (
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

func (t *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	var req dto.CreateUpdateRequestBodyDto
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
	println("before create transaction")
	data, transactionErr := t.transactionService.CreateOneTransaction(&dto.CreateUpdateTransactionDto{
		Amount:        req.Amount,
		UserID:        user.UserID,
		CategoryID:    req.CategoryID,
		TransactionAt: transactionTime,
		SubcategoryID: req.SubcategoryID,
		Description:   req.Description,
	})

	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code": fiber.StatusCreated,
		"data": data,
	})
}

func (t *TransactionHandler) DeleteTransactionByID(c *fiber.Ctx) error {
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	transactionErr := t.transactionService.DeleteTransactionByID(transactionID)
	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
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

	var req dto.CreateUpdateRequestBodyDto
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

	res, transactionErr := t.transactionService.UpdateTransactionByID(transactionID, &dto.CreateUpdateTransactionDto{
		Amount:        req.Amount,
		UserID:        user.UserID,
		CategoryID:    req.CategoryID,
		TransactionAt: transactionTime,
		SubcategoryID: req.SubcategoryID,
		Description:   req.Description,
	})

	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": fiber.StatusOK,
		"data": res,
	})

}

func (t *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	user := c.Locals("user").(*dto.ATClaims)

	var query dto.GetAllQueryParams
	if err := c.QueryParser(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	if err := t.validator.Struct(&query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "code": fiber.StatusBadRequest})
	}

	res, transactionErr := t.transactionService.FindAllTransactions(int(user.UserID), &query)
	if transactionErr != nil {
		return c.Status(transactionErr.Code).JSON(fiber.Map{"error": transactionErr.Err.Error(), "code": transactionErr.Code})
	}

	metadata, metadataErr := t.transactionService.GetTransactionPaginationMetadata(int(user.UserID), &query)

	if metadataErr != nil {
		return c.Status(metadataErr.Code).JSON(fiber.Map{"error": metadataErr.Err.Error(), "code": metadataErr.Code})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":     fiber.StatusOK,
		"data":     res,
		"metadata": metadata,
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

	user := c.Locals("user").(*dto.ATClaims)

	if user.UserID != res.UserID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UNAUTHORIZED", "code": fiber.StatusUnauthorized})
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
