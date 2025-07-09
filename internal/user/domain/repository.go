package domain

import "errors"

// Domain errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	Save(user *User) error
	FindByEmail(email *Email) (*User, error)
	FindByID(id string) (*User, error)
	FindByAuth0UserID(auth0UserID string) (*User, error)
}