package controller

import (
	"kredit-plus/app/constants"
	"kredit-plus/app/service/dto/response"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, code int, message string, err error) {
	c.Set(constants.STATUS_CODE, code)
	c.AbortWithStatusJSON(code, response.ResponseV3{Success: false, Message: message, Data: err.Error()})
}

func RespondWithSuccess(c *gin.Context, code int, message string, data interface{}, pagination *response.Pagination) {
	c.Set(constants.STATUS_CODE, code)
	response := response.ResponseV3{Success: true, Message: message, Data: data}
	if pagination != nil {
		response.Meta = pagination
	}
	c.JSON(code, response)
}
