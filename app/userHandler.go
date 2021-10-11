package app

import (
	"net/http"

	"github.com/rokafela/udemy-banking-auth/service"
)

type UserHandler struct {
	UserService service.UserService
}

func (uh UserHandler) VerifyUser(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, nil)
}
