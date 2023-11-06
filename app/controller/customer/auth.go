package customer

import (
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	customerDBModels "kredit-plus/app/db/dto/customer"
	"kredit-plus/app/service/correlation"
	"kredit-plus/app/service/logger"
	"kredit-plus/app/service/util"
	"net/http"

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
