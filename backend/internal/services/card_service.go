package services

import (
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
)

type CardService struct{}

func NewCardService() *CardService {
	return &CardService{}
}

func (s *CardService) GetAllCards() ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	result := db.Order("card_no ASC").Find(&cards)
	return cards, result.Error
}

func (s *CardService) GetCardByNumber(cardNo string) (*models.Card, error) {
	var card models.Card
	db := database.GetDB()
	result := db.Where("card_no = ?", cardNo).First(&card)
	if result.Error != nil {
		return nil, result.Error
	}
	return &card, nil
}

func (s *CardService) GetCardsByType(cardType models.CardType) ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	result := db.Where("type = ?", cardType).Find(&cards)
	return cards, result.Error
}

func (s *CardService) GetCardsByColor(color models.CardColor) ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	result := db.Where("color = ?", color).Find(&cards)
	return cards, result.Error
}

func (s *CardService) SearchCards(query string) ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	searchQuery := "%" + query + "%"
	result := db.Where("name LIKE ? OR effect LIKE ?", searchQuery, searchQuery).Find(&cards)
	return cards, result.Error
}

func (s *CardService) GetPromoCards() ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	// Promo cards have (P) in their card number or is_promo flag
	result := db.Where("card_no LIKE '%(P)%' OR is_promo = ?", true).Order("card_no ASC").Find(&cards)
	return cards, result.Error
}

func (s *CardService) GetNonPromoCards() ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	// Non-promo cards don't have (P) in number, aren't parallel (-P), and don't have is_promo flag
	result := db.Where("card_no NOT LIKE '%(P)%' AND card_no NOT LIKE '%-P' AND (is_promo = ? OR is_promo IS NULL)", false).Order("card_no ASC").Find(&cards)
	return cards, result.Error
}

func (s *CardService) GetParallelCards() ([]models.Card, error) {
	var cards []models.Card
	db := database.GetDB()
	// Parallel cards end with -P
	result := db.Where("card_no LIKE '%-P'").Order("card_no ASC").Find(&cards)
	return cards, result.Error
}