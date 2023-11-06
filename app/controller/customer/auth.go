package customer

import (
	"errors"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	customerDBModels "kredit-plus/app/db/dto/customer"
	customerTokenDBModels "kredit-plus/app/db/dto/customer_token"
	"kredit-plus/app/service/correlation"
	customerRequest "kredit-plus/app/service/dto/request/customer"
	"kredit-plus/app/service/logger"
	"kredit-plus/app/service/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u CustomerController) Signup(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Parse the request body
	var dataFromBody customerDBModels.Customer
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Validate the request body
	if err := dataFromBody.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Generate a hashed password
	hashedPassword, err := util.GenerateHash(dataFromBody.Password)
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the password with the hashed password
	dataFromBody.Password = hashedPassword

	// Generate a UUID for the user
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the UUID with the generated UUID
	dataFromBody.UUID = uuid

	// Create the user
	if err := u.CustomerDBClient.Create(ctx, &dataFromBody); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Clear the password from the response
	dataFromBody.Password = ""

	controller.RespondWithSuccess(c, http.StatusCreated, constants.CREATED_SUCCESSFULLY, dataFromBody, nil)
}

func (u CustomerController) Signin(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Parse the request body
	var dataFromBody customerRequest.SignInRequest
	if err := c.BindJSON(&dataFromBody); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Validate the request body
	if err := dataFromBody.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	filter := map[string]interface{}{}
	if dataFromBody.Email != "" {
		filter[customerDBModels.COLUMN_EMAIL] = dataFromBody.Email
	}
	if dataFromBody.Phone != "" {
		filter[customerDBModels.COLUMN_PHONE] = dataFromBody.Phone
	}

	// Get the user from the database
	user, err := u.CustomerDBClient.Get(ctx, filter)
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Check if the password matches
	if !util.ValidatePassword(dataFromBody.Password, user.Password) {
		controller.RespondWithError(c, http.StatusUnauthorized, "Wrong credentials", errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	// Generate a JWT
	token, err := u.JWT.GenerateToken(c, user)
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the user's last login time
	patcher := map[string]interface{}{
		customerDBModels.COLUMN_LAST_LOGIN: time.Now(),
	}

	if err := u.CustomerDBClient.Update(ctx, filter, patcher); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Delete the user's previous token from the database
	if err := u.CustomerTokenDBClient.Delete(ctx, map[string]interface{}{customerTokenDBModels.COLUMN_CUSTOMER_ID: user.ID}); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the user's token in the database with the generated JWT token
	tokenRecord := customerTokenDBModels.CustomerToken{
		CustomerID:            user.ID,
		AccessToken:           token.AccessToken,
		RefreshToken:          token.RefreshToken,
		UserAgent:             c.Request.UserAgent(),
		IPAddress:             c.ClientIP(),
		AccessTokenExpiredAt:  time.Unix(token.AtExpires, 0),
		RefreshTokenExpiredAt: time.Unix(token.RtExpires, 0),
	}

	if err := u.CustomerTokenDBClient.Create(ctx, &tokenRecord); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.LOGIN_SUCCESSFULLY, token, nil)
}

func (u CustomerController) Signout(c *gin.Context) {
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

	// Delete the user's token from the database
	if err := u.CustomerTokenDBClient.Delete(ctx, map[string]interface{}{customerTokenDBModels.COLUMN_CUSTOMER_ID: user.ID}); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	controller.RespondWithSuccess(c, http.StatusOK, constants.LOGOUT_SUCCESSFULLY, nil, nil)
}

func (u CustomerController) RefreshToken(c *gin.Context) {
	ctx := correlation.WithReqContext(c)
	log := logger.Logger(ctx)

	// Parse the refresh token from the request
	var refreshTokenRequest customerRequest.RefreshTokenRequest
	if err := c.BindJSON(&refreshTokenRequest); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Validate the request body
	if err := refreshTokenRequest.Validate(); err != nil {
		log.Error(constants.BAD_REQUEST, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BAD_REQUEST, err)
		return
	}

	// Get the token from the database to check if it exists and is still valid
	token, err := u.CustomerTokenDBClient.Get(ctx, map[string]interface{}{customerTokenDBModels.COLUMN_REFRESH_TOKEN: refreshTokenRequest.RefreshToken})
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	if token.ID == 0 {
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	if token.RefreshTokenExpiredAt.Before(time.Now()) {
		controller.RespondWithError(c, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New(constants.UNAUTHORIZED_ACCESS))
		return
	}

	// Parse the refresh token from the request
	refreshToken := refreshTokenRequest.RefreshToken

	// Use the JWT service to refresh the token
	tokenDetails, err := u.JWT.RefreshToken(ctx, refreshToken)
	if err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Delete the user's previous token from the database
	if err := u.CustomerTokenDBClient.Delete(ctx, map[string]interface{}{customerTokenDBModels.COLUMN_CUSTOMER_ID: token.CustomerID}); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Update the user's token in the database with the generated JWT token
	tokenRecord := customerTokenDBModels.CustomerToken{
		CustomerID:            token.CustomerID,
		AccessToken:           tokenDetails.AccessToken,
		RefreshToken:          tokenDetails.RefreshToken,
		UserAgent:             c.Request.UserAgent(),
		IPAddress:             c.ClientIP(),
		AccessTokenExpiredAt:  time.Unix(tokenDetails.AtExpires, 0),
		RefreshTokenExpiredAt: time.Unix(tokenDetails.RtExpires, 0),
	}

	if err := u.CustomerTokenDBClient.Create(ctx, &tokenRecord); err != nil {
		log.Errorf(constants.INTERNAL_SERVER_ERROR, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR, err)
		return
	}

	// Respond with the new token details
	controller.RespondWithSuccess(c, http.StatusOK, constants.UPDATED_SUCCESSFULLY, tokenDetails, nil)
}
