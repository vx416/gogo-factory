package genutil

import (
	"time"
)

func SeqInt(start int, delta int) func() int {
	return func() int {
		val := start
		start = start + delta
		return val
	}
}

func SeqUint(start uint, delta uint) func() uint {
	return func() uint {
		val := start
		start = start + delta
		return val
	}
}

func SeqFloat(start float64, delta float64) func() float64 {
	return func() float64 {
		val := start
		start = start + delta
		return val
	}
}

func SeqTime(t time.Time, delta time.Duration) func() time.Time {
	return func() time.Time {
		val := t
		t = t.Add(delta)
		return val
	}
}

func SeqIntSet(set ...int) func() int {
	index := 0
	maxIndex := len(set) - 1
	return func() int {
		val := set[index]
		index++
		if index > maxIndex {
			index = 0
		}
		return val
	}
}

func SeqStrSet(set ...string) func() string {
	index := 0
	maxIndex := len(set) - 1
	return func() string {
		val := set[index]
		index++
		if index > maxIndex {
			index = 0
		}
		return val
	}
}

func SeqUintSet(set ...uint) func() uint {
	index := 0
	maxIndex := len(set) - 1
	return func() uint {
		val := set[index]
		index++
		if index > maxIndex {
			index = 0
		}
		return val
	}
}

func SeqFloatSet(set ...float64) func() float64 {
	index := 0
	maxIndex := len(set) - 1
	return func() float64 {
		val := set[index]
		index++
		if index > maxIndex {
			index = 0
		}
		return val
	}
}

func SeqTimeSet(set ...time.Time) func() time.Time {
	index := 0
	maxIndex := len(set) - 1
	return func() time.Time {
		val := set[index]
		index++
		if index > maxIndex {
			index = 0
		}
		return val
	}
}
