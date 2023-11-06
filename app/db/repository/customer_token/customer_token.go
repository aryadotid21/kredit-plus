package customer_token

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	customerTokenDBModels "kredit-plus/app/db/dto/customer_token"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with customerToken token data.
type ICustomerTokenRepository interface {
	Create(ctx context.Context, customerToken *customerTokenDBModels.CustomerToken) error
	Get(ctx context.Context, filter map[string]interface{}) (customerTokenDBModels.CustomerToken, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]customerTokenDBModels.CustomerToken, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type CustomerTokenRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new CustomerTokenRepository.
func NewCustomerTokenRepository(dbService *db.DBService) ICustomerTokenRepository {
	return &CustomerTokenRepository{
		DBService: dbService,
	}
}

const tableName = customerTokenDBModels.TABLE_NAME

// Create a new customerToken record.
func (u *CustomerTokenRepository) Create(ctx context.Context, customerToken *customerTokenDBModels.CustomerToken) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(customerToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a customerToken based on filter criteria.
func (u *CustomerTokenRepository) Get(ctx context.Context, filter map[string]interface{}) (customerTokenDBModels.CustomerToken, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var customerToken customerTokenDBModels.CustomerToken

	if err := tx.Where(filter).First(&customerToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerToken, nil
		}
		return customerToken, err
	}

	return customerToken, nil
}

// List customers based on filtering and pagination criteria.
func (u *CustomerTokenRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []customerTokenDBModels.CustomerToken, paginationResponse response.Pagination, err error) {
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

// Update customerToken records based on filter criteria and a patch.
func (u *CustomerTokenRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
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

	var customerToken customerTokenDBModels.CustomerToken

	if err := tx.Where(filter).First(&customerToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete customerToken records based on filter criteria.
func (u *CustomerTokenRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&customerTokenDBModels.CustomerToken{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
