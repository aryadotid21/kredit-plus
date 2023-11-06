package customer_token

import (
	"errors"
	"kredit-plus/app/constants"
	"time"
)

const (
	TABLE_NAME                      = "customer_tokens"
	COLUMN_ID                       = "id"
	COLUMN_CUSTOMER_ID              = "customer_id"
	COLUMN_ACCESS_TOKEN             = "access_token"
	COLUMN_REFRESH_TOKEN            = "refresh_token"
	COLUMN_USER_AGENT               = "user_agent"
	COLUMN_IP_ADDRESS               = "ip_address"
	COLUMN_ACCESS_TOKEN_EXPIRED_AT  = "access_token_expired_at"
	COLUMN_REFRESH_TOKEN_EXPIRED_AT = "refresh_token_expired_at"
	COLUMN_CREATED_AT               = "created_at"
	COLUMN_UPDATED_AT               = "updated_at"
)

type CustomerToken struct {
	ID                    int        `json:"-"`
	CustomerID            int        `json:"customer_id" form:"customer_id"`
	AccessToken           string     `json:"access_token" form:"access_token"`
	RefreshToken          string     `json:"refresh_token" form:"refresh_token"`
	UserAgent             string     `json:"user_agent" form:"user_agent"`
	IPAddress             string     `json:"ip_address" form:"ip_address"`
	AccessTokenExpiredAt  time.Time  `json:"access_token_expired_at" form:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time  `json:"refresh_token_expired_at" form:"refresh_token_expired_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerToken.
func (u *CustomerToken) Validate() error {
	if u.CustomerID == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if len(u.AccessToken) == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if len(u.RefreshToken) == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.AccessTokenExpiredAt.IsZero() {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.RefreshTokenExpiredAt.IsZero() {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
