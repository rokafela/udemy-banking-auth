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

type RegistrationHandler struct {
	registration_service service.RegistrationService
	auth_service         service.AuthService
}

func (handler RegistrationHandler) HandleCustomerRegistration(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterCustomerRequest
	var response = dto.RegisterCustomerResponse{
		Code:    0,
		Message: "success",
		Data:    nil,
	}

	// ensure authorization token exist
	header_auth := r.Header.Get("authorization")
	logger.Info("received param|" + string(header_auth))
	if header_auth == "" {
		logger.Info("response error|token not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmp := strings.Split(header_auth, " ")
	if len(tmp) != 2 {
		logger.Info("response error|token wrong format")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	app_error := handler.auth_service.Verify(&dto.VerifyRequest{Token: tmp[1]})
	if app_error != nil {
		logger.Info("response error|token unauthorized")
		w.WriteHeader(http.StatusForbidden)
		return
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
	err := decoder.Decode(&request)
	if err != nil {
		response.Code = http.StatusUnprocessableEntity
		response.Message = "invalid parameters"
		response.Data = map[string]interface{}{
			"errors": err.Error(),
		}
		encoded_response, _ := json.Marshal(response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, http.StatusUnprocessableEntity, response)
		return
	}

	// log the received param
	received_param, _ := json.Marshal(request)
	logger.Info("received param|" + string(received_param))

	// validate param with validation lib
	validation_error := validation.ValidateStruct(request)
	if validation_error != nil {
		response.Code = http.StatusUnprocessableEntity
		response.Message = "invalid parameters"
		response.Data = map[string]interface{}{
			"errors": validation_error,
		}
		encoded_response, _ := json.Marshal(response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, http.StatusUnprocessableEntity, response)
		return
	}

	// param validated, send to service
	customer_id, app_error := handler.registration_service.RegisterCustomer(&request)
	if app_error != nil {
		response.Code = app_error.Code
		response.Message = app_error.Message
		encoded_response, _ := json.Marshal(response)
		logger.Info("response error|" + string(encoded_response))
		writeResponse(w, app_error.Code, response)
		return
	}

	response.Data = map[string]*int64{
		"customer_id": customer_id,
	}
	encoded_response, _ := json.Marshal(response)
	logger.Info("response success|" + string(encoded_response))
	writeResponse(w, http.StatusOK, response)
}
