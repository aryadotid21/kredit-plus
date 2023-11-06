package customer

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	customers_DBModels "kredit-plus/app/db/dto/customer"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with customer data.
type ICustomerRepository interface {
	Create(ctx context.Context, customer *customers_DBModels.Customer) error
	Get(ctx context.Context, filter map[string]interface{}) (customers_DBModels.Customer, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]customers_DBModels.Customer, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type CustomerRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new CustomerRepository.
func NewCustomerRepository(dbService *db.DBService) ICustomerRepository {
	return &CustomerRepository{
		DBService: dbService,
	}
}

var tableName = customers_DBModels.TABLE_NAME

// Create a new customer record.
func (u *CustomerRepository) Create(ctx context.Context, customer *customers_DBModels.Customer) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(customer).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a customer based on filter criteria.
func (u *CustomerRepository) Get(ctx context.Context, filter map[string]interface{}) (customers_DBModels.Customer, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var customer customers_DBModels.Customer

	if err := tx.Where(filter).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customer, nil
		}
		return customer, err
	}

	return customer, nil
}

// List customers based on filtering and pagination criteria.
func (u *CustomerRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []customers_DBModels.Customer, paginationResponse response.Pagination, err error) {
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

// Update customer records based on filter criteria and a patch.
func (u *CustomerRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
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

	var customer customers_DBModels.Customer

	if err := tx.Where(filter).First(&customer).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete customer records based on filter criteria.
func (u *CustomerRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&customers_DBModels.Customer{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
