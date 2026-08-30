package validator

import (
	"regexp"
	"strings"
)

// IsValidPhone performs a light normalization/validation for phone numbers.
// Accepts digits, optional leading '+', spaces and dashes. Returns true when
// it looks like a plausible phone number (at least 9 digits).
func IsValidPhone(phone string) bool {
	digits := digitsOnly(phone)
	return len(digits) >= 9 && len(digits) <= 15
}

// NormalizePhone strips non-digit characters and optionally strips a leading
// '0'/+country prefix, returning a plain E.164-style local number as provided.
func NormalizePhone(phone string) string {
	return digitsOnly(phone)
}

func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// IsValidEmail does a basic structural email check.
func IsValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	match, _ := regexp.MatchString(`^[^@\s]+@[^@\s]+\.[^@\s]+$`, email)
	return match
}
