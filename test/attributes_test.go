package test

import (
	"testing"
	"time"

	"github.com/vicxu416/gogo-factory/randutil"

	"github.com/stretchr/testify/assert"

	"github.com/Pallinder/go-randomdata"
	factory "github.com/vicxu416/gogo-factory"
	"github.com/vicxu416/gogo-factory/attr"
)

func between(t *testing.T, target, max, min int) {
	assert.LessOrEqual(t, target, max)
	assert.GreaterOrEqual(t, target, min)
}

func TestBasicAttributes(t *testing.T) {
	phoneSet := []string{"091234567", "09765432", "096789234"}
	userFactory := factory.New(
		&User{CreatedAt: time.Now(), Host: true},
		attr.Seq("ID", 1),
		attr.Str("Username", randomdata.LastName),
		attr.StrSeq("Phone", phoneSet),
		attr.Int("Gender", func() int { return int(randomdata.Number(1, 2)) }),
		attr.Attr("Age", func() interface{} { return int32(randomdata.Number(1, 100)) }),
		attr.Float("Weight", func() float64 { return randomdata.Decimal(1, 20, 1) }),
		attr.Float("Height", func() float64 { return randomdata.Decimal(1, 20, 1) }),
	)

	users := make([]*User, 0, 1)
	for i := 1; i <= 5; i++ {
		user := userFactory.MustBuild().(*User)
		assert.Equal(t, user.ID, int64(i))
		assert.NotEmpty(t, user.Username)
		assert.NotEmpty(t, user.Phone)
		assert.NotZero(t, user.Gender)
		assert.NotZero(t, user.Age)
		assert.NotZero(t, user.Weight)
		assert.NotZero(t, user.Height)
		users = append(users, user)
	}
}

func TestRandAttr(t *testing.T) {
	userNameSet := []string{"divaid", "vic", "shelly", "jason"}
	phoneSet := []string{"090123543", "0954323123", "0924325345"}
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Str("Username", randutil.StrSetRander(userNameSet...)),
		attr.Str("Phone", randutil.StrSetRander(phoneSet...)),
		attr.Int("Gender", randutil.IntRander(1, 2)),
		attr.Int("Age", randutil.IntRander(25, 100)),
		attr.Float("Height", randutil.FloatRander(55.0, 99.0)),
		attr.Float("Weight", randutil.FloatRander(155.0, 190.0)),
		attr.Bool("Host", randutil.BoolRander(0.5)),
		attr.Time("CreatedAt", randutil.TimeRander(minTime, maxTime)),
	)

	user := userFactory.MustBuild().(*User)
	assert.Subset(t, userNameSet, []string{user.Username})
	assert.Subset(t, phoneSet, []string{user.Phone})
	between(t, int(user.Gender), 2, 1)
	between(t, int(*user.Age), 100, 25)
	assert.InDelta(t, float64(user.Height), 90.0, 55.0)
	assert.InDelta(t, float64(user.Weight), 190.0, 155.0)
	between(t, int(user.CreatedAt.Unix()), int(maxTime.Unix()), int(minTime.Unix()))
}

func TestNullableFields(t *testing.T) {
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	stringSet := []string{"ptr_string_1", "ptr_string_2"}
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Time("UpdatedAt", randutil.TimeRander(minTime, maxTime)),
		attr.Str("PtrString", randutil.StrSetRander(stringSet...)),
	)
	user := userFactory.MustBuild().(*User)
	between(t, int(user.UpdatedAt.Time.Unix()), int(maxTime.Unix()), int(minTime.Unix()))
	assert.NotNil(t, user.PtrString)
	assert.Subset(t, stringSet, []string{*user.PtrString})
}
