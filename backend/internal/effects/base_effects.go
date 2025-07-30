package effects

import (
	"fmt"
	"github.com/your-org/card_game/backend/internal/models"
)

// BaseEffect provides common functionality for effects
type BaseEffect struct {
	Trigger     TriggerType
	Description string
}

func (e *BaseEffect) GetTrigger() TriggerType {
	return e.Trigger
}

func (e *BaseEffect) GetDescription() string {
	return e.Description
}

// PowerModifierEffect modifies the power of friends
type PowerModifierEffect struct {
	BaseEffect
	Scope    EffectScope
	Amount   int
	Duration string // "turn", "permanent", "while_condition"
	Condition func(game *GameContext, source *models.Card) bool
}

func (e *PowerModifierEffect) CanActivate(game *GameContext, source *models.Card) bool {
	if e.Condition != nil {
		return e.Condition(game, source)
	}
	return true
}

func (e *PowerModifierEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	switch e.Scope {
	case ScopeSelf:
		// Target only the source friend
		for pos, friend := range playerState.BattleArea {
			if friend.CardNo == source.CardNo {
				targets = append(targets, Target{
					Type:     "friend",
					ID:       friend.CardNo,
					Location: fmt.Sprintf("battle_area_%d_%s", game.ActivePlayer, pos),
				})
			}
		}
	case ScopeMyFriends:
		// Target all of the player's friends
		for pos, friend := range playerState.BattleArea {
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: fmt.Sprintf("battle_area_%d_%s", game.ActivePlayer, pos),
			})
		}
	case ScopeOppFriends:
		// Target all opponent's friends
		oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
		oppState := game.GetPlayerState(oppPlayer)
		for pos, friend := range oppState.BattleArea {
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: fmt.Sprintf("battle_area_%d_%s", oppPlayer, pos),
			})
		}
	}
	
	return targets
}

func (e *PowerModifierEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	for _, target := range targets {
		// Extract player number from location
		var player int
		fmt.Sscanf(target.Location, "battle_area_%d_", &player)
		
		if err := game.ModifyPower(player, target.ID, e.Amount); err != nil {
			return err
		}
	}
	return nil
}

// DrawCardEffect allows drawing cards
type DrawCardEffect struct {
	BaseEffect
	Count int
}

func (e *DrawCardEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *DrawCardEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// No targets needed for drawing cards
	return nil
}

func (e *DrawCardEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	return game.DrawCards(game.ActivePlayer, e.Count)
}

// ReturnToHandEffect returns a friend to hand
type ReturnToHandEffect struct {
	BaseEffect
	Scope         EffectScope
	MaxCost       int  // Maximum cost of friend that can be returned
	RequireTarget bool // Whether the player must choose a target
}

func (e *ReturnToHandEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Check if there are valid targets
	targets := e.GetTargets(game, source)
	return len(targets) > 0 || !e.RequireTarget
}

func (e *ReturnToHandEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	
	if e.Scope == ScopeOppFriends {
		oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
		oppState := game.GetPlayerState(oppPlayer)
		
		for pos, friend := range oppState.BattleArea {
			// TODO: Check friend cost against MaxCost
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: fmt.Sprintf("battle_area_%d_%s", oppPlayer, pos),
			})
		}
	}
	
	return targets
}

func (e *ReturnToHandEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) == 0 && e.RequireTarget {
		return fmt.Errorf("no target selected")
	}
	
	for _, target := range targets {
		// Extract player number from location
		var player int
		var pos string
		fmt.Sscanf(target.Location, "battle_area_%d_%s", &player, &pos)
		
		if err := game.ReturnToHand(player, target.ID); err != nil {
			return err
		}
	}
	
	return nil
}

// DestroyFriendEffect destroys friends
type DestroyFriendEffect struct {
	BaseEffect
	MaxPower      int  // Maximum power of friend that can be destroyed
	RequireTarget bool // Whether the player must choose a target
}

func (e *DestroyFriendEffect) CanActivate(game *GameContext, source *models.Card) bool {
	targets := e.GetTargets(game, source)
	return len(targets) > 0 || !e.RequireTarget
}

func (e *DestroyFriendEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	oppState := game.GetPlayerState(oppPlayer)
	
	for pos, friend := range oppState.BattleArea {
		if friend.Power <= e.MaxPower {
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: fmt.Sprintf("battle_area_%d_%s", oppPlayer, pos),
				Data:     friend,
			})
		}
	}
	
	return targets
}

func (e *DestroyFriendEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) == 0 && e.RequireTarget {
		return fmt.Errorf("no target selected")
	}
	
	for _, target := range targets {
		// Extract player number from location
		var player int
		fmt.Sscanf(target.Location, "battle_area_%d_", &player)
		
		if err := game.DestroyFriend(player, target.ID); err != nil {
			return err
		}
	}
	
	return nil
}

// RevealNegativeEnergyEffect reveals cards from negative energy area
type RevealNegativeEnergyEffect struct {
	BaseEffect
	Count int
}

func (e *RevealNegativeEnergyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if there are unrevealed cards in negative energy
	return len(playerState.NegativeEnergy) > 0
}

func (e *RevealNegativeEnergyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *RevealNegativeEnergyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	return game.RevealNegEnergy(game.ActivePlayer, e.Count)
}

// RestFriendEffect rests or activates friends
type RestFriendEffect struct {
	BaseEffect
	Action string // "rest" or "active"
}

func (e *RestFriendEffect) CanActivate(game *GameContext, source *models.Card) bool {
	targets := e.GetTargets(game, source)
	return len(targets) > 0
}

func (e *RestFriendEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	
	if e.Action == "rest" {
		// Target opponent's active friends
		oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
		oppState := game.GetPlayerState(oppPlayer)
		
		for pos, friend := range oppState.BattleArea {
			if !friend.IsRest {
				targets = append(targets, Target{
					Type:     "friend",
					ID:       friend.CardNo,
					Location: fmt.Sprintf("battle_area_%d_%s", oppPlayer, pos),
				})
			}
		}
	} else if e.Action == "active" {
		// Target own rested friends
		playerState := game.GetPlayerState(game.ActivePlayer)
		
		for pos, friend := range playerState.BattleArea {
			if friend.IsRest {
				targets = append(targets, Target{
					Type:     "friend",
					ID:       friend.CardNo,
					Location: fmt.Sprintf("battle_area_%d_%s", game.ActivePlayer, pos),
				})
			}
		}
	}
	
	return targets
}

func (e *RestFriendEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	for _, target := range targets {
		// Extract player number from location
		var player int
		fmt.Sscanf(target.Location, "battle_area_%d_", &player)
		
		if e.Action == "rest" {
			if err := game.RestFriend(player, target.ID); err != nil {
				return err
			}
		} else if e.Action == "active" {
			if err := game.ActiveFriend(player, target.ID); err != nil {
				return err
			}
		}
	}
	
	return nil
}