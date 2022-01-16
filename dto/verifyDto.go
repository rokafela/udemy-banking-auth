package dto

type VerifyRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Authorized bool        `json:"authorized"`
}
