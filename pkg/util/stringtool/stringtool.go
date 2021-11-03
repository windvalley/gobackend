package stringtool

import (
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

// Diff creates an slice of slice values not included in the other given slice.
func Diff(base, exclude []string) (result []string) {
	excludeMap := make(map[string]bool)
	for _, s := range exclude {
		excludeMap[s] = true
	}
	for _, s := range base {
		if !excludeMap[s] {
			result = append(result, s)
		}
	}
	return result
}

// Unique ...
func Unique(ss []string) (result []string) {
	smap := make(map[string]bool)
	for _, s := range ss {
		smap[s] = true
	}
	for s := range smap {
		result = append(result, s)
	}
	return result
}

// CamelCaseToUnderscore ...
func CamelCaseToUnderscore(str string) string {
	return govalidator.CamelCaseToUnderscore(str)
}

// UnderscoreToCamelCase ...
func UnderscoreToCamelCase(str string) string {
	return govalidator.UnderscoreToCamelCase(str)
}

// FindString ...
func FindString(array []string, str string) int {
	for index, s := range array {
		if str == s {
			return index
		}
	}
	return -1
}

// StringIn ...
func StringIn(str string, array []string) bool {
	return FindString(array, str) > -1
}

// Reverse ...
func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}
