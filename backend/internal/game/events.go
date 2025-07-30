package game

import (
	"fmt"
	"github.com/your-org/card_game/backend/internal/effects"
	"github.com/your-org/card_game/backend/internal/models"
)

// EventType represents the type of game event
type EventType string

const (
	EventFriendPlayed    EventType = "friend_played"
	EventFriendAttacks   EventType = "friend_attacks"
	EventFriendBlocks    EventType = "friend_blocks"
	EventFriendDestroyed EventType = "friend_destroyed"
	EventDamageDealt     EventType = "damage_dealt"
	EventPhaseStart      EventType = "phase_start"
	EventPhaseEnd        EventType = "phase_end"
	EventSupportPlayed   EventType = "support_played"
	EventFieldPlayed     EventType = "field_played"
	EventTurnStart       EventType = "turn_start"
	EventTurnEnd         EventType = "turn_end"
)

// GameEvent represents an event that occurred in the game
type GameEvent struct {
	Type     EventType
	Player   int
	CardNo   string
	Target   string
	Phase    models.GamePhase
	Data     map[string]interface{}
}

// EventHandler handles game events and triggers effects
type EventHandler struct {
	game           *models.Game
	effectRegistry *effects.EffectRegistry
	eventQueue     []GameEvent
	context        *effects.GameContext
}

// NewEventHandler creates a new event handler
func NewEventHandler(game *models.Game) *EventHandler {
	return &EventHandler{
		game:           game,
		effectRegistry: effects.GetGlobalRegistry(),
		eventQueue:     make([]GameEvent, 0),
		context:        createGameContext(game),
	}
}

// TriggerEvent adds an event to the queue and processes it
func (h *EventHandler) TriggerEvent(event GameEvent) error {
	h.eventQueue = append(h.eventQueue, event)
	
	// Process the event
	return h.processEvent(event)
}

// processEvent handles a single event and triggers relevant effects
func (h *EventHandler) processEvent(event GameEvent) error {
	var triggerType effects.TriggerType
	
	switch event.Type {
	case EventFriendPlayed:
		triggerType = effects.TriggerOnPlay
	case EventFriendAttacks:
		triggerType = effects.TriggerOnAttack
	case EventFriendBlocks:
		triggerType = effects.TriggerOnBlock
	case EventFriendDestroyed:
		triggerType = effects.TriggerOnDestroy
	case EventDamageDealt:
		triggerType = effects.TriggerOnDamageDealt
	case EventPhaseStart:
		triggerType = effects.TriggerStartPhase
	case EventPhaseEnd:
		triggerType = effects.TriggerEndPhase
	case EventSupportPlayed:
		// Support cards can have main or counter triggers
		if h.game.CurrentPhase == models.PhaseMain {
			triggerType = effects.TriggerMain
		} else {
			triggerType = effects.TriggerCounter
		}
	default:
		// No effects to trigger for this event type
		return nil
	}
	
	// Get all cards that could trigger on this event
	cardNos := h.effectRegistry.GetEffectsForTrigger(triggerType)
	
	// Check each potential trigger
	for _, cardNo := range cardNos {
		// Check if this card is in play or relevant to the event
		if h.isCardRelevant(cardNo, event) {
			effect, _ := h.effectRegistry.GetEffect(cardNo)
			
			// Load card data
			card, err := h.loadCard(cardNo)
			if err != nil {
				continue
			}
			
			// Update context with current player
			h.context.ActivePlayer = event.Player
			
			// Check if effect can activate
			if effect.CanActivate(h.context, card) {
				// Get targets
				targets := effect.GetTargets(h.context, card)
				
				// Apply effect
				if err := effect.Apply(h.context, card, targets); err != nil {
					return fmt.Errorf("failed to apply effect for %s: %w", cardNo, err)
				}
			}
		}
	}
	
	// Process persistent effects
	h.processPersistentEffects()
	
	return nil
}

// isCardRelevant checks if a card is relevant to the current event
func (h *EventHandler) isCardRelevant(cardNo string, event GameEvent) bool {
	// Check if the card triggered the event
	if event.CardNo == cardNo {
		return true
	}
	
	// Check if card is on the battlefield
	if h.game.GameState != nil {
		player1State := h.game.GameState.Player1State
		player2State := h.game.GameState.Player2State
		
		// Check battle areas
		for _, friend := range player1State.BattleArea {
			if friend.CardNo == cardNo {
				return true
			}
		}
		for _, friend := range player2State.BattleArea {
			if friend.CardNo == cardNo {
				return true
			}
		}
		
		// Check field cards
		if player1State.FieldCard != nil && *player1State.FieldCard == cardNo {
			return true
		}
		if player2State.FieldCard != nil && *player2State.FieldCard == cardNo {
			return true
		}
	}
	
	return false
}

// processPersistentEffects processes all active persistent effects
func (h *EventHandler) processPersistentEffects() {
	// Get all persistent effects
	cardNos := h.effectRegistry.GetEffectsForTrigger(effects.TriggerPersistent)
	
	for _, cardNo := range cardNos {
		if h.isCardOnField(cardNo) {
			effect, _ := h.effectRegistry.GetEffect(cardNo)
			card, err := h.loadCard(cardNo)
			if err != nil {
				continue
			}
			
			// Check which player controls this card
			player := h.getCardController(cardNo)
			h.context.ActivePlayer = player
			
			if effect.CanActivate(h.context, card) {
				targets := effect.GetTargets(h.context, card)
				effect.Apply(h.context, card, targets)
			}
		}
	}
}

// isCardOnField checks if a card is currently on the field
func (h *EventHandler) isCardOnField(cardNo string) bool {
	if h.game.GameState == nil {
		return false
	}
	
	// Check battle areas
	for _, friend := range h.game.GameState.Player1State.BattleArea {
		if friend.CardNo == cardNo {
			return true
		}
	}
	for _, friend := range h.game.GameState.Player2State.BattleArea {
		if friend.CardNo == cardNo {
			return true
		}
	}
	
	// Check field cards
	if h.game.GameState.Player1State.FieldCard != nil && 
	   *h.game.GameState.Player1State.FieldCard == cardNo {
		return true
	}
	if h.game.GameState.Player2State.FieldCard != nil && 
	   *h.game.GameState.Player2State.FieldCard == cardNo {
		return true
	}
	
	return false
}

// getCardController returns which player controls a card
func (h *EventHandler) getCardController(cardNo string) int {
	if h.game.GameState == nil {
		return 0
	}
	
	// Check player 1's cards
	for _, friend := range h.game.GameState.Player1State.BattleArea {
		if friend.CardNo == cardNo {
			return 1
		}
	}
	if h.game.GameState.Player1State.FieldCard != nil && 
	   *h.game.GameState.Player1State.FieldCard == cardNo {
		return 1
	}
	
	// Check player 2's cards
	for _, friend := range h.game.GameState.Player2State.BattleArea {
		if friend.CardNo == cardNo {
			return 2
		}
	}
	if h.game.GameState.Player2State.FieldCard != nil && 
	   *h.game.GameState.Player2State.FieldCard == cardNo {
		return 2
	}
	
	return 0
}

// loadCard loads card data (placeholder - should use actual card service)
func (h *EventHandler) loadCard(cardNo string) (*models.Card, error) {
	// TODO: Load from card service/database
	return &models.Card{
		CardNo: cardNo,
	}, nil
}

// createGameContext creates a game context for effects
func createGameContext(game *models.Game) *effects.GameContext {
	return &effects.GameContext{
		Game: game,
		
		DrawCards: func(player int, count int) error {
			// TODO: Implement
			return nil
		},
		
		DestroyFriend: func(player int, cardNo string) error {
			// TODO: Implement
			return nil
		},
		
		ReturnToHand: func(player int, cardNo string) error {
			// TODO: Implement
			return nil
		},
		
		RestFriend: func(player int, cardNo string) error {
			if game.GameState == nil {
				return fmt.Errorf("no game state")
			}
			
			var playerState *models.PlayerState
			if player == 1 {
				playerState = &game.GameState.Player1State
			} else {
				playerState = &game.GameState.Player2State
			}
			
			for pos, friend := range playerState.BattleArea {
				if friend.CardNo == cardNo {
					friend.IsRest = true
					playerState.BattleArea[pos] = friend
					return nil
				}
			}
			
			return fmt.Errorf("friend not found")
		},
		
		ActiveFriend: func(player int, cardNo string) error {
			if game.GameState == nil {
				return fmt.Errorf("no game state")
			}
			
			var playerState *models.PlayerState
			if player == 1 {
				playerState = &game.GameState.Player1State
			} else {
				playerState = &game.GameState.Player2State
			}
			
			for pos, friend := range playerState.BattleArea {
				if friend.CardNo == cardNo {
					friend.IsRest = false
					playerState.BattleArea[pos] = friend
					return nil
				}
			}
			
			return fmt.Errorf("friend not found")
		},
		
		ModifyPower: func(player int, cardNo string, amount int) error {
			if game.GameState == nil {
				return fmt.Errorf("no game state")
			}
			
			var playerState *models.PlayerState
			if player == 1 {
				playerState = &game.GameState.Player1State
			} else {
				playerState = &game.GameState.Player2State
			}
			
			for pos, friend := range playerState.BattleArea {
				if friend.CardNo == cardNo {
					friend.Power += amount
					if friend.Power < 0 {
						friend.Power = 0
					}
					playerState.BattleArea[pos] = friend
					return nil
				}
			}
			
			return fmt.Errorf("friend not found")
		},
		
		RevealNegEnergy: func(player int, count int) error {
			// TODO: Implement - need to track revealed state
			return nil
		},
		
		PlaceFieldCard: func(player int, cardNo string) error {
			// TODO: Implement
			return nil
		},
		
		MoveToTrash: func(player int, cardNo string, from string) error {
			// TODO: Implement
			return nil
		},
		
		MoveCardToDeck: func(player int, cardNo string, position string) error {
			// TODO: Implement
			return nil
		},
		
		DealDamage: func(player int, amount int) error {
			// TODO: Implement
			return nil
		},
		
		AddToEnergyArea: func(player int, cardNo string) error {
			// TODO: Implement
			return nil
		},
		
		GetPlayerState: func(player int) *models.PlayerState {
			if game.GameState == nil {
				return nil
			}
			
			if player == 1 {
				return &game.GameState.Player1State
			}
			return &game.GameState.Player2State
		},
		
		GetOpponentPlayer: func(player int) int {
			if player == 1 {
				return 2
			}
			return 1
		},
	}
}