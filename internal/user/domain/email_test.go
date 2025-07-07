package domain

import (
	"testing"
)

func TestNewEmail_ValidEmail_ShouldCreateEmail(t *testing.T) {
	// Arrange
	validEmail := "test@example.com"

	// Act
	email, err := NewEmail(validEmail)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if email == nil {
		t.Error("Expected email to be created, got nil")
	}
	if email.Value() != validEmail {
		t.Errorf("Expected email value %s, got %s", validEmail, email.Value())
	}
}

func TestNewEmail_EmptyEmail_ShouldReturnError(t *testing.T) {
	// Arrange
	emptyEmail := ""

	// Act
	email, err := NewEmail(emptyEmail)

	// Assert
	if err == nil {
		t.Error("Expected error for empty email, got nil")
	}
	if email != nil {
		t.Error("Expected nil email for empty input, got email object")
	}
	expectedError := "email cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}

func TestNewEmail_InvalidFormat_ShouldReturnError(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{"missing @", "testexample.com"},
		{"missing domain", "test@"},
		{"missing local part", "@example.com"},
		{"consecutive dots", "test..user@example.com"},
		{"starts with dot", ".test@example.com"},
		{"ends with dot", "test.@example.com"},
		{"domain starts with dot", "test@.example.com"},
		{"domain ends with dot", "test@example.com."},
		{"invalid characters", "test@exam ple.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			email, err := NewEmail(tc.email)

			// Assert
			if err == nil {
				t.Errorf("Expected error for invalid email '%s', got nil", tc.email)
			}
			if email != nil {
				t.Errorf("Expected nil email for invalid input '%s', got email object", tc.email)
			}
			expectedError := "invalid email format"
			if err.Error() != expectedError {
				t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
			}
		})
	}
}

func TestNewEmail_EmailNormalization_ShouldLowercaseAndTrim(t *testing.T) {
	// Arrange
	unnormalizedEmail := "  TEST@EXAMPLE.COM  "
	expectedEmail := "test@example.com"

	// Act
	email, err := NewEmail(unnormalizedEmail)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if email.Value() != expectedEmail {
		t.Errorf("Expected normalized email '%s', got '%s'", expectedEmail, email.Value())
	}
}

func TestEmail_Equals_SameValue_ShouldReturnTrue(t *testing.T) {
	// Arrange
	email1, _ := NewEmail("test@example.com")
	email2, _ := NewEmail("test@example.com")

	// Act
	result := email1.Equals(email2)

	// Assert
	if !result {
		t.Error("Expected emails with same value to be equal")
	}
}

func TestEmail_Equals_DifferentValue_ShouldReturnFalse(t *testing.T) {
	// Arrange
	email1, _ := NewEmail("test1@example.com")
	email2, _ := NewEmail("test2@example.com")

	// Act
	result := email1.Equals(email2)

	// Assert
	if result {
		t.Error("Expected emails with different values to not be equal")
	}
}

func TestEmail_Equals_NilEmail_ShouldReturnFalse(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")

	// Act
	result := email.Equals(nil)

	// Assert
	if result {
		t.Error("Expected email compared to nil to return false")
	}
}

func TestEmail_String_ShouldReturnEmailValue(t *testing.T) {
	// Arrange
	emailValue := "test@example.com"
	email, _ := NewEmail(emailValue)

	// Act
	result := email.String()

	// Assert
	if result != emailValue {
		t.Errorf("Expected string representation '%s', got '%s'", emailValue, result)
	}
}