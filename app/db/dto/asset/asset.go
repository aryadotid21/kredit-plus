package transaction

import (
	"errors"
	"kredit-plus/app/constants"
	"time"
)

const (
	TABLE_NAME                = "transactions"
	COLUMN_ID                 = "id"
	COLUMN_CUSTOMER_ID        = "customer_id"
	COLUMN_ASSET_ID           = "asset_id"
	COLUMN_CONTRACT_NUMBER    = "contract_number"
	COLUMN_OTR                = "otr"
	COLUMN_ADMIN_FEE          = "admin_fee"
	COLUMN_INSTALLMENT_AMOUNT = "installment_amount"
	COLUMN_INSTALLMENT_PERIOD = "installment_period"
	COLUMN_CREATED_AT         = "created_at"
	COLUMN_UPDATED_AT         = "updated_at"
)

type Asset struct {
	ID                int        `json:"id"`
	CustomerID        int        `json:"customer_id" form:"customer_id"`
	AssetID           int        `json:"asset_id" form:"asset_id"`
	ContractNumber    string     `json:"contract_number" form:"contract_number"`
	OTR               int        `json:"otr" form:"otr"`
	AdminFee          int        `json:"admin_fee" form:"admin_fee"`
	InstallmentAmount int        `json:"installment_amount" form:"installment_amount"`
	InstallmentPeriod int        `json:"installment_period" form:"installment_period"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerToken.
func (u *Asset) Validate() error {
	if u.CustomerID == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.AssetID == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.ContractNumber == "" {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.AdminFee == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	if u.InstallmentPeriod == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
