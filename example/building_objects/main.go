package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	factory "github.com/vicxu416/gogo-factory"
	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/randutil"
)

func main() {
	db, err := initSqliteDB()
	if err != nil {
		log.Panicf("err:%+v", err)
	}
	factory.DB(db, "sqlite3")
	initSqliteDB()
	// insertObject()
	omitFieldObject()
}

func insertObject() {
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1, "id"),
		attr.Str("Name", randutil.NameRander(3), "name"),
		attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
		attr.Str("Phone", randomdata.PhoneNumber, "phone"),
		attr.Str("Address", randomdata.Address, "address"),
		attr.Time("CreatedAt", randutil.NowRander(), "created_at"),
		attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
	).Table("users")

	for i := 0; i < 5; i++ {
		user := userFactory.MustInsert().(*User)
		spew.Printf("user_%d: %+v\n", user.ID, user)
	}
}

func omitFieldObject() {
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1, "id"),
		attr.Str("Name", randutil.NameRander(3), "name"),
		attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
		attr.Str("Phone", randomdata.PhoneNumber, "phone"),
		attr.Str("Address", randomdata.Address, "address"),
		attr.Time("CreatedAt", randutil.NowRander(), "created_at"),
		attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
	).Table("users")

	for i := 0; i < 5; i++ {
		user := userFactory.Omit("Address").MustInsert().(*User)
		fmt.Printf("%+v\n", user)
	}
	user := userFactory.MustInsert().(*User)
	fmt.Printf("%+v\n", user)
}

// func fixFieldObject() {
// 	userFactory := factory.New(
// 		func() interface{} { return &User{CreatedAt: time.Now()} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Name", randutil.NameRander(3), "name"),
// 		attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
// 		attr.Str("Phone", randomdata.PhoneNumber, "phone"),
// 		attr.Str("Address", randomdata.Address, "address"),
// 		attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
// 	).Table("users").Fix("CreatedAt", "created_at")
// }

// func tempFieldObject() {
// 	userFactory := factory.New(
// 		func() interface{} { return &User{CreatedAt: time.Now()} },
// 		attr.Seq("ID", 1, "id"),
// 		attr.Str("Name", randutil.NameRander(3), "name"),
// 		attr.Int("Gender", randutil.IntRander(1, 2), "gender"),
// 		attr.Str("Phone", randomdata.PhoneNumber, "phone"),
// 		attr.Str("Address", randomdata.Address, "address"),
// 		attr.Time("UpdatedAt", randutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour)), "updated_at"),
// 	).Table("users").Fix("CreatedAt", "created_at")
// }

func initSqliteDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	schema, err := readSchema("schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func readSchema(fileName string) (string, error) {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	bytes, err := ioutil.ReadFile(filepath.Join(dir, "./"+fileName))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
