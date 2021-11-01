package app

import (
	"encoding/json"
	"net/http"

	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/logger"
	"github.com/rokafela/udemy-banking-auth/service"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func (uh AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, app_error := uh.AuthService.Login(&request)
		if app_error != nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}
