package constants

import "kredit-plus/config"

var Config *config.ServiceConfig

const (
	//Header constants
	AUTHORIZATION      = "Authorization"
	BEARER             = "Bearer "
	CTK_CLAIM_KEY      = CONTEXT_KEY("claims")
	CORRELATION_KEY_ID = CORRELATION_KEY("X-Correlation-ID")
	STATUS_CODE        = "status_code"
	TIME_NOW           = "time_now"
)

type (
	ENVIRONMENT     string
	CONTEXT_KEY     string
	CORRELATION_KEY string
	ROLE            string
)

var AvailableAccountType = []string{}

var DshopiMODE bool

func (c CONTEXT_KEY) String() string {
	return string(c)
}

func (c CORRELATION_KEY) String() string {
	return string(c)
}

var (
	Local ENVIRONMENT = "local"
	Dev   ENVIRONMENT = "dev"
	Prod  ENVIRONMENT = "prod"
	Stage ENVIRONMENT = "stage"
)

func (c ENVIRONMENT) String() string {
	return string(c)
}

var (
	USER  ROLE = "user"
	ADMIN ROLE = "admin"
)

func (c ROLE) String() string {
	return string(c)
}
