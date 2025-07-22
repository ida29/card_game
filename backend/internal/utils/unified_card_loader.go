package utils

import (
	"encoding/json"
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
	"os"
	"strings"
)

// Different JSON structures from various sources
type UnifiedCardData interface{}

func LoadAllCardData() error {
	log.Println("Loading all card data from various sources...")
	
	// Load image data first
	imageData, err := loadCardJSON("data/card_data.json")
	if err != nil {
		log.Printf("Warning: Failed to load image data: %v", err)
	}
	
	// Create image map
	imageMap := make(map[string]models.CardJSON)
	for _, img := range imageData {
		cardNo := strings.Split(img.Number, " ")[0]
		imageMap[cardNo] = img
	}
	
	// Process each JSON file with appropriate parser
	dataFiles := map[string]func(string, map[string]models.CardJSON) error{
		"data/mememe_cards_f001_f020.json": loadFirstFormat,
		"data/cards_F021_F040.json": loadSecondFormat,
		"data/cards_F041_F060.json": loadSecondFormat,
		"data/cards_F061_F080.json": loadSecondFormat,
		"data/cards_F081_F102.json": loadSecondFormat,
	}
	
	for file, loader := range dataFiles {
		if err := loader(file, imageMap); err != nil {
			log.Printf("Error loading %s: %v", file, err)
		}
	}
	
	// Also load promo cards from card_data.json
	if err := loadPromoCards(imageMap); err != nil {
		log.Printf("Error loading promo cards: %v", err)
	}
	
	// Also load parallel cards from card_data.json
	if err := loadParallelCards(imageMap); err != nil {
		log.Printf("Error loading parallel cards: %v", err)
	}
	
	// Count total cards
	var count int64
	database.GetDB().Model(&models.Card{}).Count(&count)
	log.Printf("Total cards in database: %d", count)
	
	return nil
}

// First format handler (F-001 to F-020)
func loadFirstFormat(filepath string, imageMap map[string]models.CardJSON) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	var data struct {
		Cards []struct {
			CardNumber string `json:"card_number"`
			Name       string `json:"name"`
			CardType   string `json:"card_type"`
			Rarity     string `json:"rarity"`
			Cost       int    `json:"cost"`
			CostBreakdown struct {
				Colorless int `json:"colorless"`
				Red       int `json:"red"`
				Blue      int `json:"blue"`
				Yellow    int `json:"yellow"`
				Green     int `json:"green"`
			} `json:"cost_breakdown"`
			Color      string  `json:"color"`
			Power      int     `json:"power"`
			Attribute  string  `json:"attribute"`
			Emotion    string  `json:"emotion"`
			Effect     string  `json:"effect"`
			FlavorText string  `json:"flavor_text"`
			Height     string  `json:"height"`
			Weight     string  `json:"weight"`
		} `json:"cards"`
	}
	
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	
	db := database.GetDB()
	
	for _, cardData := range data.Cards {
		card := models.Card{
			CardNo:     cardData.CardNumber,
			Name:       cardData.Name,
			Cost:       cardData.Cost,
			Power:      cardData.Power,
			Effect:     cardData.Effect,
			FlavorText: cardData.FlavorText,
		}
		
		// Set type
		switch cardData.CardType {
		case "ふれんど":
			card.Type = models.CardTypeFriend
		case "サポート":
			card.Type = models.CardTypeSupport
		case "フィールド":
			card.Type = models.CardTypeField
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
		case "無色":
			card.Color = models.ColorRed // Default to red for neutral cards
		}
		
		// Set rarity
		switch cardData.Rarity {
		case "C":
			card.Rarity = models.RarityC
		case "U":
			card.Rarity = models.RarityU
		case "R":
			card.Rarity = models.RarityR
		case "SR":
			card.Rarity = models.RaritySR
		case "SEC", "UR":
			card.Rarity = models.RaritySEC
		}
		
		// Add physical stats
		if cardData.Height != "" || cardData.Weight != "" {
			stats := []string{}
			if cardData.Height != "" {
				stats = append(stats, "身長: " + cardData.Height)
			}
			if cardData.Weight != "" {
				stats = append(stats, "体重: " + cardData.Weight)
			}
			if len(stats) > 0 {
				card.FlavorText = card.FlavorText + "\n" + strings.Join(stats, " / ")
			}
		}
		
		// Add attributes
		if cardData.Attribute != "" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Attribute)
		}
		if cardData.Emotion != "" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Emotion)
		}
		
		// Check for counter abilities
		if strings.Contains(cardData.Effect, "【カウンター】") {
			card.IsCounter = true
		}
		if strings.Contains(cardData.Effect, "【メイン/カウンター】") {
			card.IsMainCounter = true
		}
		
		// Add image data
		if imgData, exists := imageMap[cardData.CardNumber]; exists {
			card.ImageURL = imgData.ImageURL
			card.LocalImagePath = imgData.LocalImagePath
		}
		
		// Check if it's a promo card
		if strings.Contains(cardData.CardNumber, "(P)") {
			card.IsPromo = true
		}
		
		// Save to database
		var existingCard models.Card
		result := db.Where("card_no = ? AND rarity = ?", card.CardNo, card.Rarity).First(&existingCard)
		if result.Error == nil {
			// Skip if already exists with same rarity
			continue
		} else {
			db.Create(&card)
		}
	}
	
	return nil
}

// Second format handler (F-021 onwards)
func loadSecondFormat(filepath string, imageMap map[string]models.CardJSON) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	var data struct {
		Cards []struct {
			Number     string `json:"number"`
			Name       string `json:"name"`
			Type       string `json:"type"`
			Rarity     string `json:"rarity"`
			Cost       struct {
				Total     int `json:"total"`
				Colored   int `json:"colored"`
				Colorless int `json:"colorless"`
				Red       int `json:"red,omitempty"`
				Blue      int `json:"blue,omitempty"`
				Yellow    int `json:"yellow,omitempty"`
				Green     int `json:"green,omitempty"`
				Multi     int `json:"multi,omitempty"`
			} `json:"cost"`
			Power      *int   `json:"power"`
			Color      string `json:"color"`
			Attribute  string `json:"attribute"`
			Emotion    string `json:"emotion"`
			Effect     string `json:"effect"`
			FlavorText string `json:"flavorText"`
			Height     string `json:"height"`
			Weight     string `json:"weight"`
		} `json:"cards"`
	}
	
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}
	
	db := database.GetDB()
	
	for _, cardData := range data.Cards {
		card := models.Card{
			CardNo:     cardData.Number,
			Name:       cardData.Name,
			Cost:       cardData.Cost.Total,
			Effect:     cardData.Effect,
			FlavorText: cardData.FlavorText,
		}
		
		if cardData.Power != nil {
			card.Power = *cardData.Power
		}
		
		// Set type
		switch cardData.Type {
		case "Friend", "ふれんど":
			card.Type = models.CardTypeFriend
		case "Support", "サポート":
			card.Type = models.CardTypeSupport
		case "Field", "フィールド":
			card.Type = models.CardTypeField
		case "Item":
			card.Type = models.CardTypeSupport // Treat Item as Support for now
		}
		
		// Set color
		switch cardData.Color {
		case "Red", "赤":
			card.Color = models.ColorRed
		case "Blue", "青":
			card.Color = models.ColorBlue
		case "Yellow", "黄":
			card.Color = models.ColorYellow
		case "Green", "緑":
			card.Color = models.ColorGreen
		case "Neutral", "無色", "Multi":
			card.Color = models.ColorRed // Default to red for neutral/multi cards
		}
		
		// Set rarity
		switch cardData.Rarity {
		case "C":
			card.Rarity = models.RarityC
		case "U":
			card.Rarity = models.RarityU
		case "R":
			card.Rarity = models.RarityR
		case "SR":
			card.Rarity = models.RaritySR
		case "SEC", "UR":
			card.Rarity = models.RaritySEC
		}
		
		// Add physical stats
		if cardData.Height != "" && cardData.Height != "null" {
			card.FlavorText = card.FlavorText + "\n身長: " + cardData.Height
		}
		if cardData.Weight != "" && cardData.Weight != "null" {
			if strings.Contains(card.FlavorText, "\n身長:") {
				card.FlavorText = card.FlavorText + " / 体重: " + cardData.Weight
			} else {
				card.FlavorText = card.FlavorText + "\n体重: " + cardData.Weight
			}
		}
		
		// Add attributes
		if cardData.Attribute != "" && cardData.Attribute != "null" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Attribute)
		}
		if cardData.Emotion != "" && cardData.Emotion != "null" {
			card.EnergyIcons = append(card.EnergyIcons, cardData.Emotion)
		}
		
		// Check for counter abilities
		if strings.Contains(cardData.Effect, "【カウンター】") {
			card.IsCounter = true
		}
		if strings.Contains(cardData.Effect, "【メイン/カウンター】") {
			card.IsMainCounter = true
		}
		
		// Add image data
		if imgData, exists := imageMap[cardData.Number]; exists {
			card.ImageURL = imgData.ImageURL
			card.LocalImagePath = imgData.LocalImagePath
		}
		
		// Check if it's a promo card
		if strings.Contains(cardData.Number, "(P)") {
			card.IsPromo = true
		}
		
		// Save to database
		var existingCard models.Card
		result := db.Where("card_no = ? AND rarity = ?", card.CardNo, card.Rarity).First(&existingCard)
		if result.Error == nil {
			// Skip if already exists with same rarity
			continue
		} else {
			db.Create(&card)
		}
	}
	
	return nil
}

// Load parallel cards from card_data.json
func loadParallelCards(imageMap map[string]models.CardJSON) error {
	jsonData, err := loadCardJSON("data/card_data.json")
	if err != nil {
		return err
	}
	
	db := database.GetDB()
	parallelCount := 0
	
	for _, jCard := range jsonData {
		// Only process -P cards
		if !strings.Contains(jCard.Number, "-P") {
			continue
		}
		
		card := models.Card{
			CardNo:         jCard.Number,
			Name:           jCard.Name + " (パラレル)",
			Type:           models.CardTypeFriend, // Default, will be overridden
			ImageURL:       jCard.ImageURL,
			LocalImagePath: jCard.LocalImagePath,
		}
		
		// Try to find the base card to copy its data
		baseCardNo := strings.Replace(jCard.Number, "-P", "", 1)
		var baseCard models.Card
		if err := db.Where("card_no = ?", baseCardNo).First(&baseCard).Error; err == nil {
			// Copy base card data
			card.Type = baseCard.Type
			card.Color = baseCard.Color
			card.Cost = baseCard.Cost
			card.Power = baseCard.Power
			card.Rarity = baseCard.Rarity
			card.Effect = baseCard.Effect
			card.FlavorText = baseCard.FlavorText
			card.EnergyIcons = baseCard.EnergyIcons
			card.IsCounter = baseCard.IsCounter
			card.IsMainCounter = baseCard.IsMainCounter
			// Add -P to rarity for parallel cards
			card.Rarity = models.CardRarity(string(baseCard.Rarity) + "-P")
		} else {
			// Try to infer from image data
			if strings.Contains(jCard.LocalImagePath, "_SR") {
				card.Rarity = "SR-P"
			} else if strings.Contains(jCard.LocalImagePath, "_R") {
				card.Rarity = "R-P"
			} else if strings.Contains(jCard.LocalImagePath, "_SEC") {
				card.Rarity = "SEC-P"
			} else if strings.Contains(jCard.LocalImagePath, "_U") {
				card.Rarity = "U-P"
			} else if strings.Contains(jCard.LocalImagePath, "_C") {
				card.Rarity = "C-P"
			}
		}
		
		// Check if card exists
		var existingCard models.Card
		result := db.Where("card_no = ? AND rarity = ?", card.CardNo, card.Rarity).First(&existingCard)
		if result.Error == nil {
			// Skip if already exists
			continue
		} else {
			// Create new card
			db.Create(&card)
			parallelCount++
		}
	}
	
	log.Printf("Loaded %d parallel cards from card_data.json", parallelCount)
	return nil
}

// Load promo cards from card_data.json
func loadPromoCards(imageMap map[string]models.CardJSON) error {
	jsonData, err := loadCardJSON("data/card_data.json")
	if err != nil {
		return err
	}
	
	db := database.GetDB()
	promoCount := 0
	
	for _, jCard := range jsonData {
		// Only process promo cards with (P) in number
		if !strings.Contains(jCard.Number, "(P)") {
			continue
		}
		
		card := models.Card{
			CardNo:         jCard.Number,
			Name:           jCard.Name,
			Type:           models.CardTypeFriend, // Default, will be overridden
			ImageURL:       jCard.ImageURL,
			LocalImagePath: jCard.LocalImagePath,
			IsPromo:        true,
		}
		
		// Try to find the base card to copy its data
		baseCardNo := strings.Replace(jCard.Number, " (P)", "", 1)
		var baseCard models.Card
		if err := db.Where("card_no = ?", baseCardNo).First(&baseCard).Error; err == nil {
			// Copy base card data
			card.Type = baseCard.Type
			card.Color = baseCard.Color
			card.Cost = baseCard.Cost
			card.Power = baseCard.Power
			card.Rarity = baseCard.Rarity
			card.Effect = baseCard.Effect
			card.FlavorText = baseCard.FlavorText
			card.EnergyIcons = baseCard.EnergyIcons
			card.IsCounter = baseCard.IsCounter
			card.IsMainCounter = baseCard.IsMainCounter
		} else {
			// Try to infer type from jCard.Type
			if jCard.Type == "サポート" {
				card.Type = models.CardTypeSupport
			} else if jCard.Type == "フィールド" {
				card.Type = models.CardTypeField
			}
			
			// Try to infer rarity from image data
			if strings.Contains(jCard.LocalImagePath, "_SR") || strings.Contains(jCard.LocalImagePath, "-P_SR") {
				card.Rarity = models.RaritySR
			} else if strings.Contains(jCard.LocalImagePath, "_R") || strings.Contains(jCard.LocalImagePath, "-P_R") {
				card.Rarity = models.RarityR
			} else if strings.Contains(jCard.LocalImagePath, "_U") || strings.Contains(jCard.LocalImagePath, "-P_U") {
				card.Rarity = models.RarityU
			} else if strings.Contains(jCard.LocalImagePath, "_C") || strings.Contains(jCard.LocalImagePath, "-P_C") {
				card.Rarity = models.RarityC
			}
		}
		
		// Check if card exists
		var existingCard models.Card
		result := db.Where("card_no = ? AND rarity = ?", card.CardNo, card.Rarity).First(&existingCard)
		if result.Error == nil {
			// Skip if already exists
			continue
		} else {
			// Create new card
			db.Create(&card)
			promoCount++
		}
	}
	
	log.Printf("Loaded %d promo cards from card_data.json", promoCount)
	return nil
}
