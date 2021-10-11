package domain

import (
	"github.com/rokafela/udemy-banking-auth/errs"
)

type User struct {
	Username   string
	Password   string
	Role       string
	CustomerId string
	Status     string
}

type UserRepository interface {
	FindUserById(*User) ([]User, *errs.AppError)
}
