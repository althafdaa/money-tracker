package transaction

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
)

type TransactionService interface {
	CreateOneTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	UpdateTransactionByID(transactionID int, transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	GetOneTransactionByID(transactionID int) (*entity.Transaction, *domain.Error)
	DeleteTransactionByID(transactionID int) *domain.Error
	FindAllTransactions(userID int, query *dto.GetAllQueryParams) (*[]entity.PgTransaction, *domain.Error)
}
type transactionService struct {
	transactionRepository TransactionRepository
}

// FindAllTransactions implements TransactionService.
func (t *transactionService) FindAllTransactions(userID int, query *dto.GetAllQueryParams) (*[]entity.PgTransaction, *domain.Error) {
	offset := (query.Page - 1) * query.Limit
	res, err := t.transactionRepository.FindAllTransactions(userID, &dto.GetAllValueRepository{
		Offset: offset,
		Limit:  query.Limit,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateTransactionByID implements TransactionService.
func (t *transactionService) UpdateTransactionByID(transactionID int, transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	if transaction.TransactionType == entity.Expense {
		transaction.Amount = transaction.Amount * -1
	}

	res, err := t.transactionRepository.UpdateTransactionByID(transactionID, transaction)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateOneTransaction implements TransactionService.
func (t *transactionService) CreateOneTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	if transaction.TransactionType == entity.Expense {
		transaction.Amount = transaction.Amount * -1
	}

	res, err := t.transactionRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionService) GetOneTransactionByID(transactionID int) (*entity.Transaction, *domain.Error) {
	return t.transactionRepository.GetOneTransactionByID(transactionID)
}

// DeleteTransactionByID implements TransactionService.
func (t *transactionService) DeleteTransactionByID(transactionID int) *domain.Error {
	return t.transactionRepository.DeleteTransactionByID(transactionID)
}

func NewTransactionService(transactionRepository TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepository,
	}
}
