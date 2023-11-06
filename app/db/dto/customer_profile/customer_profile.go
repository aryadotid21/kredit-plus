package customer_profile

import (
	"errors"
	"kredit-plus/app/constants"
	"time"
)

const (
	TABLE_NAME            = "customer_profiles"
	COLUMN_ID             = "id"
	COLUMN_CUSTOMER_ID    = "customer_id"
	COLUMN_NIK            = "nik"
	COLUMN_FULL_NAME      = "full_name"
	COLUMN_LEGAL_NAME     = "legal_name"
	COLUMN_PLACE_OF_BIRTH = "place_of_birth"
	COLUMN_DATE_OF_BIRTH  = "date_of_birth"
	COLUMN_SALARY         = "salary"
	COLUMN_KTP_IMAGE      = "ktp_image"
	COLUMN_SELFIE_IMAGE   = "selfie_image"
	COLUMN_CREATED_AT     = "created_at"
	COLUMN_UPDATED_AT     = "updated_at"
)

type CustomerProfile struct {
	ID           int        `json:"id"`
	CustomerID   int        `json:"customer_id" form:"customer_id"`
	NIK          string     `json:"nik" form:"nik"`
	FullName     string     `json:"full_name" form:"full_name"`
	LegalName    string     `json:"legal_name" form:"legal_name"`
	PlaceOfBirth string     `json:"place_of_birth" form:"place_of_birth"`
	DateOfBirth  string     `json:"date_of_birth" form:"date_of_birth"`
	Salary       float32    `json:"salary" form:"salary"`
	KtpImage     string     `json:"ktp_image" form:"ktp_image"`
	SelfieImage  string     `json:"selfie_image" form:"selfie_image"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// Validate the fields of a customerProfile.
func (u *CustomerProfile) Validate() error {
	if u.CustomerID == 0 {
		return errors.New(constants.INVALID_INPUT)
	}

	return nil
}
