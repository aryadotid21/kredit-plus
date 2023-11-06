package transaction

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"net/http"

	transactionDBModels "kredit-plus/app/db/dto/transaction"
	transactionDB "kredit-plus/app/db/repository/transaction"

	customerDBModels "kredit-plus/app/db/dto/customer"
	customerDB "kredit-plus/app/db/repository/customer"

	customerLimitDBModels "kredit-plus/app/db/dto/customer_limit"
	customerLimitDB "kredit-plus/app/db/repository/customer_limit"

	assetDBModels "kredit-plus/app/db/dto/asset"
	assetDB "kredit-plus/app/db/repository/asset"

	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/dto/request/transaction"
	"kredit-plus/app/service/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ITransactionController interface {
	Checkout(c *gin.Context)

	CreateTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
	GetTransaction(c *gin.Context)
	UpdateTransaction(c *gin.Context)
	DeleteTransaction(c *gin.Context)
}

type TransactionController struct {
	TransactionDBClient   transactionDB.ITransactionRepository
	CustomerDBClient      customerDB.ICustomerRepository
	CustomerLimitDBClient customerLimitDB.ICustomerLimitRepository
	AssetDBClient         assetDB.IAssetRepository
}

func NewTransactionController(TransactionClient transactionDB.ITransactionRepository, CustomerClient customerDB.ICustomerRepository, CustomerLimitClient customerLimitDB.ICustomerLimitRepository, AssetClient assetDB.IAssetRepository) ITransactionController {
	return &TransactionController{
		TransactionDBClient:   TransactionClient,
		CustomerDBClient:      CustomerClient,
		CustomerLimitDBClient: CustomerLimitClient,
		AssetDBClient:         AssetClient,
	}
}

func (u TransactionController) Checkout(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Get the user from the context
	userUUID, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Error(constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	// Get the user from the database
	user, err := u.CustomerDBClient.Get(ctx, map[string]interface{}{customerDBModels.COLUMN_UUID: userUUID})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Parse and validate the request body
	var dataFromBody transaction.CheckoutRequest
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	if err := dataFromBody.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Get the customer limit from the database
	customerLimit, err := u.CustomerLimitDBClient.Get(ctx, map[string]interface{}{
		customerLimitDBModels.COLUMN_CUSTOMER_ID: user.ID,
		customerLimitDBModels.COLUMN_TENOR:       dataFromBody.InstallmentPeriod,
	})

	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Check if the customer limit exists and is sufficient
	if customerLimit.ID == 0 || int(customerLimit.LimitAmount) < dataFromBody.InstallmentAmount {
		controller.RespondWithError(c, http.StatusForbidden, constants.FORBIDDEN, errors.New(constants.INSUFFICIENT_LIMIT))
		return
	}

	// Create a new UUID
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	now := time.Now()

	// Create the transaction
	transaction := transactionDBModels.Transaction{
		UUID:              uuid,
		CustomerID:        user.ID,
		ContractNumber:    dataFromBody.ContractNumber,
		OTRAmount:         dataFromBody.OTRAmount,
		AdminFee:          dataFromBody.AdminFee,
		InstallmentAmount: dataFromBody.InstallmentAmount,
		InstallmentPeriod: dataFromBody.InstallmentPeriod,
		CreatedAt:         now,
		UpdatedAt:         &now,
	}

	if err := transaction.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	if err := u.TransactionDBClient.Create(ctx, &transaction); err != nil {
		log.Error(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Create the asset
	asset := assetDBModels.Asset{
		Name:        dataFromBody.AssetName,
		Type:        dataFromBody.AssetType,
		Description: dataFromBody.AssetDescription,
		Price:       dataFromBody.AssetPrice,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	if err := asset.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	if err := u.AssetDBClient.Create(ctx, &asset); err != nil {
		log.Error(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the customer's limit
	patcher := map[string]interface{}{
		customerLimitDBModels.COLUMN_LIMIT_AMOUNT: customerLimit.LimitAmount - float32(dataFromBody.InstallmentAmount),
		customerLimitDBModels.COLUMN_UPDATED_AT:   time.Now(),
	}

	filter := map[string]interface{}{
		customerLimitDBModels.COLUMN_CUSTOMER_ID: user.ID,
		customerLimitDBModels.COLUMN_TENOR:       dataFromBody.InstallmentPeriod,
	}

	if err := u.CustomerLimitDBClient.Update(ctx, filter, patcher); err != nil {
		log.Error(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the transaction with the asset's ID
	patcher = map[string]interface{}{
		transactionDBModels.COLUMN_ASSET_ID:   asset.ID,
		transactionDBModels.COLUMN_UPDATED_AT: time.Now(),
	}

	filter = map[string]interface{}{
		transactionDBModels.COLUMN_UUID: transaction.UUID,
	}

	if err := u.TransactionDBClient.Update(ctx, filter, patcher); err != nil {
		log.Error(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	transaction.AssetID = asset.ID

	controller.RespondWithSuccess(c, http.StatusOK, constants.CREATED_SUCCESSFULLY, transaction, nil)
}

func (u TransactionController) CreateTransaction(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	var dataFromBody transactionDBModels.Transaction
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	now := time.Now()

	transaction := transactionDBModels.Transaction{
		UUID:              uuid,
		CustomerID:        dataFromBody.CustomerID,
		AssetID:           dataFromBody.AssetID,
		ContractNumber:    dataFromBody.ContractNumber,
		OTRAmount:         dataFromBody.OTRAmount,
		AdminFee:          dataFromBody.AdminFee,
		InstallmentAmount: dataFromBody.InstallmentAmount,
		InstallmentPeriod: dataFromBody.InstallmentPeriod,
		CreatedAt:         now,
		UpdatedAt:         &now,
	}

	if err := transaction.Validate(); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	if err = u.TransactionDBClient.Create(ctx, &transaction); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.CREATED_SUCCESSFULLY, transaction, nil)
}

func (u TransactionController) GetTransactions(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	var pagination request.Pagination

	if err := c.ShouldBindQuery(&pagination); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	pagination.Validate()

	f := map[string]interface{}{}

	if c.Query(transactionDBModels.COLUMN_UUID) != "" {
		f[transactionDBModels.COLUMN_UUID] = c.Query(transactionDBModels.COLUMN_UUID)
	}

	if c.Query(transactionDBModels.COLUMN_CUSTOMER_ID) != "" {
		f[transactionDBModels.COLUMN_CUSTOMER_ID] = c.Query(transactionDBModels.COLUMN_CUSTOMER_ID)
	}

	if c.Query(transactionDBModels.COLUMN_ASSET_ID) != "" {
		f[transactionDBModels.COLUMN_ASSET_ID] = c.Query(transactionDBModels.COLUMN_ASSET_ID)
	}

	if c.Query(transactionDBModels.COLUMN_CONTRACT_NUMBER) != "" {
		f[transactionDBModels.COLUMN_CONTRACT_NUMBER] = c.Query(transactionDBModels.COLUMN_CONTRACT_NUMBER)
	}

	if c.Query(transactionDBModels.COLUMN_CONTRACT_NUMBER) != "" {
		f[transactionDBModels.COLUMN_CONTRACT_NUMBER] = c.Query(transactionDBModels.COLUMN_CONTRACT_NUMBER)
	}

	if c.Query(transactionDBModels.COLUMN_OTR_AMOUNT) != "" {
		f[transactionDBModels.COLUMN_OTR_AMOUNT] = c.Query(transactionDBModels.COLUMN_OTR_AMOUNT)
	}

	if c.Query(transactionDBModels.COLUMN_ADMIN_FEE) != "" {
		f[transactionDBModels.COLUMN_ADMIN_FEE] = c.Query(transactionDBModels.COLUMN_ADMIN_FEE)
	}

	if c.Query(transactionDBModels.COLUMN_INSTALLMENT_AMOUNT) != "" {
		f[transactionDBModels.COLUMN_INSTALLMENT_AMOUNT] = c.Query(transactionDBModels.COLUMN_INSTALLMENT_AMOUNT)
	}

	if c.Query(transactionDBModels.COLUMN_INSTALLMENT_PERIOD) != "" {
		f[transactionDBModels.COLUMN_INSTALLMENT_PERIOD] = c.Query(transactionDBModels.COLUMN_INSTALLMENT_PERIOD)
	}

	transactions, paginationResponse, err := u.TransactionDBClient.List(ctx, pagination, f)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, transactions, &paginationResponse)
}

func (u TransactionController) GetTransaction(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param(transactionDBModels.COLUMN_UUID)
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		transactionDBModels.COLUMN_UUID: id,
	}

	r, err := u.TransactionDBClient.Get(ctx, filter)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if r.UUID == uuid.Nil {
		controller.RespondWithError(c, http.StatusNotFound, constants.NOT_FOUND, errors.New(constants.RESOURCE_NOT_FOUND))
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, r, nil)
}

func (u TransactionController) UpdateTransaction(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param(transactionDBModels.COLUMN_UUID)
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	var dataFromBody transactionDBModels.Transaction
	if err := c.ShouldBindJSON(&dataFromBody); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	patcher := make(map[string]interface{})

	if dataFromBody.ContractNumber != "" {
		patcher[transactionDBModels.COLUMN_CONTRACT_NUMBER] = dataFromBody.ContractNumber
	}

	if dataFromBody.OTRAmount != 0 {
		patcher[transactionDBModels.COLUMN_OTR_AMOUNT] = dataFromBody.OTRAmount
	}

	if dataFromBody.AdminFee != 0 {
		patcher[transactionDBModels.COLUMN_ADMIN_FEE] = dataFromBody.AdminFee
	}

	if dataFromBody.InstallmentAmount != 0 {
		patcher[transactionDBModels.COLUMN_INSTALLMENT_AMOUNT] = dataFromBody.InstallmentAmount
	}

	if dataFromBody.InstallmentPeriod != 0 {
		patcher[transactionDBModels.COLUMN_INSTALLMENT_PERIOD] = dataFromBody.InstallmentPeriod
	}

	patcher[transactionDBModels.COLUMN_UPDATED_AT] = time.Now()

	filter := map[string]interface{}{
		transactionDBModels.COLUMN_UUID: id,
	}

	if err := u.TransactionDBClient.Update(ctx, filter, patcher); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	transaction, err := u.TransactionDBClient.Get(ctx, filter)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if transaction.UUID == uuid.Nil {
		controller.RespondWithError(c, http.StatusNotFound, constants.NOT_FOUND, errors.New(constants.RESOURCE_NOT_FOUND))
		return
	}

	controller.RespondWithSuccess(c, http.StatusAccepted, constants.UPDATED_SUCCESSFULLY, transaction, nil)
}

func (u TransactionController) DeleteTransaction(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param(transactionDBModels.COLUMN_UUID)
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		transactionDBModels.COLUMN_UUID: id,
	}

	if err := u.TransactionDBClient.Delete(ctx, filter); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.DELETED_SUCCESSFULLY, nil, nil)
}
