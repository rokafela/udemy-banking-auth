package domain

import "github.com/rokafela/udemy-banking-auth/errs"

//go:generate mockgen -destination=../mocks/domain/mockRegistrationRepository.go -package=domain github.com/rokafela/udemy-banking-auth/domain RegistrationRepository
type CustomerRepository interface {
	SaveCustomer(*Customer) *errs.AppError
}
