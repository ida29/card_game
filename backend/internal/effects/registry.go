package effects

import (
	"sync"
)

var (
	globalRegistry *EffectRegistry
	once          sync.Once
)

// GetGlobalRegistry returns the global effect registry
func GetGlobalRegistry() *EffectRegistry {
	once.Do(func() {
		globalRegistry = NewEffectRegistry()
		InitializeEffects(globalRegistry)
	})
	return globalRegistry
}

// InitializeEffects registers all card effects
func InitializeEffects(registry *EffectRegistry) {
	// Friend card effects
	registry.Register("F-002", NewNamidabukuronEffect()) // なみだぶくろん - Main phase power boost
	
	registry.Register("F-003", NewFurafuraEffect()) // フラフラ - Power boost based on hand size
	
	registry.Register("F-004", NewTiranoEffect()) // ハシルシト - Can attack on turn played
	
	registry.Register("F-006", NewHiyakeratopusEffect())
	registry.Register("F-006 (P)", NewHiyakeratopusEffect())
	
	registry.Register("F-008", NewBoyEffect())
	registry.Register("F-008 (P)", NewBoyEffect())
	
	registry.Register("F-011", NewPochiEffect())
	registry.Register("F-011 (P)", NewPochiEffect())
	
	registry.Register("F-013", NewRukusoEffect())
	registry.Register("F-013 (P)", NewRukusoEffect())
	
	registry.Register("F-015", NewTiranoEffect())
	registry.Register("F-015 (P)", NewTiranoEffect())
	
	registry.Register("F-016", NewKurageboEffect())
	registry.Register("F-016 (P)", NewKurageboEffect())
	
	registry.Register("F-020", NewMarukaniEffect())
	registry.Register("F-020 (P)", NewMarukaniEffect())
	
	registry.Register("F-022", NewJohnnyEffect())
	registry.Register("F-022 (P)", NewJohnnyEffect())
	
	registry.Register("F-023", NewYupiEffect())
	registry.Register("F-023 (P)", NewYupiEffect())
	
	registry.Register("F-025", NewShimonEffect())
	registry.Register("F-025 (P)", NewShimonEffect())
	
	registry.Register("F-034", NewMegarokkoEffect())
	registry.Register("F-034 (P)", NewMegarokkoEffect())
	
	registry.Register("F-041", NewHayaoEffect())
	registry.Register("F-041 (P)", NewHayaoEffect())
	
	registry.Register("F-042", NewUkkiEffect())
	registry.Register("F-042 (P)", NewUkkiEffect())
	
	registry.Register("F-044", NewUkkiAttackEffect())
	registry.Register("F-044 (P)", NewUkkiAttackEffect())
	
	registry.Register("F-055", NewKo2Effect())
	registry.Register("F-055 (P)", NewKo2Effect())
	
	registry.Register("F-056", NewShiranEffect())
	registry.Register("F-056 (P)", NewShiranEffect())
	
	// Support card effects
	registry.Register("F-065", NewBardonEffect())
	registry.Register("F-065 (P)", NewBardonEffect())
	
	registry.Register("F-066", NewMasashiKurageboEffect())
	registry.Register("F-066 (P)", NewMasashiKurageboEffect())
	
	registry.Register("F-067", NewDaikoubutsuEffect())
	registry.Register("F-067 (P)", NewDaikoubutsuEffect())
	
	registry.Register("F-068", NewDecorationEffect())
	registry.Register("F-068 (P)", NewDecorationEffect())
	
	registry.Register("F-069", NewTokuiTenEffect())
	registry.Register("F-069 (P)", NewTokuiTenEffect())
	
	registry.Register("F-070", NewBlueDragonKickEffect())
	registry.Register("F-070 (P)", NewBlueDragonKickEffect())
	
	registry.Register("F-071", NewZettaiUragiranaiFriendEffect())
	
	registry.Register("F-072", NewRyuyaYupiEffect())
	
	registry.Register("F-073", NewFuruikeDivingEffect())
	
	registry.Register("F-080", NewNazonoYoninEffect())
	registry.Register("F-080 (P)", NewNazonoYoninEffect())
	
	// Field card effects
	registry.Register("F-089", NewMiharashidaiEffect())
	registry.Register("F-090", NewJinjaEffect())
	registry.Register("F-091", NewMasashiHouseEffect())
	registry.Register("F-092", NewFushigiKyoshitsuEffect())
	registry.Register("F-093", NewGakuenPoolEffect())
	registry.Register("F-094", NewKenkyujoEffect())
	registry.Register("F-095", NewTokumoUniversityEffect())
	registry.Register("F-096", NewTenchiKyukaiEffect())
	registry.Register("F-097", NewTokumoFestivalEffect())
	registry.Register("F-098", NewMoguraHouseEffect())
	registry.Register("F-099", NewGomisutebaEffect())
	
	// Transformed cards
	registry.Register("F-102", NewKurageboTransformEffect())
	
	// TODO: Add more card effects as they are discovered
}