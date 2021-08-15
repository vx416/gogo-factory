package model

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
	"gopkg.in/guregu/null.v4"
)

type Product struct {
	UID       string
	Buyer     *User
	Price     decimal.Decimal
	Quantity  decimal.Decimal
	Discount  null.String
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
