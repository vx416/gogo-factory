package genutil

import "time"

func Now() func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}

func FixInt(val int) func() int {
	return func() int {
		return val
	}
}

func FixUint(val uint) func() uint {
	return func() uint {
		return val
	}
}

func FixStr(str string) func() string {
	return func() string {
		return str
	}
}

func FixFloat(val float64) func() float64 {
	return func() float64 {
		return val
	}
}

func FixTime(val time.Time) func() time.Time {
	return func() time.Time {
		return val
	}
}
