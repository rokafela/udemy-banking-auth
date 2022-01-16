package domain

type User struct {
	Username   string `db:"username"`
	Password   string `db:"password"`
	Role       string `db:"role"`
	CustomerId int64  `db:"customer_id"`
}
