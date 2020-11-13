package test

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Home belongs To User
// User has a Home

type Gender int8

type User struct {
	ID        int64
	Username  string
	Phone     string
	Gender    Gender
	Age       *int32
	Host      bool
	Height    float32
	Weight    float32
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	PtrString *string
	Home      *Home
	Rented    []*Home
}

type Home struct {
	ID        int64
	HostID    int64
	CreatedAt time.Time
	Location  *Location
}

type Location struct {
	ID      int64
	Address string
}
