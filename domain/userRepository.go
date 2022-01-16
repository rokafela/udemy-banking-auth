package domain

import "github.com/rokafela/udemy-banking-auth/errs"

//go:generate mockgen -destination=../mocks/domain/mockUserRepository.go -package=domain github.com/rokafela/udemy-banking-auth/domain UserRepository
type UserRepository interface {
	SaveUser(*User) *errs.AppError
}
