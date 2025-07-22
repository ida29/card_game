package utils

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
	"os"
	"strconv"
	"strings"
)

func LoadCardsFromFiles() error {
	// Load card data from JSON
	cardJSONData, err := loadCardJSON("data/card_data.json")
	if err != nil {
		return err
	}

	// Load card data from CSV
	cardCSVData, err := loadCardCSV("data/mememe_cards_parsed.csv")
	if err != nil {
		return err
	}

	// Merge and save cards
	cards := mergeCardData(cardJSONData, cardCSVData)
	
	db := database.GetDB()
	for _, card := range cards {
		// Check if card already exists
		var existingCard models.Card
		result := db.Where("card_no = ?", card.CardNo).First(&existingCard)
		if result.Error == nil {
			// Update existing card
			db.Model(&existingCard).Updates(&card)
		} else {
			// Create new card
			db.Create(&card)
		}
	}

	log.Printf("Loaded %d cards into database", len(cards))
	return nil
}

func loadCardJSON(filepath string) ([]models.CardJSON, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cards []models.CardJSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func loadCardCSV(filepath string) ([]models.CardCSV, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var cards []models.CardCSV
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		card := models.CardCSV{
			ID:              record[0],
			Name:            record[1],
			Type:            record[2],
			Owner:           record[3],
			Description:     record[4],
			Abilities:       record[5],
			Characteristics: record[6],
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func mergeCardData(jsonCards []models.CardJSON, csvCards []models.CardCSV) []models.Card {
	var cards []models.Card
	
	// Only process cards that have images (from JSON file)
	for _, jCard := range jsonCards {
		// Use the full number including (P) if present as card ID
		cardNo := jCard.Number
		baseCardNo := strings.Split(jCard.Number, " ")[0]
		
		card := models.Card{
			CardNo:         cardNo,
			Name:           jCard.Name,
			Type:           models.CardTypeFriend, // Most cards are friends
			ImageURL:       jCard.ImageURL,
			LocalImagePath: jCard.LocalImagePath,
		}
		
		// Extract card info from base ID  
		parseCardID(&card, baseCardNo)
		
		// Add version indicator based on number
		if strings.Contains(jCard.Number, "(P)") {
			card.Name = card.Name + " (プロモ)"
		} else if strings.Contains(jCard.LocalImagePath, "_P") || strings.Contains(jCard.LocalImagePath, "-P") {
			// Some promotional cards might not have (P) in number but have _P in filename
			if !strings.Contains(card.Name, "(プロモ)") {
				card.Name = card.Name + " (特別版)"
			}
		}
		
		// Try to find additional data from CSV using base card number
		for _, csvCard := range csvCards {
			// Extract number from CSV ID (e.g., C-009 -> 009, F-009 -> 009)
			csvParts := strings.Split(csvCard.ID, "-")
			csvNumber := ""
			if len(csvParts) >= 2 {
				csvNumber = csvParts[1]
			}
			
			// Extract number from card number (e.g., F-009 -> 009)
			cardParts := strings.Split(baseCardNo, "-")
			cardNumber := ""
			if len(cardParts) >= 2 {
				cardNumber = cardParts[1]
			}
			
			// Try to match by number only (ignoring prefix)
			if csvNumber == cardNumber && csvNumber != "" {
				// Found matching CSV data
				if csvCard.Description != "" {
					card.Effect = csvCard.Description
				}
				if csvCard.Characteristics != "" {
					card.FlavorText = csvCard.Characteristics
				}
				
				// Determine card type
				switch csvCard.Type {
				case "ふれんど":
					card.Type = models.CardTypeFriend
				case "サポート":
					card.Type = models.CardTypeSupport
				case "フィールド":
					card.Type = models.CardTypeField
				}
				break
			}
		}
		
		cards = append(cards, card)
	}

	return cards
}

func parseCardID(card *models.Card, id string) {
	// Parse format like "F-001", "F-013-P", etc.
	parts := strings.Split(id, "-")
	if len(parts) < 2 {
		return
	}

	// Extract number
	numStr := parts[1]
	if strings.Contains(numStr, "_") {
		numParts := strings.Split(numStr, "_")
		numStr = numParts[0]
	}

	num, _ := strconv.Atoi(numStr)

	// Determine color based on number range (example ranges)
	switch {
	case num >= 1 && num <= 22:
		card.Color = models.ColorRed
	case num >= 23 && num <= 32:
		card.Color = models.ColorBlue
	case num >= 33 && num <= 48:
		card.Color = models.ColorYellow
	case num >= 49 && num <= 64:
		card.Color = models.ColorGreen
	default:
		card.Color = models.ColorRed // Default
	}

	// Extract rarity from ID
	if strings.Contains(id, "_") {
		rarityPart := strings.Split(id, "_")[1]
		switch {
		case strings.HasPrefix(rarityPart, "SEC"):
			card.Rarity = models.RaritySEC
		case strings.HasPrefix(rarityPart, "SR"):
			card.Rarity = models.RaritySR
		case strings.HasPrefix(rarityPart, "R"):
			card.Rarity = models.RarityR
		case strings.HasPrefix(rarityPart, "U"):
			card.Rarity = models.RarityU
		default:
			card.Rarity = models.RarityC
		}
	} else {
		card.Rarity = models.RarityC // Default to common
	}

	// Set example costs and power (would be from actual card data)
	card.Cost = (num % 5) + 1
	if card.Type == models.CardTypeFriend {
		card.Power = (num % 8 + 2) * 1000
	}
}