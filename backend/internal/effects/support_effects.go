package effects

import (
	"github.com/your-org/card_game/backend/internal/models"
)

// Support card specific effects

// F-065 バードン - Main/Counter: +2000 power to one friend this turn
func NewBardonEffect() Effect {
	return &PowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】このターン中、自分のふれんど1体のパワー+2000。",
		},
		Amount:   2000,
		Duration: "turn",
		Target:   "one_my_friend",
	}
}

// F-066/F-066(P) 正志とくらげ坊 - Draw 2 cards. If you have くらげ坊, destroy 1 opponent's friend with 3000 power or less
func NewMasashiKurageboEffect() Effect {
	return &CompositeEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain,
			Description: "【メイン】自分はデッキから2枚ドローする。自分の場に「くらげ坊」がいるなら、さらにパワー3000以下の相手のふれんど1体を破壊する。",
		},
		Effects: []Effect{
			&DrawCardEffect{
				BaseEffect: BaseEffect{Trigger: TriggerMain},
				Count:      2,
			},
			&ConditionalDestroyEffect{
				BaseEffect: BaseEffect{Trigger: TriggerMain},
				MaxPower:   3000,
				Condition: func(game *GameContext, source *models.Card) bool {
					playerState := game.GetPlayerState(game.ActivePlayer)
					// Check if player has くらげ坊 on field
					for _, friend := range playerState.BattleArea {
						if friend.CardNo == "F-016" || friend.CardNo == "F-016 (P)" {
							return true
						}
					}
					return false
				},
			},
		},
	}
}

// F-067/F-067(P) 大好物！ - Main/Counter: Return 1 card from energy area to hand
func NewDaikoubutsuEffect() Effect {
	return &ReturnEnergyToHandEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】自分のエネルギーエリアのカード1枚を手札に戻せる。",
		},
	}
}

// F-068/F-068(P) デコレーション - Main/Counter: Destroy 1 field card
func NewDecorationEffect() Effect {
	return &DestroyFieldCardEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】相手の場のフィールドカード1枚を破壊する。",
		},
	}
}

// F-069/F-069(P) 特異点が開く扉 - Main/Counter: Destroy 1 opponent's friend with 5000 power or less
func NewTokuiTenEffect() Effect {
	return &DestroyFriendEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】パワー5000以下の相手のふれんど1体を破壊する。",
		},
		MaxPower:      5000,
		RequireTarget: true,
	}
}

// F-070/F-070(P) ブルードラゴン飛連蹴 - Main: Destroy 1 opponent's friend with 10000 power or less
func NewBlueDragonKickEffect() Effect {
	return &DestroyFriendEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain,
			Description: "【メイン】パワー10000以下の相手のふれんど1体を破壊する。",
		},
		MaxPower:      10000,
		RequireTarget: true,
	}
}

// F-071 絶対に裏切らない友達 - Main/Counter: Revive 1 friend from trash rested, destroy at end of turn
func NewZettaiUragiranaiFriendEffect() Effect {
	return &ReviveFriendEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】自分のトラッシュから、ふれんどカード1枚をコストを支払わずにレストして登場させる。この効果で登場したふれんどは、エンドフェイズに破壊される。",
		},
		EnterRested: true,
		DestroyAtEnd: true,
	}
}

// F-072 竜也とユピ - Main/Counter: Look at top 3 cards, add 1 to hand, discard rest. If you have ユピ, can return it to hand
func NewRyuyaYupiEffect() Effect {
	return &LookAndDrawEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】自分はデッキの上から3枚オープンする。その中の1枚を手札に加え、残りは破棄する。自分の場に「ユピ」がいるなら、さらに自分の「ユピ」1体を手札に戻せる。",
		},
		LookCount: 3,
		DrawCount: 1,
		AdditionalEffect: func(game *GameContext, source *models.Card) error {
			playerState := game.GetPlayerState(game.ActivePlayer)
			// Check if player has ユピ on field
			for pos, friend := range playerState.BattleArea {
				if friend.CardNo == "F-023" || friend.CardNo == "F-023 (P)" {
					// Allow returning ユピ to hand
					_ = pos
					return game.ReturnToHand(game.ActivePlayer, friend.CardNo)
				}
			}
			return nil
		},
	}
}

// F-073 古池ダイビング - Main/Counter: Return 1 friend to hand
func NewFuruikeDivingEffect() Effect {
	return &ReturnToHandEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain, // Can also be TriggerCounter
			Description: "【メイン/カウンター】ふれんど1体を手札に戻す。",
		},
		Scope:         ScopeTarget, // Can target any friend
		RequireTarget: true,
	}
}

// F-080/F-080(P) 謎の四人衆 - Main: Both players discard to 4 cards, then draw to 4 cards
func NewNazonoYoninEffect() Effect {
	return &HandResetEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerMain,
			Description: "【メイン】お互いは手札を4枚になるまで破棄する。その後、お互いはデッキから、手札が4枚になるまでドローする。",
		},
		TargetHandSize: 4,
		AffectBoth: true,
	}
}

// Custom effect types for support cards

type PowerBoostEffect struct {
	BaseEffect
	Amount   int
	Duration string // "turn" or "permanent"
	Target   string // "one_my_friend", "all_my_friends", etc.
}

func (e *PowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return e.GetTargets(game, source) != nil
}

func (e *PowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	for pos, friend := range playerState.BattleArea {
		targets = append(targets, Target{
			Type:     "friend",
			ID:       friend.CardNo,
			Location: pos,
		})
	}
	
	return targets
}

func (e *PowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) > 0 {
		return game.ModifyPower(game.ActivePlayer, targets[0].ID, e.Amount)
	}
	return nil
}

type CompositeEffect struct {
	BaseEffect
	Effects []Effect
}

func (e *CompositeEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// All sub-effects must be activatable
	for _, effect := range e.Effects {
		if !effect.CanActivate(game, source) {
			return false
		}
	}
	return true
}

func (e *CompositeEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// Collect targets from all sub-effects
	var allTargets []Target
	for _, effect := range e.Effects {
		targets := effect.GetTargets(game, source)
		allTargets = append(allTargets, targets...)
	}
	return allTargets
}

func (e *CompositeEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Apply each sub-effect in order
	for _, effect := range e.Effects {
		if err := effect.Apply(game, source, targets); err != nil {
			return err
		}
	}
	return nil
}

type ConditionalDestroyEffect struct {
	BaseEffect
	MaxPower  int
	Condition func(game *GameContext, source *models.Card) bool
}

func (e *ConditionalDestroyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return e.Condition(game, source)
}

func (e *ConditionalDestroyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	if !e.CanActivate(game, source) {
		return nil
	}
	
	var targets []Target
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	oppState := game.GetPlayerState(oppPlayer)
	
	for pos, friend := range oppState.BattleArea {
		if friend.Power <= e.MaxPower {
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: pos,
			})
		}
	}
	
	return targets
}

func (e *ConditionalDestroyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) > 0 {
		oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
		return game.DestroyFriend(oppPlayer, targets[0].ID)
	}
	return nil
}

type ReturnEnergyToHandEffect struct {
	BaseEffect
}

func (e *ReturnEnergyToHandEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	return len(playerState.EnergyArea) > 0
}

func (e *ReturnEnergyToHandEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	for i, energy := range playerState.EnergyArea {
		targets = append(targets, Target{
			Type:     "energy",
			ID:       energy.CardNo,
			Location: string(i),
		})
	}
	
	return targets
}

func (e *ReturnEnergyToHandEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement returning energy to hand
	return nil
}

type DestroyFieldCardEffect struct {
	BaseEffect
}

func (e *DestroyFieldCardEffect) CanActivate(game *GameContext, source *models.Card) bool {
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	oppState := game.GetPlayerState(oppPlayer)
	return oppState.FieldCard != nil
}

func (e *DestroyFieldCardEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	oppState := game.GetPlayerState(oppPlayer)
	
	if oppState.FieldCard != nil {
		return []Target{{
			Type:     "field",
			ID:       *oppState.FieldCard,
			Location: "field",
		}}
	}
	
	return nil
}

func (e *DestroyFieldCardEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	// TODO: Implement destroying field card
	_ = oppPlayer
	return nil
}

type ReviveFriendEffect struct {
	BaseEffect
	EnterRested  bool
	DestroyAtEnd bool
}

func (e *ReviveFriendEffect) CanActivate(game *GameContext, source *models.Card) bool {
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if there are friend cards in trash
	// TODO: Filter for friend cards
	return len(playerState.Trash) > 0
}

func (e *ReviveFriendEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	// TODO: Filter for friend cards
	for _, cardNo := range playerState.Trash {
		targets = append(targets, Target{
			Type:     "card",
			ID:       cardNo,
			Location: "trash",
		})
	}
	
	return targets
}

func (e *ReviveFriendEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement reviving friend from trash
	// Need to mark for destruction at end of turn if DestroyAtEnd is true
	return nil
}

type LookAndDrawEffect struct {
	BaseEffect
	LookCount        int
	DrawCount        int
	AdditionalEffect func(game *GameContext, source *models.Card) error
}

func (e *LookAndDrawEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *LookAndDrawEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	// TODO: Show top cards of deck for selection
	return nil
}

func (e *LookAndDrawEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement looking at top cards and selecting
	
	// Apply additional effect if present
	if e.AdditionalEffect != nil {
		return e.AdditionalEffect(game, source)
	}
	
	return nil
}

type HandResetEffect struct {
	BaseEffect
	TargetHandSize int
	AffectBoth     bool
}

func (e *HandResetEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true
}

func (e *HandResetEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *HandResetEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement hand reset for both players
	return nil
}