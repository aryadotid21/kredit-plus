package customer

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	customerTokenDBModels "kredit-plus/app/db/dto/customer_token"
	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u CustomerController) GetCustomerTokens(c *gin.Context) {
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

	if c.Query(customerTokenDBModels.COLUMN_CUSTOMER_ID) != "" {
		f[customerTokenDBModels.COLUMN_CUSTOMER_ID] = c.Query(customerTokenDBModels.COLUMN_CUSTOMER_ID)
	}

	if c.Query(customerTokenDBModels.COLUMN_ACCESS_TOKEN) != "" {
		f[customerTokenDBModels.COLUMN_ACCESS_TOKEN] = c.Query(customerTokenDBModels.COLUMN_ACCESS_TOKEN)
	}

	if c.Query(customerTokenDBModels.COLUMN_REFRESH_TOKEN) != "" {
		f[customerTokenDBModels.COLUMN_REFRESH_TOKEN] = c.Query(customerTokenDBModels.COLUMN_REFRESH_TOKEN)
	}

	if c.Query(customerTokenDBModels.COLUMN_USER_AGENT) != "" {
		f[customerTokenDBModels.COLUMN_USER_AGENT] = c.Query(customerTokenDBModels.COLUMN_USER_AGENT)
	}

	if c.Query(customerTokenDBModels.COLUMN_IP_ADDRESS) != "" {
		f[customerTokenDBModels.COLUMN_IP_ADDRESS] = c.Query(customerTokenDBModels.COLUMN_IP_ADDRESS)
	}

	customerTokens, paginationResponse, err := u.CustomerTokenDBClient.List(ctx, pagination, f)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, customerTokens, &paginationResponse)
}

func (u CustomerController) GetCustomerToken(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param("uuid")
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		customerTokenDBModels.COLUMN_ID: id,
	}

	r, err := u.CustomerDBClient.Get(ctx, filter)
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

func (u CustomerController) DeleteCustomerToken(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param(customerTokenDBModels.COLUMN_ID)
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		customerTokenDBModels.COLUMN_ID: id,
	}

	if err := u.CustomerDBClient.Delete(ctx, filter); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.DELETED_SUCCESSFULLY, nil, nil)
}
