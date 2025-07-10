package domain

import (
	"errors"
	"time"
)

// ExperienceContent represents the content of an experience record
type ExperienceContent struct {
	session        *MeditationSession
	emotionalState *EmotionalState
	createdAt      time.Time
	updatedAt      time.Time
}

// Domain errors for ExperienceContent
var (
	ErrNilSession        = errors.New("session cannot be nil")
	ErrNilEmotionalState = errors.New("emotional state cannot be nil")
	ErrInvalidTimestamp  = errors.New("updated at must be equal to or after created at")
)

// NewExperienceContent creates a new ExperienceContent value object
func NewExperienceContent(session *MeditationSession, emotionalState *EmotionalState, createdAt, updatedAt time.Time) *ExperienceContent {
	return &ExperienceContent{
		session:        session,
		emotionalState: emotionalState,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// NewExperienceContentWithValidation creates a new ExperienceContent with validation
func NewExperienceContentWithValidation(session *MeditationSession, emotionalState *EmotionalState, createdAt, updatedAt time.Time) (*ExperienceContent, error) {
	if session == nil {
		return nil, ErrNilSession
	}
	
	if emotionalState == nil {
		return nil, ErrNilEmotionalState
	}
	
	if updatedAt.Before(createdAt) {
		return nil, ErrInvalidTimestamp
	}
	
	return &ExperienceContent{
		session:        session,
		emotionalState: emotionalState,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}, nil
}

// Getter methods
func (ec *ExperienceContent) Session() *MeditationSession {
	return ec.session
}

func (ec *ExperienceContent) EmotionalState() *EmotionalState {
	return ec.emotionalState
}

func (ec *ExperienceContent) CreatedAt() time.Time {
	return ec.createdAt
}

func (ec *ExperienceContent) UpdatedAt() time.Time {
	return ec.updatedAt
}

// UpdateTimestamp creates a new ExperienceContent with updated timestamp
func (ec *ExperienceContent) UpdateTimestamp(newTimestamp time.Time) *ExperienceContent {
	return &ExperienceContent{
		session:        ec.session,
		emotionalState: ec.emotionalState,
		createdAt:      ec.createdAt,
		updatedAt:      newTimestamp,
	}
}

// HasEmotionalImprovement returns true if the emotional state has improved
func (ec *ExperienceContent) HasEmotionalImprovement() bool {
	return ec.emotionalState.IsImproved()
}

// IsRecent returns true if the content was created within the specified duration from the reference time
func (ec *ExperienceContent) IsRecent(referenceTime time.Time, duration time.Duration) bool {
	return referenceTime.Sub(ec.createdAt) <= duration
}

// Equals checks if two ExperienceContent instances are equal
func (ec *ExperienceContent) Equals(other *ExperienceContent) bool {
	if other == nil {
		return false
	}
	
	return ec.session.Equals(other.session) &&
		ec.emotionalState.Equals(other.emotionalState) &&
		ec.createdAt.Equal(other.createdAt) &&
		ec.updatedAt.Equal(other.updatedAt)
}