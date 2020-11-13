package test

import (
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"
	factory "github.com/vicxu416/gogo-factory"
	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/randutil"
)

func TestBuildN(t *testing.T) {
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

	usersData, err := userFactory.BuildN(10)
	assert.NoError(t, err)
	users := usersData.([]*User)
	assert.Len(t, users, 10)
	for i, user := range users {
		assert.Equal(t, user.ID, int64(i+1))
		assert.NotEmpty(t, user.Username)
		assert.NotEmpty(t, user.Phone)
		assert.NotZero(t, user.Gender)
		assert.NotZero(t, user.Age)
		assert.NotZero(t, user.Weight)
		assert.NotZero(t, user.Height)
	}
}

func TestBuildNBelongsTo(t *testing.T) {
	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
	)

	locFactory := factory.New(
		&Location{},
		attr.Seq("ID", 1),
	)

	homesData, err := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).BuildN(5)
	assert.NoError(t, err)
	homes := homesData.([]*Home)
	assert.Len(t, homes, 5)
	for i := range homes {
		testBelongsTo(t, homes[i])
	}
}

func TestBuildNHasOneAndMany(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Int("Gender", randutil.IntRander(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("HostID")
	usersData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", homeAss, 5).BuildN(6)
	assert.NoError(t, err)
	users := usersData.([]*User)
	assert.Len(t, users, 6)
	for i := range users {
		testHasOne(t, users[i])
		testHasMany(t, users[i], 5)
	}
}

func TestBuildNManyAndMany(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Int("Gender", randutil.IntRander(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
	)

	country := factory.New(
		&Country{},
		attr.Seq("ID", 1),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("CountryID")
	countryAss := country.HasMany("Homes", homeAss, 10).ToAssociation().ReferField("ID").ForeignField("HostID")
	usersData, err := userFactory.HasMany("Countries", countryAss, 5).BuildN(10)
	assert.NoError(t, err)
	users := usersData.([]*User)
	for _, user := range users {
		testManyAndMany(t, user, false)
	}
}

func TestBuildNAll(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Int("Gender", randutil.IntRander(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
	)

	locFactory := factory.New(
		&Location{},
		attr.Seq("ID", 1),
	)

	country := factory.New(
		&Country{},
		attr.Seq("ID", 1),
	)

	homeAss := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).ToAssociation().ReferField("ID").ForeignField("HostID")
	countryHomeAss := homeAss.ReferField("ID").ForeignField("CountryID")
	countryAss := country.HasMany("Homes", countryHomeAss, 10).ToAssociation().ReferField("ID").ForeignField("HostID")
	usersData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", homeAss, 5).HasMany("Countries", countryAss, 5).BuildN(15)
	assert.NoError(t, err)
	users := usersData.([]*User)
	assert.Len(t, users, 15)
	for i := range users {
		testHasOne(t, users[i])
		testHasMany(t, users[i], 5)
		testBelongsTo(t, users[i].Home)
		for i := range users[i].Rented {
			testBelongsTo(t, users[i].Rented[i])
		}
	}
}
