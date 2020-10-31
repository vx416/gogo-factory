package randutil

import (
	"time"

	"github.com/Pallinder/go-randomdata"

	"github.com/google/uuid"
)

func IntRander(min, max int) func() int {
	return func() int {
		return RandInts(min, max, 1)[0]
	}
}

func UintRander(min, max uint) func() uint {
	return func() uint {
		return RandUints(min, max, 1)[0]
	}
}

func FloatRander(min, max float64) func() float64 {
	return func() float64 {
		return RandFloats(min, max, 1)[0]
	}
}

func IntSetRander(set ...int) func() int {
	maxIndex := len(set) - 1
	return func() int {
		index := RandInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func FloatSetRander(set ...float64) func() float64 {
	maxIndex := len(set) - 1
	return func() float64 {
		index := RandInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func UintSetRander(set ...uint) func() uint {
	maxIndex := len(set) - 1
	return func() uint {
		index := RandInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func AlphRander(n int) func() string {
	return func() string {
		return RandString(n)
	}
}

func UUIDRander() func() string {
	return func() string {
		return uuid.New().String()
	}
}

func StrSetRander(set ...string) func() string {
	maxIndex := len(set) - 1
	return func() string {
		index := RandInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func BoolRander(ratio float64) func() bool {
	return func() bool {
		return RandBool(ratio)
	}
}

func TimeRander(min, max time.Time) func() time.Time {
	minUnix := int(min.Unix())
	maxUnix := int(min.Unix())
	return func() time.Time {
		timeUnix := int64(RandInts(minUnix, maxUnix, 1)[0])
		return time.Unix(timeUnix, 0)
	}
}

func FirstNameRander(gender int) func() string {
	return func() string {
		return randomdata.FirstName(gender)
	}
}

func NameRander(gender int) func() string {
	return func() string {
		return randomdata.FirstName(gender) + ", " + randomdata.LastName()
	}
}

func NowRander() func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}
