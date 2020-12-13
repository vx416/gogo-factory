package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gofactory "github.com/vx416/gogo-factory"
)

func TestGormTagInsert(t *testing.T) {
	gofactory.Opt().SetTagProcess(gofactory.GormTagProcess)
	db, err := initSqliteDB()
	if err != nil {
		t.Fatalf("db init failed, err:%+v", err)
	}
	gofactory.Opt().SetDB(db, "sqlite3")
	specFactory := SpecialtyFactory.BelongsToDomain(DomainFactory)
	employees := EmployeeFactory.HasOneSpecialty(specFactory).
		HasManySecondSpecialties(specFactory, 5).MustInsertN(5).([]*Employee)
	employees2, _, err := AllEmployees(db, "sqlite3")
	assert.Len(t, employees2, len(employees))
	gofactory.Opt().SetTagProcess(nil)
}

func TestDBTagInsert(t *testing.T) {
	gofactory.Opt().SetTagProcess(gofactory.GormTagProcess)
	db, err := initSqliteDB()
	if err != nil {
		t.Fatalf("db init failed, err:%+v", err)
	}
	gofactory.Opt().SetDB(db, "sqlite3")
	specFactory := SpecialtyFactory.BelongsToDomain(DomainFactory)
	employees := EmployeeFactory.HasOneSpecialty(specFactory).
		HasManySecondSpecialties(specFactory, 5).MustInsertN(5).([]*Employee)
	employees2, _, err := AllEmployees(db, "sqlite3")
	assert.Len(t, employees2, len(employees))
	gofactory.Opt().SetTagProcess(nil)
}
