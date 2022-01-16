package service

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	realdomain "github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/mocks/domain"
)

func Test_find_user_by_credential(t *testing.T) {
	// ----- Arrange -----
	// dummy_login_request adalah data yang harus diterima oleh service dari handler
	dummy_login_request := dto.LoginRequest{
		Username: "admin",
		Password: "abc123",
	}

	// dummy_login
	dummy_login := realdomain.Login{
		Username: "admin",
		Role:     "admin",
	}

	// dummy_auth_token
	dummy_token := "token"

	// setup the service with repo
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := domain.NewMockLoginRepository(ctrl)
	service := NewAuthService(mockRepository)

	// ----- Act -----
	mockRepository.EXPECT().FindUserByCredential(&dummy_login_request.Username, &dummy_login_request.Password).Return(&dummy_login, nil)
	mockRepository.EXPECT().SaveToken(&dummy_token).Return(nil)
	str, err := service.Login(&dummy_login_request)

	// ----- Assert -----
	fmt.Println(*str)
	fmt.Println(err)

}
