package dto

import "github.com/rokafela/udemy-banking-auth/errs"

type LoginRequest struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

type LoginResponse struct {
	Token string
}

func (login_request LoginRequest) Validate() *errs.AppError {
	if login_request.Username == "" || login_request.Password == "" {
		return errs.NewValidationError("Username and password are required")
	}
	return nil
}
