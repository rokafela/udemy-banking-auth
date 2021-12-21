package dto

import "github.com/rokafela/udemy-banking-auth/errs"

type LoginRequest struct {
	Username string `validate:"required,alphanum"`
	Password string `validate:"required,alphanum"`
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
