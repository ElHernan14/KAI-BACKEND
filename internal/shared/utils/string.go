package utils

import "strings"

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func ToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
