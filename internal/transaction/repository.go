package transaction

import (
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	DeleteTransactionByID(transactionID int64) *domain.Error
	UpdateTransaction(transactionID int64, transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
}
type transactionRepository struct {
	db *gorm.DB
}

// CreateTransaction implements TransactionRepository.
func (t *transactionRepository) CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	var data entity.Transaction
	result := t.db.Raw("insert into \"transaction\" (amount, user_id, category_id, subcategory_id, transaction_type, transaction_at, description) values (?, ?, ?, ?, ?, ? ,?) returning *", transaction.Amount, transaction.UserID, transaction.CategoryID, transaction.SubcategoryID, transaction.TransactionType, transaction.TransactionAt, transaction.Description).Scan(&data)

	if result.Error != nil {
		return nil, &domain.Error{Code: 500, Err: result.Error}
	}

	return &data, nil
}

// DeleteTransactionByID implements TransactionRepository.
func (t *transactionRepository) DeleteTransactionByID(transactionID int64) *domain.Error {
	now := time.Now()

	err := t.db.Exec("update transaction set deleted_at = ? where id = ?", &now, transactionID).Error
	if err != nil {
		return &domain.Error{Code: 500, Err: err}
	}
	return nil
}

// UpdateTransaction implements TransactionRepository.
func (t *transactionRepository) UpdateTransaction(transactionID int64, transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
	var data entity.Transaction
	result := t.db.Raw("update \"transaction\" set amount = ?, user_id = ?, category_id = ?, subcategory_id = ?, transaction_type = ?, description = ? where id = ? returning *", transaction.Amount, transaction.UserID, transaction.CategoryID, transaction.SubcategoryID, transaction.TransactionType, transaction.Description, transactionID).Scan(&data)

	if result.Error != nil {
		return nil, &domain.Error{Code: 500, Err: result.Error}
	}

	return &data, nil
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}
