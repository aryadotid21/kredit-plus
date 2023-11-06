package customer

import (
	"errors"
	"fmt"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	customerDBModels "kredit-plus/app/db/dto/customer"
	customerProfileDBModels "kredit-plus/app/db/dto/customer_profile"
	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (u CustomerController) CreateCustomerProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	userUUID, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Error(constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	user, err := u.CustomerDBClient.Get(ctx, map[string]interface{}{customerDBModels.COLUMN_UUID: userUUID})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	var dataFromBody customerProfileDBModels.CustomerProfile
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	now := time.Now()

	customerProfile := customerProfileDBModels.CustomerProfile{
		CustomerID:   user.ID,
		NIK:          dataFromBody.NIK,
		FullName:     dataFromBody.FullName,
		LegalName:    dataFromBody.LegalName,
		PlaceOfBirth: dataFromBody.PlaceOfBirth,
		DateOfBirth:  dataFromBody.DateOfBirth,
		Salary:       dataFromBody.Salary,
		KtpImage:     dataFromBody.KtpImage,
		SelfieImage:  dataFromBody.SelfieImage,
		CreatedAt:    now,
		UpdatedAt:    &now,
	}

	if err := customerProfile.Validate(); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	if err := u.CustomerProfileDBClient.Delete(ctx, map[string]interface{}{customerProfileDBModels.COLUMN_CUSTOMER_ID: user.ID}); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if err := u.CustomerProfileDBClient.Create(ctx, &customerProfile); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.CREATED_SUCCESSFULLY, customerProfile, nil)
}

func (u CustomerController) GetCustomerProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	userUUID, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Error(constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	user, err := u.CustomerDBClient.Get(ctx, map[string]interface{}{customerDBModels.COLUMN_UUID: userUUID})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	filter := map[string]interface{}{
		customerProfileDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	r, err := u.CustomerProfileDBClient.Get(ctx, filter)
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

func (u CustomerController) UpdateCustomerProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	userUUID, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Error(constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	user, err := u.CustomerDBClient.Get(ctx, map[string]interface{}{customerDBModels.COLUMN_UUID: userUUID})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	var dataFromBody customerProfileDBModels.CustomerProfile
	if err := c.ShouldBindJSON(&dataFromBody); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.BAD_REQUEST, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusBadRequest, errorMsg, err)
		return
	}

	patcher := make(map[string]interface{})

	if dataFromBody.NIK != "" {
		patcher[customerProfileDBModels.COLUMN_NIK] = dataFromBody.NIK
	}

	if dataFromBody.FullName != "" {
		patcher[customerProfileDBModels.COLUMN_FULL_NAME] = dataFromBody.FullName
	}

	if dataFromBody.LegalName != "" {
		patcher[customerProfileDBModels.COLUMN_LEGAL_NAME] = dataFromBody.LegalName
	}

	if dataFromBody.PlaceOfBirth != "" {
		patcher[customerProfileDBModels.COLUMN_PLACE_OF_BIRTH] = dataFromBody.PlaceOfBirth
	}

	if dataFromBody.DateOfBirth != "" {
		patcher[customerProfileDBModels.COLUMN_DATE_OF_BIRTH] = dataFromBody.DateOfBirth
	}

	if dataFromBody.Salary != 0 {
		patcher[customerProfileDBModels.COLUMN_SALARY] = dataFromBody.Salary
	}

	if dataFromBody.KtpImage != "" {
		patcher[customerProfileDBModels.COLUMN_KTP_IMAGE] = dataFromBody.KtpImage
	}

	if dataFromBody.SelfieImage != "" {
		patcher[customerProfileDBModels.COLUMN_SELFIE_IMAGE] = dataFromBody.SelfieImage
	}

	patcher[customerDBModels.COLUMN_UPDATED_AT] = time.Now()

	filter := map[string]interface{}{
		customerProfileDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	if err := u.CustomerProfileDBClient.Update(ctx, filter, patcher); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	customer, err := u.CustomerProfileDBClient.Get(ctx, filter)
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

func (u CustomerController) DeleteCustomerProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	userUUID, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Error(constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	user, err := u.CustomerDBClient.Get(ctx, map[string]interface{}{customerDBModels.COLUMN_UUID: userUUID})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	filter := map[string]interface{}{
		customerProfileDBModels.COLUMN_CUSTOMER_ID: user.ID,
	}

	if err := u.CustomerProfileDBClient.Delete(ctx, filter); err != nil {
		errorMsg := fmt.Sprintf("%s: %v", constants.INTERNAL_SERVER_ERROR, err)
		log.Error(errorMsg)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.DELETED_SUCCESSFULLY, nil, nil)
}
