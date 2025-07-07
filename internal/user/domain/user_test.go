package domain

import (
	"testing"
	"time"
)

// TODO: 以下以降のテストは全部見直す必要がある。
func TestNewUnverifiedUser_ValidInputs_ShouldCreateUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")

	// Act
	user, err := NewUnverifiedUser(email, password)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user to be created, got nil")
	}
	if user.ID() == "" {
		t.Error("Expected user ID to be generated")
	}
	if !user.Email().Equals(email) {
		t.Error("Expected user email to match input email")
	}
	if user.PasswordHash() == "" {
		t.Error("Expected password hash to be generated")
	}
	if user.IsVerified() {
		t.Error("Expected new user to be unverified")
	}
	if !user.IsUnverified() {
		t.Error("Expected new user to be unverified")
	}
}

func TestNewUnverifiedUser_ShouldGenerateUserRegisteredEvent(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")

	// Act
	user, _ := NewUnverifiedUser(email, password)

	// Assert
	events := user.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventName() != "UserRegistered" {
		t.Errorf("Expected UserRegistered event, got %s", event.EventName())
	}
	if event.AggregateID() != user.ID() {
		t.Error("Expected event aggregate ID to match user ID")
	}
}

func TestUnverifiedUser_VerifyEmail_ShouldReturnVerifiedUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	unverifiedUser, _ := NewUnverifiedUser(email, password)

	// Act
	verifiedUser := unverifiedUser.VerifyEmail()

	// Assert
	if verifiedUser == nil {
		t.Error("Expected verified user to be created")
	}
	if !verifiedUser.IsVerified() {
		t.Error("Expected user to be verified")
	}
	if verifiedUser.IsUnverified() {
		t.Error("Expected user to not be unverified")
	}
	if verifiedUser.ID() != unverifiedUser.ID() {
		t.Error("Expected verified user to have same ID as unverified user")
	}
	if !verifiedUser.Email().Equals(unverifiedUser.Email()) {
		t.Error("Expected verified user to have same email as unverified user")
	}
	if verifiedUser.PasswordHash() != unverifiedUser.PasswordHash() {
		t.Error("Expected verified user to have same password hash as unverified user")
	}
}

func TestUnverifiedUser_VerifyEmail_ShouldGenerateEmailVerifiedEvent(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	unverifiedUser, _ := NewUnverifiedUser(email, password)

	// Act
	verifiedUser := unverifiedUser.VerifyEmail()

	// Assert
	events := verifiedUser.Events()
	if len(events) != 2 {
		t.Errorf("Expected 2 events (UserRegistered + EmailVerified), got %d", len(events))
	}

	emailVerifiedEvent := events[1]
	if emailVerifiedEvent.EventName() != "EmailVerified" {
		t.Errorf("Expected EmailVerified event, got %s", emailVerifiedEvent.EventName())
	}
	if emailVerifiedEvent.AggregateID() != verifiedUser.ID() {
		t.Error("Expected event aggregate ID to match user ID")
	}
}

func TestUnverifiedUser_VerifyPassword_CorrectPassword_ShouldReturnTrue(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	user, _ := NewUnverifiedUser(email, password)

	// Act
	result := user.VerifyPassword(password)

	// Assert
	if !result {
		t.Error("Expected password verification to return true for correct password")
	}
}

func TestUnverifiedUser_VerifyPassword_IncorrectPassword_ShouldReturnFalse(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	wrongPassword, _ := NewPassword("WrongPass456#")
	user, _ := NewUnverifiedUser(email, password)

	// Act
	result := user.VerifyPassword(wrongPassword)

	// Assert
	if result {
		t.Error("Expected password verification to return false for incorrect password")
	}
}

func TestUnverifiedUser_ChangePassword_ShouldUpdatePasswordHash(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	oldPassword, _ := NewPassword("OldPass123!")
	newPassword, _ := NewPassword("NewPass456#")
	user, _ := NewUnverifiedUser(email, oldPassword)
	oldHash := user.PasswordHash()

	// Act
	err := user.ChangePassword(newPassword)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.PasswordHash() == oldHash {
		t.Error("Expected password hash to change")
	}
	if !user.VerifyPassword(newPassword) {
		t.Error("Expected new password to be verified")
	}
	if user.VerifyPassword(oldPassword) {
		t.Error("Expected old password to not be verified")
	}
}

func TestVerifiedUser_VerifiedAt_ShouldReturnVerificationTime(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	unverifiedUser, _ := NewUnverifiedUser(email, password)
	beforeVerification := time.Now()

	// Act
	verifiedUser := unverifiedUser.VerifyEmail()
	afterVerification := time.Now()

	// Assert
	verifiedAt := verifiedUser.VerifiedAt()
	if verifiedAt.Before(beforeVerification) || verifiedAt.After(afterVerification) {
		t.Error("Expected verified at time to be between before and after verification")
	}
}

func TestUser_ClearEvents_ShouldRemoveAllEvents(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	user, _ := NewUnverifiedUser(email, password)

	// Act
	user.ClearEvents()

	// Assert
	events := user.Events()
	if len(events) != 0 {
		t.Errorf("Expected 0 events after clearing, got %d", len(events))
	}
}

func TestIsUnverified_WithUnverifiedUser_ShouldReturnTrue(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	user, _ := NewUnverifiedUser(email, password)

	// Act
	unverifiedUser, ok := IsUnverified(user)

	// Assert
	if !ok {
		t.Error("Expected IsUnverified to return true for unverified user")
	}
	if unverifiedUser == nil {
		t.Error("Expected to get unverified user instance")
	}
}

func TestIsUnverified_WithVerifiedUser_ShouldReturnFalse(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	unverifiedUser, _ := NewUnverifiedUser(email, password)
	verifiedUser := unverifiedUser.VerifyEmail()

	// Act
	_, ok := IsUnverified(verifiedUser)

	// Assert
	if ok {
		t.Error("Expected IsUnverified to return false for verified user")
	}
}

func TestIsVerified_WithVerifiedUser_ShouldReturnTrue(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	unverifiedUser, _ := NewUnverifiedUser(email, password)
	verifiedUser := unverifiedUser.VerifyEmail()

	// Act
	actualVerifiedUser, ok := IsVerified(verifiedUser)

	// Assert
	if !ok {
		t.Error("Expected IsVerified to return true for verified user")
	}
	if actualVerifiedUser == nil {
		t.Error("Expected to get verified user instance")
	}
}

func TestIsVerified_WithUnverifiedUser_ShouldReturnFalse(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	password, _ := NewPassword("ValidPass123!")
	user, _ := NewUnverifiedUser(email, password)

	// Act
	_, ok := IsVerified(user)

	// Assert
	if ok {
		t.Error("Expected IsVerified to return false for unverified user")
	}
}

func TestFromSnapshot_UnverifiedUser_ShouldRestoreUnverifiedUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	passwordHash := "hashed_password"
	createdAt := time.Now()

	// Act
	user := FromSnapshot("user-123", email, passwordHash, false, createdAt, nil)

	// Assert
	if user.IsVerified() {
		t.Error("Expected restored user to be unverified")
	}
	if !user.IsUnverified() {
		t.Error("Expected restored user to be unverified")
	}
	if user.ID() != "user-123" {
		t.Error("Expected restored user to have correct ID")
	}
}

func TestFromSnapshot_VerifiedUser_ShouldRestoreVerifiedUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	passwordHash := "hashed_password"
	createdAt := time.Now()
	verifiedAt := time.Now().Add(time.Hour)

	// Act
	user := FromSnapshot("user-123", email, passwordHash, true, createdAt, &verifiedAt)

	// Assert
	if !user.IsVerified() {
		t.Error("Expected restored user to be verified")
	}
	if user.IsUnverified() {
		t.Error("Expected restored user to not be unverified")
	}

	verifiedUser, ok := IsVerified(user)
	if !ok {
		t.Error("Expected to be able to cast to VerifiedUser")
	}
	if !verifiedUser.VerifiedAt().Equal(verifiedAt) {
		t.Error("Expected restored user to have correct verified at time")
	}
}
