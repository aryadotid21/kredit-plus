package transaction

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"net/http"

	transactionDBModels "kredit-plus/app/db/dto/transaction"
	transactionDB "kredit-plus/app/db/repository/transaction"

	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ITransactionController interface {
	CreateTransaction(c *gin.Context)
	GetTransactions(c *gin.Context)
	GetTransaction(c *gin.Context)
	UpdateTransaction(c *gin.Context)
	DeleteTransaction(c *gin.Context)
}

type TransactionController struct {
	TransactionDBClient transactionDB.ITransactionRepository
}

func NewTransactionController(TransactionClient transactionDB.ITransactionRepository) ITransactionController {
	return &TransactionController{
		TransactionDBClient: TransactionClient,
	}
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
		UUID:        uuid,
		Name:        dataFromBody.Name,
		Type:        dataFromBody.Type,
		Description: dataFromBody.Description,
		Price:       dataFromBody.Price,
		CreatedAt:   now,
		UpdatedAt:   &now,
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

	if c.Query(transactionDBModels.COLUMN_NAME) != "" {
		f[transactionDBModels.COLUMN_NAME] = c.Query(transactionDBModels.COLUMN_NAME)
	}

	if c.Query(transactionDBModels.COLUMN_TYPE) != "" {
		f[transactionDBModels.COLUMN_TYPE] = c.Query(transactionDBModels.COLUMN_TYPE)
	}

	if c.Query(transactionDBModels.COLUMN_DESCRIPTION) != "" {
		f[transactionDBModels.COLUMN_DESCRIPTION] = c.Query(transactionDBModels.COLUMN_DESCRIPTION)
	}

	if c.Query(transactionDBModels.COLUMN_PRICE) != "" {
		f[transactionDBModels.COLUMN_PRICE] = c.Query(transactionDBModels.COLUMN_PRICE)
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

	if dataFromBody.Name != "" {
		patcher[transactionDBModels.COLUMN_NAME] = dataFromBody.Name
	}

	if dataFromBody.Type != "" {
		patcher[transactionDBModels.COLUMN_TYPE] = dataFromBody.Type
	}

	if dataFromBody.Description != "" {
		patcher[transactionDBModels.COLUMN_DESCRIPTION] = dataFromBody.Description
	}

	if dataFromBody.Price != 0 {
		patcher[transactionDBModels.COLUMN_PRICE] = dataFromBody.Price
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
