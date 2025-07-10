package domain

import (
	"errors"
	"strings"
)

// EmotionalState represents the emotional state before and after meditation
type EmotionalState struct {
	before string
	after  string
}

// Domain errors for EmotionalState
var (
	ErrEmptyBeforeState = errors.New("before state cannot be empty")
	ErrEmptyAfterState  = errors.New("after state cannot be empty")
)

// Positive emotional states for comparison
var positiveStates = map[string]int{
	"非常に穏やか": 5,
	"穏やか":      4,
	"やや穏やか":   3,
	"普通":       2,
	"やや不安":    1,
	"不安":       0,
	"非常に不安":   -1,
}

// NewEmotionalState creates a new EmotionalState value object
func NewEmotionalState(before, after string) *EmotionalState {
	return &EmotionalState{
		before: before,
		after:  after,
	}
}

// NewEmotionalStateWithValidation creates a new EmotionalState with validation
func NewEmotionalStateWithValidation(before, after string) (*EmotionalState, error) {
	if strings.TrimSpace(before) == "" {
		return nil, ErrEmptyBeforeState
	}
	
	if strings.TrimSpace(after) == "" {
		return nil, ErrEmptyAfterState
	}
	
	return &EmotionalState{
		before: before,
		after:  after,
	}, nil
}

// Getter methods
func (es *EmotionalState) Before() string {
	return es.before
}

func (es *EmotionalState) After() string {
	return es.after
}

// HasChanged returns true if the emotional state has changed
func (es *EmotionalState) HasChanged() bool {
	return es.before != es.after
}

// IsImproved returns true if the emotional state has improved
func (es *EmotionalState) IsImproved() bool {
	beforeScore, beforeExists := positiveStates[es.before]
	afterScore, afterExists := positiveStates[es.after]
	
	// If we don't have predefined scores, use simple string comparison
	if !beforeExists || !afterExists {
		return es.isImprovedByStringComparison()
	}
	
	return afterScore > beforeScore
}

// isImprovedByStringComparison provides a simple fallback for unknown states
func (es *EmotionalState) isImprovedByStringComparison() bool {
	// Simple heuristic: if "before" contains negative words and "after" contains positive words
	negativeWords := []string{"不安", "憂鬱", "心配", "ストレス", "イライラ", "疲れ"}
	positiveWords := []string{"穏やか", "リラックス", "落ち着", "平和", "満足", "幸せ"}
	
	beforeIsNegative := containsAny(es.before, negativeWords)
	afterIsPositive := containsAny(es.after, positiveWords)
	
	return beforeIsNegative && afterIsPositive
}

// containsAny checks if a string contains any of the given substrings
func containsAny(str string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(str, substring) {
			return true
		}
	}
	return false
}

// Equals checks if two EmotionalState instances are equal
func (es *EmotionalState) Equals(other *EmotionalState) bool {
	if other == nil {
		return false
	}
	
	return es.before == other.before && es.after == other.after
}