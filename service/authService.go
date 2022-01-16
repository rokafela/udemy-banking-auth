package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/logger"
)

type DefaultAuthService struct {
	UserRepo domain.LoginRepository
}

func NewAuthService(userRepo domain.LoginRepository) DefaultAuthService {
	return DefaultAuthService{
		UserRepo: userRepo,
	}
}

//go:generate mockgen -destination=../mocks/service/mockAuthService.go -package=service github.com/rokafela/udemy-banking-auth/service AuthService
type AuthService interface {
	Login(*dto.LoginRequest) (*string, *errs.AppError)
	Verify(*dto.VerifyRequest) *errs.AppError
}

func (dus DefaultAuthService) Login(login_request *dto.LoginRequest) (*string, *errs.AppError) {
	// validate credetential
	found_user, app_error := dus.UserRepo.FindUserByCredential(&login_request.Username, &login_request.Password)
	if app_error != nil {
		return nil, app_error
	}

	// user found, determine role
	claims := domain.GetClaimsForRole(found_user)
	// role determined, generate token
	authToken, app_error := domain.GenerateAuthToken(&claims)
	if app_error != nil {
		return nil, app_error
	}

	// token generated, save to db
	app_error = dus.UserRepo.SaveToken(&authToken)
	if app_error != nil {
		return nil, app_error
	}

	return &authToken, nil
}

func (dus DefaultAuthService) Verify(verify_request *dto.VerifyRequest) *errs.AppError {
	if jwtToken, err := jwtTokenFromString(verify_request.Token); err != nil {
		return errs.NewAuthenticationError(err.Error())
	} else {
		if jwtToken.Valid {
			return nil
		} else {
			return errs.NewAuthenticationError("Invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}
