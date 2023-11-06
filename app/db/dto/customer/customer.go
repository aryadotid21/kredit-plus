package customer

import (
	"errors"
	"kredit-plus/app/service/util"
	"time"

	"github.com/google/uuid"
)

const (
	TABLE_NAME        = "customers"
	COLUMN_ID         = "id"
	COLUMN_UUID       = "uuid"
	COLUMN_EMAIL      = "email"
	COLUMN_PHONE      = "phone"
	COLUMN_PASSWORD   = "password"
	COLUMN_LAST_LOGIN = "last_login"
	COLUMN_CREATED_AT = "created_at"
	COLUMN_UPDATED_AT = "updated_at"
)

type Customer struct {
	ID        int        `json:"-"`
	UUID      uuid.UUID  `json:"uuid" form:"uuid"`
	Email     string     `json:"email" form:"email"`
	Phone     string     `json:"phone" form:"phone"`
	Password  string     `json:"password,omitempty"`
	LastLogin time.Time  `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (Customer) TableName() string {
	return TABLE_NAME
}

func (f Customer) Validate() error {
	if f.Email == "" {
		return errors.New("email is required")
	}
	if !util.IsValidEmail(f.Email) {
		return errors.New("email is invalid")
	}
	if f.Phone == "" {
		return errors.New("phone is required")
	}
	if f.Password == "" {
		return errors.New("password is required")
	}

	return nil
}
