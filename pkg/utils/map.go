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

// DeepMerge merges src into dst recursively.
func DeepMerge(dst, src map[string]interface{}) map[string]interface{} {
	for key, val := range src {
		if vMap, ok := val.(map[string]interface{}); ok {
			// If dst has a nested map at this key, merge recursively
			if dMap, ok := dst[key].(map[string]interface{}); ok {
				dst[key] = DeepMerge(dMap, vMap)
			} else {
				// Otherwise, replace entirely
				dst[key] = vMap
			}
		} else {
			// Replace scalar values
			dst[key] = val
		}
	}
	return dst
}
