package gofactory

import (
	"regexp"
	"strings"
)

// DBTagProcess db tag process
func DBTagProcess(tagGetter TagGetter) string {
	return tagGetter.Get("db")
}

// GormTagProcess gorm tag process
func GormTagProcess(tagGetter TagGetter) string {
	regex := regexp.MustCompile(`column:(.*?(;|$))`)
	gormTag := tagGetter.Get("gorm")

	subMatch := regex.FindAllStringSubmatch(gormTag, -1)

	if len(subMatch) == 0 {
		return ""
	}

	firstMatch := subMatch[0]

	if len(firstMatch) < 1 {
		return ""
	}

	trimed := strings.TrimSpace(firstMatch[1])
	return strings.Trim(trimed, ";")
}
