package utils

import (
	"regexp"
	"strings"
	"time"
)

func CompileDate(date string) (time.Time, error) {
	// loc, _ := time.LoadLocation("Australia/Melbourne")
	re := regexp.MustCompile("[TZ]")

	cleanedTimestamp := re.ReplaceAllString(date, " ")

	datetime, err := time.Parse(time.DateTime, strings.TrimSpace(cleanedTimestamp))

	return datetime, err
}
