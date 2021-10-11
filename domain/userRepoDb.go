package domain

import (
	"github.com/jmoiron/sqlx"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/logger"
)

type UserRepositoryDb struct {
	Client *sqlx.DB
}

func NewUserRepositoryDb(client *sqlx.DB) UserRepositoryDb {
	return UserRepositoryDb{
		Client: client,
	}
}

func (db UserRepositoryDb) FindUserById(searched *User) ([]User, *errs.AppError) {
	var result []User
	selectQuery := "SELECT * FROM users WHERE status = 1 AND username = ? AND password = ?"
	err := db.Client.Select(&result, selectQuery, searched.Username, searched.Password)
	if err != nil {
		logger.Error(err.Error())
		return nil, errs.NewNotFoundError("User not found")
	}
	return result, nil
}
