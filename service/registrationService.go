package service

import (
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/dto"
	"github.com/rokafela/udemy-banking-auth/errs"
)

//go:generate mockgen -destination=../mocks/service/mockRegistrationService.go -package=service github.com/rokafela/udemy-banking-auth/service RegistrationService
type RegistrationService interface {
	RegisterCustomer(*dto.RegisterCustomerRequest) (*int64, *errs.AppError)
}

type DefaultRegistrationService struct {
	RepoCustomer domain.CustomerRepository
}

func NewRegistrationService(repo domain.CustomerRepository) DefaultRegistrationService {
	return DefaultRegistrationService{
		RepoCustomer: repo,
	}
}

func (drs DefaultRegistrationService) RegisterCustomer(register_request *dto.RegisterCustomerRequest) (*int64, *errs.AppError) {
	customer_data := register_request.TransformToCustomer()
	customer_id, app_error := drs.RepoCustomer.SaveCustomer(&customer_data)
	if app_error != nil {
		return nil, app_error
	}
	return customer_id, nil
}
