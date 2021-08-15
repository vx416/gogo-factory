package model

import (
	"database/sql"
	"time"
)

type Gender int8

const (
	Male Gender = iota + 1
	Female
)

type Phone string
type Hash []byte

type User struct {
	ID        int64
	Name      string
	Gender    Gender
	Phone     Phone
	Address   sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Password  Hash
}
