package models

import (
	"gorm.io/gorm"
)

type CardType string
type CardColor string
type CardRarity string

const (
	CardTypeFriend  CardType = "ふれんど"
	CardTypeSupport CardType = "サポート"
	CardTypeField   CardType = "フィールド"
)

const (
	ColorRed    CardColor = "赤"
	ColorBlue   CardColor = "青"
	ColorYellow CardColor = "黄"
	ColorGreen  CardColor = "緑"
	ColorNone   CardColor = "無"
)

const (
	RarityC   CardRarity = "C"
	RarityU   CardRarity = "U"
	RarityR   CardRarity = "R"
	RaritySR  CardRarity = "SR"
	RaritySEC CardRarity = "SEC"
)

type Card struct {
	gorm.Model
	CardNo          string     `json:"card_no" gorm:"uniqueIndex:idx_card_no_rarity"`
	Name            string     `json:"name"`
	Type            CardType   `json:"type"`
	Color           CardColor  `json:"color"`
	Cost            int        `json:"cost"`
	CostRed         int        `json:"cost_red"`
	CostBlue        int        `json:"cost_blue"`
	CostYellow      int        `json:"cost_yellow"`
	CostGreen       int        `json:"cost_green"`
	CostColorless   int        `json:"cost_colorless"`
	Power           int        `json:"power,omitempty"`
	Rarity          CardRarity `json:"rarity" gorm:"uniqueIndex:idx_card_no_rarity"`
	Effect          string     `json:"effect,omitempty"`
	FlavorText      string     `json:"flavor_text,omitempty"`
	ImageURL        string     `json:"image_url"`
	LocalImagePath  string     `json:"local_image_path"`
	EnergyIcons     []string   `json:"energy_icons,omitempty" gorm:"serializer:json"`
	IsCounter       bool       `json:"is_counter"`
	IsMainCounter   bool       `json:"is_main_counter"`
	IsPromo         bool       `json:"is_promo"` // True if card number contains (P)
}

type CardCSV struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Owner           string `json:"owner"`
	Description     string `json:"description"`
	Abilities       string `json:"abilities"`
	Characteristics string `json:"characteristics"`
}

type CardJSON struct {
	Name           string `json:"name"`
	Number         string `json:"number"`
	Type           string `json:"type"`
	ImageURL       string `json:"image_url"`
	LocalImagePath string `json:"local_image_path"`
}