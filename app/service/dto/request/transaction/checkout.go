package transaction

type CheckoutRequest struct {
	ContractNumber    string  `json:"contract_number" form:"contract_number"`
	OTRAmount         float32 `json:"otr_amount" form:"otr_amount"`
	AdminFee          float32 `json:"admin_fee" form:"admin_fee"`
	InstallmentAmount float32 `json:"installment_amount" form:"installment_amount"`
	InstallmentPeriod int     `json:"installment_period" form:"installment_period"`
	InterestAmount    float32 `json:"interest_amount" form:"interest_amount"`
	SalesChannel      string  `json:"sales_channel" form:"sales_channel"`
	AssetName         string  `json:"asset_name" form:"asset_name"`
	AssetType         string  `json:"asset_type" form:"asset_type"`
	AssetDescription  string  `json:"asset_description" form:"asset_description"`
	AssetPrice        float32 `json:"asset_price" form:"asset_price"`
}

func (u *CheckoutRequest) Validate() error {
	return nil
}
