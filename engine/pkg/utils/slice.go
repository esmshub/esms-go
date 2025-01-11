package utils

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

func Each[T any](s []T, f func(T)) {
	for _, v := range s {
		f(v)
	}
}
