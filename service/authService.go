package service

import (
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/errs"
)

type DefaultAuthService struct {
	UserRepo domain.AuthRepositoryDb
}

func NewAuthService(userRepo domain.AuthRepositoryDb) DefaultAuthService {
	return DefaultAuthService{
		UserRepo: userRepo,
	}
}

type AuthService interface {
	Login(*dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

func (dus DefaultAuthService) Login(login_request *dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	// validate param
	err := login_request.Validate()
	if err != nil {
		return nil, err
	}

	// validate credetential
	found_user, err := dus.UserRepo.FindUserByCredential(&login_request.Username, &login_request.Password)
	if err != nil {
		return nil, err
	}

	// user found, determine role
	claims := domain.GetClaimsForRole(found_user)
	// role determined, generate token
	authToken, err := domain.GenerateAuthToken(&claims)
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponse{
		Token: authToken,
	}
	return &response, nil
}
