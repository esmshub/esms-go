package utils

import (
	"reflect"
	"regexp"
)

func IsNumericalStr(s string) bool {
	match, _ := regexp.MatchString(`^-?\d+(\.\d+)?$`, s)
	return match
}

func IsNumber(value interface{}) bool {
	// Get the reflect.Kind of the value
	kind := reflect.TypeOf(value).Kind()

	// Check if the kind is one of the numeric kinds
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
