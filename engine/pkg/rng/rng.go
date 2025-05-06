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

// Float64 returns, as a float64, a pseudo-random number in the half-open interval [0.0,1.0).
func RandomF() float64 {
	return r.Float64()
}

// RandomRange generates a random number in the range [min, max)
func RandomRange(min, max int) int {
	return r.IntN(max-min) + min
}

// Shuffle pseudo-randomizes the order of elements.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) {
	r.Shuffle(n, swap)
}
