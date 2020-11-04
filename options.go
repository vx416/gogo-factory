package factory

import (
	"database/sql"
	"strings"

	"github.com/vicxu416/gogo-factory/dbutil"
)

var options = &Options{}

type Options struct {
	DB         *sql.DB
	Driver     string
	TagName    string
	InsertFunc dbutil.InsertFunc
}

func (opt *Options) SetDB(db *sql.DB, driver string) *Options {
	opt.DB = db
	opt.Driver = strings.ToLower(driver)
	return opt
}

func (opt *Options) SetInsertFunc(fn dbutil.InsertFunc) *Options {
	opt.InsertFunc = fn
	return opt
}

func Opt() *Options {
	return options
}
