package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	factory "github.com/vicxu416/gogo-factory"
	"github.com/vicxu416/gogo-factory/attr"
	"github.com/vicxu416/gogo-factory/randutil"
)

func testBelongsTo(t *testing.T, home *Home) {
	assert.NotZero(t, home.ID)
	assert.NotNil(t, home.Location)
	assert.NotZero(t, home.Location.ID)
}

func testHasOne(t *testing.T, user *User) {
	assert.NotZero(t, user.ID)
	assert.NotNil(t, user.Home)
	assert.NotZero(t, user.Home.ID)
	assert.NotZero(t, user.Home.HostID)
	assert.Equal(t, user.ID, user.Home.HostID)
}

func testHasMany(t *testing.T, user *User, num int) {
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Rented)
	assert.Len(t, user.Rented, num)
	for _, home := range user.Rented {
		assert.NotZero(t, home.ID)
		assert.NotZero(t, home.HostID)
		assert.Equal(t, user.ID, home.HostID)
	}
}

func testManyAndMany(t *testing.T, user *User, belongsTo bool) {
	assert.NotZero(t, user.ID)
	assert.NotEmpty(t, user.Countries)
	for _, country := range user.Countries {
		assert.NotZero(t, country.ID)
		assert.NotEmpty(t, country.Homes)
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
		attr.Seq("ID", 1),
	)

	locFactory := factory.New(
		&Location{},
		attr.Seq("ID", 1),
	)

	homeData, err := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).Build()
	assert.NoError(t, err)
	home := homeData.(*Home)
	testBelongsTo(t, home)
}

func TestHasOne(t *testing.T) {
	userFactory := factory.New(
		&User{},
		attr.Seq("ID", 1),
		attr.Int("Gender", randutil.IntRander(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
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
		attr.Seq("ID", 1),
		attr.Int("Gender", randutil.IntRander(1, 2)),
	)

	homeFactory := factory.New(
		&Home{},
		attr.Seq("ID", 1),
	)

	homeAss := homeFactory.ToAssociation().ReferField("ID").ForeignField("HostID")
	userData, err := userFactory.HasMany("Rented", homeAss, 5).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasMany(t, user, 5)
}

func TestHasOneAndMany(t *testing.T) {
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
	userData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", homeAss, 5).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	testHasMany(t, user, 5)
}

func TestManyAndMany(t *testing.T) {
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
	user := userFactory.HasMany("Countries", countryAss, 5).MustBuild().(*User)
	testManyAndMany(t, user, false)
}

func TestAllAssociations(t *testing.T) {
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
