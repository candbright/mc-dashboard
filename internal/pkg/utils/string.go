package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToCamelCase converts a kebab-case string to camelCase
func ToCamelCase(s string) string {
	parts := strings.Split(s, "-")
	for i := 1; i < len(parts); i++ {
		parts[i] = cases.Title(language.English).String(parts[i])
	}
	return strings.Join(parts, "")
}
