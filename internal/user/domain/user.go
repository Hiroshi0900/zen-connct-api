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

	// ProvisionalUser represents a user that hasn't completed Auth0 registration
	ProvisionalUser struct {
		id        string
		email     *Email
		createdAt time.Time
		events    []DomainEvent
	}

	// ActiveUser represents a user that has completed Auth0 registration
	ActiveUser struct {
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

	// User interface for both user types
	User interface {
		ID() string
		Email() *Email
		CreatedAt() time.Time
		Events() []DomainEvent
		ClearEvents()
		IsProvisional() bool
		IsActive() bool
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

// Getters for Profile
func (p *Profile) DisplayName() string     { return p.displayName }
func (p *Profile) Bio() string            { return p.bio }
func (p *Profile) ProfileImageURL() string { return p.profileImageURL }

// UpdateDisplayName updates the display name
func (p *Profile) UpdateDisplayName(displayName string) {
	p.displayName = displayName
}

// UpdateBio updates the bio
func (p *Profile) UpdateBio(bio string) {
	p.bio = bio
}

// UpdateProfileImageURL updates the profile image URL
func (p *Profile) UpdateProfileImageURL(url string) {
	p.profileImageURL = url
}

// NewProvisionalUser creates a new provisional user
func NewProvisionalUser(email *Email) *ProvisionalUser {
	id := uuid.New().String()
	createdAt := time.Now()
	userRegisteredEvent := NewUserRegistered(id, email.Value(), createdAt)

	return &ProvisionalUser{
		id:        id,
		email:     email,
		createdAt: createdAt,
		events:    []DomainEvent{userRegisteredEvent},
	}
}

// NewActiveUser creates a new active user from Auth0 information
func NewActiveUser(auth0UserID string, email *Email, displayName string, emailVerified bool) *ActiveUser {
	id := uuid.New().String()
	now := time.Now()
	profile := NewProfile(displayName, "", "")

	var verifiedAt *time.Time
	if emailVerified {
		verifiedAt = &now
	}

	userRegisteredEvent := NewUserRegistered(id, email.Value(), now)

	return &ActiveUser{
		id:            id,
		auth0UserID:   auth0UserID,
		email:         email,
		profile:       profile,
		emailVerified: emailVerified,
		createdAt:     now,
		verifiedAt:    verifiedAt,
		updatedAt:     now,
		events:        []DomainEvent{userRegisteredEvent},
	}
}

// ActivateUser transitions a ProvisionalUser to ActiveUser with Auth0 information
func (u *ProvisionalUser) ActivateUser(auth0UserID string, displayName string, emailVerified bool) *ActiveUser {
	now := time.Now()
	profile := NewProfile(displayName, "", "")

	var verifiedAt *time.Time
	if emailVerified {
		verifiedAt = &now
	}

	emailVerifiedEvent := NewEmailVerified(u.id, u.email.Value(), now)

	return &ActiveUser{
		id:            u.id,
		auth0UserID:   auth0UserID,
		email:         u.email,
		profile:       profile,
		emailVerified: emailVerified,
		createdAt:     u.createdAt,
		verifiedAt:    verifiedAt,
		updatedAt:     now,
		events:        append(u.events, emailVerifiedEvent),
	}
}

// ProvisionalUser methods
func (u *ProvisionalUser) ID() string            { return u.id }
func (u *ProvisionalUser) Email() *Email         { return u.email }
func (u *ProvisionalUser) CreatedAt() time.Time  { return u.createdAt }
func (u *ProvisionalUser) Events() []DomainEvent { return u.events }
func (u *ProvisionalUser) IsProvisional() bool   { return true }
func (u *ProvisionalUser) IsActive() bool        { return false }

func (u *ProvisionalUser) ClearEvents() {
	u.events = []DomainEvent{}
}

// ActiveUser methods
func (a *ActiveUser) ID() string              { return a.id }
func (a *ActiveUser) Auth0UserID() string     { return a.auth0UserID }
func (a *ActiveUser) Email() *Email           { return a.email }
func (a *ActiveUser) Profile() *Profile       { return a.profile }
func (a *ActiveUser) EmailVerified() bool     { return a.emailVerified }
func (a *ActiveUser) CreatedAt() time.Time    { return a.createdAt }
func (a *ActiveUser) VerifiedAt() *time.Time  { return a.verifiedAt }
func (a *ActiveUser) UpdatedAt() time.Time    { return a.updatedAt }
func (a *ActiveUser) Events() []DomainEvent   { return a.events }
func (a *ActiveUser) IsProvisional() bool     { return false }
func (a *ActiveUser) IsActive() bool          { return true }

func (a *ActiveUser) ClearEvents() {
	a.events = []DomainEvent{}
}

// UpdateProfile updates the user's profile information
func (a *ActiveUser) UpdateProfile(displayName, bio, profileImageURL string) {
	a.profile.UpdateDisplayName(displayName)
	a.profile.UpdateBio(bio)
	a.profile.UpdateProfileImageURL(profileImageURL)
	a.updatedAt = time.Now()
}

// UpdateEmail updates the user's email
func (a *ActiveUser) UpdateEmail(email *Email) {
	a.email = email
	a.updatedAt = time.Now()
}

// VerifyEmail marks the email as verified
func (a *ActiveUser) VerifyEmail() {
	if !a.emailVerified {
		a.emailVerified = true
		now := time.Now()
		a.verifiedAt = &now
		a.updatedAt = now
		
		emailVerifiedEvent := NewEmailVerified(a.id, a.email.Value(), now)
		a.events = append(a.events, emailVerifiedEvent)
	}
}

// Type guards
func IsProvisional(user User) (*ProvisionalUser, bool) {
	u, ok := user.(*ProvisionalUser)
	return u, ok
}

func IsActive(user User) (*ActiveUser, bool) {
	a, ok := user.(*ActiveUser)
	return a, ok
}

// FromSnapshot restores a user from persistence data
func FromSnapshot(id, auth0UserID string, email *Email, displayName, bio, profileImageURL string, emailVerified bool, createdAt time.Time, verifiedAt *time.Time, updatedAt time.Time) User {
	if auth0UserID != "" {
		// Active user with Auth0 integration
		profile := NewProfile(displayName, bio, profileImageURL)
		return &ActiveUser{
			id:            id,
			auth0UserID:   auth0UserID,
			email:         email,
			profile:       profile,
			emailVerified: emailVerified,
			createdAt:     createdAt,
			verifiedAt:    verifiedAt,
			updatedAt:     updatedAt,
			events:        []DomainEvent{},
		}
	}

	// Provisional user without Auth0 integration
	return &ProvisionalUser{
		id:        id,
		email:     email,
		createdAt: createdAt,
		events:    []DomainEvent{},
	}
}
