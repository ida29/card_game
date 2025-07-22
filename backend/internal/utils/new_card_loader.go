package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
	"path/filepath"
	"strings"
)

// CompleteCardData represents the full card data from the new JSON
type CompleteCardData struct {
	Number      string      `json:"number"`
	Name        string      `json:"name"`
	Rarity      string      `json:"rarity"`
	Cost        CostData    `json:"cost"`
	Color       string      `json:"color"`
	Type        string      `json:"type"`
	Attribute   string      `json:"attribute,omitempty"`
	Emotion     string      `json:"emotion,omitempty"`
	Power       interface{} `json:"power,omitempty"`
	Ability     string      `json:"ability,omitempty"`
	FlavorText  string      `json:"flavor_text,omitempty"`
	ImageURL    string      `json:"image_url"`
	IsPromo     bool        `json:"is_promo"`
	IsParallel  bool        `json:"is_parallel"`
	Height      string      `json:"height,omitempty"`
	Weight      string      `json:"weight,omitempty"`
}

// CostData represents the cost breakdown
type CostData struct {
	Total     int `json:"total"`
	Red       int `json:"red"`
	Blue      int `json:"blue"`
	Yellow    int `json:"yellow"`
	Green     int `json:"green"`
	Colorless int `json:"colorless"`
}

// LoadCompleteCardData loads all card data from the new complete JSON file
func LoadCompleteCardData() error {
	log.Println("Loading complete card data from mememe_cards_complete.json...")
	
	// Read the JSON file
	data, err := ioutil.ReadFile("data/mememe_cards_complete.json")
	if err != nil {
		return fmt.Errorf("failed to read card data: %v", err)
	}
	
	var cards []CompleteCardData
	if err := json.Unmarshal(data, &cards); err != nil {
		return fmt.Errorf("failed to parse card data: %v", err)
	}
	
	db := database.GetDB()
	
	// Clear existing cards
	if err := db.Exec("DELETE FROM cards").Error; err != nil {
		log.Printf("Warning: Failed to clear existing cards: %v", err)
	}
	
	// Process each card
	successCount := 0
	for _, cardData := range cards {
		card := models.Card{
			CardNo:         cardData.Number,
			Name:           cardData.Name,
			Cost:           cardData.Cost.Total,
			CostRed:        cardData.Cost.Red,
			CostBlue:       cardData.Cost.Blue,
			CostYellow:     cardData.Cost.Yellow,
			CostGreen:      cardData.Cost.Green,
			CostColorless:  cardData.Cost.Colorless,
			Power:          0, // Will be set below
			Effect:         cardData.Ability,
			FlavorText:     cardData.FlavorText,
			ImageURL:       cardData.ImageURL,
			IsPromo:        cardData.IsPromo,
		}
		
		// Set card type
		switch cardData.Type {
		case "ふれんど":
			card.Type = models.CardTypeFriend
		case "サポート":
			card.Type = models.CardTypeSupport
		case "フィールド":
			card.Type = models.CardTypeField
		default:
			card.Type = models.CardTypeFriend
		}
		
		// Set color
		switch cardData.Color {
		case "赤":
			card.Color = models.ColorRed
		case "青":
			card.Color = models.ColorBlue
		case "黄":
			card.Color = models.ColorYellow
		case "緑":
			card.Color = models.ColorGreen
		default:
			card.Color = models.ColorNone
		}
		
		// Set rarity
		rarityStr := cardData.Rarity
		if strings.Contains(rarityStr, "パラレル") {
			// Handle parallel cards
			if strings.Contains(rarityStr, "SEC") {
				card.Rarity = "SEC-P"
			} else if strings.Contains(rarityStr, "SR") {
				card.Rarity = "SR-P"
			} else if strings.Contains(rarityStr, "R") {
				card.Rarity = "R-P"
			} else if strings.Contains(rarityStr, "U") {
				card.Rarity = "U-P"
			} else if strings.Contains(rarityStr, "C") {
				card.Rarity = "C-P"
			}
		} else {
			// Regular rarities
			switch rarityStr {
			case "C":
				card.Rarity = models.RarityC
			case "U":
				card.Rarity = models.RarityU
			case "R":
				card.Rarity = models.RarityR
			case "SR":
				card.Rarity = models.RaritySR
			case "SEC":
				card.Rarity = models.RaritySEC
			default:
				card.Rarity = models.RarityC
			}
		}
		
		// Set power value - handle both int and string values
		switch v := cardData.Power.(type) {
		case int:
			card.Power = v
		case float64:
			card.Power = int(v)
		case string:
			// For support/field cards with "ー" power
			card.Power = 0
		default:
			card.Power = 0
		}
		
		// Set local image path
		filename := filepath.Base(cardData.ImageURL)
		card.LocalImagePath = "card_images/" + filename
		
		// Add energy icons (attribute and emotion)
		if cardData.Attribute != "" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Attribute)
		}
		if cardData.Emotion != "" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Emotion)
		}
		
		// Add physical stats to flavor text if available
		if cardData.Height != "" || cardData.Weight != "" {
			stats := []string{}
			if cardData.Height != "" {
				stats = append(stats, "身長: " + cardData.Height)
			}
			if cardData.Weight != "" {
				stats = append(stats, "体重: " + cardData.Weight)
			}
			if len(stats) > 0 && card.FlavorText != "" {
				card.FlavorText = card.FlavorText + "\n" + strings.Join(stats, " / ")
			}
		}
		
		// Check for counter abilities
		if strings.Contains(card.Effect, "【カウンター】") {
			card.IsCounter = true
		}
		if strings.Contains(card.Effect, "【メイン/カウンター】") {
			card.IsMainCounter = true
		}
		
		// Save to database
		if err := db.Create(&card).Error; err != nil {
			log.Printf("Failed to create card %s: %v", cardData.Number, err)
		} else {
			successCount++
		}
		
		// Also save cost breakdown to a separate table if needed
		// For now, we'll store it in the card's effect or a JSON field
	}
	
	log.Printf("Successfully loaded %d/%d cards", successCount, len(cards))
	return nil
}