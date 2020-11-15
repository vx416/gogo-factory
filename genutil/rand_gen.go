package genutil

import (
	"time"

	"github.com/Pallinder/go-randomdata"

	"github.com/google/uuid"
)

func RandInt(min, max int) func() int {
	return func() int {
		return randInts(min, max, 1)[0]
	}
}

func RandUint(min, max uint) func() uint {
	return func() uint {
		return randUints(min, max, 1)[0]
	}
}

func RandFloat(min, max float64) func() float64 {
	return func() float64 {
		return randFloats(min, max, 1)[0]
	}
}

func RandIntSet(set ...int) func() int {
	maxIndex := len(set) - 1
	return func() int {
		index := randInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func RandFloatSet(set ...float64) func() float64 {
	maxIndex := len(set) - 1
	return func() float64 {
		index := randInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func RandUintSet(set ...uint) func() uint {
	maxIndex := len(set) - 1
	return func() uint {
		index := randInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func RandAlph(n int) func() string {
	return func() string {
		return randString(n)
	}
}

func RandAlphSet(n int, alphs string) func() string {
	maxIndex := len(alphs) - 1
	return func() string {
		res := ""
		indexs := randInts(0, maxIndex, n)
		for _, index := range indexs {
			res += string(alphs[index])
		}
		return res
	}
}

func RandUUID() func() string {
	return func() string {
		return uuid.New().String()
	}
}

func RandStrSet(set ...string) func() string {
	maxIndex := len(set) - 1
	return func() string {
		index := randInts(0, maxIndex, 1)[0]
		return set[index]
	}
}

func RandBool(ratio float64) func() bool {
	return func() bool {
		return randBool(ratio)
	}
}

func RandTime(min, max time.Time) func() time.Time {
	minUnix := int(min.Unix())
	maxUnix := int(min.Unix())
	return func() time.Time {
		timeUnix := int64(randInts(minUnix, maxUnix, 1)[0])
		return time.Unix(timeUnix, 0)
	}
}

func RandFirstName(gender int) func() string {
	return func() string {
		return randomdata.FirstName(gender)
	}
}

func RandName(gender int) func() string {
	return func() string {
		return randomdata.FirstName(gender) + ", " + randomdata.LastName()
	}
}
