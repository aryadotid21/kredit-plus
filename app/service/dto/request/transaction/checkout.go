package transaction

type CheckoutRequest struct {
	ContractNumber    string `json:"contract_number" form:"contract_number"`
	OTRAmount         int    `json:"otr_amount" form:"otr_amount"`
	AdminFee          int    `json:"admin_fee" form:"admin_fee"`
	InstallmentAmount int    `json:"installment_amount" form:"installment_amount"`
	InstallmentPeriod int    `json:"installment_period" form:"installment_period"`
	SalesChannel      string `json:"sales_channel" form:"sales_channel"`
	AssetName         string `json:"asset_name" form:"asset_name"`
	AssetType         string `json:"asset_type" form:"asset_type"`
	AssetDescription  string `json:"asset_description" form:"asset_description"`
	AssetPrice        int    `json:"asset_price" form:"asset_price"`
}

func (u *CheckoutRequest) Validate() error {
	return nil
}
