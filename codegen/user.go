package codegen

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Gender int8

const (
	Male Gender = iota + 1
	Female
)

type Phone string

type User struct {
	ID        int64
	Name      string
	Gender    Gender
	Phone     Phone
	Address   sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Price     decimal.Decimal
	Amount    decimal.NullDecimal
	Timestamp
}

type Timestamp struct {
	CreatedAti uint64
}
