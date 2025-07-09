package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Profile represents user profile information
	Profile struct {
		displayName     string
		bio            string
		profileImageURL string
	}

	// User represents an authenticated Auth0 user
	User struct {
		id            string
		auth0UserID   string
		email         *Email
		profile       *Profile
		emailVerified bool
		createdAt     time.Time
		verifiedAt    *time.Time
		updatedAt     time.Time
		events        []DomainEvent
	}
)

// NewProfile creates a new Profile value object
func NewProfile(displayName, bio, profileImageURL string) *Profile {
	return &Profile{
		displayName:     displayName,
		bio:            bio,
		profileImageURL: profileImageURL,
	}
}

// Profile getter methods
func (p *Profile) DisplayName() string     { return p.displayName }
func (p *Profile) Bio() string             { return p.bio }
func (p *Profile) ProfileImageURL() string { return p.profileImageURL }

// NewUser creates a new Auth0 authenticated user
func NewUser(auth0UserID string, email *Email, displayName string, emailVerified bool) *User {
	now := time.Now()
	user := &User{
		id:            uuid.New().String(),
		auth0UserID:   auth0UserID,
		email:         email,
		profile:       NewProfile(displayName, "", ""),
		emailVerified: emailVerified,
		createdAt:     now,
		updatedAt:     now,
		events:        []DomainEvent{},
	}

	if emailVerified {
		user.verifiedAt = &now
	}

	// Emit domain event
	user.events = append(user.events, NewUserRegistered(user.id, email.String(), now))

	return user
}

// User getter methods
func (u *User) ID() string              { return u.id }
func (u *User) Auth0UserID() string     { return u.auth0UserID }
func (u *User) Email() *Email           { return u.email }
func (u *User) Profile() *Profile       { return u.profile }
func (u *User) EmailVerified() bool     { return u.emailVerified }
func (u *User) CreatedAt() time.Time    { return u.createdAt }
func (u *User) VerifiedAt() *time.Time  { return u.verifiedAt }
func (u *User) UpdatedAt() time.Time    { return u.updatedAt }
func (u *User) Events() []DomainEvent   { return u.events }

// ClearEvents clears domain events after they have been processed
func (u *User) ClearEvents() {
	u.events = []DomainEvent{}
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(displayName, bio, profileImageURL string) {
	u.profile = NewProfile(displayName, bio, profileImageURL)
	u.updatedAt = time.Now()
	u.events = append(u.events, NewUserProfileUpdated(u.id, u.updatedAt))
}

// VerifyEmail verifies the user's email
func (u *User) VerifyEmail() {
	if !u.emailVerified {
		u.emailVerified = true
		now := time.Now()
		u.verifiedAt = &now
		u.updatedAt = now
		u.events = append(u.events, NewEmailVerified(u.id, u.email.String(), now))
	}
}

// FromSnapshot recreates a user from persisted data
func FromSnapshot(
	id string,
	auth0UserID string,
	email string,
	displayName string,
	bio string,
	profileImageURL string,
	emailVerified bool,
	createdAt time.Time,
	verifiedAt *time.Time,
	updatedAt time.Time,
) (*User, error) {
	emailObj, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:            id,
		auth0UserID:   auth0UserID,
		email:         emailObj,
		profile:       NewProfile(displayName, bio, profileImageURL),
		emailVerified: emailVerified,
		createdAt:     createdAt,
		verifiedAt:    verifiedAt,
		updatedAt:     updatedAt,
		events:        []DomainEvent{},
	}, nil
}