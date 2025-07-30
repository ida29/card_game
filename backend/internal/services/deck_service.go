package services

import (
	"errors"
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"

	"gorm.io/gorm"
)

type DeckService struct{}

func NewDeckService() *DeckService {
	return &DeckService{}
}

func (s *DeckService) CreateDeck(deck *models.Deck) error {
	// Skip validation for new empty decks
	if len(deck.Cards) > 0 {
		if err := deck.Validate(); err != nil {
			return err
		}
	}

	db := database.GetDB()
	return db.Create(deck).Error
}

func (s *DeckService) GetDeck(deckID uint) (*models.Deck, error) {
	var deck models.Deck
	db := database.GetDB()
	
	// Preload Cards with a condition to exclude soft-deleted records
	result := db.Preload("Cards", "deleted_at IS NULL").Preload("Cards.Card").First(&deck, deckID)
	if result.Error != nil {
		return nil, result.Error
	}
	
	log.Printf("GetDeck returned deck ID %d with %d cards", deck.ID, len(deck.Cards))
	for i, card := range deck.Cards {
		log.Printf("  Card %d: ID=%d, CardNo=%s, Quantity=%d", i, card.ID, card.CardNo, card.Quantity)
	}
	
	return &deck, nil
}

func (s *DeckService) GetUserDecks(userID uint) ([]models.Deck, error) {
	var decks []models.Deck
	db := database.GetDB()
	
	result := db.Preload("Cards", "deleted_at IS NULL").Preload("Cards.Card").Where("user_id = ?", userID).Find(&decks)
	return decks, result.Error
}

func (s *DeckService) UpdateDeck(deck *models.Deck) error {
	// Log incoming deck data
	log.Printf("UpdateDeck called for deck ID %d with %d cards, main_card_no=%s", deck.ID, len(deck.Cards), deck.MainCardNo)
	for i, card := range deck.Cards {
		log.Printf("  Card %d: CardNo=%s, Quantity=%d", i, card.CardNo, card.Quantity)
	}

	// Only validate if the deck has cards (allow saving work-in-progress decks)
	// The frontend will handle the validation for completeness
	if len(deck.Cards) > 0 {
		// Only validate card quantity limits, not total count
		cardQuantities := make(map[string]int)
		for _, dc := range deck.Cards {
			cardQuantities[dc.CardNo] += dc.Quantity
			if cardQuantities[dc.CardNo] > 4 {
				return errors.New("同じカードは4枚までしか入れることができません")
			}
		}
	}

	db := database.GetDB()
	
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete existing deck cards (soft delete)
		var deleteResult *gorm.DB
		deleteResult = tx.Where("deck_id = ?", deck.ID).Delete(&models.DeckCard{})
		if deleteResult.Error != nil {
			return deleteResult.Error
		}
		log.Printf("Deleted %d existing deck cards", deleteResult.RowsAffected)
		
		// Save new deck cards
		for i := range deck.Cards {
			// Create a new DeckCard without the ID to ensure GORM creates new records
			newCard := models.DeckCard{
				DeckID:   deck.ID,
				CardNo:   deck.Cards[i].CardNo,
				Quantity: deck.Cards[i].Quantity,
			}
			if err := tx.Create(&newCard).Error; err != nil {
				return err
			}
		}
		log.Printf("Created %d new deck cards", len(deck.Cards))
		
		// Update deck info including main_card_no
		return tx.Model(deck).Updates(map[string]interface{}{
			"name":         deck.Name,
			"is_active":    deck.IsActive,
			"main_card_no": deck.MainCardNo,
		}).Error
	})
}

func (s *DeckService) DeleteDeck(deckID uint) error {
	db := database.GetDB()
	
	return db.Transaction(func(tx *gorm.DB) error {
		// Delete deck cards first
		if err := tx.Where("deck_id = ?", deckID).Delete(&models.DeckCard{}).Error; err != nil {
			return err
		}
		
		// Delete deck
		return tx.Delete(&models.Deck{}, deckID).Error
	})
}

func (s *DeckService) ValidateDeckCards(cards []models.DeckCard) error {
	cardCount := 0
	cardQuantities := make(map[string]int)

	for _, dc := range cards {
		cardCount += dc.Quantity
		cardQuantities[dc.CardNo] += dc.Quantity

		if cardQuantities[dc.CardNo] > 4 {
			return errors.New("同じカードは4枚までしか入れることができません")
		}
	}

	if cardCount != 50 {
		return errors.New("デッキは正確に50枚でなければなりません")
	}

	return nil
}