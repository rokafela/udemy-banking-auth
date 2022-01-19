package domain

import (
	"github.com/jmoiron/sqlx"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/logger"
)

type CustomerRepositoryDb struct {
	Client *sqlx.DB
}

func NewCustomerRepositoryDb(client *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{
		Client: client,
	}
}

func (db CustomerRepositoryDb) SaveCustomer(customer_data *Customer) (*int64, *errs.AppError) {
	// var result int64
	query_insert := "INSERT INTO customers (name, date_of_birth, city, zipcode, status) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Client.Exec(query_insert, customer_data.Name, customer_data.DateOfBirth, customer_data.City, customer_data.Zipcode, 1)
	if err != nil {
		logger.Error("Error while creating new customer: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	return &id, nil
}
