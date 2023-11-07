package transaction

import (
	"kredit-plus/app/db/dto/asset"
	"kredit-plus/app/db/dto/transaction"
)

type TransactionDetailResponse struct {
	Transaction transaction.Transaction `json:"transaction"`
	Asset       asset.Asset             `json:"asset"`
}
