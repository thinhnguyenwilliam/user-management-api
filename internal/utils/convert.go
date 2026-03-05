// user-management-api/internal/utils/convert.go
package utils

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	spaceRegex    = regexp.MustCompile(`\s+`)
)

func NormalizeEmail(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	str = spaceRegex.ReplaceAllString(str, " ")
	return str
}

func CamelToSnake(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
