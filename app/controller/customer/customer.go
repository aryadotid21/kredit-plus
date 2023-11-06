package customer

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"net/http"

	customerDBModels "kredit-plus/app/db/dto/customer"
	customerDB "kredit-plus/app/db/repository/customer"

	customerProfileDB "kredit-plus/app/db/repository/customer_profile"
	customerTokenDB "kredit-plus/app/db/repository/customer_token"

	"kredit-plus/app/api/middleware/jwt"
	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/dto/request"
	"kredit-plus/app/service/logger"
	"kredit-plus/app/service/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ICustomerController interface {
	CreateCustomer(c *gin.Context)
	GetCustomers(c *gin.Context)
	GetCustomer(c *gin.Context)
	UpdateCustomer(c *gin.Context)
	DeleteCustomer(c *gin.Context)

	CreateCustomerProfile(c *gin.Context)
	GetCustomerProfile(c *gin.Context)
	UpdateCustomerProfile(c *gin.Context)
	DeleteCustomerProfile(c *gin.Context)

	GetCustomerTokens(c *gin.Context)
	GetCustomerToken(c *gin.Context)
	DeleteCustomerToken(c *gin.Context)

	Signup(c *gin.Context)
	Signin(c *gin.Context)
}

type CustomerController struct {
	CustomerDBClient        customerDB.ICustomerRepository
	CustomerProfileDBClient customerProfileDB.ICustomerProfileRepository
	CustomerTokenDBClient   customerTokenDB.ICustomerTokenRepository

	JWT jwt.IJWTService
}

func NewCustomerController(CustomerClient customerDB.ICustomerRepository, CustomerProfileClient customerProfileDB.ICustomerProfileRepository, CustomerTokenClient customerTokenDB.ICustomerTokenRepository, JWT jwt.IJWTService) ICustomerController {
	return &CustomerController{
		CustomerDBClient:        CustomerClient,
		CustomerProfileDBClient: CustomerProfileClient,
		CustomerTokenDBClient:   CustomerTokenClient,
		JWT:                     JWT,
	}
}

func (u CustomerController) CreateCustomer(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	var dataFromBody customerDBModels.Customer
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	hashedPassword, err := util.GenerateHash(dataFromBody.Password)
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	now := time.Now()

	customer := customerDBModels.Customer{
		UUID:      uuid,
		Email:     dataFromBody.Email,
		Phone:     dataFromBody.Phone,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	if err := customer.Validate(); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	if err = u.CustomerDBClient.Create(ctx, &customer); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	customer.Password = ""

	controller.RespondWithSuccess(c, http.StatusOK, constants.CREATED_SUCCESSFULLY, customer, nil)
}

func (u CustomerController) GetCustomers(c *gin.Context) {
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

	if c.Query(customerDBModels.COLUMN_EMAIL) != "" {
		f[customerDBModels.COLUMN_EMAIL] = c.Query(customerDBModels.COLUMN_EMAIL)
	}

	if c.Query(customerDBModels.COLUMN_PHONE) != "" {
		f[customerDBModels.COLUMN_PHONE] = c.Query(customerDBModels.COLUMN_PHONE)
	}

	customers, paginationResponse, err := u.CustomerDBClient.List(ctx, pagination, f)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.GET_SUCCESSFULLY, customers, &paginationResponse)
}

func (u CustomerController) GetCustomer(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param("uuid")
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		customerDBModels.COLUMN_UUID: id,
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

func (u CustomerController) UpdateCustomer(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param("uuid")
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	var dataFromBody customerDBModels.Customer
	if err := c.ShouldBindJSON(&dataFromBody); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	patcher := make(map[string]interface{})

	if dataFromBody.Email != "" {
		patcher[customerDBModels.COLUMN_EMAIL] = dataFromBody.Email
	}

	if dataFromBody.Phone != "" {
		patcher[customerDBModels.COLUMN_PHONE] = dataFromBody.Phone
	}

	patcher[customerDBModels.COLUMN_UPDATED_AT] = time.Now()

	filter := map[string]interface{}{
		customerDBModels.COLUMN_UUID: id,
	}

	if err := u.CustomerDBClient.Update(ctx, filter, patcher); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	customer, err := u.CustomerDBClient.Get(ctx, filter)
	if err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if customer.UUID == uuid.Nil {
		controller.RespondWithError(c, http.StatusNotFound, constants.NOT_FOUND, errors.New(constants.RESOURCE_NOT_FOUND))
		return
	}

	customer.Password = ""

	controller.RespondWithSuccess(c, http.StatusAccepted, constants.UPDATED_SUCCESSFULLY, customer, nil)
}

func (u CustomerController) DeleteCustomer(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	id := c.Param(customerDBModels.COLUMN_UUID)
	if id == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	if _, err := uuid.Parse(id); err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, errors.New(constants.INVALID_INPUT))
		return
	}

	filter := map[string]interface{}{
		customerDBModels.COLUMN_UUID: id,
	}

	if err := u.CustomerDBClient.Delete(ctx, filter); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.DELETED_SUCCESSFULLY, nil, nil)
}
