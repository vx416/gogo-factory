package factory

import (
	"time"

	gofactory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/example/gencode/model"
	"github.com/vx416/gogo-factory/genutil"
)

type ProductFactory struct {
	*gofactory.Factory
}

var Product = &ProductFactory{gofactory.New(
	&model.Product{},
	attr.Str("UID", genutil.RandName(3)),
	attr.Attr("Buyer", genutil.FixInterface(nil)),
	attr.Str("Price", genutil.RandName(3)),
	attr.Str("Quantity", genutil.RandName(3)),
	attr.Str("Discount", genutil.RandName(3)),
	attr.Time("CreatedAt", genutil.Now(time.UTC)),
	attr.Time("DeletedAt", genutil.Now(time.UTC)),
)}
