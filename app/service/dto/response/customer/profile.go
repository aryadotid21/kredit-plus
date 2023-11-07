package customer

import (
	customerDBModels "kredit-plus/app/db/dto/customer"
	customerLimitDBModels "kredit-plus/app/db/dto/customer_limit"
	customerProfileDBModels "kredit-plus/app/db/dto/customer_profile"
	customerTokenDBModels "kredit-plus/app/db/dto/customer_token"
)

type CustomerProfileDetail struct {
	Customer customerDBModels.Customer               `json:"customer"`
	Profile  customerProfileDBModels.CustomerProfile `json:"profile"`
	Token    customerTokenDBModels.CustomerToken     `json:"token"`
	Limit    []customerLimitDBModels.CustomerLimit   `json:"limit"`
}
