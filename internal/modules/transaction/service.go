package transaction

import (
	"errors"
	"math"
	"money-tracker/internal/database/entity"
	"money-tracker/internal/domain"
	"money-tracker/internal/dto"
	"money-tracker/internal/modules/category"
	"money-tracker/internal/modules/category/subcategory"
	"time"

	"gorm.io/gorm"
)

type TransactionService interface {
	CreateOneTransaction(transaction *dto.CreateUpdateTransaction) (*dto.TransactionOneResponse, *domain.Error)
	UpdateTransactionByID(transactionID int, transaction *dto.CreateUpdateTransaction) (*dto.TransactionOneResponse, *domain.Error)
	GetOneTransactionByID(transactionID int) (*dto.TransactionOneResponse, *domain.Error)
	DeleteTransactionByID(transactionID int) *domain.Error
	GetAllPaginatedTransactions(userID int, query *dto.GetAllQueryParams) (*dto.Pagination[dto.TransactionsWithTotalResponse], *domain.Error)
	GetAllTransactionTotal(userID int, query *dto.GetAllQueryParams) (*entity.TotalTransaction, *domain.Error)
	findAllTransactions(userID int, query *dto.GetAllQueryParams) (*dto.TransactionsResponse, *domain.Error)
	getTransactionPaginationMetadata(totalDocs int, query *dto.GetAllQueryParams) *dto.PaginationMetadata
	generateGetTransactionFilter(query *dto.GetAllQueryParams) (*dto.GetAllQueryParams, *domain.Error)
}
type transactionService struct {
	transactionRepository TransactionRepository
	categoryService       category.CategoryService
	subcategoryService    subcategory.SubcategoryService
}

// GetAllTransactionTotal implements TransactionService.
func (t *transactionService) GetAllTransactionTotal(userID int, query *dto.GetAllQueryParams) (*entity.TotalTransaction, *domain.Error) {
	total, err := t.transactionRepository.GetAllTransactionTotal(userID, query)
	if err != nil {
		return nil, err
	}
	return total, nil
}

// GetAllPaginatedTransactions implements TransactionService.
func (t *transactionService) GetAllPaginatedTransactions(userID int, query *dto.GetAllQueryParams) (*dto.Pagination[dto.TransactionsWithTotalResponse], *domain.Error) {
	totals, totalsErr := t.GetAllTransactionTotal(userID, query)
	if totalsErr != nil {
		return nil, totalsErr
	}
	metadata := t.getTransactionPaginationMetadata(totals.Count, query)

	res, transactionErr := t.findAllTransactions(userID, query)

	if transactionErr != nil {
		return nil, transactionErr
	}

	return &dto.Pagination[dto.TransactionsWithTotalResponse]{
		Code: 200,
		Data: dto.TransactionsWithTotalResponse{
			Transactions: *res,
			Total: dto.Total{
				Total:        totals.Total,
				TotalIncome:  totals.TotalIncome,
				TotalExpense: totals.TotalExpense,
			},
		},
		Metadata: *metadata,
	}, nil
}

// generateGetTransactionFilter implements TransactionService.
func (t *transactionService) generateGetTransactionFilter(query *dto.GetAllQueryParams) (*dto.GetAllQueryParams, *domain.Error) {
	if query.StartedAt != "" && query.EndedAt == "" {
		return nil, &domain.Error{
			Code: 400,
			Err:  errors.New("END_DATE_REQUIRED"),
		}
	}
	if query.EndedAt != "" && query.StartedAt == "" {
		return nil, &domain.Error{
			Code: 400,
			Err:  errors.New("START_DATE_REQUIRED"),
		}
	}

	if query.StartedAt != "" && query.EndedAt != "" {
		startedAt, _ := time.Parse("2006-01-02", query.StartedAt)
		endedAt, _ := time.Parse("2006-01-02", query.EndedAt)

		if startedAt.After(endedAt) {
			return nil, &domain.Error{
				Code: 400,
				Err:  errors.New("START_DATE_AFTER_END_DATE"),
			}
		}
	}

	offset := (query.Page - 1) * query.Limit
	query.Offset = &offset

	return query, nil
}

// getTransactionPaginationMetadata implements TransactionService.
func (t *transactionService) getTransactionPaginationMetadata(totalDocs int, query *dto.GetAllQueryParams) *dto.PaginationMetadata {
	totalPages := int(math.Ceil(float64(totalDocs) / float64(query.Limit)))

	return &dto.PaginationMetadata{
		Page:        query.Page,
		Limit:       query.Limit,
		TotalDocs:   totalDocs,
		TotalPages:  totalPages,
		HasNextPage: query.Page < totalPages,
	}
}

// FindAllTransactions implements TransactionService.
func (t *transactionService) findAllTransactions(userID int, query *dto.GetAllQueryParams) (*dto.TransactionsResponse, *domain.Error) {
	filter, filterErr := t.generateGetTransactionFilter(query)
	if filterErr != nil {
		return nil, filterErr
	}
	res, err := t.transactionRepository.FindAllTransactions(userID, filter)
	if err != nil {
		return nil, err
	}

	var transactions dto.TransactionsResponse
	for _, data := range *res {
		cat := entity.Category{
			ID:   data.CategoryID,
			Name: data.CategoryName,
			Slug: data.CategorySlug,
			Type: data.CategoryType,
		}

		transactions = append(transactions, dto.TransactionListResponse{
			ID:              data.ID,
			Amount:          data.Amount,
			UserID:          data.UserID,
			TransactionType: data.TransactionType,
			TransactionAt:   data.TransactionAt,
			Category:        cat,
			Description:     data.Description,
			CreatedAt:       data.CreatedAt,
			UpdatedAt:       data.UpdatedAt,
			DeletedAt:       data.DeletedAt,
			SubcategoryID:   data.SubcategoryID,
		})

	}

	return &transactions, nil
}

// UpdateTransactionByID implements TransactionService.
func (t *transactionService) UpdateTransactionByID(transactionID int, transaction *dto.CreateUpdateTransaction) (*dto.TransactionOneResponse, *domain.Error) {
	cat, catErr := t.categoryService.GetOneCategoryAndSubcategoryByID(transaction.CategoryID, transaction.SubcategoryID)

	if catErr != nil {
		return nil, catErr
	}

	if entity.TransactionType(cat.Category.Type) == entity.Expense {
		transaction.Amount = transaction.Amount * -1
	}

	res, err := t.transactionRepository.UpdateTransactionByID(transactionID, &entity.Transaction{
		Amount:          transaction.Amount,
		UserID:          transaction.UserID,
		CategoryID:      transaction.CategoryID,
		TransactionAt:   transaction.TransactionAt,
		SubcategoryID:   transaction.SubcategoryID,
		Description:     transaction.Description,
		TransactionType: entity.TransactionType(cat.Category.Type),
	})
	if err != nil {
		return nil, err
	}

	return &dto.TransactionOneResponse{
		ID:              res.ID,
		Amount:          res.Amount,
		UserID:          res.UserID,
		TransactionType: res.TransactionType,
		TransactionAt:   res.TransactionAt,
		Category:        cat.Category,
		Subcategory:     cat.Subcategory,
		Description:     res.Description,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		DeletedAt:       res.DeletedAt,
	}, nil
}

// CreateOneTransaction implements TransactionService.
func (t *transactionService) CreateOneTransaction(transaction *dto.CreateUpdateTransaction) (*dto.TransactionOneResponse, *domain.Error) {
	cat_subcat, catErr := t.categoryService.GetOneCategoryAndSubcategoryByID(transaction.CategoryID, transaction.SubcategoryID)

	if catErr != nil {
		return nil, catErr
	}

	if entity.TransactionType(cat_subcat.Category.Type) == entity.Expense {
		transaction.Amount = transaction.Amount * -1
	}

	res, err := t.transactionRepository.CreateTransaction(&entity.Transaction{
		Amount:          transaction.Amount,
		UserID:          transaction.UserID,
		CategoryID:      transaction.CategoryID,
		TransactionType: entity.TransactionType(cat_subcat.Category.Type),
		TransactionAt:   transaction.TransactionAt,
		SubcategoryID:   transaction.SubcategoryID,
		Description:     transaction.Description,
	})
	if err != nil {
		return nil, err
	}

	return &dto.TransactionOneResponse{
		ID:              res.ID,
		Amount:          res.Amount,
		UserID:          res.UserID,
		TransactionType: res.TransactionType,
		TransactionAt:   res.TransactionAt,
		Category:        cat_subcat.Category,
		Subcategory:     cat_subcat.Subcategory,
		Description:     res.Description,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		DeletedAt:       res.DeletedAt,
	}, nil
}

func (t *transactionService) GetOneTransactionByID(transactionID int) (*dto.TransactionOneResponse, *domain.Error) {
	data, err := t.transactionRepository.GetOneTransactionByID(transactionID)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, &domain.Error{
				Code: 404,
				Err:  errors.New("TRANSACTION_NOT_FOUND"),
			}
		}
		return nil, err

	}
	cat := entity.Category{
		ID:        data.CategoryID,
		Name:      data.CategoryName,
		Slug:      data.CategorySlug,
		Type:      data.CategoryType,
		CreatedAt: data.CategoryCreatedAt,
		UpdatedAt: data.CategoryUpdatedAt,
	}

	var subCat *entity.Subcategory
	if data.SubcategoryID != nil {
		subCat = &entity.Subcategory{
			ID:        *data.SubcategoryID,
			Name:      *data.SubcategoryName,
			CreatedAt: data.SubcategoryCreatedAt,
			UpdatedAt: data.SubcategoryUpdatedAt,
			UserID:    data.UserID,
		}
	}

	return &dto.TransactionOneResponse{
		ID:              data.ID,
		Amount:          data.Amount,
		UserID:          data.UserID,
		TransactionType: data.TransactionType,
		TransactionAt:   data.TransactionAt,
		Category:        cat,
		Subcategory:     subCat,
		Description:     data.Description,
		CreatedAt:       data.CreatedAt,
		UpdatedAt:       data.UpdatedAt,
		DeletedAt:       data.DeletedAt,
	}, nil
}

// DeleteTransactionByID implements TransactionService.
func (t *transactionService) DeleteTransactionByID(transactionID int) *domain.Error {
	err := t.transactionRepository.DeleteTransactionByID(transactionID)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return &domain.Error{
				Code: 404,
				Err:  errors.New("TRANSACTION_NOT_FOUND"),
			}
		}

		return err
	}
	return nil
}

func NewTransactionService(
	transactionRepository TransactionRepository,
	categoryService category.CategoryService,
	subcategoryService subcategory.SubcategoryService,
) TransactionService {
	return &transactionService{
		transactionRepository,
		categoryService,
		subcategoryService,
	}
}
