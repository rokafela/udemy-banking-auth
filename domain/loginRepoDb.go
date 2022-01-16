package domain

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/logger"
)

//go:generate mockgen -destination=../mocks/domain/mockLoginRepository.go -package=domain github.com/rokafela/udemy-banking-auth/domain LoginRepository
type LoginRepository interface {
	FindUserByCredential(*string, *string) (*Login, *errs.AppError)
	SaveToken(*string) *errs.AppError
}

type AuthRepositoryDb struct {
	Client *sqlx.DB
}

func NewAuthRepositoryDb(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{
		Client: client,
	}
}

func (db AuthRepositoryDb) FindUserByCredential(username *string, password *string) (*Login, *errs.AppError) {
	var result Login
	selectQuery := `select u.username, u.role, u.customer_id, group_concat(a.account_id) as account_id
	from users u
	left join accounts a on a.customer_id = u.customer_id
	where u.username = ?
	and u.password = ?
	group by u.username, u.role, u.customer_id;`
	err := db.Client.Get(&result, selectQuery, username, password)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("Invalid credentials")
		} else {
			logger.Error("Error while verifying credentials: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &result, nil
}

func (db AuthRepositoryDb) SaveToken(token *string) *errs.AppError {
	insertQuery := `insert into refresh_token_store (refresh_token, created_on) values (?, ?)`
	tx := db.Client.MustBegin()
	tx.MustExec(insertQuery, token, time.Now())
	err := tx.Commit()
	if err != nil {
		return errs.NewUnexpectedError("Error while saving token")
	}
	return nil
}
