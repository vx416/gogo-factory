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
	Countries []*Country
}

type Country struct {
	ID     int64
	HostID int64
	Homes  []*Home
}

type Home struct {
	ID        int64
	HostID    int64
	CountryID int64
	CreatedAt time.Time
	Location  *Location
}

type Location struct {
	ID      int64
	Address string
}

type Employee struct {
	ID                int64        `db:"id" gorm:"column:id"`
	Name              string       `db:"name" gorm:"column:name"`
	Gender            Gender       `db:"gender" gorm:"column:gender"`
	Age               *int32       `db:"age" gorm:"column:age"`
	Phone             string       `db:"phone" gorm:"column:phone"`
	Salary            float64      `db:"salary" gorm:"column:salary"`
	Specialty         *Specialty   `gorm:"-"`
	SecondSpecialties []*Specialty `gorm:"-"`
	Projects          []*Project   `gorm:"-"`
	CreatedAt         time.Time    `db:"created_at" gorm:"column:created_at"`
	UpdatedAt         sql.NullTime `db:"updated_at" gorm:"column:updated_at"`
}

type Project struct {
	ID        int64       `db:"id" gorm:"column:id"`
	Name      string      `db:"name" gorm:"column:name"`
	Employees []*Employee `gorm:"-"`
	Tasks     []*Task     `gorm:"-"`
	Deadline  time.Time   `db:"deadline" gorm:"column:deadline"`
}

type EmployeesProjects struct {
	ID         int64 `db:"id" gorm:"column:id"`
	ProjectID  int64 `db:"project_id" gorm:"column:project_id"`
	EmployeeID int64 `db:"employee_id" gorm:"column:employee_id"`
}

type Task struct {
	ID        int64     `db:"id" gorm:"column:id"`
	Name      string    `db:"name" gorm:"column:name"`
	ProjectID int64     `db:"project_id" gorm:"column:project_id"`
	Deadline  time.Time `db:"deadline" gorm:"column:deadline"`
}

type Specialty struct {
	ID       int64         `db:"id" gorm:"column:id"`
	Name     string        `db:"name" gorm:"column:name"`
	OwnerID  sql.NullInt64 `db:"owner_id" gorm:"column:owner_id"`
	Domain   *Domain
	DomainID sql.NullInt64 `db:"domain_id" gorm:"column:domain_id"`
}

type Domain struct {
	ID   int64  `db:"id" gorm:"column:id"`
	Name string `db:"name" gorm:"column:name"`
}
