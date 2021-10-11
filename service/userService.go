package service

import "github.com/rokafela/udemy-banking-auth/domain"

type DefaultUserService struct {
	UserRepo domain.UserRepositoryDb
}

func NewUserService(userRepo domain.UserRepositoryDb) DefaultUserService {
	return DefaultUserService{
		UserRepo: userRepo,
	}
}

type UserService interface {
}
