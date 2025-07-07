package domain

import (
	"strings"
	"testing"
)

func TestNewPassword_ValidPassword_ShouldCreatePassword(t *testing.T) {
	// Arrange
	validPassword := "ValidPass123!"

	// Act
	password, err := NewPassword(validPassword)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if password == nil {
		t.Error("Expected password to be created, got nil")
	}
	if password.Value() != validPassword {
		t.Errorf("Expected password value %s, got %s", validPassword, password.Value())
	}
}

func TestNewPassword_TooShort_ShouldReturnError(t *testing.T) {
	// Arrange
	shortPassword := "Short1!"

	// Act
	password, err := NewPassword(shortPassword)

	// Assert
	if err == nil {
		t.Error("Expected error for short password, got nil")
	}
	if password != nil {
		t.Error("Expected nil password for short input, got password object")
	}
	expectedError := "password must be at least 8 characters long"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewPassword_MissingUppercase_ShouldReturnError(t *testing.T) {
	// Arrange
	noUppercasePassword := "validpass123!"

	// Act
	password, err := NewPassword(noUppercasePassword)

	// Assert
	if err == nil {
		t.Error("Expected error for password without uppercase, got nil")
	}
	if password != nil {
		t.Error("Expected nil password, got password object")
	}
	expectedError := "password must contain at least one uppercase letter"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewPassword_MissingLowercase_ShouldReturnError(t *testing.T) {
	// Arrange
	noLowercasePassword := "VALIDPASS123!"

	// Act
	password, err := NewPassword(noLowercasePassword)

	// Assert
	if err == nil {
		t.Error("Expected error for password without lowercase, got nil")
	}
	if password != nil {
		t.Error("Expected nil password, got password object")
	}
	expectedError := "password must contain at least one lowercase letter"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewPassword_MissingDigit_ShouldReturnError(t *testing.T) {
	// Arrange
	noDigitPassword := "ValidPass!"

	// Act
	password, err := NewPassword(noDigitPassword)

	// Assert
	if err == nil {
		t.Error("Expected error for password without digit, got nil")
	}
	if password != nil {
		t.Error("Expected nil password, got password object")
	}
	expectedError := "password must contain at least one number"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewPassword_MissingSpecialChar_ShouldReturnError(t *testing.T) {
	// Arrange
	noSpecialCharPassword := "ValidPass123"

	// Act
	password, err := NewPassword(noSpecialCharPassword)

	// Assert
	if err == nil {
		t.Error("Expected error for password without special character, got nil")
	}
	if password != nil {
		t.Error("Expected nil password, got password object")
	}
	expectedError := "password must contain at least one special character"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestPassword_Hash_ShouldGenerateValidHash(t *testing.T) {
	// Arrange
	password, _ := NewPassword("ValidPass123!")

	// Act
	hash, err := password.Hash()

	// Assert
	if err != nil {
		t.Errorf("Expected no error when hashing, got %v", err)
	}
	if hash == "" {
		t.Error("Expected hash to be generated, got empty string")
	}
	
	// Hash should contain salt and hash separated by ":"
	parts := strings.Split(hash, ":")
	if len(parts) != 2 {
		t.Errorf("Expected hash format 'salt:hash', got %s", hash)
	}
	
	// Salt should be hex encoded (32 characters for 16 bytes)
	if len(parts[0]) != 32 {
		t.Errorf("Expected salt length 32, got %d", len(parts[0]))
	}
	
	// Hash should be hex encoded (128 characters for 64 bytes)
	if len(parts[1]) != 128 {
		t.Errorf("Expected hash length 128, got %d", len(parts[1]))
	}
}

func TestPassword_Hash_DifferentPasswordsDifferentHashes(t *testing.T) {
	// Arrange
	password1, _ := NewPassword("ValidPass123!")
	password2, _ := NewPassword("AnotherPass456#")

	// Act
	hash1, _ := password1.Hash()
	hash2, _ := password2.Hash()

	// Assert
	if hash1 == hash2 {
		t.Error("Expected different passwords to produce different hashes")
	}
}

func TestPassword_Hash_SamePasswordDifferentHashes(t *testing.T) {
	// Arrange
	password, _ := NewPassword("ValidPass123!")

	// Act
	hash1, _ := password.Hash()
	hash2, _ := password.Hash()

	// Assert
	if hash1 == hash2 {
		t.Error("Expected same password to produce different hashes due to random salt")
	}
}

func TestVerifyHash_CorrectPassword_ShouldReturnTrue(t *testing.T) {
	// Arrange
	passwordValue := "ValidPass123!"
	password, _ := NewPassword(passwordValue)
	hash, _ := password.Hash()

	// Act
	result := VerifyHash(passwordValue, hash)

	// Assert
	if !result {
		t.Error("Expected verification of correct password to return true")
	}
}

func TestVerifyHash_IncorrectPassword_ShouldReturnFalse(t *testing.T) {
	// Arrange
	password, _ := NewPassword("ValidPass123!")
	hash, _ := password.Hash()
	wrongPassword := "WrongPass456#"

	// Act
	result := VerifyHash(wrongPassword, hash)

	// Assert
	if result {
		t.Error("Expected verification of incorrect password to return false")
	}
}

func TestVerifyHash_InvalidHashFormat_ShouldReturnFalse(t *testing.T) {
	// Arrange
	invalidHash := "invalid-hash-format"

	// Act
	result := VerifyHash("ValidPass123!", invalidHash)

	// Assert
	if result {
		t.Error("Expected verification with invalid hash format to return false")
	}
}

func TestPassword_Equals_SameValue_ShouldReturnTrue(t *testing.T) {
	// Arrange
	password1, _ := NewPassword("ValidPass123!")
	password2, _ := NewPassword("ValidPass123!")

	// Act
	result := password1.Equals(password2)

	// Assert
	if !result {
		t.Error("Expected passwords with same value to be equal")
	}
}

func TestPassword_Equals_DifferentValue_ShouldReturnFalse(t *testing.T) {
	// Arrange
	password1, _ := NewPassword("ValidPass123!")
	password2, _ := NewPassword("DifferentPass456#")

	// Act
	result := password1.Equals(password2)

	// Assert
	if result {
		t.Error("Expected passwords with different values to not be equal")
	}
}

func TestPassword_Equals_NilPassword_ShouldReturnFalse(t *testing.T) {
	// Arrange
	password, _ := NewPassword("ValidPass123!")

	// Act
	result := password.Equals(nil)

	// Assert
	if result {
		t.Error("Expected password compared to nil to return false")
	}
}