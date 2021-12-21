package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/logger"
	"github.com/rokafela/udemy-banking-auth/service"
)

type AuthHandler struct {
	AuthService service.AuthService
	Validate    *validator.Validate
}

func (ah AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var login_request dto.LoginRequest
	// ensure param in json format
	err := json.NewDecoder(r.Body).Decode(&login_request)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// log the received param
		received_param, err := json.Marshal(login_request)
		if err != nil {
			logger.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Info(string(received_param))

		// validate param indexes
		validation_error := ah.Validate.Struct(login_request)
		if validation_error != nil {
			if _, ok := validation_error.(*validator.InvalidValidationError); ok {
				fmt.Println(validation_error)
				return
			}

			for _, err := range validation_error.(validator.ValidationErrors) {
				// fmt.Println(err.Namespace())
				// fmt.Println(err.Field())
				// fmt.Println(err.StructNamespace())
				// fmt.Println(err.StructField())
				// fmt.Println(err.Tag())
				// fmt.Println(err.ActualTag())
				// fmt.Println(err.Kind())
				// fmt.Println(err.Type())
				// fmt.Println(err.Value())
				// fmt.Println(err.Param())
				fmt.Println(err.Translate(translator))
			}

			// log.Println("validation_error2")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		app_error := login_request.Validate()
		if app_error != nil {
			writeResponse(w, app_error.Code, app_error.Message)
		}

		token, app_error := ah.AuthService.Login(&login_request)
		if app_error != nil {
			writeResponse(w, app_error.Code, app_error.Message)
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}

func (ah AuthHandler) HandleVerify(w http.ResponseWriter, r *http.Request) {
	var verify_request dto.VerifyRequest
	err := json.NewDecoder(r.Body).Decode(&verify_request)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if verify_request.Token != "" {
			appErr := ah.AuthService.Verify(&verify_request)
			if appErr != nil {
				writeResponse(w, appErr.Code, notAuthorizedResponse(appErr.Message))
			} else {
				writeResponse(w, http.StatusOK, authorizedResponse())
			}
		} else {
			writeResponse(w, http.StatusForbidden, notAuthorizedResponse("Missing token"))
		}
	}
}

func notAuthorizedResponse(msg string) dto.VerifyResponse {
	return dto.VerifyResponse{
		Authorized: false,
		Message:    msg,
	}
}

func authorizedResponse() dto.VerifyResponse {
	return dto.VerifyResponse{
		Authorized: true,
	}
}
