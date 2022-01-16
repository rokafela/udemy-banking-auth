package service

import (
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/errs"
)

//go:generate mockgen -destination=../mocks/service/mockRegistrationService.go -package=service github.com/rokafela/udemy-banking-auth/service RegistrationService
type RegistrationService interface {
	RegisterCustomer(*dto.RegisterCustomerRequest) *errs.AppError
}

type DefaultRegistrationService struct {
	RepoCustomer domain.CustomerRepository
}

func NewRegistrationService(repo domain.CustomerRepository) DefaultRegistrationService {
	return DefaultRegistrationService{
		RepoCustomer: repo,
	}
}
