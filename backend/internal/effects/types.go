package effects

import (
	"github.com/your-org/card_game/backend/internal/models"
)

// TriggerType represents when an effect triggers
type TriggerType string

const (
	TriggerOnPlay        TriggerType = "on_play"        // 登場時
	TriggerOnAttack      TriggerType = "on_attack"      // アタック時
	TriggerOnBlock       TriggerType = "on_block"       // ブロック時
	TriggerOnDestroy     TriggerType = "on_destroy"     // 破壊時
	TriggerOnDamageDealt TriggerType = "on_damage_dealt" // ダメージを与えた時
	TriggerPersistent    TriggerType = "persistent"     // 常在
	TriggerStartPhase    TriggerType = "start_phase"    // フェイズ開始時
	TriggerEndPhase      TriggerType = "end_phase"      // フェイズ終了時
	TriggerMain          TriggerType = "main"           // メイン
	TriggerCounter       TriggerType = "counter"        // カウンター
)

// EffectScope represents the scope of an effect
type EffectScope string

const (
	ScopeSelf         EffectScope = "self"          // 自分
	ScopeOpponent     EffectScope = "opponent"      // 相手
	ScopeAllFriends   EffectScope = "all_friends"   // 全てのふれんど
	ScopeMyFriends    EffectScope = "my_friends"    // 自分のふれんど
	ScopeOppFriends   EffectScope = "opp_friends"   // 相手のふれんど
	ScopeTarget       EffectScope = "target"         // 対象
)

// Effect represents a card effect
type Effect interface {
	// GetTrigger returns when this effect triggers
	GetTrigger() TriggerType
	
	// CanActivate checks if the effect can be activated in the current game state
	CanActivate(game *GameContext, source *models.Card) bool
	
	// GetTargets returns valid targets for the effect
	GetTargets(game *GameContext, source *models.Card) []Target
	
	// Apply applies the effect to the game state
	Apply(game *GameContext, source *models.Card, targets []Target) error
	
	// GetDescription returns a human-readable description of the effect
	GetDescription() string
}

// Target represents a valid target for an effect
type Target struct {
	Type     string      // "friend", "card", "player", etc.
	ID       string      // Unique identifier for the target
	Location string      // Where the target is located (battle_area, hand, etc.)
	Data     interface{} // Additional target data
}

// GameContext provides access to the game state and methods to modify it
type GameContext struct {
	Game         *models.Game
	ActivePlayer int // 1 or 2
	
	// Helper methods to modify game state
	DrawCards         func(player int, count int) error
	DestroyFriend     func(player int, cardNo string) error
	ReturnToHand      func(player int, cardNo string) error
	RestFriend        func(player int, cardNo string) error
	ActiveFriend      func(player int, cardNo string) error
	ModifyPower       func(player int, cardNo string, amount int) error
	RevealNegEnergy   func(player int, count int) error
	PlaceFieldCard    func(player int, cardNo string) error
	MoveToTrash       func(player int, cardNo string, from string) error
	MoveCardToDeck    func(player int, cardNo string, position string) error // "top" or "bottom"
	DealDamage        func(player int, amount int) error
	AddToEnergyArea   func(player int, cardNo string) error
	GetPlayerState    func(player int) *models.PlayerState
	GetOpponentPlayer func(player int) int
}

// EffectRegistry holds all registered effects
type EffectRegistry struct {
	effects map[string]Effect
}

// NewEffectRegistry creates a new effect registry
func NewEffectRegistry() *EffectRegistry {
	return &EffectRegistry{
		effects: make(map[string]Effect),
	}
}

// Register registers an effect for a card
func (r *EffectRegistry) Register(cardNo string, effect Effect) {
	r.effects[cardNo] = effect
}

// GetEffect returns the effect for a card
func (r *EffectRegistry) GetEffect(cardNo string) (Effect, bool) {
	effect, exists := r.effects[cardNo]
	return effect, exists
}

// GetEffectsForTrigger returns all effects that trigger on a specific event
func (r *EffectRegistry) GetEffectsForTrigger(trigger TriggerType) []string {
	var cardNos []string
	for cardNo, effect := range r.effects {
		if effect.GetTrigger() == trigger {
			cardNos = append(cardNos, cardNo)
		}
	}
	return cardNos
}