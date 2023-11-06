package customer_limit

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	customerLimitDBModels "kredit-plus/app/db/dto/customer_limit"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with customerLimit limit data.
type ICustomerLimitRepository interface {
	Create(ctx context.Context, customerLimit *customerLimitDBModels.CustomerLimit) error
	Get(ctx context.Context, filter map[string]interface{}) (customerLimitDBModels.CustomerLimit, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]customerLimitDBModels.CustomerLimit, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type CustomerLimitRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new CustomerLimitRepository.
func NewCustomerLimitRepository(dbService *db.DBService) ICustomerLimitRepository {
	return &CustomerLimitRepository{
		DBService: dbService,
	}
}

const tableName = customerLimitDBModels.TABLE_NAME

// Create a new customerLimit record.
func (u *CustomerLimitRepository) Create(ctx context.Context, customerLimit *customerLimitDBModels.CustomerLimit) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a customerLimit based on filter criteria.
func (u *CustomerLimitRepository) Get(ctx context.Context, filter map[string]interface{}) (customerLimitDBModels.CustomerLimit, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var customerLimit customerLimitDBModels.CustomerLimit

	if err := tx.Where(filter).First(&customerLimit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerLimit, nil
		}
		return customerLimit, err
	}

	return customerLimit, nil
}

// List customers based on filtering and pagination criteria.
func (u *CustomerLimitRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []customerLimitDBModels.CustomerLimit, paginationResponse response.Pagination, err error) {
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

// Update customerLimit records based on filter criteria and a patch.
func (u *CustomerLimitRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
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

	var customerLimit customerLimitDBModels.CustomerLimit

	if err := tx.Where(filter).First(&customerLimit).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete customerLimit records based on filter criteria.
func (u *CustomerLimitRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&customerLimitDBModels.CustomerLimit{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
