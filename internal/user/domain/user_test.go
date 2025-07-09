package domain

import (
	"testing"
	"time"
)

// TODO: 以下以降のテストは全部見直す必要がある。
func TestNewUser_ValidInputs_ShouldCreateUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")

	// Act
	user := NewUser("auth0-user-id", email, "Test User", false)

	// Assert
	if user == nil {
		t.Error("Expected user to be created")
	}
	if user.ID() == "" {
		t.Error("Expected user ID to be generated")
	}
	if user.Email().String() != email.String() {
		t.Error("Expected user email to match input email")
	}
	if user.EmailVerified() {
		t.Error("Expected new user to be unverified")
	}
}

func TestNewUser_ShouldGenerateUserRegisteredEvent(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")

	// Act
	user := NewUser("auth0-user-id", email, "Test User", false)

	// Assert
	events := user.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventName() != "UserRegistered" {
		t.Errorf("Expected event name 'UserRegistered', got '%s'", event.EventName())
	}
	if event.AggregateID() != user.ID() {
		t.Errorf("Expected event aggregate ID to match user ID")
	}
}

func TestVerifyEmail_ShouldUpdateUserAndGenerateEvent(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	user := NewUser("auth0-user-id", email, "Test User", false)
	user.ClearEvents() // Clear the initial registration event

	// Act
	user.VerifyEmail()

	// Assert
	if !user.EmailVerified() {
		t.Error("Expected user to be verified")
	}
	if user.VerifiedAt() == nil {
		t.Error("Expected verified at timestamp to be set")
	}

	events := user.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventName() != "EmailVerified" {
		t.Errorf("Expected event name 'EmailVerified', got '%s'", event.EventName())
	}
}

func TestVerifyEmail_AlreadyVerified_ShouldNotChangeUser(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	user := NewUser("auth0-user-id", email, "Test User", true)
	user.ClearEvents() // Clear the initial registration event
	originalVerifiedAt := user.VerifiedAt()

	// Act
	user.VerifyEmail()

	// Assert
	if user.VerifiedAt() != originalVerifiedAt {
		t.Error("Expected verified at timestamp to remain unchanged")
	}

	events := user.Events()
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

func TestUpdateProfile_ShouldUpdateUserAndGenerateEvent(t *testing.T) {
	// Arrange
	email, _ := NewEmail("test@example.com")
	user := NewUser("auth0-user-id", email, "Test User", false)
	user.ClearEvents() // Clear the initial registration event

	// Act
	user.UpdateProfile("New Name", "New bio", "new-image.jpg")

	// Assert
	if user.Profile().DisplayName() != "New Name" {
		t.Error("Expected display name to be updated")
	}
	if user.Profile().Bio() != "New bio" {
		t.Error("Expected bio to be updated")
	}
	if user.Profile().ProfileImageURL() != "new-image.jpg" {
		t.Error("Expected profile image URL to be updated")
	}

	events := user.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event := events[0]
	if event.EventName() != "UserProfileUpdated" {
		t.Errorf("Expected event name 'UserProfileUpdated', got '%s'", event.EventName())
	}
}

func TestFromSnapshot_ValidInputs_ShouldCreateUser(t *testing.T) {
	// Arrange
	id := "user-123"
	auth0UserID := "auth0-user-id"
	email := "test@example.com"
	displayName := "Test User"
	bio := "Test bio"
	profileImageURL := "test-image.jpg"
	emailVerified := true
	createdAt := time.Now()
	verifiedAt := &createdAt
	updatedAt := time.Now()

	// Act
	user, err := FromSnapshot(id, auth0UserID, email, displayName, bio, profileImageURL, emailVerified, createdAt, verifiedAt, updatedAt)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.ID() != id {
		t.Error("Expected user ID to match snapshot")
	}
	if user.Auth0UserID() != auth0UserID {
		t.Error("Expected Auth0 user ID to match snapshot")
	}
	if user.Email().String() != email {
		t.Error("Expected email to match snapshot")
	}
	if user.Profile().DisplayName() != displayName {
		t.Error("Expected display name to match snapshot")
	}
	if user.EmailVerified() != emailVerified {
		t.Error("Expected email verified status to match snapshot")
	}
}

func TestFromSnapshot_InvalidEmail_ShouldReturnError(t *testing.T) {
	// Arrange
	invalidEmail := "invalid-email"

	// Act
	_, err := FromSnapshot("id", "auth0-id", invalidEmail, "name", "bio", "image", false, time.Now(), nil, time.Now())

	// Assert
	if err == nil {
		t.Error("Expected error for invalid email")
	}
}