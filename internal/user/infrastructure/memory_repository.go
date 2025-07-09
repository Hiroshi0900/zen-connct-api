package infrastructure

import (
	"sync"
	"zen-connect/internal/user/domain"
)

// InMemoryUserRepository implements UserRepository interface using in-memory storage
type InMemoryUserRepository struct {
	users  map[string]*domain.User // userID -> User
	emails map[string]string       // email -> userID
	mutex  sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[string]*domain.User),
		emails: make(map[string]string),
	}
}

// Save stores a user in memory
func (r *InMemoryUserRepository) Save(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Store user by ID
	r.users[user.ID()] = user

	// Store email mapping
	r.emails[user.Email().String()] = user.ID()

	return nil
}

// FindByEmail finds a user by email address
func (r *InMemoryUserRepository) FindByEmail(email *domain.Email) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, exists := r.emails[email.String()]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}