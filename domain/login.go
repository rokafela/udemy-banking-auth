package domain

import (
	"database/sql"
)

type Login struct {
	Username   string         `db:"username"`
	Role       string         `db:"role"`
	CustomerId sql.NullString `db:"customer_id"`
	AccountId  sql.NullString `db:"account_id"`
}
