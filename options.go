package factory

import (
	"strings"
	"database/sql"
)

var options = &Options{}

type Options struct {
	db     *sql.DB
	driver string
}

func DB(db *sql.DB, driver string) {
	options.db = db
	options.driver = strings.ToLower(driver)
}
