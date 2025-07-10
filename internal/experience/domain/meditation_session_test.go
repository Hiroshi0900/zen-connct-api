package domain

import (
	"testing"
	"time"
)

func TestNewMeditationSession_ShouldCreateSessionWhenValidInputs(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	meditationType := "mindfulness"
	note := "感じたことを記録"
	
	// when
	session := NewMeditationSession(startTime, endTime, meditationType, note)
	
	// then
	if session == nil {
		t.Error("Expected session to be created")
	}
	if session.StartTime() != startTime {
		t.Error("Expected start time to match input")
	}
	if session.EndTime() != endTime {
		t.Error("Expected end time to match input")
	}
	if session.MeditationType() != meditationType {
		t.Error("Expected meditation type to match input")
	}
	if session.Note() != note {
		t.Error("Expected note to match input")
	}
}

func TestNewMeditationSession_ShouldReturnErrorWhenEndTimeBeforeStartTime(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(-10 * time.Minute)
	meditationType := "mindfulness"
	note := "test note"
	
	// when
	session, err := NewMeditationSessionWithValidation(startTime, endTime, meditationType, note)
	
	// then
	if err == nil {
		t.Error("Expected error for invalid time range")
	}
	if session != nil {
		t.Error("Expected session to be nil when error occurs")
	}
}

func TestNewMeditationSession_ShouldReturnErrorWhenMeditationTypeIsEmpty(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	meditationType := ""
	note := "test note"
	
	// when
	session, err := NewMeditationSessionWithValidation(startTime, endTime, meditationType, note)
	
	// then
	if err == nil {
		t.Error("Expected error for empty meditation type")
	}
	if session != nil {
		t.Error("Expected session to be nil when error occurs")
	}
}

func TestMeditationSession_Duration_ShouldReturnCorrectDuration(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(45 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	expected := 45 * time.Minute
	
	// when
	actual := session.Duration()
	
	// then
	if actual != expected {
		t.Errorf("Expected duration %v, got %v", expected, actual)
	}
}

func TestMeditationSession_IsValidDuration_ShouldReturnTrueForValidDuration(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	
	// when
	actual := session.IsValidDuration()
	
	// then
	if !actual {
		t.Error("Expected valid duration to return true")
	}
}

func TestMeditationSession_IsValidDuration_ShouldReturnFalseForNegativeDuration(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(-10 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	
	// when
	actual := session.IsValidDuration()
	
	// then
	if actual {
		t.Error("Expected invalid duration to return false")
	}
}

func TestMeditationSession_Equals_ShouldReturnTrueWhenSameValues(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session1 := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	session2 := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	
	// when
	actual := session1.Equals(session2)
	
	// then
	if !actual {
		t.Error("Expected sessions with same values to be equal")
	}
}

func TestMeditationSession_Equals_ShouldReturnFalseWhenDifferentValues(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session1 := NewMeditationSession(startTime, endTime, "mindfulness", "test")
	session2 := NewMeditationSession(startTime, endTime, "breathing", "test")
	
	// when
	actual := session1.Equals(session2)
	
	// then
	if actual {
		t.Error("Expected sessions with different values to be not equal")
	}
}