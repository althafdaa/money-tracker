package transaction

import (
	"errors"
	"fmt"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	FindAllTransactions(userID int, values *dto.GetAllValueRepository) (*[]entity.PgTransaction, *domain.Error)
	FindAllTransactionsV2(userID int, values *dto.GetAllValueRepository) (*[]entity.Transaction, *domain.Error)
	DeleteTransactionByID(transactionID int) *domain.Error
	GetOneTransactionByID(transactionID int) (*entity.Transaction, *domain.Error)
	UpdateTransactionByID(transactionID int, transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
}
type transactionRepository struct {
	db *gorm.DB
}

// FindAllTransactionsV2 implements TransactionRepository.
func (t *transactionRepository) FindAllTransactionsV2(userID int, values *dto.GetAllValueRepository) (*[]entity.Transaction, *domain.Error) {
	args := make([]interface{}, 0)
	args = append(args, userID)
	args = append(args, values.Limit)
	args = append(args, values.Offset)

	whereClause := "where t.user_id = ? and t.deleted_at is null"

	query := fmt.Sprintf(`
		select
			t.*
			c.*
			s.*
		from
			"transaction" t
		join
			category c on t.category_id = c.id
		left join
			subcategory s on t.subcategory_id = s.id
		%s
		limit ?
		offset ?
	`, whereClause)

	rows, err := t.db.Raw(
		query, args...,
	).Rows()

	if err != nil {
		return nil, &domain.Error{Code: 500, Err: err}
	}

	defer rows.Close()

	var txs []entity.PgTransaction
	for rows.Next() {
		var tx entity.PgTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.Amount,
			&tx.UserID,
			&tx.CreatedAt,
			&tx.UpdatedAt,
			&tx.Category,
		)

		if err != nil {
			println(err)
			return nil, &domain.Error{Code: 500, Err: err}
		}

		txs = append(txs, tx)
	}

	return nil, nil
}

// FindAllTransactions implements TransactionRepository.
func (t *transactionRepository) FindAllTransactions(userID int, values *dto.GetAllValueRepository) (*[]entity.PgTransaction, *domain.Error) {

	var raw []entity.TransactionRaw
	args := make([]interface{}, 0)
	args = append(args, userID)
	args = append(args, values.Limit)
	args = append(args, values.Offset)

	whereClause := "where t.user_id = ? and t.deleted_at is null"

	query := fmt.Sprintf(`
		select
			t.*
			c.name as category_name,
			c.slug as category_slug,
			c.type as category_type,

			s.name as subcategory_name,
			s.slug as subcategory_slug
		from
			"transaction" t
		join
			category c on t.category_id = c.id
		left join
			subcategory s on t.subcategory_id = s.id
		%s
		limit ?
		offset ?
	`, whereClause)

	result := t.db.Exec(query, args...).Scan(&raw)

	if result.Error != nil {
		return nil, &domain.Error{Code: 500, Err: result.Error}
	}
	var txs []entity.PgTransaction
	for _, r := range raw {
		txs = append(txs, entity.PgTransaction{
			ID:              r.ID,
			Amount:          r.Amount,
			UserID:          r.UserID,
			TransactionAt:   r.TransactionAt,
			TransactionType: r.TransactionType,
			Description:     r.Description,
			CreatedAt:       r.CreatedAt,
			UpdatedAt:       r.UpdatedAt,
			DeletedAt:       r.DeletedAt,

			Category: entity.PgCategory{
				ID:   r.CategoryID,
				Name: r.CategoryName,
				Slug: r.CategorySlug,
				Type: r.CategoryType,
			},
			Subcategory: &entity.PgSubcategory{
				ID:   *r.SubcategoryID,
				Name: *r.SubcategoryName,
				Slug: *r.SubcategorySlug,
			},
		})
	}

	return &txs, nil
}

// GetOneTransactionByID implements TransactionRepository.
func (t *transactionRepository) GetOneTransactionByID(transactionID int) (*entity.Transaction, *domain.Error) {
	var data entity.Transaction
	result := t.db.Exec("select * from \"transaction\" where id = ?", transactionID).Scan(&data)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, &domain.Error{Code: 404, Err: result.Error, Info: "RECORD_NOT_FOUND"}
		}
		return nil, &domain.Error{Code: 500, Err: result.Error}
	}

	return &data, nil
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
func (t *transactionRepository) DeleteTransactionByID(transactionID int) *domain.Error {
	now := time.Now()

	err := t.db.Exec("update transaction set deleted_at = ? where id = ?", &now, transactionID).Error
	if err != nil {
		return &domain.Error{Code: 500, Err: err}
	}
	return nil
}

// UpdateTransaction implements TransactionRepository.
func (t *transactionRepository) UpdateTransactionByID(transactionID int, transaction *entity.Transaction) (*entity.Transaction, *domain.Error) {
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
