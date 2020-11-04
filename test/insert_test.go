package test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"

	factory "github.com/vicxu416/gogo-factory"
	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/randutil"

	"github.com/stretchr/testify/suite"
)

func TestSqlite(t *testing.T) {
	suite.Run(t, new(sqliteSuite))
}

type sqliteSuite struct {
	db *sql.DB
	suite.Suite
}

func (suite *sqliteSuite) SetupSuite() {
	db, err := initSqliteDB()
	suite.Require().NoError(err)
	suite.db = db
	factory.Opt().SetDB(db, "sqlite3")
}

func (suite *sqliteSuite) TestInsertWithFullAssociation() {
	locationFactory := factory.New(
		&Location{},
		attr.Seq("ID", 1, "id"),
		attr.Str("Address", randomdata.Address, "address"),
	).Table("locations")

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1, "id"),
	).Columns(factory.Col("HostID", "host_id"), factory.Col("Location", "ID", "location_id")).
		FAssociate("Location", locationFactory, 1, true, nil, "location_id").Table("homes")

	userFactory := factory.New(
		&User{CreatedAt: time.Now()},
		attr.Seq("ID", 1, "id"),
		attr.Str("Username", randutil.NameRander(3), "username"),
		attr.Int("Age", randutil.IntRander(25, 50), "age"),
	).Columns(factory.Col("CreatedAt", "created_at")).
		FAssociate("Home", homeFactory, 1, false, func(data, depend interface{}) error {
			user := data.(*User)
			home := depend.(*Home)
			home.HostID = user.ID
			return nil
		}).
		FAssociate("Rented", homeFactory, 5, false, func(data, depend interface{}) error {
			user := data.(*User)
			home := depend.(*Home)
			home.HostID = user.ID
			return nil
		}).Table("users")

	_, err := userFactory.Attrs(
		attr.Time("UpdatedAt", randutil.NowRander, "updated_at"),
	).Insert()
	suite.Require().NoError(err)
	_, err = userFactory.Insert()
	suite.Require().NoError(err)
}
