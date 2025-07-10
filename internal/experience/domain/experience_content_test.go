package domain

import (
	"testing"
	"time"
)

func TestNewExperienceContent_ShouldCreateContentWhenValidInputs(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt
	
	// when
	content := NewExperienceContent(session, emotionalState, createdAt, updatedAt)
	
	// then
	if content == nil {
		t.Error("Expected content to be created")
	}
	if !content.Session().Equals(session) {
		t.Error("Expected session to match input")
	}
	if !content.EmotionalState().Equals(emotionalState) {
		t.Error("Expected emotional state to match input")
	}
	if content.CreatedAt() != createdAt {
		t.Error("Expected created at to match input")
	}
	if content.UpdatedAt() != updatedAt {
		t.Error("Expected updated at to match input")
	}
}

func TestNewExperienceContent_ShouldReturnErrorWhenSessionIsNil(t *testing.T) {
	// given
	var session *MeditationSession = nil
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt
	
	// when
	content, err := NewExperienceContentWithValidation(session, emotionalState, createdAt, updatedAt)
	
	// then
	if err == nil {
		t.Error("Expected error for nil session")
	}
	if content != nil {
		t.Error("Expected content to be nil when error occurs")
	}
}

func TestNewExperienceContent_ShouldReturnErrorWhenEmotionalStateIsNil(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	var emotionalState *EmotionalState = nil
	createdAt := time.Now()
	updatedAt := createdAt
	
	// when
	content, err := NewExperienceContentWithValidation(session, emotionalState, createdAt, updatedAt)
	
	// then
	if err == nil {
		t.Error("Expected error for nil emotional state")
	}
	if content != nil {
		t.Error("Expected content to be nil when error occurs")
	}
}

func TestNewExperienceContent_ShouldReturnErrorWhenUpdatedAtBeforeCreatedAt(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt.Add(-10 * time.Minute)
	
	// when
	content, err := NewExperienceContentWithValidation(session, emotionalState, createdAt, updatedAt)
	
	// then
	if err == nil {
		t.Error("Expected error for invalid timestamp order")
	}
	if content != nil {
		t.Error("Expected content to be nil when error occurs")
	}
}

func TestExperienceContent_UpdateTimestamp_ShouldUpdateTimestamp(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	content := NewExperienceContent(session, emotionalState, createdAt, createdAt)
	newTimestamp := createdAt.Add(10 * time.Minute)
	
	// when
	updatedContent := content.UpdateTimestamp(newTimestamp)
	
	// then
	if updatedContent == nil {
		t.Error("Expected updated content to be created")
	}
	if updatedContent.UpdatedAt() != newTimestamp {
		t.Error("Expected updated at to be updated")
	}
	if updatedContent.CreatedAt() != createdAt {
		t.Error("Expected created at to remain unchanged")
	}
}

func TestExperienceContent_HasEmotionalImprovement_ShouldReturnTrueWhenImproved(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	content := NewExperienceContent(session, emotionalState, time.Now(), time.Now())
	
	// when
	actual := content.HasEmotionalImprovement()
	
	// then
	if !actual {
		t.Error("Expected HasEmotionalImprovement to return true when improved")
	}
}

func TestExperienceContent_HasEmotionalImprovement_ShouldReturnFalseWhenNotImproved(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("穏やか", "穏やか")
	content := NewExperienceContent(session, emotionalState, time.Now(), time.Now())
	
	// when
	actual := content.HasEmotionalImprovement()
	
	// then
	if actual {
		t.Error("Expected HasEmotionalImprovement to return false when not improved")
	}
}

func TestExperienceContent_IsRecent_ShouldReturnTrueWhenRecent(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	recentTime := time.Now().Add(-1 * time.Hour)
	content := NewExperienceContent(session, emotionalState, recentTime, recentTime)
	
	// when
	actual := content.IsRecent(time.Now(), 24*time.Hour)
	
	// then
	if !actual {
		t.Error("Expected IsRecent to return true for recent content")
	}
}

func TestExperienceContent_IsRecent_ShouldReturnFalseWhenNotRecent(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	oldTime := time.Now().Add(-25 * time.Hour)
	content := NewExperienceContent(session, emotionalState, oldTime, oldTime)
	
	// when
	actual := content.IsRecent(time.Now(), 24*time.Hour)
	
	// then
	if actual {
		t.Error("Expected IsRecent to return false for old content")
	}
}

func TestExperienceContent_Equals_ShouldReturnTrueWhenSameValues(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt
	content1 := NewExperienceContent(session, emotionalState, createdAt, updatedAt)
	content2 := NewExperienceContent(session, emotionalState, createdAt, updatedAt)
	
	// when
	actual := content1.Equals(content2)
	
	// then
	if !actual {
		t.Error("Expected contents with same values to be equal")
	}
}

func TestExperienceContent_Equals_ShouldReturnFalseWhenDifferentValues(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session1 := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	session2 := NewMeditationSession(startTime, endTime, "breathing", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	createdAt := time.Now()
	updatedAt := createdAt
	content1 := NewExperienceContent(session1, emotionalState, createdAt, updatedAt)
	content2 := NewExperienceContent(session2, emotionalState, createdAt, updatedAt)
	
	// when
	actual := content1.Equals(content2)
	
	// then
	if actual {
		t.Error("Expected contents with different values to be not equal")
	}
}

func TestExperienceContent_Equals_ShouldReturnFalseWhenComparingWithNil(t *testing.T) {
	// given
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)
	session := NewMeditationSession(startTime, endTime, "mindfulness", "集中できた")
	emotionalState := NewEmotionalState("不安", "穏やか")
	content := NewExperienceContent(session, emotionalState, time.Now(), time.Now())
	
	// when
	actual := content.Equals(nil)
	
	// then
	if actual {
		t.Error("Expected content to not be equal to nil")
	}
}