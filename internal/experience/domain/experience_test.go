package domain

import (
	"testing"
	"time"
)

func TestNewExperience_ShouldCreateExperienceWhenValidInputs(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	
	// when
	experience := NewExperience(userID, content)
	
	// then
	if experience == nil {
		t.Error("Expected experience to be created")
	}
	if experience.ID() == "" {
		t.Error("Expected experience ID to be generated")
	}
	if experience.UserID() != userID {
		t.Error("Expected user ID to match input")
	}
	if !experience.Content().Equals(content) {
		t.Error("Expected content to match input")
	}
	if experience.IsPublic() {
		t.Error("Expected new experience to be private by default")
	}
}

func TestNewExperience_ShouldGenerateExperienceCreatedEvent(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	
	// when
	experience := NewExperience(userID, content)
	
	// then
	events := experience.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	
	event := events[0]
	if event.EventName() != "ExperienceCreated" {
		t.Errorf("Expected event name 'ExperienceCreated', got '%s'", event.EventName())
	}
	if event.AggregateID() != experience.ID() {
		t.Error("Expected event aggregate ID to match experience ID")
	}
}

func TestNewExperience_ShouldReturnErrorWhenUserIDIsEmpty(t *testing.T) {
	// given
	userID := ""
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	
	// when
	experience, err := NewExperienceWithValidation(userID, content)
	
	// then
	if err == nil {
		t.Error("Expected error for empty user ID")
	}
	if experience != nil {
		t.Error("Expected experience to be nil when error occurs")
	}
}

func TestNewExperience_ShouldReturnErrorWhenContentIsNil(t *testing.T) {
	// given
	userID := "user-123"
	var content *ExperienceContent = nil
	
	// when
	experience, err := NewExperienceWithValidation(userID, content)
	
	// then
	if err == nil {
		t.Error("Expected error for nil content")
	}
	if experience != nil {
		t.Error("Expected experience to be nil when error occurs")
	}
}

func TestExperience_MakePublic_ShouldMakeExperiencePublicAndGenerateEvent(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	experience.ClearEvents()
	
	// when
	experience.MakePublic()
	
	// then
	if !experience.IsPublic() {
		t.Error("Expected experience to be public")
	}
	
	events := experience.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	
	event := events[0]
	if event.EventName() != "ExperienceVisibilityChanged" {
		t.Errorf("Expected event name 'ExperienceVisibilityChanged', got '%s'", event.EventName())
	}
}

func TestExperience_MakePrivate_ShouldMakeExperiencePrivateAndGenerateEvent(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	experience.MakePublic()
	experience.ClearEvents()
	
	// when
	experience.MakePrivate()
	
	// then
	if experience.IsPublic() {
		t.Error("Expected experience to be private")
	}
	
	events := experience.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	
	event := events[0]
	if event.EventName() != "ExperienceVisibilityChanged" {
		t.Errorf("Expected event name 'ExperienceVisibilityChanged', got '%s'", event.EventName())
	}
}

func TestExperience_UpdateContent_ShouldUpdateContentAndGenerateEvent(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	experience.ClearEvents()
	
	newSession := NewMeditationSession(startTime, endTime, "breathing", "呼吸に集中")
	newEmotionalState := NewEmotionalState("ストレス", "リラックス")
	newContent := NewExperienceContent(newSession, newEmotionalState, createdAt, time.Now())
	
	// when
	experience.UpdateContent(newContent)
	
	// then
	if !experience.Content().Equals(newContent) {
		t.Error("Expected content to be updated")
	}
	
	events := experience.Events()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}
	
	event := events[0]
	if event.EventName() != "ExperienceUpdated" {
		t.Errorf("Expected event name 'ExperienceUpdated', got '%s'", event.EventName())
	}
}

func TestExperience_ClearEvents_ShouldClearAllEvents(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	
	// when
	experience.ClearEvents()
	
	// then
	events := experience.Events()
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}

func TestExperience_BelongsToUser_ShouldReturnTrueWhenUserMatches(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	
	// when
	actual := experience.BelongsToUser(userID)
	
	// then
	if !actual {
		t.Error("Expected BelongsToUser to return true when user matches")
	}
}

func TestExperience_BelongsToUser_ShouldReturnFalseWhenUserDoesNotMatch(t *testing.T) {
	// given
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	experience := NewExperience(userID, content)
	
	// when
	actual := experience.BelongsToUser("different-user")
	
	// then
	if actual {
		t.Error("Expected BelongsToUser to return false when user does not match")
	}
}

func TestFromSnapshot_ShouldCreateExperienceFromSnapshot(t *testing.T) {
	// given
	id := "experience-123"
	userID := "user-123"
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt
	content := NewExperienceContent(session, emotionalState, createdAt, updatedAt)
	isPublic := true
	
	// when
	experience := FromSnapshot(id, userID, content, isPublic, createdAt, updatedAt)
	
	// then
	if experience == nil {
		t.Error("Expected experience to be created from snapshot")
	}
	if experience.ID() != id {
		t.Error("Expected experience ID to match snapshot")
	}
	if experience.UserID() != userID {
		t.Error("Expected user ID to match snapshot")
	}
	if !experience.Content().Equals(content) {
		t.Error("Expected content to match snapshot")
	}
	if experience.IsPublic() != isPublic {
		t.Error("Expected public status to match snapshot")
	}
	if experience.CreatedAt() != createdAt {
		t.Error("Expected created at to match snapshot")
	}
	if experience.UpdatedAt() != updatedAt {
		t.Error("Expected updated at to match snapshot")
	}
	if len(experience.Events()) != 0 {
		t.Error("Expected no events for experience created from snapshot")
	}
}