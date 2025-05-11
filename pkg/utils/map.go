package utils

import (
	"fmt"
	"reflect"
)

func MustGetKey[V any](m map[string]any, key string) V {
	v, ok := m[key]
	if !ok {
		panic(fmt.Sprintf("key %s not found in map", key))
	}

	if result, ok := v.(V); ok {
		return result
	} else {
		panic(fmt.Sprintf("key %s is not of type %s", key, reflect.TypeOf(v).Name()))
	}
}
