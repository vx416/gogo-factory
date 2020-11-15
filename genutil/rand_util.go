package genutil

import (
	"math/rand"
	"time"
)

var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randFloats(min, max float64, n int) []float64 {
	rand.Seed(time.Now().UnixNano())
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func randInts(min, max int, n int) []int {
	rand.Seed(time.Now().UnixNano())
	res := make([]int, n)
	for i := range res {
		res[i] = rand.Intn(max-min+1) + min
	}
	return res
}

func randUints(min, max uint, n int) []uint {
	rand.Seed(time.Now().UnixNano())
	res := make([]uint, n)
	for i := range res {
		res[i] = uint(rand.Intn(int(max)-int(min)+1) + int(min))
	}
	return res
}

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func randBool(ratio float64) bool {
	rand.Seed(time.Now().UnixNano())
	if ratio > rand.Float64() {
		return true
	}
	return false
}
