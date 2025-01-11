package utils

import (
	"fmt"
	"strings"
)

func Substring(src string, from string, to string) (string, error) {
	startIndex := strings.Index(src, from)
	endIndex := strings.Index(src, to)
	if startIndex == -1 {
		return src, fmt.Errorf("substring '%s' was not found in the source string", from)
	} else if endIndex == -1 {
		return src, fmt.Errorf("substring '%s' was not found in the source string", to)
	} else {
		return src[startIndex+1 : endIndex], nil
	}
}
