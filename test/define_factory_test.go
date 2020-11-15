package test

import (
	"testing"
	"time"

	"github.com/vx416/gogo-factory/genutil"

	"github.com/stretchr/testify/assert"

	"github.com/Pallinder/go-randomdata"
	factory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
)

func between(t *testing.T, target, max, min int) {
	assert.LessOrEqual(t, target, max)
	assert.GreaterOrEqual(t, target, min)
}

func TestBasicAttributes(t *testing.T) {
	phoneSet := []string{"091234567", "09765432", "096789234"}
	userFactory := factory.New(
		&User{CreatedAt: time.Now(), Host: true},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Str("Username", randomdata.LastName),
		attr.Str("Phone", genutil.SeqStrSet(phoneSet...)),
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
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Str("Username", genutil.RandStrSet(userNameSet...)),
		attr.Str("Phone", genutil.RandStrSet(phoneSet...)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
		attr.Int("Age", genutil.RandInt(25, 100)),
		attr.Float("Height", genutil.RandFloat(55.0, 99.0)),
		attr.Float("Weight", genutil.RandFloat(155.0, 190.0)),
		attr.Bool("Host", genutil.RandBool(0.5)),
		attr.Time("CreatedAt", genutil.RandTime(minTime, maxTime)),
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
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Time("UpdatedAt", genutil.RandTime(minTime, maxTime)),
		attr.Str("PtrString", genutil.RandStrSet(stringSet...)),
	)
	user := userFactory.MustBuild().(*User)
	between(t, int(user.UpdatedAt.Time.Unix()), int(maxTime.Unix()), int(minTime.Unix()))
	assert.NotNil(t, user.PtrString)
	assert.Subset(t, stringSet, []string{*user.PtrString})
}

func TestOmit(t *testing.T) {
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	employeeFactory := factory.New(
		&Employee{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandIntSet(1, 2)),
		attr.Time("UpdatedAt", genutil.RandTime(minTime, maxTime)),
	)

	omited := employeeFactory.Omit("Gender")
	e1 := employeeFactory.MustBuild().(*Employee)
	e2 := omited.MustBuild().(*Employee)
	assert.Equal(t, e1.ID, int64(1))
	assert.Equal(t, e2.ID, int64(2))
	assert.NotZero(t, e1.Gender)
	assert.Zero(t, e2.Gender)

}

func TestOverWrite(t *testing.T) {
	minTime, maxTime := time.Now().Add(-30*24*time.Hour), time.Now()
	employeeFactory := factory.New(
		&Employee{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandIntSet(1, 2)),
		attr.Time("UpdatedAt", genutil.RandTime(minTime, maxTime)),
	)

	newFactory := employeeFactory.Attrs(
		attr.Int("Gender", genutil.FixInt(3)),
		attr.Int("Age", genutil.SeqInt(20, 5)),
	)

	e1 := employeeFactory.MustBuild().(*Employee)
	assert.Equal(t, e1.ID, int64(1))
	assert.Nil(t, e1.Age)
	assert.Less(t, int(e1.Gender), 3)

	e2 := newFactory.MustBuild().(*Employee)
	assert.Equal(t, e2.ID, int64(2))
	assert.NotNil(t, e2.Age)
	assert.Equal(t, e2.Gender, Gender(3))
}
