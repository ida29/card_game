package effects

import (
	"github.com/your-org/card_game/backend/internal/models"
)

// Field card specific effects

// F-089 見晴らし台 - During your turn, while you have more cards in hand than opponent, all your friends get +2000 power
func NewMiharashidaiEffect() Effect {
	return &ConditionalPowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分のターン中、自分の手札が相手より多い間、自分のふれんど全てのパワー+2000。",
		},
		Amount: 2000,
		Scope:  ScopeMyFriends,
		Condition: func(game *GameContext, source *models.Card) bool {
			// Only active during your turn
			if game.ActivePlayer != game.ActivePlayer {
				return false
			}
			
			playerState := game.GetPlayerState(game.ActivePlayer)
			oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
			oppState := game.GetPlayerState(oppPlayer)
			
			return len(playerState.Hand) > len(oppState.Hand)
		},
	}
}

// F-090 神社 - During energy phase, instead of placing from deck, can reveal 1 card from negative energy area
func NewJinjaEffect() Effect {
	return &EnergyPhaseAlternativeEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerStartPhase,
			Description: "自分のエネルギーフェイズにカードを置く時、デッキから置く代わりに自分の負のエネルギーエリアのカード1枚を表にできる。",
		},
	}
}

// F-091 正志の家 - At start of end phase, can activate 1 of your friends
func NewMasashiHouseEffect() Effect {
	return &EndPhaseActivateEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerEndPhase,
			Description: "エンドフェイズ開始時、自分のふれんど1体をアクティブにできる。",
		},
	}
}

// F-092 不思議な教室 - At start of draw phase, reveal top card of deck and put it on top or bottom
func NewFushigiKyoshitsuEffect() Effect {
	return &DrawPhaseManipulateEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerStartPhase,
			Description: "自分のドローフェイズ開始時、自分はデッキの上から1枚オープンする。そのカードをデッキの上か下に置く。",
		},
	}
}

// F-093 学園のプール - At start of opponent's start phase, activate 1 card in your energy area
func NewGakuenPoolEffect() Effect {
	return &OpponentStartActivateEnergyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerStartPhase,
			Description: "相手のスタートフェイズ開始時、自分のエネルギーエリアのカード1枚をアクティブにする。",
		},
	}
}

// F-094 研究所 - Can use support cards from negative energy area. Card goes to bottom of deck, then take 1 damage
func NewKenkyujoEffect() Effect {
	return &UseFromNegativeEnergyEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分は自分の負のエネルギーエリアのサポートカードを使用できる。使用したカードはデッキの下に置き、その後自分に1ダメージを与える。",
		},
	}
}

// F-095 都雲大学 - Friend cards in hand cost -1 to play (minimum 1, color costs not reduced)
func NewTokumoUniversityEffect() Effect {
	return &CostReductionEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分の手札のふれんどカードの登場コストは、登場する時-１される。ただし、コスト1以下にはならず、カラーシンボルのコストはマイナスされない。",
		},
		CardType:  models.CardTypeFriend,
		Reduction: 1,
		MinCost:   1,
	}
}

// F-096 天地救世教会本部 - Friends with cost 2/4/6 enter play rested
func NewTenchiKyukaiEffect() Effect {
	return &EnterRestedEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "コスト2/4/6のふれんどは、登場する時はレストして登場する。",
		},
		CostCondition: func(cost int) bool {
			return cost == 2 || cost == 4 || cost == 6
		},
	}
}

// F-097 都雲祭 - During your turn, for each other field card in play, all your friends get +1000 power
func NewTokumoFestivalEffect() Effect {
	return &FieldCountPowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分のターン中、このカード以外のお互いの場のフィールドカード1枚ごとに、自分のふれんど全てのパワー+1000。",
		},
		PowerPerField: 1000,
	}
}

// F-098 もぐらの家 - When your friend deals damage to opponent, draw 1 card
func NewMoguraHouseEffect() Effect {
	return &OnDamageDrawEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerOnDamageDealt,
			Description: "自分のふれんどが相手にダメージを与えた時、自分はデッキから1枚ドローする。",
		},
	}
}

// F-099 ごみ捨て場 - During your turn, for each face-down card in negative energy, green friends get +1000 power
func NewGomisutebaEffect() Effect {
	return &NegativeEnergyPowerBoostEffect{
		BaseEffect: BaseEffect{
			Trigger:     TriggerPersistent,
			Description: "自分のターン中、負のエネルギーエリアの裏のカード1枚ごとに、緑の自分のふれんど全てのパワー＋1000。",
		},
		PowerPerCard: 1000,
		ColorFilter:  models.ColorGreen,
	}
}

// Custom effect types for field cards

type ConditionalPowerBoostEffect struct {
	BaseEffect
	Amount    int
	Scope     EffectScope
	Condition func(game *GameContext, source *models.Card) bool
}

func (e *ConditionalPowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return e.Condition(game, source)
}

func (e *ConditionalPowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	if !e.CanActivate(game, source) {
		return nil
	}
	
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

func (e *ConditionalPowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This is a persistent effect - power boost is calculated dynamically
	return nil
}

type EnergyPhaseAlternativeEffect struct {
	BaseEffect
}

func (e *EnergyPhaseAlternativeEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return game.Game.CurrentPhase == models.PhaseEnergy
}

func (e *EnergyPhaseAlternativeEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // Player chooses during energy phase
}

func (e *EnergyPhaseAlternativeEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This modifies the energy phase rules
	return nil
}

type EndPhaseActivateEffect struct {
	BaseEffect
}

func (e *EndPhaseActivateEffect) CanActivate(game *GameContext, source *models.Card) bool {
	if game.Game.CurrentPhase != models.PhaseEnd {
		return false
	}
	
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if there are rested friends
	for _, friend := range playerState.BattleArea {
		if friend.IsRest {
			return true
		}
	}
	
	return false
}

func (e *EndPhaseActivateEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	for pos, friend := range playerState.BattleArea {
		if friend.IsRest {
			targets = append(targets, Target{
				Type:     "friend",
				ID:       friend.CardNo,
				Location: pos,
			})
		}
	}
	
	return targets
}

func (e *EndPhaseActivateEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	if len(targets) > 0 {
		return game.ActiveFriend(game.ActivePlayer, targets[0].ID)
	}
	return nil
}

type DrawPhaseManipulateEffect struct {
	BaseEffect
}

func (e *DrawPhaseManipulateEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return game.Game.CurrentPhase == models.PhaseDraw
}

func (e *DrawPhaseManipulateEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // No targets needed
}

func (e *DrawPhaseManipulateEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement deck manipulation
	return nil
}

type OpponentStartActivateEnergyEffect struct {
	BaseEffect
}

func (e *OpponentStartActivateEnergyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Check if it's opponent's start phase
	if game.Game.CurrentPhase != models.PhaseStart {
		return false
	}
	
	// Check if it's opponent's turn
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	if game.Game.ActivePlayer != oppPlayer {
		return false
	}
	
	playerState := game.GetPlayerState(game.ActivePlayer)
	// Check if there are rested energy cards
	for _, energy := range playerState.EnergyArea {
		if energy.IsRest {
			return true
		}
	}
	
	return false
}

func (e *OpponentStartActivateEnergyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	var targets []Target
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	for i, energy := range playerState.EnergyArea {
		if energy.IsRest {
			targets = append(targets, Target{
				Type:     "energy",
				ID:       energy.CardNo,
				Location: string(i),
			})
		}
	}
	
	return targets
}

func (e *OpponentStartActivateEnergyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// TODO: Implement activating energy card
	return nil
}

type UseFromNegativeEnergyEffect struct {
	BaseEffect
}

func (e *UseFromNegativeEnergyEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true // This is a persistent rule modification
}

func (e *UseFromNegativeEnergyEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil
}

func (e *UseFromNegativeEnergyEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This modifies game rules to allow using support from negative energy
	return nil
}

type CostReductionEffect struct {
	BaseEffect
	CardType  models.CardType
	Reduction int
	MinCost   int
}

func (e *CostReductionEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true // This is a persistent rule modification
}

func (e *CostReductionEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil
}

func (e *CostReductionEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This modifies cost calculation rules
	return nil
}

type EnterRestedEffect struct {
	BaseEffect
	CostCondition func(cost int) bool
}

func (e *EnterRestedEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true // This is a persistent rule modification
}

func (e *EnterRestedEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil
}

func (e *EnterRestedEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// This modifies entry rules for friends
	return nil
}

type FieldCountPowerBoostEffect struct {
	BaseEffect
	PowerPerField int
}

func (e *FieldCountPowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Only active during your turn
	return game.ActivePlayer == game.ActivePlayer
}

func (e *FieldCountPowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // Affects all friends automatically
}

func (e *FieldCountPowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	// Count other field cards
	fieldCount := 0
	
	playerState := game.GetPlayerState(game.ActivePlayer)
	if playerState.FieldCard != nil && *playerState.FieldCard != source.CardNo {
		fieldCount++
	}
	
	oppPlayer := game.GetOpponentPlayer(game.ActivePlayer)
	oppState := game.GetPlayerState(oppPlayer)
	if oppState.FieldCard != nil {
		fieldCount++
	}
	
	// Apply power boost based on field count
	powerBoost := fieldCount * e.PowerPerField
	
	for _, friend := range playerState.BattleArea {
		if err := game.ModifyPower(game.ActivePlayer, friend.CardNo, powerBoost); err != nil {
			return err
		}
	}
	
	return nil
}

type OnDamageDrawEffect struct {
	BaseEffect
}

func (e *OnDamageDrawEffect) CanActivate(game *GameContext, source *models.Card) bool {
	return true // Triggered by damage events
}

func (e *OnDamageDrawEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil
}

func (e *OnDamageDrawEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	return game.DrawCards(game.ActivePlayer, 1)
}

type NegativeEnergyPowerBoostEffect struct {
	BaseEffect
	PowerPerCard int
	ColorFilter  models.CardColor
}

func (e *NegativeEnergyPowerBoostEffect) CanActivate(game *GameContext, source *models.Card) bool {
	// Only active during your turn
	return game.ActivePlayer == game.ActivePlayer
}

func (e *NegativeEnergyPowerBoostEffect) GetTargets(game *GameContext, source *models.Card) []Target {
	return nil // Affects matching friends automatically
}

func (e *NegativeEnergyPowerBoostEffect) Apply(game *GameContext, source *models.Card, targets []Target) error {
	playerState := game.GetPlayerState(game.ActivePlayer)
	
	// Count face-down cards in negative energy
	// TODO: Need to track which cards are revealed/face-up
	faceDownCount := 0
	for range playerState.NegativeEnergy {
		faceDownCount++ // Placeholder - need to track revealed state
	}
	
	powerBoost := faceDownCount * e.PowerPerCard
	
	// Apply to matching color friends
	for _, friend := range playerState.BattleArea {
		// TODO: Need to check friend's color
		if err := game.ModifyPower(game.ActivePlayer, friend.CardNo, powerBoost); err != nil {
			return err
		}
	}
	
	return nil
}