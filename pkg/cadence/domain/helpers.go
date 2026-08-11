package domain

import (
	"strings"
	"unicode/utf8"
)

func checkStringLimits(str string, limit int, emptyErr, tooLongErr error) error {
	str = strings.TrimSpace(str)
	if str == "" {
		return emptyErr
	}
	if utf8.RuneCountInString(str) > limit {
		return tooLongErr
	}
	return nil
}
