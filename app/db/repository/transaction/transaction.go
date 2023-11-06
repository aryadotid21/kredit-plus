package transaction

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	transactions_DBModels "kredit-plus/app/db/dto/transaction"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with transaction data.
type ITransactionRepository interface {
	Create(ctx context.Context, transaction *transactions_DBModels.Transaction) error
	Get(ctx context.Context, filter map[string]interface{}) (transactions_DBModels.Transaction, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]transactions_DBModels.Transaction, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type TransactionRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new TransactionRepository.
func NewTransactionRepository(dbService *db.DBService) ITransactionRepository {
	return &TransactionRepository{
		DBService: dbService,
	}
}

var tableName = transactions_DBModels.TABLE_NAME

// Create a new transaction record.
func (u *TransactionRepository) Create(ctx context.Context, transaction *transactions_DBModels.Transaction) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a transaction based on filter criteria.
func (u *TransactionRepository) Get(ctx context.Context, filter map[string]interface{}) (transactions_DBModels.Transaction, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var transaction transactions_DBModels.Transaction

	if err := tx.Where(filter).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return transaction, nil
		}
		return transaction, err
	}

	return transaction, nil
}

// List transactions based on filtering and pagination criteria.
func (u *TransactionRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []transactions_DBModels.Transaction, paginationResponse response.Pagination, err error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	tx, err = util.ApplyFilterCondition(tx, filter)
	if err != nil {
		return nil, paginationResponse, err
	}

	if err := tx.Count(&paginationResponse.TotalCount).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, paginationResponse, nil
		}
		return nil, paginationResponse, err
	}

	if !paginationRequest.GetAllData {
		offset := (*paginationRequest.Page - 1) * *paginationRequest.Limit
		tx = tx.Limit(*paginationRequest.Limit).Offset(offset)
		paginationResponse.Page = *paginationRequest.Page
		paginationResponse.PerPage = *paginationRequest.Limit
		paginationResponse.TotalPages = (paginationResponse.TotalCount + *paginationRequest.Limit - 1) / *paginationRequest.Limit
	}

	tx = tx.Order(fmt.Sprintf("%s %s", paginationRequest.Sort, paginationRequest.Order))

	if err := tx.Find(&record).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, paginationResponse, nil
		}
		return record, paginationResponse, err
	}

	return record, paginationResponse, nil
}

// Update transaction records based on filter criteria and a patch.
func (u *TransactionRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Updates(patch).Error; err != nil {
		tx.Rollback()
		return err
	}

	var transaction transactions_DBModels.Transaction

	if err := tx.Where(filter).First(&transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete transaction records based on filter criteria.
func (u *TransactionRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&transactions_DBModels.Transaction{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
