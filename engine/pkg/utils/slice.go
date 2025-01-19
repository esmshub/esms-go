package utils

import "slices"

func Map[T comparable](s []T, f func(T) T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func Reduce[T any, U any](s []T, f func(U, T) U, init U) U {
	result := init
	for _, v := range s {
		result = f(result, v)
	}
	return result
}

func EachFunc[T any](s []T, f func(T)) {
	for _, v := range s {
		f(v)
	}
}

func FindFunc[S ~[]E, E any](s S, f func(E) bool) E {
	var def E
	if i := slices.IndexFunc(s, f); i != -1 {
		return s[i]
	} else {
		return def
	}
}

func CountFunc[S ~[]E, E any](s S, f func(E) bool) int {
	count := 0
	for _, v := range s {
		if f(v) {
			count++
		}
	}
	return count
}

func SumFunc[S ~[]E, E any, K int | float64](s S, f func(E) K) K {
	var sum K = 0
	for _, v := range s {
		sum += f(v)
	}
	return sum
}
