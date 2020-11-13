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
	userData, err := userFactory.HasMany("Rented", 5, homeAss).Build()
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
	userData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", 5, homeAss).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	testHasMany(t, user, 5)
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

	homeAss := homeFactory.BelongsTo("Location", locFactory.ToAssociation()).ToAssociation().ReferField("ID").ForeignField("HostID")
	userData, err := userFactory.HasOne("Home", homeAss).HasMany("Rented", 5, homeAss).Build()
	assert.NoError(t, err)
	user := userData.(*User)
	testHasOne(t, user)
	testHasMany(t, user, 5)
	testBelongsTo(t, user.Home)
}
