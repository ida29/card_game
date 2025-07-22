package models

import (
	"gorm.io/gorm"
	"time"
)

type GamePhase string
type GameStatus string

const (
	PhaseStart  GamePhase = "start"
	PhaseDraw   GamePhase = "draw"
	PhaseEnergy GamePhase = "energy"
	PhaseMain   GamePhase = "main"
	PhaseEnd    GamePhase = "end"
)

const (
	StatusWaiting  GameStatus = "waiting"
	StatusPlaying  GameStatus = "playing"
	StatusFinished GameStatus = "finished"
)

type Game struct {
	gorm.Model
	GameID        string       `json:"game_id" gorm:"uniqueIndex"`
	Player1ID     uint         `json:"player1_id"`
	Player2ID     uint         `json:"player2_id"`
	CurrentTurn   int          `json:"current_turn"`
	CurrentPhase  GamePhase    `json:"current_phase"`
	ActivePlayer  int          `json:"active_player"`
	Status        GameStatus   `json:"status"`
	WinnerID      *uint        `json:"winner_id,omitempty"`
	StartedAt     *time.Time   `json:"started_at,omitempty"`
	FinishedAt    *time.Time   `json:"finished_at,omitempty"`
	GameState     *GameState   `json:"game_state,omitempty" gorm:"serializer:json"`
}

type GameState struct {
	Player1State PlayerState `json:"player1_state"`
	Player2State PlayerState `json:"player2_state"`
}

type PlayerState struct {
	Deck             []string          `json:"deck"`
	Hand             []string          `json:"hand"`
	BattleArea       map[string]Friend `json:"battle_area"`
	EnergyArea       []EnergyCard      `json:"energy_area"`
	NegativeEnergy   []string          `json:"negative_energy"`
	Trash            []string          `json:"trash"`
	FieldCard        *string           `json:"field_card,omitempty"`
}

type Friend struct {
	CardNo     string `json:"card_no"`
	Power      int    `json:"power"`
	IsRest     bool   `json:"is_rest"`
	TurnPlayed int    `json:"turn_played"`
}

type EnergyCard struct {
	CardNo string `json:"card_no"`
	Color  CardColor `json:"color"`
	IsRest bool   `json:"is_rest"`
}