package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/davecgh/go-spew/spew"
	factory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/genutil"
)

func main() {
	fmt.Println("random attributes")
	randomAttributes()
	fmt.Println("----------------------------")
	fmt.Println("customize attributes")
	customizeAttributes()
}

func randomAttributes() {
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Int("ID", genutil.SeqInt(1, 1), "id"),
		attr.Str("Name", genutil.RandName(3)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
		attr.Str("Phone", randomdata.PhoneNumber),
		attr.Str("Address", randomdata.Address),
		attr.Time("CreatedAt", genutil.Now(time.UTC)),
		attr.Time("UpdatedAt", genutil.RandTime(time.Now(), time.Now().Add(30*time.Hour))),
	)

	for i := 0; i < 5; i++ {
		user := userFactory.MustBuild().(*User)
		spew.Printf("user_%d: %+v\n", user.ID, user)
	}
}

func customizeAttributes() {
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Int("ID", genutil.SeqInt(1, 1), "id"),
		attr.Str("Name", genutil.RandName(3)).Process(func(a attr.Attributer) error {
			user := a.GetObject().(*User)
			name := a.GetVal().(string)
			name = "username-" + strconv.Itoa(int(user.ID))
			return a.SetVal(name)
		}),
		attr.Int("Gender", genutil.RandInt(1, 2)),
		attr.Attr("Phone", func() interface{} { return Phone(randomdata.PhoneNumber()) }),
		attr.Str("Address", randomdata.Address),
		attr.Time("CreatedAt", genutil.Now(time.UTC)),
		attr.Time("UpdatedAt", genutil.RandTime(time.Now(), time.Now().Add(30*time.Hour))),
	)

	for i := 0; i < 5; i++ {
		user := userFactory.MustBuild().(*User)
		spew.Printf("user_%d: %+v\n", user.ID, user)
	}
}
