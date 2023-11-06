package customer

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	customerDBModels "kredit-plus/app/db/dto/customer"
	customerLimitDBModels "kredit-plus/app/db/dto/customer_limit"
	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (u CustomerController) CreateCustomerLimit(c *gin.Context) {
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

	var dataFromBody customerLimitDBModels.CustomerLimit
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	now := time.Now()

	customerLimit := customerLimitDBModels.CustomerLimit{
		CustomerID:  user.ID,
		Tenor:       dataFromBody.Tenor,
		LimitAmount: dataFromBody.LimitAmount,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	if err := customerLimit.Validate(); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	if err := u.CustomerLimitDBClient.Create(ctx, &customerLimit); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.CREATED_SUCCESSFULLY, customerLimit, nil)
}

func (u CustomerController) GetCustomerLimits(c *gin.Context) {
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

	var pagination request.Pagination

	if err := c.ShouldBindQuery(&pagination); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	pagination.Validate()

	f := map[string]interface{}{}

	f[customerLimitDBModels.COLUMN_CUSTOMER_ID] = user.ID

	if c.Query(customerLimitDBModels.COLUMN_TENOR) != "" {
		f[customerLimitDBModels.COLUMN_TENOR] = c.Query(customerLimitDBModels.COLUMN_TENOR)
	}

	if c.Query(customerLimitDBModels.COLUMN_LIMIT_AMOUNT) != "" {
		f[customerLimitDBModels.COLUMN_LIMIT_AMOUNT] = c.Query(customerLimitDBModels.COLUMN_LIMIT_AMOUNT)
	}

	customerLimits, paginationResponse, err := u.CustomerLimitDBClient.List(ctx, pagination, f)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, customerLimits, &paginationResponse)
}

func (u CustomerController) GetCustomerLimit(c *gin.Context) {
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

	filter := map[string]interface{}{
		customerLimitDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	r, err := u.CustomerLimitDBClient.Get(ctx, filter)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if r.ID == 0 {
		controller.RespondWithError(c, http.StatusNotFound, constants.NOT_FOUND, errors.New(constants.RESOURCE_NOT_FOUND))
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, r, nil)
}

func (u CustomerController) UpdateCustomerLimit(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Get the id param from the url
	id := c.Param("id")
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

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

	var dataFromBody customerLimitDBModels.CustomerLimit
	if err := c.ShouldBindJSON(&dataFromBody); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	patcher := make(map[string]interface{})

	if dataFromBody.LimitAmount != 0 {
		patcher[customerLimitDBModels.COLUMN_LIMIT_AMOUNT] = dataFromBody.LimitAmount
	}

	if dataFromBody.Tenor != 0 {
		patcher[customerLimitDBModels.COLUMN_TENOR] = dataFromBody.Tenor
	}

	patcher[customerDBModels.COLUMN_UPDATED_AT] = time.Now()

	filter := map[string]interface{}{
		customerLimitDBModels.COLUMN_ID:          id,
		customerLimitDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	if err := u.CustomerLimitDBClient.Update(ctx, filter, patcher); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	customer, err := u.CustomerLimitDBClient.Get(ctx, filter)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if customer.ID == 0 {
		controller.RespondWithError(c, http.StatusNotFound, constants.NOT_FOUND, errors.New(constants.RESOURCE_NOT_FOUND))
		return
	}

	controller.RespondWithSuccess(c, http.StatusAccepted, constants.UPDATED_SUCCESSFULLY, customer, nil)
}

func (u CustomerController) DeleteCustomerLimit(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Get the id param from the url
	id := c.Param("id")
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

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

	filter := map[string]interface{}{
		customerLimitDBModels.COLUMN_ID:          id,
		customerLimitDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	if err := u.CustomerLimitDBClient.Delete(ctx, filter); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.DELETED_SUCCESSFULLY, nil, nil)
}
