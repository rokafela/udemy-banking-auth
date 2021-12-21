package dto

type VerifyRequest struct {
	Token  string `json:token`
	Action string `json:action`
	Param  string `json:param`
}

type VerifyResponse struct {
	Authorized bool        `json:"authorized"`
	Message    interface{} `json:"message,omitempty"`
}
