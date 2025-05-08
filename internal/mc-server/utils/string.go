package utils

import (
	"strings"
)

// ToCamelCase converts a kebab-case string to camelCase
func ToCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}
