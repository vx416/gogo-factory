package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	factory "github.com/vx416/gogo-factory"
	"github.com/vx416/gogo-factory/attr"
	"github.com/vx416/gogo-factory/genutil"
)

func testBelongsTo(t *testing.T, home *Home) {
	assert.NotZero(t, home.ID)
	assert.NotNil(t, home.Location, "home's location not nil")
	assert.NotZero(t, home.Location.ID)
}

func testHasOne(t *testing.T, user *User) {
	assert.NotZero(t, user.ID)
	assert.NotNil(t, user.Home, "user's home is nil")
	assert.NotZero(t, user.Home.ID)
	assert.NotZero(t, user.Home.HostID)
	assert.Equal(t, user.ID, user.Home.HostID)
}

func testHasMany(t *testing.T, user *User, num int) {
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Rented, "user's rented is empty")
	assert.Len(t, user.Rented, num)
	for _, home := range user.Rented {
		assert.NotZero(t, home.ID)
		assert.NotZero(t, home.HostID)
		assert.Equal(t, user.ID, home.HostID)
	}
}

func testManyAndMany(t *testing.T, user *User, belongsTo bool) {
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Countries, "user's country is empty")
	for _, country := range user.Countries {
		assert.NotZero(t, country.ID)
		assert.NotEmpty(t, country.Homes, "country's home is empty")
		for _, home := range country.Homes {
			assert.NotZero(t, home.ID)
			if belongsTo {
				testBelongsTo(t, home)
			}
		}
	}
}

func TestBelongsTo(t *testing.T) {
	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	locFactory := factory.New(
		&Location{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	homeData, err := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).Build()
	assert.NoError(t, err)
	assert.NotNil(t, homeData, "home not nil")
	home := homeData.(*Home)
	testBelongsTo(t, home)
}

func TestHasOne(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	userData, err := userFactory.HasOne("Home",
		homeFactory.ToAssociation().ReferField("ID").ForeignField("HostID")).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	user2 := userFactory.MustBuild().(*User)
	assert.Nil(t, user2.Home)
}

func TestHasMany(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("HostID")
	userData, err := userFactory.HasMany("Rented", homeAss, 5).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasMany(t, user, 5)
	userData, err = userFactory.HasMany("Rented", homeAss, 1).Build()
	assert.NoError(t, err)
	user = userData.(*User)
	testHasMany(t, user, 1)
}

func TestHasOneAndMany(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("HostID")
	userData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", homeAss, 5).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	testHasMany(t, user, 5)
}

func TestManyAndMany(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	country := factory.New(
		&Country{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("CountryID")
	countryAss := country.HasMany("Homes", homeAss, 10).ToAssociation().ReferField("ID").ForeignField("HostID")
	user := userFactory.HasMany("Countries", countryAss, 5).MustBuild().(*User)
	testManyAndMany(t, user, false)
}

func TestAllAssociations(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
		attr.Int("Gender", genutil.RandInt(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	locFactory := factory.New(
		&Location{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	country := factory.New(
		&Country{},
		attr.Int("ID", genutil.SeqInt(1, 1)),
	)

	homeAss := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).ToAssociation().ReferField("ID").ForeignField("HostID")
	countryHomeAss := homeAss.ReferField("ID").ForeignField("CountryID")
	countryAss := country.HasMany("Homes", countryHomeAss, 10).ToAssociation().ReferField("ID").ForeignField("HostID")
	userData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", homeAss, 5).HasMany("Countries", countryAss, 5).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	testHasMany(t, user, 5)
	testBelongsTo(t, user.Home)
	for i := range user.Rented {
		testBelongsTo(t, user.Rented[i])
	}
	testManyAndMany(t, user, true)
	home := homeFactory.MustBuild().(*Home)
	assert.Nil(t, home.Location)
}
