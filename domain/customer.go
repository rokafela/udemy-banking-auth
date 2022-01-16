package domain

type Customer struct {
	CustomerId  int64  `db:"customer_id"`
	Name        string `db:"name"`
	DateOfBirth string `db:"date_of_birth"`
	City        string `db:"city"`
	Zipcode     string `db:"zipcode"`
	Status      int8   `db:"status"`
}
