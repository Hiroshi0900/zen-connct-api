package domain

import "time"

// DomainEvent represents a domain event interface
type DomainEvent interface {
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}

// ExperienceCreated event fired when a new experience is created
type ExperienceCreated struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	userID      string
	createdAt   time.Time
}

func NewExperienceCreated(aggregateID, userID string, createdAt time.Time) *ExperienceCreated {
	return &ExperienceCreated{
		eventName:   "ExperienceCreated",
		aggregateID: aggregateID,
		occurredAt:  createdAt,
		userID:      userID,
		createdAt:   createdAt,
	}
}

func (e *ExperienceCreated) EventName() string     { return e.eventName }
func (e *ExperienceCreated) AggregateID() string   { return e.aggregateID }
func (e *ExperienceCreated) OccurredAt() time.Time { return e.occurredAt }
func (e *ExperienceCreated) UserID() string        { return e.userID }
func (e *ExperienceCreated) CreatedAt() time.Time  { return e.createdAt }

// ExperienceUpdated event fired when an experience is updated
type ExperienceUpdated struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	updatedAt   time.Time
}

func NewExperienceUpdated(aggregateID string, updatedAt time.Time) *ExperienceUpdated {
	return &ExperienceUpdated{
		eventName:   "ExperienceUpdated",
		aggregateID: aggregateID,
		occurredAt:  updatedAt,
		updatedAt:   updatedAt,
	}
}

func (e *ExperienceUpdated) EventName() string     { return e.eventName }
func (e *ExperienceUpdated) AggregateID() string   { return e.aggregateID }
func (e *ExperienceUpdated) OccurredAt() time.Time { return e.occurredAt }
func (e *ExperienceUpdated) UpdatedAt() time.Time  { return e.updatedAt }

// ExperienceVisibilityChanged event fired when experience visibility is changed
type ExperienceVisibilityChanged struct {
	eventName   string
	aggregateID string
	occurredAt  time.Time
	isPublic    bool
	changedAt   time.Time
}

func NewExperienceVisibilityChanged(aggregateID string, isPublic bool, changedAt time.Time) *ExperienceVisibilityChanged {
	return &ExperienceVisibilityChanged{
		eventName:   "ExperienceVisibilityChanged",
		aggregateID: aggregateID,
		occurredAt:  changedAt,
		isPublic:    isPublic,
		changedAt:   changedAt,
	}
}

func (e *ExperienceVisibilityChanged) EventName() string     { return e.eventName }
func (e *ExperienceVisibilityChanged) AggregateID() string   { return e.aggregateID }
func (e *ExperienceVisibilityChanged) OccurredAt() time.Time { return e.occurredAt }
func (e *ExperienceVisibilityChanged) IsPublic() bool        { return e.isPublic }
func (e *ExperienceVisibilityChanged) ChangedAt() time.Time  { return e.changedAt }