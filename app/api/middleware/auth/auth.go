package auth

import (
	"errors"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	customerTokenDBModels "kredit-plus/app/db/dto/customer_token"
	customerTokenDBClient "kredit-plus/app/db/repository/customer_token"

	"kredit-plus/app/api/middleware/jwt"
)

func Authenticated(JWT jwt.IJWTService, customerTokenDBClient customerTokenDBClient.ICustomerTokenRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getHeaderToken(ctx)
		if err != nil {
			controller.RespondWithError(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, err)
			return
		}

		claims, err := JWT.ParseToken(ctx, token)
		if err != nil {
			controller.RespondWithError(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, err)
			return
		}

		filter := map[string]interface{}{
			customerTokenDBModels.COLUMN_ACCESS_TOKEN: token,
		}

		customerToken, err := customerTokenDBClient.Get(ctx, filter)
		if err != nil {
			controller.RespondWithError(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, err)
			return
		}

		if customerToken.ID == 0 {
			controller.RespondWithError(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New("token not found"))
			return
		}

		if customerToken.AccessTokenExpiredAt.Before(time.Now()) {
			controller.RespondWithError(ctx, http.StatusUnauthorized, constants.UNAUTHORIZED_ACCESS, errors.New("token expired"))
			return
		}

		ctx.Set(constants.CTK_CLAIM_KEY.String(), claims.UserUUID)
		ctx.Next()
	}
}

func getHeaderToken(ctx *gin.Context) (string, error) {
	header := string(ctx.GetHeader(constants.AUTHORIZATION))
	return extractToken(header)
}

func extractToken(header string) (string, error) {
	if strings.HasPrefix(header, constants.BEARER) {
		return header[len(constants.BEARER):], nil
	}
	return "", errors.New("token not found")
}
