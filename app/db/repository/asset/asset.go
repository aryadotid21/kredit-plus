package asset

import (
	"context"
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/db"
	assets_DBModels "kredit-plus/app/db/dto/asset"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/response"
	"kredit-plus/app/service/util"

	"github.com/jinzhu/gorm"
)

// Interface methods for interacting with asset data.
type IAssetRepository interface {
	Create(ctx context.Context, asset *assets_DBModels.Asset) error
	Get(ctx context.Context, filter map[string]interface{}) (assets_DBModels.Asset, error)
	List(ctx context.Context, pagination request.Pagination, filter map[string]interface{}) ([]assets_DBModels.Asset, response.Pagination, error)
	Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error
	Delete(ctx context.Context, filter map[string]interface{}) error
}

type AssetRepository struct {
	DBService *db.DBService
}

// Constructor for creating a new AssetRepository.
func NewAssetRepository(dbService *db.DBService) IAssetRepository {
	return &AssetRepository{
		DBService: dbService,
	}
}

var tableName = assets_DBModels.TABLE_NAME

// Create a new asset record.
func (u *AssetRepository) Create(ctx context.Context, asset *assets_DBModels.Asset) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	if err := tx.Create(asset).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Retrieve a asset based on filter criteria.
func (u *AssetRepository) Get(ctx context.Context, filter map[string]interface{}) (assets_DBModels.Asset, error) {
	tx := u.DBService.GetDB().Table(tableName)
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var asset assets_DBModels.Asset

	if err := tx.Where(filter).First(&asset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return asset, nil
		}
		return asset, err
	}

	return asset, nil
}

// List assets based on filtering and pagination criteria.
func (u *AssetRepository) List(ctx context.Context, paginationRequest request.Pagination, filter map[string]interface{}) (record []assets_DBModels.Asset, paginationResponse response.Pagination, err error) {
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

// Update asset records based on filter criteria and a patch.
func (u *AssetRepository) Update(ctx context.Context, filter map[string]interface{}, patch map[string]interface{}) error {
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

	var asset assets_DBModels.Asset

	if err := tx.Where(filter).First(&asset).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// Delete asset records based on filter criteria.
func (u *AssetRepository) Delete(ctx context.Context, filter map[string]interface{}) error {
	tx := u.DBService.GetDB().Table(tableName).Begin()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where(filter).Delete(&assets_DBModels.Asset{}).Error; err != nil {
		return err
	}

	return tx.Commit().Error
}
