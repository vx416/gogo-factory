package factory

import (
	"database/sql"
	"strings"
)

var options = &Options{}

type Options struct {
	db      *sql.DB
	driver  string
	tagName string
}

func DB(db *sql.DB, driver string) {
	options.db = db
	options.driver = strings.ToLower(driver)
}
