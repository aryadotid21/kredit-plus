package asset

import (
	"errors"
	"kredit-plus/app/constants"
	"time"
)

const (
	TABLE_NAME         = "assets"
	COLUMN_ID          = "id"
	COLUMN_NAME        = "name"
	COLUMN_TYPE        = "type"
	COLUMN_DESCRIPTION = "description"
	COLUMN_PRICE       = "price"
	COLUMN_CREATED_AT  = "created_at"
	COLUMN_UPDATED_AT  = "updated_at"
)

type Asset struct {
	ID          int        `json:"id"`
	Name        string     `json:"name" form:"name"`
	Type        string     `json:"type" form:"type"`
	Description string     `json:"description" form:"description"`
	Price       int        `json:"price" form:"price"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerToken.
func (u *Asset) Validate() error {
	if u.Name == "" {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.Type == "" {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.Price == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
