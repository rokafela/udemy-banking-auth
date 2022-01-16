package app

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/logger"
	"github.com/rokafela/udemy-banking-auth/service"
	"github.com/rokafela/udemy-banking-auth/validation"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func (ah AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var login_request dto.LoginRequest
	var login_response = dto.LoginResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	}

	// ensure content type is json
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		logger.Info("response error|unsupported content type")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// ensure param in json format
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&login_request)
	if err != nil {
		login_response.Code = http.StatusUnprocessableEntity
		login_response.Message = "invalid parameters"
		login_response.Data = map[string]interface{}{
			"errors": err.Error(),
		}
		encoded_response, _ := json.Marshal(login_response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, http.StatusUnprocessableEntity, login_response)
		return
	}

	// log the received param
	received_param, _ := json.Marshal(login_request)
	logger.Info("received param|" + string(received_param))

	// validate param with validation lib
	validation_error := validation.ValidateStruct(login_request)
	if validation_error != nil {
		login_response.Code = http.StatusUnprocessableEntity
		login_response.Message = "invalid parameters"
		login_response.Data = map[string]interface{}{
			"errors": validation_error,
		}
		encoded_response, _ := json.Marshal(login_response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, http.StatusUnprocessableEntity, login_response)
		return
	}

	// param validated, send to service
	token, app_error := ah.AuthService.Login(&login_request)
	if app_error != nil {
		login_response.Code = app_error.Code
		login_response.Message = app_error.Message
		encoded_response, _ := json.Marshal(login_response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, app_error.Code, login_response)
		return
	}

	login_response.Data = map[string]*string{
		"token": token,
	}
	encoded_response, _ := json.Marshal(login_response)
	logger.Info("response success|" + string(encoded_response))
	writeResponse(w, http.StatusOK, login_response)
}

func (ah AuthHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	var verify_request = dto.VerifyRequest{}
	var verify_response = dto.VerifyResponse{
		Code:       0,
		Message:    "success",
		Authorized: true,
	}

	// ensure authorization token exist
	header_auth := r.Header.Get("authorization")
	logger.Info("received param|" + string(header_auth))
	if header_auth == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmp := strings.Split(header_auth, " ")
	if len(tmp) != 2 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	verify_request.Token = tmp[1]
	app_error := ah.AuthService.Verify(&verify_request)
	if app_error != nil {
		verify_response.Code = app_error.Code
		verify_response.Message = app_error.Message
		verify_response.Authorized = false
		encoded_response, _ := json.Marshal(verify_response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, app_error.Code, verify_response)
		return
	}

	encoded_response, _ := json.Marshal(verify_response)
	logger.Info("response success|" + string(encoded_response))
	writeResponse(w, http.StatusOK, verify_response)
}
