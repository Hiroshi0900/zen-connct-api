package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Experience represents an experience record entity (aggregate root)
type Experience struct {
	id        string
	userID    string
	content   *ExperienceContent
	isPublic  bool
	createdAt time.Time
	updatedAt time.Time
	events    []DomainEvent
}

// Domain errors for Experience
var (
	ErrEmptyUserID   = errors.New("user ID cannot be empty")
	ErrNilContent    = errors.New("content cannot be nil")
)

// NewExperience creates a new Experience entity
func NewExperience(userID string, content *ExperienceContent) *Experience {
	now := time.Now()
	experience := &Experience{
		id:        uuid.New().String(),
		userID:    userID,
		content:   content,
		isPublic:  false, // Default to private
		createdAt: now,
		updatedAt: now,
		events:    []DomainEvent{},
	}
	
	// Emit domain event
	experience.events = append(experience.events, NewExperienceCreated(experience.id, userID, now))
	
	return experience
}

// NewExperienceWithValidation creates a new Experience with validation
func NewExperienceWithValidation(userID string, content *ExperienceContent) (*Experience, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, ErrEmptyUserID
	}
	
	if content == nil {
		return nil, ErrNilContent
	}
	
	now := time.Now()
	experience := &Experience{
		id:        uuid.New().String(),
		userID:    userID,
		content:   content,
		isPublic:  false, // Default to private
		createdAt: now,
		updatedAt: now,
		events:    []DomainEvent{},
	}
	
	// Emit domain event
	experience.events = append(experience.events, NewExperienceCreated(experience.id, userID, now))
	
	return experience, nil
}

// Getter methods
func (e *Experience) ID() string {
	return e.id
}

func (e *Experience) UserID() string {
	return e.userID
}

func (e *Experience) Content() *ExperienceContent {
	return e.content
}

func (e *Experience) IsPublic() bool {
	return e.isPublic
}

func (e *Experience) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Experience) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Experience) Events() []DomainEvent {
	return e.events
}

// ClearEvents clears domain events after they have been processed
func (e *Experience) ClearEvents() {
	e.events = []DomainEvent{}
}

// MakePublic makes the experience public
func (e *Experience) MakePublic() {
	if !e.isPublic {
		e.isPublic = true
		e.updatedAt = time.Now()
		e.events = append(e.events, NewExperienceVisibilityChanged(e.id, true, e.updatedAt))
	}
}

// MakePrivate makes the experience private
func (e *Experience) MakePrivate() {
	if e.isPublic {
		e.isPublic = false
		e.updatedAt = time.Now()
		e.events = append(e.events, NewExperienceVisibilityChanged(e.id, false, e.updatedAt))
	}
}

// UpdateContent updates the experience content
func (e *Experience) UpdateContent(content *ExperienceContent) {
	if content != nil {
		e.content = content
		e.updatedAt = time.Now()
		e.events = append(e.events, NewExperienceUpdated(e.id, e.updatedAt))
	}
}

// BelongsToUser checks if the experience belongs to the specified user
func (e *Experience) BelongsToUser(userID string) bool {
	return e.userID == userID
}

// FromSnapshot recreates an experience from persisted data
func FromSnapshot(
	id string,
	userID string,
	content *ExperienceContent,
	isPublic bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Experience {
	return &Experience{
		id:        id,
		userID:    userID,
		content:   content,
		isPublic:  isPublic,
		createdAt: createdAt,
		updatedAt: updatedAt,
		events:    []DomainEvent{},
	}
}