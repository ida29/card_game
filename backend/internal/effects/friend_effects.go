package effects

import (
	"mememe-tcg/internal/models"
)

// Friend card specific effects

// F-002 なみだぶくろん - Main phase: This turn, this friend gets +1000 power
func NewNamidabukuronEffect() Effect {
	return &MainPhasePowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain,
			Description: "【メイン：コスト〇】このターン中、このふれんどのパワー+1000。",
		},
		PowerBoost: 1000,
		Duration:   "turn",
	}
}

// F-003 フラフラ - For every 2 cards in hand, +1000 power
func NewFurafuraEffect() Effect {
	return &HandSizePowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分の手札2枚ごとに、このふれんどのパワー+1000。",
		},
		CardsPerBoost: 2,
		PowerBoost:    1000,
	}
}

// F-006/F-006(P) ヒヤケラトプス - When this friend attacks, draw 1 card
func NewHiyakeratopusEffect() Effect {
	return &DrawCardEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分はデッキから1枚ドローする。",
		},
		Amount: 1,
		Scope:  ScopeSelf,
	}
}

// F-008/F-008(P) ボーイ - When this friend attacks, destroy 1 opponent's friend with power 3000 or less
func NewBoyEffect() Effect {
	return &DestroyFriendEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、パワー3000以下の相手のふれんど1体を破壊する。",
		},
		Scope:         ScopeOppFriends,
		MaxPower:      3000,
		RequireTarget: true,
	}
}

// F-011/F-011(P) ポチ - When this friend attacks, draw 1 card
func NewPochiEffect() Effect {
	return &DrawCardEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分はデッキから1枚ドローする。",
		},
		Amount: 1,
		Scope:  ScopeSelf,
	}
}

// F-020/F-020(P) マルカニ - When this friend attacks, reveal top deck card and place on top or bottom
func NewMarukaniEffect() Effect {
	return &RevealAndPlaceDeckTopEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分か相手のデッキを上から1枚オープンする。そのカードをデッキの上か下に置く。",
		},
	}
}

// F-022/F-022(P) ジョニー - When this friend attacks, discard top card of deck
func NewJohnnyEffect() Effect {
	return &DiscardDeckTopEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分のデッキを上から1枚破棄する。",
		},
	}
}

// F-044/F-044(P) うっきー - When this friend attacks, may discard 1 negative energy
func NewUkkiAttackEffect() Effect {
	return &DiscardNegativeEnergyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分の負のエネルギーエリアのカード1枚を破棄できる。",
		},
		Count:    1,
		Optional: true,
	}
}

// F-102 くらげ坊(変身) - When this friend attacks, turn 2 negative energy face down to activate
func NewKurageboTransformEffect() Effect {
	return &ActivateByFlippingNegativeEnergyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、自分の負のエネルギーエリアのカード2枚を裏にすることで、このふれんどをアクティブにする。",
		},
		Count: 2,
	}
}

// F-013/F-013(P) るくそー - When this friend blocks, reveal 1 card from your negative energy area
func NewRukusoEffect() Effect {
	return &RevealNegativeEnergyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnBlock,
			Description: "このふれんどがブロックした時、自分の負のエネルギーエリアのカード1枚を表にする。",
		},
		Count: 1,
	}
}

// F-015/F-015(P) ティラノちゃん - This friend can attack on the turn it's played
func NewTiranoEffect() Effect {
	return &CanAttackImmediatelyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "このふれんどは登場したターンにアタックできる。",
		},
	}
}

// F-016/F-016(P) くらげ坊 - While you have 3 or more revealed cards in negative energy, this friend deals +1 damage
func NewKurageboEffect() Effect {
	return &DamageModifierEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分の負のエネルギーエリアの表のカードが3枚以上の間、このふれんどの与えるダメージ+1。",
		},
		Amount: 1,
		Condition: func(game *GameContext, source *models.Card) bool {
			playerState := game.GetPlayerState(game.ActivePlayer)
			// TODO: Count revealed cards in negative energy
			revealedCount := 0
			for _, cardNo := range playerState.NegativeEnergy {
				// Need to track which cards are revealed
				_ = cardNo
				revealedCount++
			}
			return revealedCount >= 3
		},
	}
}

// F-023/F-023(P) ユピ - When this friend enters play, return 1 opponent's friend to hand
func NewYupiEffect() Effect {
	return &ReturnToHandEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnPlay,
			Description: "このふれんどが登場した時、相手のふれんど1体を手札に戻す。",
		},
		Scope:         ScopeOppFriends,
		RequireTarget: true,
	}
}

// F-025/F-025(P) しもん - When played, return a cost 3 or less support card from trash to hand
func NewShimonEffect() Effect {
	return &ReturnSupportFromTrashEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnPlay,
			Description: "このふれんどが登場した時、コスト3以下の自分のトラッシュのサポートカード1枚を手札に加えられる。",
		},
		MaxCost: 3,
	}
}

// F-034/F-034(P) メガロッコ - When this friend blocks, activate it
func NewMegarokkoEffect() Effect {
	return &ActivateOnBlockEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnBlock,
			Description: "このふれんどがブロックした時、このふれんどをアクティブにする。",
		},
	}
}

// F-041/F-041(P) ハヤオ - When played, can play a cost 2 or less field card without paying cost
func NewHayaoEffect() Effect {
	return &PlayFieldCardEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnPlay,
			Description: "このふれんどが登場した時、自分の手札から、コスト2以下のフィールドカード1枚をコストを支払わずに置ける。",
		},
		MaxCost: 2,
		NoCost:  true,
	}
}

// F-042/F-042(P) うっきー - When played, return 3 cards from trash to top or bottom of deck
func NewUkkiEffect() Effect {
	return &ReturnToCardsToDeckEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnPlay,
			Description: "このふれんどが登場した時、自分のトラッシュのカード3枚を1枚ずつ選んでデッキの上か下に置ける。",
		},
		Count: 3,
	}
}

// F-055/F-055(P) Ko2 - When this friend attacks, rest 1 opponent's friend
func NewKo2Effect() Effect {
	return &RestFriendEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnAttack,
			Description: "このふれんどがアタックした時、相手のふれんど1体をレストする。",
		},
		Action: "rest",
	}
}

// F-056/F-056(P) シーラン - This friend can block even while rested
func NewShiranEffect() Effect {
	return &CanBlockWhileRestedEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "このふれんどはレストしていてもブロックできる。",
		},
	}
}

// Custom effect types for specific abilities

type CanAttackImmediatelyEffect struct {
	BaseEffect
}

func (e *CanAttackImmediatelyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *CanAttackImmediatelyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *CanAttackImmediatelyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This is a persistent effect that modifies game rules
	// It should be checked when determining if a friend can attack
	return nil
}

type DamageModifierEffect struct {
	BaseEffect
	Amount    int
	Condition func(game *GameContext, source *models.Card) bool
}

func (e *DamageModifierEffect) CanActivate(game *GameContext, source *models.Card) bool {
	if e.Condition != nil {
		return e.Condition(game, source)
	}
	return true
}

func (e *DamageModifierEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *DamageModifierEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This is a persistent effect that modifies damage calculation
	// It should be checked when calculating damage
	return nil
}

type ReturnSupportFromTrashEffect struct {
	BaseEffect
	MaxCost int
}

func (e *ReturnSupportFromTrashEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if there are support cards in trash
	return len(playerState.Trash) > 0 // TODO: Filter for support cards with cost <= MaxCost
}

func (e *ReturnSupportFromTrashEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	// TODO: Need to load card data to check type and cost
	for _, cardNo := range playerState.Trash {
		// Check if card is support and cost <= MaxCost
		targets = append(targets, Target{
			Type:     "card",
			ID:       cardNo,
			Location: "trash",
		})
	}
	
	return targets
}

func (e *ReturnSupportFromTrashEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Move selected card from trash to hand
	// TODO: Implement
	return nil
}

type ActivateOnBlockEffect struct {
	BaseEffect
}

func (e *ActivateOnBlockEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *ActivateOnBlockEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// Target is always the blocking friend itself
	return []Target{{
		Type:     "friend",
		ID:       source.CardNo,
		Location: "self",
	}}
}

func (e *ActivateOnBlockEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	return game.ActiveFriend(game.ActivePlayer, source.CardNo)
}

type PlayFieldCardEffect struct {
	BaseEffect
	MaxCost int
	NoCost  bool
}

func (e *PlayFieldCardEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if player has field cards in hand
	return len(playerState.Hand) > 0 // TODO: Filter for field cards with cost <= MaxCost
}

func (e *PlayFieldCardEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	// TODO: Need to load card data to check type and cost
	for _, cardNo := range playerState.Hand {
		// Check if card is field and cost <= MaxCost
		targets = append(targets, Target{
			Type:     "card",
			ID:       cardNo,
			Location: "hand",
		})
	}
	
	return targets
}

func (e *PlayFieldCardEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) > 0 {
		return game.PlaceFieldCard(game.ActivePlayer, targets[0].ID)
	}
	return nil
}

type ReturnToCardsToDeckEffect struct {
	BaseEffect
	Count int
}

func (e *ReturnToCardsToDeckEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	return len(playerState.Trash) >= e.Count
}

func (e *ReturnToCardsToDeckEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	for _, cardNo := range playerState.Trash {
		targets = append(targets, Target{
			Type:     "card",
			ID:       cardNo,
			Location: "trash",
		})
	}
	
	return targets
}

func (e *ReturnToCardsToDeckEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Player chooses position (top or bottom) for each card
	for _, target := range targets {
		// TODO: Get player choice for position
		position := "top" // or "bottom"
		if err := game.MoveCardToDeck(game.ActivePlayer, target.ID, position); err != nil {
			return err
		}
	}
	return nil
}

type CanBlockWhileRestedEffect struct {
	BaseEffect
}

func (e *CanBlockWhileRestedEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *CanBlockWhileRestedEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *CanBlockWhileRestedEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This is a persistent effect that modifies game rules
	// It should be checked when determining if a friend can block
	return nil
}

// Custom effect types for new attack effects

type RevealAndPlaceDeckTopEffect struct {
	BaseEffect
}

func (e *RevealAndPlaceDeckTopEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Can always activate if there's a deck
	return true
}

func (e *RevealAndPlaceDeckTopEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// Player chooses whose deck to reveal
	return []Target{
		{Type: "deck", ID: "self", Location: "deck"},
		{Type: "deck", ID: "opponent", Location: "deck"},
	}
}

func (e *RevealAndPlaceDeckTopEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Implementation would reveal top card and let player choose placement
	return nil
}

type DiscardDeckTopEffect struct {
	BaseEffect
}

func (e *DiscardDeckTopEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	return len(playerState.Deck) > 0
}

func (e *DiscardDeckTopEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *DiscardDeckTopEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	return game.DiscardFromDeckTop(game.ActivePlayer, 1)
}

type DiscardNegativeEnergyEffect struct {
	BaseEffect
	Count    int
	Optional bool
}

func (e *DiscardNegativeEnergyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	return len(playerState.NegativeEnergy) > 0
}

func (e *DiscardNegativeEnergyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	if !e.Optional {
		return nil // Mandatory effects might not need target selection
	}
	
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	for i := range playerState.NegativeEnergy {
		targets = append(targets, Target{
			Type:     "negative_energy",
			ID:       playerState.NegativeEnergy[i],
			Location: "negative_energy",
		})
	}
	return targets
}

func (e *DiscardNegativeEnergyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Discard selected negative energy
	return nil
}

type ActivateByFlippingNegativeEnergyEffect struct {
	BaseEffect
	Count int
}

func (e *ActivateByFlippingNegativeEnergyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Need at least Count face-up negative energy cards
	// TODO: Count face-up negative energy
	return len(playerState.NegativeEnergy) >= e.Count
}

func (e *ActivateByFlippingNegativeEnergyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	// TODO: Only include face-up negative energy
	for i := range playerState.NegativeEnergy {
		targets = append(targets, Target{
			Type:     "negative_energy",
			ID:       playerState.NegativeEnergy[i],
			Location: "negative_energy",
		})
	}
	return targets
}

func (e *ActivateByFlippingNegativeEnergyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Flip selected negative energy and activate the attacker
	return nil
}

type HandSizePowerBoostEffect struct {
	BaseEffect
	CardsPerBoost int
	PowerBoost    int
}

func (e *HandSizePowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true // Always active
}

func (e *HandSizePowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *HandSizePowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This is a persistent effect that modifies power calculation
	// It should be checked when calculating power
	return nil
}

func (e *HandSizePowerBoostEffect) GetPowerBoost(game *GameContext) int {
	playerState := game.GetPlayerState(game.ActivePlayer)
	handSize := len(playerState.Hand)
	return (handSize / e.CardsPerBoost) * e.PowerBoost
}

type MainPhasePowerBoostEffect struct {
	BaseEffect
	PowerBoost int
	Duration   string // "turn" or "permanent"
}

func (e *MainPhasePowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Can activate during main phase if not already activated this turn
	return game.Phase == "main" && game.ActivePlayer == game.CurrentPlayer
}

func (e *MainPhasePowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// Targets self
	return []Target{{
		Type:     "friend",
		ID:       source.CardNo,
		Location: "self",
	}}
}

func (e *MainPhasePowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This would apply a temporary power boost
	// Implementation would track the boost and duration
	return nil
}