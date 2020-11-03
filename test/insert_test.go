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
	factory.DB(db, "sqlite3")
}

// func (suite *sqliteSuite) TestInsert() {
// 	userFactory := factory.New(
// 		func() interface{} { return &User{CreatedAt: time.Now()} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Username", randutil.NameRander(3), "username"),
// 		attr.Int("Age", randutil.IntRander(25, 50), "age"),
// 	).Fix("CreatedAt", "created_at").Table("users")

// 	_, err := userFactory.Insert()
// 	suite.Require().NoError(err)

// 	for i := 1; i <= 10; i++ {
// 		_, err := userFactory.Insert()
// 		suite.Require().NoError(err)
// 	}
// }

// func (suite *sqliteSuite) TestInsertWithAfterAssociation() {
// 	homeFactory := factory.New(
// 		func() interface{} { return &Home{} },
// 		attr.Seq("ID", 1, "id"),
// 	).Fix("HostID", "host_id").Table("homes")

// 	userFactory := factory.New(
// 		func() interface{} { return &User{CreatedAt: time.Now()} },
// 		attr.Seq("ID", 50, "id"),
// 		attr.Str("Username", randutil.NameRander(3), "username"),
// 		attr.Int("Age", randutil.IntRander(25, 50), "age"),
// 	).FAssociate("Home", homeFactory, 1, false, func(data, depend interface{}) error {
// 		user := data.(*User)
// 		home := depend.(*Home)
// 		home.HostID = user.ID
// 		return nil
// 	}).Table("users")

// 	_, err := userFactory.Insert()
// 	suite.Require().NoError(err)
// 	_, err = homeFactory.Insert()
// 	suite.Require().NoError(err)
// }

// func (suite *sqliteSuite) TestInsertWithBeforeAssociation() {
// 	locationFactory := factory.New(
// 		func() interface{} { return &Location{} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Address", randomdata.Address, "address"),
// 	).Table("locations")

// 	homeFactory := factory.New(
// 		func() interface{} { return &Home{} },
// 		attr.Seq("ID", 1, "id"),
// 	).FAssociate("Location", locationFactory, 1, true, nil, "location_id").Table("homes")

// 	homeData, err := homeFactory.Insert()
// 	suite.Require().NoError(err)
// 	home := homeData.(*Home)
// 	suite.Assert().NotNil(home)
// 	suite.Assert().NotNil(home.Location)
// }

func (suite *sqliteSuite) TestInsertWithFullAssociation() {
	locationFactory := factory.New(
		&Location{},
		attr.Seq("ID", 1, "id"),
		attr.Str("Address", randomdata.Address, "address"),
	).Table("locations")

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1, "id"),
	).Fix("HostID", "host_id").FAssociate("Location", locationFactory, 1, true, nil, "location_id").Table("homes")

	userFactory := factory.New(
		&User{CreatedAt: time.Now()},
		attr.Seq("ID", 1, "id"),
		attr.Str("Username", randutil.NameRander(3), "username"),
		attr.Int("Age", randutil.IntRander(25, 50), "age"),
	).FAssociate("Home", homeFactory, 1, false, func(data, depend interface{}) error {
		user := data.(*User)
		home := depend.(*Home)
		home.HostID = user.ID
		return nil
	}).Table("users")

	_, err := userFactory.Insert()
	suite.Require().NoError(err)
}
