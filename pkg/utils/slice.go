package utils

import "slices"

func Map[T comparable, K any](s []T, f func(T) K) []K {
	result := make([]K, len(s))
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

func FilterFunc[S ~[]E, E any](s S, f func(E) bool) []E {
	results := []E{}
	for _, v := range s {
		if f(v) {
			results = append(results, v)
		}
	}
	return results
}

func Range(start, end int) []int {
	result := make([]int, end-start)
	for i := start; i < end; i++ {
		result[i-start] = i
	}
	return result
}
