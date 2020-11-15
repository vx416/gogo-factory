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
		attr.Seq("ID", 1),
		attr.Str("Name", genutil.NameRander(3)),
		attr.Int("Gender", genutil.IntRander(1, 2)),
		attr.Str("Phone", randomdata.PhoneNumber),
		attr.Str("Address", randomdata.Address),
		attr.Time("CreatedAt", genutil.NowRander()),
		attr.Time("UpdatedAt", genutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour))),
	)

	for i := 0; i < 5; i++ {
		user := userFactory.MustBuild().(*User)
		spew.Printf("user_%d: %+v\n", user.ID, user)
	}
}

func customizeAttributes() {
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1),
		attr.Str("Name", genutil.NameRander(3)).Process(func(a attr.Attributer) error {
			user := a.GetObject().(*User)
			name := a.GetVal().(string)
			name = "username-" + strconv.Itoa(int(user.ID))
			return a.SetVal(name)
		}),
		attr.Int("Gender", genutil.IntRander(1, 2)),
		attr.Attr("Phone", func() interface{} { return Phone(randomdata.PhoneNumber()) }),
		attr.Str("Address", randomdata.Address),
		attr.Time("CreatedAt", genutil.NowRander()),
		attr.Time("UpdatedAt", genutil.TimeRander(time.Now(), time.Now().Add(30*time.Hour))),
	)

	for i := 0; i < 5; i++ {
		user := userFactory.MustBuild().(*User)
		spew.Printf("user_%d: %+v\n", user.ID, user)
	}
}
