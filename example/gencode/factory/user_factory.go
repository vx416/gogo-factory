package factory

import (
	"time"

	gofactory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/example/gencode/model"
	"github.com/vx416/gogo-factory/genutil"
)

type UserFactory struct {
	*gofactory.Factory
}

var User = &UserFactory{gofactory.New(
	&model.User{},
	attr.Int("ID", genutil.SeqInt(1, 1)),
	attr.Str("Name", genutil.RandName(3)),
	attr.Int("Gender", genutil.SeqInt(1, 1)),
	attr.Str("Phone", genutil.RandName(3)),
	attr.Str("Address", genutil.RandName(3)),
	attr.Time("CreatedAt", genutil.Now(time.UTC)),
	attr.Time("UpdatedAt", genutil.Now(time.UTC)),
	attr.Bytes("Password", genutil.FixBytes([]byte("test"))),
)}
