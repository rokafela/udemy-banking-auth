package dto

type LoginRequest struct {
	Username string `validate:"nonzero"`
	Password string `validate:"nonzero"`
}

type LoginResponse struct {
	Username string
	Role     string
	Token    string
}
