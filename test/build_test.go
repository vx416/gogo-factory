package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Pallinder/go-randomdata"
	factory "github.com/vicxu416/seed-factory"
	"github.com/vicxu416/seed-factory/attr"
)

func between(t *testing.T, target, max, min int) {
	assert.LessOrEqual(t, target, max)
	assert.GreaterOrEqual(t, target, min)
}

func TestBasicAttributes(t *testing.T) {
	phoneSet := []string{"091234567", "09765432", "096789234"}
	userFactory := factory.New(
		func() interface{} { return &User{CreatedAt: time.Now(), Host: true} },
		attr.Seq("ID", 1),
		attr.Str("Username", randomdata.LastName),
		attr.StrSeq("Phone", phoneSet),
		attr.Attr("Gender", func() interface{} { return int8(randomdata.Number(1, 2)) }),
		attr.Int("Age", func() int { return randomdata.Number(1, 100) }),
		attr.Float("Weight", func() float64 { return randomdata.Decimal(1, 20, 1) }),
		attr.Float("Height", func() float64 { return randomdata.Decimal(1, 20, 1) }),
	)

	for i := 1; i <= 5; i++ {
		user := userFactory.MustBuild().(*User)
		assert.Equal(t, user.ID, int64(i))
		assert.NotEmpty(t, user.Username)
		assert.NotEmpty(t, user.Phone)
		assert.NotZero(t, user.Gender)
		assert.NotZero(t, user.Age)
		assert.NotZero(t, user.Weight)
		assert.NotZero(t, user.Height)
	}
}

func TestRandAttr(t *testing.T) {
	userNameSet := []string{"divaid", "vic", "shelly", "jason"}
	phoneSet := []string{"090123543", "0954323123", "0924325345"}
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1),
		attr.RandStr("Username", userNameSet),
		attr.RandStr("Phone", phoneSet),
		attr.RandInt("Gender", 1, 2),
		attr.RandInt("Age", 25, 100),
		attr.RandFloat("Height", 55.0, 90.0),
		attr.RandFloat("Weight", 155.0, 190.0),
		attr.RandBool("Host", 0.5),
		attr.RandTime("CreatedAt", minTime, maxTime),
	)

	user := userFactory.MustBuild().(*User)
	assert.Subset(t, userNameSet, []string{user.Username})
	assert.Subset(t, phoneSet, []string{user.Phone})
	between(t, int(user.Gender), 2, 1)
	between(t, int(user.Age), 100, 25)
	assert.InDelta(t, float64(user.Height), 90.0, 55.0)
	assert.InDelta(t, float64(user.Weight), 190.0, 155.0)
	between(t, int(user.CreatedAt.Unix()), int(maxTime.Unix()), int(minTime.Unix()))
}

func TestNullableFields(t *testing.T) {
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	stringSet := []string{"ptr_string_1", "ptr_string_2"}
	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1),
		attr.RandTime("UpdatedAt", minTime, maxTime),
		attr.RandStr("PtrString", stringSet),
	)
	user := userFactory.MustBuild().(*User)
	between(t, int(user.UpdatedAt.Time.Unix()), int(maxTime.Unix()), int(minTime.Unix()))
	assert.NotNil(t, user.PtrString)
	assert.Subset(t, stringSet, []string{*user.PtrString})
}

func TestFactoryAttr(t *testing.T) {
	locFactory := factory.New(
		func() interface{} { return &Location{} },
		attr.Seq("ID", 1),
		attr.Str("Loc", randomdata.Address),
	)

	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1),
		attr.Factory("Location", locFactory, false),
	)
	user := userFactory.MustBuild().(*User)
	assert.NotNil(t, user.Location)
	assert.NotZero(t, user.Location.ID)
	assert.NotEmpty(t, user.Location.Loc)
}

func TestProcess(t *testing.T) {
	locFactory := factory.New(
		func() interface{} { return &Location{} },
		attr.Seq("ID", 1),
		attr.Str("Loc", randomdata.Address),
	)

	userFactory := factory.New(
		func() interface{} { return &User{} },
		attr.Seq("ID", 1),
		attr.Factory("Location", locFactory, false).Process(
			func(attrGen attr.Attributer, data interface{}) error {
				loc := attrGen.GetVal().(*Location)
				user := data.(*User)
				loc.HostID = user.ID
				attrGen.SetVal(loc)
				return nil
			},
		),
	)

	user := userFactory.MustBuild().(*User)
	assert.NotNil(t, user.Location)
	assert.NotZero(t, user.Location.ID)
	assert.NotEmpty(t, user.Location.Loc)
	assert.Equal(t, user.Location.HostID, user.ID)
}
