package customer

import "errors"

type SignInRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
}

func (s *SignInRequest) Validate() error {
	if s.Password == "" {
		return errors.New("password is required")
	}

	count := 0
	if s.Email != "" {
		count++
	}
	if s.Phone != "" {
		count++
	}

	// Ensure only one of email, username, phone, or user is provided
	if count > 1 {
		return errors.New("exactly one of email, phone, or user should be provided")
	}

	return nil
}
