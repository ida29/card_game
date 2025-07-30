package game

import (
	"fmt"
	"github.com/your-org/card_game/backend/internal/effects"
	"github.com/your-org/card_game/backend/internal/models"
)

// EffectStack manages the resolution of effects in a Last-In-First-Out manner
type EffectStack struct {
	items []EffectStackItem
}

// EffectStackItem represents an effect waiting to be resolved
type EffectStackItem struct {
	Effect   effects.Effect
	Source   *models.Card
	Targets  []effects.Target
	Player   int
	Priority int // Higher priority effects resolve first at the same stack level
}

// NewEffectStack creates a new effect stack
func NewEffectStack() *EffectStack {
	return &EffectStack{
		items: make([]EffectStackItem, 0),
	}
}

// Push adds an effect to the stack
func (s *EffectStack) Push(item EffectStackItem) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top effect from the stack
func (s *EffectStack) Pop() (EffectStackItem, bool) {
	if len(s.items) == 0 {
		return EffectStackItem{}, false
	}
	
	// Get the last item
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	
	return item, true
}

// IsEmpty returns true if the stack has no items
func (s *EffectStack) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *EffectStack) Size() int {
	return len(s.items)
}

// EffectResolver handles the resolution of effects with proper ordering
type EffectResolver struct {
	stack        *EffectStack
	context      *effects.GameContext
	eventHandler *EventHandler
}

// NewEffectResolver creates a new effect resolver
func NewEffectResolver(context *effects.GameContext, eventHandler *EventHandler) *EffectResolver {
	return &EffectResolver{
		stack:        NewEffectStack(),
		context:      context,
		eventHandler: eventHandler,
	}
}

// QueueEffect adds an effect to be resolved
func (r *EffectResolver) QueueEffect(effect effects.Effect, source *models.Card, player int) error {
	// Check if effect can activate
	r.context.ActivePlayer = player
	if !effect.CanActivate(r.context, source) {
		return fmt.Errorf("effect cannot activate")
	}
	
	// Get targets
	targets := effect.GetTargets(r.context, source)
	
	// Add to stack
	r.stack.Push(EffectStackItem{
		Effect:  effect,
		Source:  source,
		Targets: targets,
		Player:  player,
	})
	
	return nil
}

// ResolveAll resolves all effects on the stack
func (r *EffectResolver) ResolveAll() error {
	for !r.stack.IsEmpty() {
		item, ok := r.stack.Pop()
		if !ok {
			break
		}
		
		// Set active player for this effect
		r.context.ActivePlayer = item.Player
		
		// Apply the effect
		if err := item.Effect.Apply(r.context, item.Source, item.Targets); err != nil {
			return fmt.Errorf("failed to apply effect %s: %w", item.Effect.GetDescription(), err)
		}
		
		// Check for any triggered effects
		// This is where chain reactions can occur
		r.checkTriggeredEffects(item)
	}
	
	return nil
}

// checkTriggeredEffects checks if resolving an effect triggers any other effects
func (r *EffectResolver) checkTriggeredEffects(item EffectStackItem) {
	// For example, if a friend was destroyed, check for "when destroyed" triggers
	// This would require tracking what happened during effect resolution
	
	// TODO: Implement trigger checking based on effect results
}

// InteractionController manages player interactions during effect resolution
type InteractionController struct {
	resolver     *EffectResolver
	pendingChoices map[string]PendingChoice
}

// PendingChoice represents a choice waiting for player input
type PendingChoice struct {
	Player      int
	ChoiceType  string // "target", "option", "order", etc.
	Options     []interface{}
	Required    bool
	Description string
}

// NewInteractionController creates a new interaction controller
func NewInteractionController(resolver *EffectResolver) *InteractionController {
	return &InteractionController{
		resolver:       resolver,
		pendingChoices: make(map[string]PendingChoice),
	}
}

// RequestTargetSelection requests the player to select targets
func (c *InteractionController) RequestTargetSelection(player int, validTargets []effects.Target, minTargets, maxTargets int, description string) string {
	choiceID := generateChoiceID()
	
	c.pendingChoices[choiceID] = PendingChoice{
		Player:      player,
		ChoiceType:  "target",
		Options:     targetsToInterface(validTargets),
		Required:    minTargets > 0,
		Description: description,
	}
	
	return choiceID
}

// RequestOptionSelection requests the player to select from options
func (c *InteractionController) RequestOptionSelection(player int, options []string, description string) string {
	choiceID := generateChoiceID()
	
	optionsInterface := make([]interface{}, len(options))
	for i, opt := range options {
		optionsInterface[i] = opt
	}
	
	c.pendingChoices[choiceID] = PendingChoice{
		Player:      player,
		ChoiceType:  "option",
		Options:     optionsInterface,
		Required:    true,
		Description: description,
	}
	
	return choiceID
}

// SubmitChoice submits a player's choice
func (c *InteractionController) SubmitChoice(choiceID string, selection []int) error {
	choice, exists := c.pendingChoices[choiceID]
	if !exists {
		return fmt.Errorf("choice not found")
	}
	
	// Validate selection
	if len(selection) == 0 && choice.Required {
		return fmt.Errorf("selection required")
	}
	
	for _, idx := range selection {
		if idx < 0 || idx >= len(choice.Options) {
			return fmt.Errorf("invalid selection index")
		}
	}
	
	// Process the choice based on type
	switch choice.ChoiceType {
	case "target":
		// Convert back to targets
		selectedTargets := make([]effects.Target, len(selection))
		for i, idx := range selection {
			if target, ok := choice.Options[idx].(effects.Target); ok {
				selectedTargets[i] = target
			}
		}
		// TODO: Apply the selected targets to the pending effect
		
	case "option":
		// Handle option selection
		// TODO: Apply the selected option
	}
	
	// Remove the pending choice
	delete(c.pendingChoices, choiceID)
	
	return nil
}

// Helper functions

func targetsToInterface(targets []effects.Target) []interface{} {
	result := make([]interface{}, len(targets))
	for i, target := range targets {
		result[i] = target
	}
	return result
}

func generateChoiceID() string {
	// TODO: Implement proper ID generation
	return fmt.Sprintf("choice_%d", len(choices))
}

var choices = make(map[string]bool) // Temporary tracking