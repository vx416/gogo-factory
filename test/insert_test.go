package test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/suite"
	factory "github.com/vicxu416/gogo-factory"
)

func TestSqlite(t *testing.T) {
	db, err := initSqliteDB()
	if err != nil {
		t.Fatalf("db init failed, err:%+v", err)
	}
	s := &insertSuite{
		db:     db,
		dbType: "sqlite3",
	}
	suite.Run(t, s)
}

type insertSuite struct {
	db     *sql.DB
	dbType string
	suite.Suite
}

func (suite *insertSuite) SetupSuite() {
	factory.Opt().SetDB(suite.db, suite.dbType)
}

func (suite *insertSuite) AfterTest() {

}
