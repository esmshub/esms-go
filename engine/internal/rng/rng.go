package rng

import (
	"math/rand"
	"time"
)

var prng rand.Rand = *rand.New(rand.NewSource(time.Now().UnixNano()))

func SeedRandom(src rand.Source) {
	prng = *rand.New(src)
}

func Randomp(threshold int) bool {
	value := Random(10000)
	return value < threshold
}

func Random(value int) int {
	return prng.Intn(value)
}

func RandomRange(min, max int) int {
	return rand.Intn(max-min) + min
}
