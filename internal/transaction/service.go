package transaction

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
)

type TransactionService interface {
	CreateOneIncomeTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	CreateOneExpenseTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	DeleteTransactionByID(transactionID int64) *domain.Error
}
type transactionService struct {
	transactionRepository TransactionRepository
}

// CreateOneExpenseTransaction implements TransactionService.
func (t *transactionService) CreateOneExpenseTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	res, err := t.transactionRepository.CreateTransaction(&entity.Transaction{
		Amount:          transaction.Amount,
		UserID:          transaction.UserID,
		CategoryID:      transaction.CategoryID,
		SubcategoryID:   transaction.SubcategoryID,
		TransactionType: entity.Expense,
		Description:     transaction.Description,
		TransactionAt:   transaction.TransactionAt,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateOneIncomeTransaction implements TransactionService.
func (t *transactionService) CreateOneIncomeTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	res, err := t.transactionRepository.CreateTransaction(&entity.Transaction{
		Amount:          transaction.Amount,
		UserID:          transaction.UserID,
		CategoryID:      transaction.CategoryID,
		SubcategoryID:   transaction.SubcategoryID,
		TransactionType: entity.Income,
		Description:     transaction.Description,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteTransactionByID implements TransactionService.
func (t *transactionService) DeleteTransactionByID(transactionID int64) *domain.Error {
	return t.transactionRepository.DeleteTransactionByID(transactionID)
}

func NewTransactionService(transactionRepository TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepository,
	}
}
