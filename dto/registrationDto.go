package dto

import "github.com/rokafela/udemy-banking-auth/domain"

type RegisterCustomerRequest struct {
	Name        string `json:"name" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
	City        string `json:"city" validate:"required"`
	Zipcode     string `json:"zipcode" validate:"required"`
}

type RegisterCustomerResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (rcr RegisterCustomerRequest) TransformToCustomer() domain.Customer {
	return domain.Customer{
		Name:        rcr.Name,
		DateOfBirth: rcr.DateOfBirth,
		City:        rcr.City,
		Zipcode:     rcr.Zipcode,
	}
}
