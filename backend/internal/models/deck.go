package models

import (
	"gorm.io/gorm"
)

type Deck struct {
	gorm.Model
	Name     string     `json:"name"`
	UserID   uint       `json:"user_id"`
	Cards    []DeckCard `json:"cards"`
	IsActive bool       `json:"is_active"`
}

type DeckCard struct {
	gorm.Model
	DeckID   uint   `json:"deck_id"`
	CardNo   string `json:"card_no"`
	Quantity int    `json:"quantity"`
	Card     Card   `json:"card" gorm:"foreignKey:CardNo;references:CardNo"`
}

func (d *Deck) Validate() error {
	cardCount := 0
	cardQuantities := make(map[string]int)

	for _, dc := range d.Cards {
		cardCount += dc.Quantity
		cardQuantities[dc.CardNo] += dc.Quantity

		if cardQuantities[dc.CardNo] > 4 {
			return &ValidationError{Message: "同じカードは4枚までしか入れることができません"}
		}
	}

	if cardCount != 50 {
		return &ValidationError{Message: "デッキは正確に50枚でなければなりません"}
	}

	return nil
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}