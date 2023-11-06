package customer_limit

import (
	"errors"
	"kredit-plus/app/constants"
	"time"
)

const (
	TABLE_NAME          = "customer_limits"
	COLUMN_ID           = "id"
	COLUMN_CUSTOMER_ID  = "customer_id"
	COLUMN_TENOR        = "tenor"
	COLUMN_LIMIT_AMOUNT = "limit_amount"
	COLUMN_CREATED_AT   = "created_at"
	COLUMN_UPDATED_AT   = "updated_at"
)

type CustomerLimit struct {
	ID          int        `json:"id"`
	CustomerID  int        `json:"customer_id" form:"customer_id"`
	Tenor       int        `json:"tenor" form:"tenor"`
	LimitAmount float32    `json:"limit_amount" form:"limit_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerLimit.
func (u *CustomerLimit) Validate() error {
	if u.CustomerID == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.Tenor == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.LimitAmount == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
