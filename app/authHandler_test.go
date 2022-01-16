package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/mocks/service"
)

var mockAuthHandler AuthHandler
var mockRouter *mux.Router
var mockAuthService *service.MockAuthService

func setupMock(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockAuthService = service.NewMockAuthService(ctrl)
	mockAuthHandler = AuthHandler{mockAuthService}
	mockRouter = mux.NewRouter()
	mockRouter.HandleFunc("/login", mockAuthHandler.HandleLogin)

	return func() {
		mockRouter = nil
		defer ctrl.Finish()
	}
}

func Test_login_success(t *testing.T) {
	// ----- Arrange -----
	// dummy_body adalah data yang akan dikirim melalui http request
	dummy_body := []byte(`{"username":"admin","password":"abc123"}`)

	// dummy_login_request adalah data yang harus diterima oleh service dari handler
	dummy_login_request := dto.LoginRequest{
		Username: "admin",
		Password: "abc123",
	}

	// dummy token adalah data yang akan dikembalikan oleh service ke handler
	dummy_token := "token"

	// dummy_login_response adalah format data yang harus dikirim oleh handler sebagai http response
	dummy_login_response := dto.LoginResponse{}

	// setup mock handler, router, and service
	teardown := setupMock(t)
	defer teardown()

	// menentukan data yang harus diterima oleh service dan data yang akan dikembalikan oleh service
	mockAuthService.EXPECT().Login(&dummy_login_request).Return(&dummy_token, nil)

	// menentukan http request dan header
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(dummy_body))
	request.Header.Set("Content-Type", "application/json")

	// ----- Act -----
	recorder := httptest.NewRecorder()
	mockRouter.ServeHTTP(recorder, request)

	// ----- Assert -----
	// cek http status code 200
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

	// cek response body adalah json dan dapat di-decode
	err := json.NewDecoder(recorder.Body).Decode(&dummy_login_response)
	if err != nil {
		t.Error("Failed while decoding json response")
	}

	// cek status code di dalam response body
	if dummy_login_response.Code != 0 {
		t.Error("Failed while testing the response code")
	}

	// cek token di dalam response body
	if dummy_login_response.Data.(map[string]interface{})["token"] == "" {
		t.Error("Failed while testing the response data token")
	}
}

func Test_login_tanpa_header_content_type(t *testing.T) {
	// ----- Arrange -----
	// dummy_body adalah data yang akan dikirim melalui http request
	dummy_body := []byte(`{"username":"admin","password":"abc123"}`)

	// setup mock handler, router, and service
	teardown := setupMock(t)
	defer teardown()

	// menentukan http request dan header
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(dummy_body))

	// ----- Act -----
	recorder := httptest.NewRecorder()
	mockRouter.ServeHTTP(recorder, request)

	// ----- Assert -----
	// cek status code 415
	if recorder.Code != http.StatusUnsupportedMediaType {
		t.Error("Failed while testing the status code")
	}
}

func Test_login_tanpa_body_json(t *testing.T) {
	// ----- Arrange -----
	// dummy_body adalah data yang akan dikirim melalui http request
	dummy_body := []byte(``)

	// dummy_login_response adalah format data yang harus dikirim oleh handler sebagai http response
	dummy_login_response := dto.LoginResponse{}

	// setup mock handler, router, and service
	teardown := setupMock(t)
	defer teardown()

	// menentukan http request dan header
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(dummy_body))
	request.Header.Set("Content-Type", "application/json")

	// ----- Act -----
	recorder := httptest.NewRecorder()
	mockRouter.ServeHTTP(recorder, request)

	// ----- Assert -----
	// cek http status code 422
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Error("Failed while testing the status code")
	}

	// cek response body adalah json dan dapat di-decode
	err := json.NewDecoder(recorder.Body).Decode(&dummy_login_response)
	if err != nil {
		t.Error("Failed while decoding json response")
	}

	// cek status code di dalam response body
	if dummy_login_response.Code != http.StatusUnprocessableEntity {
		t.Error("Failed while testing the response code")
	}

	// cek errors di dalam response body
	if dummy_login_response.Data.(map[string]interface{})["errors"] == "" || dummy_login_response.Data.(map[string]interface{})["errors"] == nil {
		t.Error("Failed while testing the response data error")
	}
}

func Test_login_gagal_validate_param(t *testing.T) {
	// ----- Arrange -----
	// dummy_body adalah data yang akan dikirim melalui http request
	dummy_body := []byte(`{"username":"admin","passwords":"abc123"}`)

	// dummy_login_response adalah format data yang harus dikirim oleh handler sebagai http response
	dummy_login_response := dto.LoginResponse{}

	// setup mock handler, router, and service
	teardown := setupMock(t)
	defer teardown()

	// menentukan http request dan header
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(dummy_body))
	request.Header.Set("Content-Type", "application/json")

	// ----- Act -----
	recorder := httptest.NewRecorder()
	mockRouter.ServeHTTP(recorder, request)

	// ----- Assert -----
	// cek http status code 422
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Error("Failed while testing the status code")
	}

	// cek response body adalah json dan dapat di-decode
	err := json.NewDecoder(recorder.Body).Decode(&dummy_login_response)
	if err != nil {
		t.Error("Failed while decoding json response")
	}

	// cek status code di dalam response body
	if dummy_login_response.Code != http.StatusUnprocessableEntity {
		t.Error("Failed while testing the response code")
	}

	// cek errors di dalam response body
	fmt.Println(dummy_login_response.Data.(map[string]interface{})["errors"])
	if dummy_login_response.Data.(map[string]interface{})["errors"] == "" || dummy_login_response.Data.(map[string]interface{})["errors"] == nil {
		t.Error("Failed while testing the response data error")
	}
}

func Test_login_gagal(t *testing.T) {
	// ----- Arrange -----
	// dummy_body adalah data yang akan dikirim melalui http request
	dummy_body := []byte(`{"username":"admin","password":"abc123"}`)

	// dummy_login_request adalah data yang harus diterima oleh service dari handler
	dummy_login_request := dto.LoginRequest{
		Username: "admin",
		Password: "abc123",
	}

	// dummy token adalah data yang akan dikembalikan oleh service ke handler
	dummy_error := errs.NewAuthenticationError("Invalid token")

	// dummy_login_response adalah format data yang harus dikirim oleh handler sebagai http response
	dummy_login_response := dto.LoginResponse{}

	// setup mock handler, router, and service
	teardown := setupMock(t)
	defer teardown()

	// menentukan data yang harus diterima oleh service dan data yang akan dikembalikan oleh service
	mockAuthService.EXPECT().Login(&dummy_login_request).Return(nil, dummy_error)

	// menentukan http request dan header
	request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(dummy_body))
	request.Header.Set("Content-Type", "application/json")

	// ----- Act -----
	recorder := httptest.NewRecorder()
	mockRouter.ServeHTTP(recorder, request)

	// ----- Assert -----
	// cek http status code 403
	if recorder.Code != http.StatusForbidden {
		t.Error("Failed while testing the status code")
	}

	// cek response body adalah json dan dapat di-decode
	err := json.NewDecoder(recorder.Body).Decode(&dummy_login_response)
	if err != nil {
		t.Error("Failed while decoding json response")
	}

	// cek status code di dalam response body
	if dummy_login_response.Code != http.StatusForbidden {
		t.Error("Failed while testing the response code")
	}
}
