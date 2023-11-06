package customer_profile

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	customerProfileDBModels "kredit-plus/app/db/dto/customer_profile"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with customerProfile profile data.
type ICustomerProfileRepository interface {
	Create(ctx context.Context, customerProfile *customerProfileDBModels.CustomerProfile) error
	Get(ctx context.Context, filter map[string]interface{}) (customerProfileDBModels.CustomerProfile, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]customerProfileDBModels.CustomerProfile, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type CustomerProfileRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new CustomerProfileRepository.
func NewCustomerProfileRepository(dbService *db.DBService) ICustomerProfileRepository {
	return &CustomerProfileRepository{
		DBService: dbService,
	}
}

const tableName = customerProfileDBModels.TABLE_NAME

// Create a new customerProfile record.
func (u *CustomerProfileRepository) Create(ctx context.Context, customerProfile *customerProfileDBModels.CustomerProfile) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(customerProfile).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a customerProfile based on filter criteria.
func (u *CustomerProfileRepository) Get(ctx context.Context, filter map[string]interface{}) (customerProfileDBModels.CustomerProfile, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var customerProfile customerProfileDBModels.CustomerProfile

	if err := tx.Where(filter).First(&customerProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerProfile, nil
		}
		return customerProfile, err
	}

	return customerProfile, nil
}

// List customers based on filtering and pagination criteria.
func (u *CustomerProfileRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []customerProfileDBModels.CustomerProfile, paginationResponse response.Pagination, err error) {
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

// Update customerProfile records based on filter criteria and a patch.
func (u *CustomerProfileRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
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

	var customerProfile customerProfileDBModels.CustomerProfile

	if err := tx.Where(filter).First(&customerProfile).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete customerProfile records based on filter criteria.
func (u *CustomerProfileRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&customerProfileDBModels.CustomerProfile{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
