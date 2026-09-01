package validator

import (
	"regexp"
	"unicode"
)

// PasswordRequirements defines password validation rules following NIST SP 800-63B.
type PasswordRequirements struct {
	MinLength         int
	RequireUppercase  bool
	RequireLowercase  bool
	RequireNumbers    bool
	RequireSpecial    bool
	ProhibitCommon    bool
}

// DefaultPasswordRequirements returns recommended password requirements.
func DefaultPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:        14, // NIST recommendation: 14+ characters
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   false, // NIST says special chars not required
		ProhibitCommon:   true,
	}
}

// LegacyPasswordRequirements returns more relaxed requirements (minimum 8 characters).
// Use only for backward compatibility.
func LegacyPasswordRequirements() PasswordRequirements {
	return PasswordRequirements{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSpecial:   false,
		ProhibitCommon:   false,
	}
}

// ValidatePassword validates a password against the requirements.
// Returns a list of validation error messages. Empty list means password is valid.
func ValidatePassword(password string, req PasswordRequirements) []string {
	var errors []string
	
	// Check minimum length
	if len(password) < req.MinLength {
		errors = append(errors, "password must be at least "+string(rune(req.MinLength))+" characters long")
	}
	
	// Check for uppercase
	if req.RequireUppercase && !hasUppercase(password) {
		errors = append(errors, "password must contain at least one uppercase letter")
	}
	
	// Check for lowercase
	if req.RequireLowercase && !hasLowercase(password) {
		errors = append(errors, "password must contain at least one lowercase letter")
	}
	
	// Check for numbers
	if req.RequireNumbers && !hasNumber(password) {
		errors = append(errors, "password must contain at least one number")
	}
	
	// Check for special characters
	if req.RequireSpecial && !hasSpecialChar(password) {
		errors = append(errors, "password must contain at least one special character (!@#$%^&*)")
	}
	
	// Check against common passwords
	if req.ProhibitCommon && isCommonPassword(password) {
		errors = append(errors, "password is too common. Please use a more unique password")
	}
	
	return errors
}

// ValidatePasswordDefault validates against default requirements.
func ValidatePasswordDefault(password string) []string {
	return ValidatePassword(password, DefaultPasswordRequirements())
}

// HasValidEmail checks if an email is valid.
func HasValidEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailPattern)
	return re.MatchString(email)
}

// Helper functions

func hasUppercase(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

func hasNumber(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, r := range s {
		for _, sc := range specialChars {
			if r == sc {
				return true
			}
		}
	}
	return false
}

// isCommonPassword checks against a basic list of common passwords.
// In production, consider using a more comprehensive list or an API service.
func isCommonPassword(password string) bool {
	commonPasswords := map[string]bool{
		"password":       true,
		"123456":         true,
		"12345678":       true,
		"qwerty":         true,
		"abc123":         true,
		"monkey":         true,
		"1234567":        true,
		"letmein":        true,
		"trustno1":       true,
		"dragon":         true,
		"baseball":       true,
		"iloveyou":       true,
		"sunshine":       true,
		"princess":       true,
		"starlight":      true,
		"admin":          true,
		"test":           true,
		"guest":          true,
		"aaaaaa":         true,
		"123123":         true,
	}
	
	return commonPasswords[password]
}
