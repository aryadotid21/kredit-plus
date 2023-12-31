package server

import (
	"context"
	"kredit-plus/app/constants"
	"kredit-plus/app/controller/healthcheck"
	"kredit-plus/app/db"
	"kredit-plus/app/service/logger"
	"strings"
	"time"

	"kredit-plus/app/api/middleware/auth"
	"kredit-plus/app/api/middleware/jwt"
	loggerMiddleware "kredit-plus/app/api/middleware/log"

	customerController "kredit-plus/app/controller/customer"
	customerDBClient "kredit-plus/app/db/repository/customer"
	customerLimitDBClient "kredit-plus/app/db/repository/customer_limit"
	customerProfileDBClient "kredit-plus/app/db/repository/customer_profile"
	customerTokenDBClient "kredit-plus/app/db/repository/customer_token"

	transactionController "kredit-plus/app/controller/transaction"
	transactionDBClient "kredit-plus/app/db/repository/transaction"

	assetDBClient "kredit-plus/app/db/repository/asset"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Init(ctx context.Context, dbConnection *db.DBService) *gin.Engine {
	if strings.EqualFold(constants.Config.Environment, "prod") {
		gin.SetMode(gin.ReleaseMode)
	}
	return NewRouter(ctx, dbConnection)

}
func addCSPHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func addReferrerPolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

func addPermissionsPolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Permissions-Policy", "default-src 'none'")
		c.Next()
	}
}

func addFeaturePolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Feature-Policy", "none")
		c.Next()
	}
}

func NewRouter(ctx context.Context, dbConnection *db.DBService) *gin.Engine {
	log := logger.Logger(ctx)

	log.Info("setting up service and controllers")

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(helmet.Default())
	//Content-Security-Policy
	router.Use(addCSPHeader())
	//Referrer-Policy
	router.Use(addReferrerPolicyHeader())
	//Permissions-Policy
	router.Use(addPermissionsPolicyHeader())
	//Feature-Policy
	router.Use(addFeaturePolicyHeader())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", constants.AUTHORIZATION, constants.CORRELATION_KEY_ID.String()},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(func(ctx *gin.Context) {
		ctx.Set(constants.TIME_NOW, time.Now())
		ctx.Next()
	})

	router.Use(uuidInjectionMiddleware())
	router.Use(loggerMiddleware.LoggerMiddleware())

	// DB Clients
	var (
		customerDBClient        = customerDBClient.NewCustomerRepository(dbConnection)
		customerProfileDBClient = customerProfileDBClient.NewCustomerProfileRepository(dbConnection)
		customerTokenDBClient   = customerTokenDBClient.NewCustomerTokenRepository(dbConnection)
		customerLimitDBClient   = customerLimitDBClient.NewCustomerLimitRepository(dbConnection)

		transactionDBClient = transactionDBClient.NewTransactionRepository(dbConnection)
		assetDBClient       = assetDBClient.NewAssetRepository(dbConnection)
	)

	// SERVICES
	var (
		JWT = jwt.NewJWTService()
	)

	// Controller
	var (
		healthCheckController = healthcheck.NewHealthCheckController()

		customerController    = customerController.NewCustomerController(customerDBClient, customerProfileDBClient, customerTokenDBClient, customerLimitDBClient, JWT)
		transactionController = transactionController.NewTransactionController(transactionDBClient, customerDBClient, customerLimitDBClient, assetDBClient)
	)

	v1 := router.Group("/kredit-plus/v1")
	{
		v1.GET(HEALTH_CHECK, healthCheckController.HealthCheck)

		// Customer
		customer := v1.Group(CUSTOMER)
		{
			v1.POST(CUSTOMER+SIGNUP, customerController.Signup)
			v1.POST(CUSTOMER+SIGNIN, customerController.Signin)
			v1.POST(CUSTOMER+SIGNOUT, auth.Authenticated(JWT, customerTokenDBClient), customerController.Signout)
			v1.POST(CUSTOMER+REFRESH_TOKEN, customerController.RefreshToken)

			customer.Use(auth.Authenticated(JWT, customerTokenDBClient))

			customer.GET(PROFILE+DETAIL, customerController.Profile)
			customer.POST(PROFILE, customerController.CreateCustomerProfile)
			customer.GET(PROFILE, customerController.GetCustomerProfile)
			customer.PATCH(PROFILE, customerController.UpdateCustomerProfile)
			customer.DELETE(PROFILE, customerController.DeleteCustomerProfile)

			customer.POST(LIMIT, customerController.CreateCustomerLimit)
			customer.GET(LIMIT, customerController.GetCustomerLimits)
			customer.GET(LIMIT+ID, customerController.GetCustomerLimit)
			customer.PATCH(LIMIT+ID, customerController.UpdateCustomerLimit)
			customer.DELETE(LIMIT+ID, customerController.DeleteCustomerLimit)

			customer.GET(TOKEN, customerController.GetCustomerTokens)
			customer.GET(TOKEN+ID, customerController.GetCustomerToken)
			customer.DELETE(TOKEN, customerController.DeleteCustomerToken)

			customer.POST("", customerController.CreateCustomer)
			customer.GET("", customerController.GetCustomers)
			customer.GET(UUID, customerController.GetCustomer)
			customer.PATCH(UUID, customerController.UpdateCustomer)
			customer.DELETE(UUID, customerController.DeleteCustomer)
		}

		// Transaction
		transaction := v1.Group(TRANSACTION)
		{
			transaction.Use(auth.Authenticated(JWT, customerTokenDBClient))

			transaction.POST("", transactionController.CreateTransaction)
			transaction.GET("", transactionController.GetTransactions)
			transaction.GET(DETAIL, transactionController.GetTransactionsDetail)
			transaction.GET(UUID, transactionController.GetTransaction)
			transaction.PATCH(UUID, transactionController.UpdateTransaction)
			transaction.DELETE(UUID, transactionController.DeleteTransaction)

			transaction.POST(CHECKOUT, transactionController.Checkout)
		}
	}

	return router
}

// uuidInjectionMiddleware injects the request context with a correlation id of type uuid
func uuidInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.GetHeader(string(constants.CORRELATION_KEY_ID))
		if len(correlationId) == 0 {
			correlationID, _ := uuid.NewUUID()
			correlationId = correlationID.String()
			c.Request.Header.Set(constants.CORRELATION_KEY_ID.String(), correlationId)
		}
		c.Writer.Header().Set(constants.CORRELATION_KEY_ID.String(), correlationId)

		c.Next()
	}
}
