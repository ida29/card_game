package game

import (
	"fmt"
	"github.com/your-org/card_game/backend/internal/effects"
	"github.com/your-org/card_game/backend/internal/models"
)

// TargetValidator validates and filters targets for effects
type TargetValidator struct {
	context *effects.GameContext
}

// NewTargetValidator creates a new target validator
func NewTargetValidator(context *effects.GameContext) *TargetValidator {
	return &TargetValidator{
		context: context,
	}
}

// ValidateTargets checks if the selected targets are valid for an effect
func (v *TargetValidator) ValidateTargets(effect effects.Effect, source *models.Card, selectedTargets []effects.Target) error {
	// Get all valid targets
	validTargets := effect.GetTargets(v.context, source)
	
	// Create a map for quick lookup
	validMap := make(map[string]bool)
	for _, target := range validTargets {
		key := fmt.Sprintf("%s_%s_%s", target.Type, target.ID, target.Location)
		validMap[key] = true
	}
	
	// Check each selected target
	for _, selected := range selectedTargets {
		key := fmt.Sprintf("%s_%s_%s", selected.Type, selected.ID, selected.Location)
		if !validMap[key] {
			return fmt.Errorf("invalid target: %s", key)
		}
	}
	
	return nil
}

// FilterTargetsByType filters targets by their type
func (v *TargetValidator) FilterTargetsByType(targets []effects.Target, targetType string) []effects.Target {
	filtered := make([]effects.Target, 0)
	for _, target := range targets {
		if target.Type == targetType {
			filtered = append(filtered, target)
		}
	}
	return filtered
}

// FilterTargetsByLocation filters targets by their location
func (v *TargetValidator) FilterTargetsByLocation(targets []effects.Target, location string) []effects.Target {
	filtered := make([]effects.Target, 0)
	for _, target := range targets {
		if target.Location == location {
			filtered = append(filtered, target)
		}
	}
	return filtered
}

// FilterFriendsByPower filters friend targets by power threshold
func (v *TargetValidator) FilterFriendsByPower(targets []effects.Target, maxPower int) []effects.Target {
	filtered := make([]effects.Target, 0)
	for _, target := range targets {
		if target.Type == "friend" {
			if friend, ok := target.Data.(models.Friend); ok {
				if friend.Power <= maxPower {
					filtered = append(filtered, target)
				}
			}
		}
	}
	return filtered
}

// FilterFriendsByCost filters friend targets by cost threshold
func (v *TargetValidator) FilterFriendsByCost(targets []effects.Target, maxCost int) []effects.Target {
	filtered := make([]effects.Target, 0)
	for _, target := range targets {
		if target.Type == "friend" {
			// TODO: Load card data to check cost
			filtered = append(filtered, target)
		}
	}
	return filtered
}

// GetTargetRequirements analyzes an effect to determine targeting requirements
func (v *TargetValidator) GetTargetRequirements(effect effects.Effect) TargetRequirements {
	// Analyze the effect type to determine requirements
	// This is a simplified version - in practice, effects should provide this info
	
	switch effect.(type) {
	case *effects.DestroyFriendEffect:
		return TargetRequirements{
			MinTargets:  1,
			MaxTargets:  1,
			TargetTypes: []string{"friend"},
			Mandatory:   true,
			Description: "相手のふれんど1体を選択",
		}
		
	case *effects.ReturnToHandEffect:
		return TargetRequirements{
			MinTargets:  1,
			MaxTargets:  1,
			TargetTypes: []string{"friend"},
			Mandatory:   true,
			Description: "ふれんど1体を選択",
		}
		
	case *effects.PowerModifierEffect:
		return TargetRequirements{
			MinTargets:  0,
			MaxTargets:  -1, // -1 means all valid targets
			TargetTypes: []string{"friend"},
			Mandatory:   false,
			Description: "対象のふれんど",
		}
		
	default:
		return TargetRequirements{
			MinTargets:  0,
			MaxTargets:  0,
			TargetTypes: []string{},
			Mandatory:   false,
			Description: "",
		}
	}
}

// TargetRequirements describes the targeting requirements for an effect
type TargetRequirements struct {
	MinTargets  int      // Minimum number of targets required
	MaxTargets  int      // Maximum number of targets allowed (-1 for unlimited)
	TargetTypes []string // Valid target types
	Mandatory   bool     // Whether targeting is mandatory
	Description string   // Human-readable description
}

// TargetSelector provides methods for interactive target selection
type TargetSelector struct {
	validator *TargetValidator
}

// NewTargetSelector creates a new target selector
func NewTargetSelector(validator *TargetValidator) *TargetSelector {
	return &TargetSelector{
		validator: validator,
	}
}

// GetSelectableTargets returns targets that can be selected for an effect
func (s *TargetSelector) GetSelectableTargets(effect effects.Effect, source *models.Card) SelectableTargets {
	targets := effect.GetTargets(s.validator.context, source)
	requirements := s.validator.GetTargetRequirements(effect)
	
	// Group targets by type and location for easier UI display
	grouped := make(map[string][]effects.Target)
	for _, target := range targets {
		key := fmt.Sprintf("%s_%s", target.Type, getTargetLocationGroup(target.Location))
		grouped[key] = append(grouped[key], target)
	}
	
	return SelectableTargets{
		Targets:      targets,
		Grouped:      grouped,
		Requirements: requirements,
	}
}

// SelectableTargets contains information about targets that can be selected
type SelectableTargets struct {
	Targets      []effects.Target
	Grouped      map[string][]effects.Target
	Requirements TargetRequirements
}

// Helper function to group locations
func getTargetLocationGroup(location string) string {
	// Extract the general location type
	if len(location) > 12 && location[:12] == "battle_area_" {
		return "battle_area"
	}
	return location
}