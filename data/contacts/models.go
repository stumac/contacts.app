// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package data

import (
	"database/sql"
)

type Contact struct {
	ID        int64
	FirstName string
	Phone     sql.NullString
	Email     sql.NullString
}
