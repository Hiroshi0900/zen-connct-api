package domain

import (
	"time"

	"github.com/google/uuid"
)

type (
	// baseUser はユーザーの基本的な情報を表すインターフェース
	baseUser struct {
		id           string
		email        *Email
		passwordHash string
	}
	// UnverifiedUser は表す未検証のユーザー
	UnverifiedUser struct {
		baseUser
		createdAt time.Time
		events    []DomainEvent
	}

	// VerifiedUser は表す検証済みのユーザー
	VerifiedUser struct {
		baseUser
		createdAt  time.Time
		verifiedAt time.Time
		events     []DomainEvent
	}

	// User はユーザーのインターフェース
	User interface {
		ID() string
		Email() *Email
		PasswordHash() string
		CreatedAt() time.Time
		Events() []DomainEvent
		ClearEvents()
		VerifyPassword(password *Password) bool
		ChangePassword(newPassword *Password) error
		IsUnverified() bool
		IsVerified() bool
	}
)

func (u baseUser) ID() string           { return u.id }
func (u baseUser) Email() *Email        { return u.email }
func (u baseUser) PasswordHash() string { return u.passwordHash }

// NewUnverifiedUser creates a new unverified user
func NewUnverifiedUser(email *Email, password *Password) (*UnverifiedUser, error) {
	id := uuid.New().String()
	passwordHash, err := password.Hash()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	userRegisteredEvent := NewUserRegistered(id, email.Value(), createdAt)

	return &UnverifiedUser{
		baseUser: baseUser{
			id:           id,
			email:        email,
			passwordHash: passwordHash,
		},
		createdAt: createdAt,
		events:    []DomainEvent{userRegisteredEvent},
	}, nil
}

// VerifyEmail transitions an UnverifiedUser to VerifiedUser
func (u *UnverifiedUser) VerifyEmail() *VerifiedUser {
	verifiedAt := time.Now()
	emailVerifiedEvent := NewEmailVerified(u.id, u.email.Value(), verifiedAt)

	return &VerifiedUser{
		baseUser: baseUser{
			id:           u.id,
			email:        u.email,
			passwordHash: u.passwordHash,
		},
		createdAt:  u.createdAt,
		verifiedAt: verifiedAt,
		events:     append(u.events, emailVerifiedEvent),
	}
}

// UnverifiedUser methods
func (u *UnverifiedUser) CreatedAt() time.Time  { return u.createdAt }
func (u *UnverifiedUser) Events() []DomainEvent { return u.events }
func (u *UnverifiedUser) IsUnverified() bool    { return true }
func (u *UnverifiedUser) IsVerified() bool      { return false }

func (u *UnverifiedUser) ClearEvents() {
	u.events = []DomainEvent{}
}

func (u *UnverifiedUser) VerifyPassword(password *Password) bool {
	return VerifyHash(password.Value(), u.passwordHash)
}

func (u *UnverifiedUser) ChangePassword(newPassword *Password) error {
	newPasswordHash, err := newPassword.Hash()
	if err != nil {
		return err
	}
	u.passwordHash = newPasswordHash

	changedAt := time.Now()
	passwordChangedEvent := NewPasswordChanged(u.id, changedAt)
	u.events = append(u.events, passwordChangedEvent)

	return nil
}

// VerifiedUser methods
func (v *VerifiedUser) CreatedAt() time.Time  { return v.createdAt }
func (v *VerifiedUser) VerifiedAt() time.Time { return v.verifiedAt }
func (v *VerifiedUser) Events() []DomainEvent { return v.events }
func (v *VerifiedUser) IsUnverified() bool    { return false }
func (v *VerifiedUser) IsVerified() bool      { return true }

func (v *VerifiedUser) ClearEvents() {
	v.events = []DomainEvent{}
}

func (v *VerifiedUser) VerifyPassword(password *Password) bool {
	return VerifyHash(password.Value(), v.passwordHash)
}

func (v *VerifiedUser) ChangePassword(newPassword *Password) error {
	newPasswordHash, err := newPassword.Hash()
	if err != nil {
		return err
	}
	v.passwordHash = newPasswordHash

	changedAt := time.Now()
	passwordChangedEvent := NewPasswordChanged(v.id, changedAt)
	v.events = append(v.events, passwordChangedEvent)

	return nil
}

// Type guards
func IsUnverified(user User) (*UnverifiedUser, bool) {
	u, ok := user.(*UnverifiedUser)
	return u, ok
}

func IsVerified(user User) (*VerifiedUser, bool) {
	v, ok := user.(*VerifiedUser)
	return v, ok
}

// FromSnapshot restores a user from persistence data
func FromSnapshot(id string, email *Email, passwordHash string, emailVerified bool, createdAt time.Time, verifiedAt *time.Time) User {
	if emailVerified && verifiedAt != nil {
		return &VerifiedUser{
			baseUser: baseUser{
				id:           id,
				email:        email,
				passwordHash: passwordHash,
			},
			createdAt:  createdAt,
			verifiedAt: *verifiedAt,
			events:     []DomainEvent{},
		}
	}

	return &UnverifiedUser{
		baseUser: baseUser{
			id:           id,
			email:        email,
			passwordHash: passwordHash,
		},
		createdAt: createdAt,
		events:    []DomainEvent{},
	}
}
