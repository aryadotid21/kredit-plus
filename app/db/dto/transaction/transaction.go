package transaction

import (
	"errors"
	"kredit-plus/app/constants"
	"time"

	"github.com/google/uuid"
)

const (
	TABLE_NAME                = "transactions"
	COLUMN_ID                 = "id"
	COLUMN_UUID               = "uuid"
	COLUMN_CUSTOMER_ID        = "customer_id"
	COLUMN_ASSET_ID           = "asset_id"
	COLUMN_CONTRACT_NUMBER    = "contract_number"
	COLUMN_OTR_AMOUNT         = "otr_amount"
	COLUMN_ADMIN_FEE          = "admin_fee"
	COLUMN_INSTALLMENT_AMOUNT = "installment_amount"
	COLUMN_INSTALLMENT_PERIOD = "installment_period"
	COLUMN_INTEREST_AMOUNT    = "interest_amount"
	COLUMN_CREATED_AT         = "created_at"
	COLUMN_UPDATED_AT         = "updated_at"
)

type Transaction struct {
	ID                int        `json:"-"`
	UUID              uuid.UUID  `json:"uuid" form:"uuid"`
	CustomerID        int        `json:"customer_id" form:"customer_id"`
	AssetID           *int       `json:"asset_id" form:"asset_id"`
	ContractNumber    string     `json:"contract_number" form:"contract_number"`
	OTRAmount         float32    `json:"otr_amount" form:"otr_amount"`
	AdminFee          float32    `json:"admin_fee" form:"admin_fee"`
	InstallmentAmount float32    `json:"installment_amount" form:"installment_amount"`
	InstallmentPeriod int        `json:"installment_period" form:"installment_period"`
	InterestAmount    float32    `json:"interest_amount" form:"interest_amount"`
	SalesChannel      string     `json:"sales_channel" form:"sales_channel"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerToken.
func (u *Transaction) Validate() error {
	if u.CustomerID == 0 {
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

	if u.SalesChannel == "" {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
