package validator

import (
	"testing"
)

func TestIsValidPhone(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"123", false},
		{"123456789", true},
		{"+254712345678", true},
		{"254712345678", true},
		{"0712345678", true},
		{"abc", false},
		{"+254-712-345-678", true},
	}

	for _, tc := range tests {
		got := IsValidPhone(tc.input)
		if got != tc.expected {
			t.Errorf("IsValidPhone(%q) = %v; want %v", tc.input, got, tc.expected)
		}
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+254712345678", "254712345678"},
		{"0712-345-678", "0712345678"},
		{"712345678", "712345678"},
		{"", ""},
		{"abc", ""},
	}

	for _, tc := range tests {
		got := NormalizePhone(tc.input)
		if got != tc.expected {
			t.Errorf("NormalizePhone(%q) = %q; want %q", tc.input, got, tc.expected)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"user@", false},
		{"@domain.com", false},
		{"user@domain.com", true},
		{"user.name+tag@domain.co.ke", true},
		{"user@domain", false},
		{"user domain@test.com", false},
	}

	for _, tc := range tests {
		got := IsValidEmail(tc.input)
		if got != tc.expected {
			t.Errorf("IsValidEmail(%q) = %v; want %v", tc.input, got, tc.expected)
		}
	}
}
