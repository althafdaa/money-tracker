package transaction

import (
	"fmt"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
	FindAllTransactions(userID int, values *dto.GetAllQueryParams) (*[]entity.TransactionRaw, *domain.Error)
	FindAllTransactionsCount(userID int, values *dto.GetAllQueryParams) (totalDocs int, err *domain.Error)
	findAllTransactionBuildQuery(userID int, values *dto.GetAllQueryParams) (string, string, []interface{})
	GetAllTransactionTotal(userID int, values *dto.GetAllQueryParams) (*entity.TotalTransaction, *domain.Error)
	DeleteTransactionByID(transactionID int) *domain.Error
	GetOneTransactionByID(transactionID int) (*entity.TransactionRaw, *domain.Error)
	UpdateTransactionByID(transactionID int, transaction *entity.Transaction) (*entity.Transaction, *domain.Error)
}
type transactionRepository struct {
	db *gorm.DB
}

// GetAllTransactionTotal implements TransactionRepository.
func (t *transactionRepository) GetAllTransactionTotal(userID int, values *dto.GetAllQueryParams) (*entity.TotalTransaction, *domain.Error) {

	var total *entity.TotalTransaction
	_, whereClause, args := t.findAllTransactionBuildQuery(userID, values)

	query := fmt.Sprintf(`
		select
			sum(t.amount) as total,
			sum(case when t.transaction_type = 'income' then t.amount else 0 end) as total_income,
			sum(case when t.transaction_type = 'expense' then t.amount else 0 end) as total_expense
		from
			"transaction" t
		%s
	`, whereClause)

	err := t.db.Raw(query, args...).Scan(&total).Error

	if err != nil {
		return nil, &domain.Error{Code: 500, Err: err}
	}

	return total, nil
}

// findAllTransactionBuildQuery implements TransactionRepository.
func (t *transactionRepository) findAllTransactionBuildQuery(userID int, values *dto.GetAllQueryParams) (string, string, []interface{}) {
	args := make([]interface{}, 0)
	args = append(args, userID)

	whereClause := "where t.user_id = ? and t.deleted_at is null"

	if values.Type != "" {
		whereClause += " and t.transaction_type = ?"
		args = append(args, values.Type)
	}

	if values.CategoryID != nil {
		whereClause += " and t.category_id in (?)"
		args = append(args, values.CategoryID)
	}

	if values.Search != "" {
		whereClause += " and t.description ilike ?"
		args = append(args, fmt.Sprintf("%%%s%%", values.Search))
	}

	if values.StartedAt != "" && values.EndedAt != "" {
		whereClause += " and t.transaction_at between ? and ?"
		args = append(args, values.StartedAt)
		args = append(args, values.EndedAt)
	}

	query := fmt.Sprintf(`
		SELECT
			t.*,
			c.name as category_name,
			c.slug as category_slug,
			c.type as category_type,
			
			s."name" AS subcategory_name,
			s.slug as subcategory_slug
		FROM
			"transaction" t
		JOIN
			category c ON t.category_id = c.id
		LEFT JOIN
			subcategory s ON t.subcategory_id = s.id
		%s
		ORDER BY
			t.transaction_at DESC
	`, whereClause)

	return query, whereClause, args
}

// findAllTransactionsCount implements TransactionRepository.
func (t *transactionRepository) FindAllTransactionsCount(userID int, values *dto.GetAllQueryParams) (totalDocs int, err *domain.Error) {
	var count int64
	_, whereClause, args := t.findAllTransactionBuildQuery(userID, values)
	query := fmt.Sprintf(`
		select
			count(*)
		from
			"transaction" t
		%s
	`, whereClause)

	result := t.db.Raw(query, args...).Count(&count)

	if result.Error != nil {
		return 0, &domain.Error{Code: 500, Err: result.Error}
	}

	return int(count), nil
}

// FindAllTransactions implements TransactionRepository.
func (t *transactionRepository) FindAllTransactions(userID int, values *dto.GetAllQueryParams) (*[]entity.TransactionRaw, *domain.Error) {

	var raw []entity.TransactionRaw
	query, _, args := t.findAllTransactionBuildQuery(userID, values)
	query += ` 
		limit ?
		offset ?
	`
	args = append(args, values.Limit)
	args = append(args, values.Offset)
	result := t.db.Raw(query, args...).Scan(&raw)

	if result.Error != nil {
		return nil, &domain.Error{Code: 500, Err: result.Error}
	}

	return &raw, nil
}

// GetOneTransactionByID implements TransactionRepository.
func (t *transactionRepository) GetOneTransactionByID(transactionID int) (*entity.TransactionRaw, *domain.Error) {
	var data entity.TransactionRaw
	query := `
	SELECT
		t.*,
		c.name as category_name,
		c.slug as category_slug,
		c.type as category_type,
		c.created_at as category_created_at,
		c.updated_at as category_updated_at,
		
		s."name" AS subcategory_name,
		s.slug as subcategory_slug,
		s.created_at as subcategory_created_at,
		s.updated_at as subcategory_updated_at
	FROM
		"transaction" t
	JOIN
		category c ON t.category_id = c.id
	LEFT JOIN
		subcategory s ON t.subcategory_id = s.id
	WHERE
		t.deleted_at is null and t.id = ?
`
	err := t.db.Raw(query, transactionID).First(&data).Error

	if err != nil {
		return nil, &domain.Error{Code: 500, Err: err}
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

	err := t.db.Raw("update transaction set deleted_at = ? where id = ? and deleted_at is null", &now, transactionID).Error
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
