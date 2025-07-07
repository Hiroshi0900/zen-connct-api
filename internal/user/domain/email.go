package domain

import (
	"errors"
	"regexp"
	"strings"
)

// Email represents an email address value object
type Email struct {
	value string
}

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (*Email, error) {
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	if !isValidEmail(normalizedEmail) {
		return nil, errors.New("invalid email format")
	}

	return &Email{value: normalizedEmail}, nil
}

// isValidEmail validates email format with comprehensive checks
func isValidEmail(email string) bool {
	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return false
	}

	// Split and validate local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Validate local part
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}

	// Validate domain part
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") || strings.Contains(domainPart, "..") {
		return false
	}

	return true
}

// Value returns the email value
func (e *Email) Value() string {
	return e.value
}

// Equals checks if two Email values are equal
func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}

// String returns the string representation
func (e *Email) String() string {
	return e.value
}