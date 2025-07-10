package domain

import (
	"testing"
)

func TestNewEmotionalState_ShouldCreateStateWhenValidInputs(t *testing.T) {
	// given
	before := "不安"
	after := "穏やか"
	
	// when
	state := NewEmotionalState(before, after)
	
	// then
	if state == nil {
		t.Error("Expected emotional state to be created")
	}
	if state.Before() != before {
		t.Error("Expected before state to match input")
	}
	if state.After() != after {
		t.Error("Expected after state to match input")
	}
}

func TestNewEmotionalState_ShouldReturnErrorWhenBeforeStateIsEmpty(t *testing.T) {
	// given
	before := ""
	after := "穏やか"
	
	// when
	state, err := NewEmotionalStateWithValidation(before, after)
	
	// then
	if err == nil {
		t.Error("Expected error for empty before state")
	}
	if state != nil {
		t.Error("Expected state to be nil when error occurs")
	}
}

func TestNewEmotionalState_ShouldReturnErrorWhenAfterStateIsEmpty(t *testing.T) {
	// given
	before := "不安"
	after := ""
	
	// when
	state, err := NewEmotionalStateWithValidation(before, after)
	
	// then
	if err == nil {
		t.Error("Expected error for empty after state")
	}
	if state != nil {
		t.Error("Expected state to be nil when error occurs")
	}
}

func TestNewEmotionalState_ShouldAllowSameBeforeAndAfterStates(t *testing.T) {
	// given
	before := "穏やか"
	after := "穏やか"
	
	// when
	state, err := NewEmotionalStateWithValidation(before, after)
	
	// then
	if err != nil {
		t.Errorf("Expected no error for same before and after states, got %v", err)
	}
	if state == nil {
		t.Error("Expected state to be created")
	}
	if state.Before() != before {
		t.Error("Expected before state to match input")
	}
	if state.After() != after {
		t.Error("Expected after state to match input")
	}
}

func TestEmotionalState_IsImproved_ShouldReturnTrueWhenImproved(t *testing.T) {
	// given
	state := NewEmotionalState("不安", "穏やか")
	
	// when
	actual := state.IsImproved()
	
	// then
	if !actual {
		t.Error("Expected IsImproved to return true for improved state")
	}
}

func TestEmotionalState_IsImproved_ShouldReturnFalseWhenNotImproved(t *testing.T) {
	// given
	state := NewEmotionalState("穏やか", "不安")
	
	// when
	actual := state.IsImproved()
	
	// then
	if actual {
		t.Error("Expected IsImproved to return false for worsened state")
	}
}

func TestEmotionalState_IsImproved_ShouldReturnFalseWhenUnchanged(t *testing.T) {
	// given
	state := NewEmotionalState("穏やか", "穏やか")
	
	// when
	actual := state.IsImproved()
	
	// then
	if actual {
		t.Error("Expected IsImproved to return false for unchanged state")
	}
}

func TestEmotionalState_HasChanged_ShouldReturnTrueWhenChanged(t *testing.T) {
	// given
	state := NewEmotionalState("不安", "穏やか")
	
	// when
	actual := state.HasChanged()
	
	// then
	if !actual {
		t.Error("Expected HasChanged to return true when states are different")
	}
}

func TestEmotionalState_HasChanged_ShouldReturnFalseWhenUnchanged(t *testing.T) {
	// given
	state := NewEmotionalState("穏やか", "穏やか")
	
	// when
	actual := state.HasChanged()
	
	// then
	if actual {
		t.Error("Expected HasChanged to return false when states are same")
	}
}

func TestEmotionalState_Equals_ShouldReturnTrueWhenSameValues(t *testing.T) {
	// given
	state1 := NewEmotionalState("不安", "穏やか")
	state2 := NewEmotionalState("不安", "穏やか")
	
	// when
	actual := state1.Equals(state2)
	
	// then
	if !actual {
		t.Error("Expected states with same values to be equal")
	}
}

func TestEmotionalState_Equals_ShouldReturnFalseWhenDifferentValues(t *testing.T) {
	// given
	state1 := NewEmotionalState("不安", "穏やか")
	state2 := NewEmotionalState("穏やか", "不安")
	
	// when
	actual := state1.Equals(state2)
	
	// then
	if actual {
		t.Error("Expected states with different values to be not equal")
	}
}

func TestEmotionalState_Equals_ShouldReturnFalseWhenComparingWithNil(t *testing.T) {
	// given
	state := NewEmotionalState("不安", "穏やか")
	
	// when
	actual := state.Equals(nil)
	
	// then
	if actual {
		t.Error("Expected state to not be equal to nil")
	}
}