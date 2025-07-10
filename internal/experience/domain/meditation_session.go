package domain

import (
	"errors"
	"time"
)

// MeditationSession represents a meditation session value object
type MeditationSession struct {
	startTime      time.Time
	endTime        time.Time
	meditationType string
	note           string
}

// Domain errors for MeditationSession
var (
	ErrInvalidTimeRange     = errors.New("end time must be after start time")
	ErrEmptyMeditationType = errors.New("meditation type cannot be empty")
)

// NewMeditationSession creates a new MeditationSession value object
func NewMeditationSession(startTime, endTime time.Time, meditationType, note string) *MeditationSession {
	return &MeditationSession{
		startTime:      startTime,
		endTime:        endTime,
		meditationType: meditationType,
		note:           note,
	}
}

// NewMeditationSessionWithValidation creates a new MeditationSession with validation
func NewMeditationSessionWithValidation(startTime, endTime time.Time, meditationType, note string) (*MeditationSession, error) {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return nil, ErrInvalidTimeRange
	}
	
	if meditationType == "" {
		return nil, ErrEmptyMeditationType
	}
	
	return &MeditationSession{
		startTime:      startTime,
		endTime:        endTime,
		meditationType: meditationType,
		note:           note,
	}, nil
}

// Getter methods
func (ms *MeditationSession) StartTime() time.Time {
	return ms.startTime
}

func (ms *MeditationSession) EndTime() time.Time {
	return ms.endTime
}

func (ms *MeditationSession) MeditationType() string {
	return ms.meditationType
}

func (ms *MeditationSession) Note() string {
	return ms.note
}

// Duration returns the duration of the meditation session
func (ms *MeditationSession) Duration() time.Duration {
	return ms.endTime.Sub(ms.startTime)
}

// IsValidDuration checks if the session has a valid duration
func (ms *MeditationSession) IsValidDuration() bool {
	return ms.endTime.After(ms.startTime)
}

// Equals checks if two MeditationSession instances are equal
func (ms *MeditationSession) Equals(other *MeditationSession) bool {
	if other == nil {
		return false
	}
	
	return ms.startTime.Equal(other.startTime) &&
		ms.endTime.Equal(other.endTime) &&
		ms.meditationType == other.meditationType &&
		ms.note == other.note
}