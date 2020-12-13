package gofactory

import (
	"database/sql"
	"strings"

	"github.com/vx416/gogo-factory/dbutil"
)

var options = &Options{}

// Options global option for factory context
type Options struct {
	DB         *sql.DB
	Driver     string
	InsertFunc dbutil.InsertFunc
	TagProcess TagProcess
}

// SetDB setup db instance
func (opt *Options) SetDB(db *sql.DB, driver string) *Options {
	opt.DB = db
	opt.Driver = strings.ToLower(driver)
	return opt
}

// SetInsertFunc setup global insert function
func (opt *Options) SetInsertFunc(fn dbutil.InsertFunc) *Options {
	opt.InsertFunc = fn
	return opt
}

// SetTagProcess setup tag process
func (opt *Options) SetTagProcess(tp TagProcess) *Options {
	opt.TagProcess = tp
	return opt
}

// Opt get global options
func Opt() *Options {
	return options
}
