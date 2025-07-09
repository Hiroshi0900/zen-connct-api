package domain

import "time"

// DomainEvent represents a domain event interface
type DomainEvent interface {
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}

// UserRegistered event fired when a new user is registered
type UserRegistered struct {
	eventName    string
	aggregateID  string
	occurredAt   time.Time
	email        string
	registeredAt time.Time
}

func NewUserRegistered(aggregateID, email string, occurredAt time.Time) *UserRegistered {
	return &UserRegistered{
		eventName:    "UserRegistered",
		aggregateID:  aggregateID,
		occurredAt:   occurredAt,
		email:        email,
		registeredAt: occurredAt,
	}
}

func (e *UserRegistered) EventName() string      { return e.eventName }
func (e *UserRegistered) AggregateID() string    { return e.aggregateID }
func (e *UserRegistered) OccurredAt() time.Time  { return e.occurredAt }
func (e *UserRegistered) Email() string          { return e.email }
func (e *UserRegistered) RegisteredAt() time.Time { return e.registeredAt }

// EmailVerified event fired when user's email is verified
type EmailVerified struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	email       string
	verifiedAt  time.Time
}

func NewEmailVerified(aggregateID, email string, verifiedAt time.Time) *EmailVerified {
	return &EmailVerified{
		eventName:   "EmailVerified",
		aggregateID: aggregateID,
		occurredAt:  verifiedAt,
		email:       email,
		verifiedAt:  verifiedAt,
	}
}

func (e *EmailVerified) EventName() string     { return e.eventName }
func (e *EmailVerified) AggregateID() string   { return e.aggregateID }
func (e *EmailVerified) OccurredAt() time.Time { return e.occurredAt }
func (e *EmailVerified) Email() string         { return e.email }
func (e *EmailVerified) VerifiedAt() time.Time { return e.verifiedAt }

// PasswordChanged event fired when user's password is changed
type PasswordChanged struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	changedAt   time.Time
}

func NewPasswordChanged(aggregateID string, changedAt time.Time) *PasswordChanged {
	return &PasswordChanged{
		eventName:   "PasswordChanged",
		aggregateID: aggregateID,
		occurredAt:  changedAt,
		changedAt:   changedAt,
	}
}

func (e *PasswordChanged) EventName() string     { return e.eventName }
func (e *PasswordChanged) AggregateID() string   { return e.aggregateID }
func (e *PasswordChanged) OccurredAt() time.Time { return e.occurredAt }
func (e *PasswordChanged) ChangedAt() time.Time  { return e.changedAt }

// UserProfileUpdated event fired when user's profile is updated
type UserProfileUpdated struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	updatedAt   time.Time
}

func NewUserProfileUpdated(aggregateID string, updatedAt time.Time) *UserProfileUpdated {
	return &UserProfileUpdated{
		eventName:   "UserProfileUpdated",
		aggregateID: aggregateID,
		occurredAt:  updatedAt,
		updatedAt:   updatedAt,
	}
}

func (e *UserProfileUpdated) EventName() string     { return e.eventName }
func (e *UserProfileUpdated) AggregateID() string   { return e.aggregateID }
func (e *UserProfileUpdated) OccurredAt() time.Time { return e.occurredAt }
func (e *UserProfileUpdated) UpdatedAt() time.Time  { return e.updatedAt }