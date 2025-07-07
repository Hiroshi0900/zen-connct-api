package domain

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"regexp"
	"strings"
)

// Password policy constants matching frontend implementation
const (
	MinLength      = 8
	SaltBytes      = 16
	HashIterations = 10000
	HashKeyLength  = 64
)

// Password patterns for validation
var (
	uppercasePattern   = regexp.MustCompile(`[A-Z]`)
	lowercasePattern   = regexp.MustCompile(`[a-z]`)
	digitPattern       = regexp.MustCompile(`\d`)
	specialCharPattern = regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`)
)

// Password represents a password value object
type Password struct {
	value string
}

// NewPassword creates a new Password value object with validation
func NewPassword(password string) (*Password, error) {
	if len(password) < MinLength {
		return nil, fmt.Errorf("password must be at least %d characters long", MinLength)
	}

	if !uppercasePattern.MatchString(password) {
		return nil, errors.New("password must contain at least one uppercase letter")
	}

	if !lowercasePattern.MatchString(password) {
		return nil, errors.New("password must contain at least one lowercase letter")
	}

	if !digitPattern.MatchString(password) {
		return nil, errors.New("password must contain at least one number")
	}

	if !specialCharPattern.MatchString(password) {
		return nil, errors.New("password must contain at least one special character")
	}

	return &Password{value: password}, nil
}

// Value returns the password value
func (p *Password) Value() string {
	return p.value
}

// Hash generates a hash of the password using PBKDF2
func (p *Password) Hash() (string, error) {
	// Generate random salt
	salt := make([]byte, SaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Create hash using PBKDF2
	hash := pbkdf2.Key([]byte(p.value), salt, HashIterations, HashKeyLength, sha512.New)

	// Return salt:hash format
	return fmt.Sprintf("%s:%s", hex.EncodeToString(salt), hex.EncodeToString(hash)), nil
}

// VerifyHash verifies if the password matches the given hash
func VerifyHash(password, hash string) bool {
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	originalHash, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Generate hash from password
	verifyHash := pbkdf2.Key([]byte(password), salt, HashIterations, HashKeyLength, sha512.New)

	// Compare hashes
	return string(originalHash) == string(verifyHash)
}

// Equals checks if two Password values are equal
func (p *Password) Equals(other *Password) bool {
	if other == nil {
		return false
	}
	return p.value == other.value
}