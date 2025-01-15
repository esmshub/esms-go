package rng

import (
	"math/rand/v2"
)

// Initial seed value from the current time
var s1 = rand.Uint64()
var s2 = rand.Uint64()

// Create and seed the generator.
// Note - using a fixed seed will produce the same output on every run.
var r = rand.New(rand.NewPCG(s1, s1))

func GetSeed() uint64 {
	return (s1 << 32) | (s2 & 0xFFFFFFFF)
}

// Seed updates the random generator with a new seed
func Seed(seed uint64) {
	s1 = seed >> 32
	s2 = seed & 0xFFFFFFFF
	r = rand.New(rand.NewPCG(s1, s2))
}

// Randomp returns true if a randomly generated value is below the threshold
func Randomp(threshold int) bool {
	value := Random(10000)
	return value < threshold
}

// Random generates a random number in the range [0, value)
func Random(value int) int {
	return r.IntN(value)
}

// RandomRange generates a random number in the range [min, max)
func RandomRange(min, max int) int {
	return r.IntN(max-min) + min
}
