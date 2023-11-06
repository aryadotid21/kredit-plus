package server

const (
	// General
	API  = "/api"
	UUID = "/:uuid"
	ID   = "/:id"

	// Auth
	AUTH          = "/auth"
	TOKEN         = "/token"
	SIGNUP        = "/signup"
	SIGNIN        = "/signin"
	SIGNOUT       = "/signout"
	REFRESH_TOKEN = "/refresh-token"

	// Profile
	PROFILE = "/profile"

	// Health Check
	HEALTH_CHECK = "/health-check"

	// Customer
	CUSTOMER = "/customer"
	LIMIT    = "/limit"
)
