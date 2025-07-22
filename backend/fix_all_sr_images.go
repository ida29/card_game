package main

import (
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
	"strings"
)

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := database.GetDB()
	
	// Get all cards
	var cards []models.Card
	db.Find(&cards)
	
	fixedCount := 0
	for _, card := range cards {
		needsFix := false
		newPath := ""
		
		// Check if it's a regular SR card using SR-P image
		if card.Rarity == "SR" && !strings.Contains(card.CardNo, "-P") && !strings.Contains(card.CardNo, "(P)") {
			if strings.Contains(card.LocalImagePath, "_SR_P") || strings.Contains(card.LocalImagePath, "_P_") || strings.Contains(card.LocalImagePath, "-P_") {
				// Should use regular SR image
				cardNoClean := strings.TrimSuffix(card.CardNo, "-P")
				newPath = "card_images/" + cardNoClean + "_SR.jpg"
				needsFix = true
			}
		}
		
		// Check if it's a regular R card using R-P image
		if card.Rarity == "R" && !strings.Contains(card.CardNo, "-P") && !strings.Contains(card.CardNo, "(P)") {
			if strings.Contains(card.LocalImagePath, "_R_P") || strings.Contains(card.LocalImagePath, "-P_") {
				// Should use regular R image
				cardNoClean := strings.TrimSuffix(card.CardNo, "-P")
				newPath = "card_images/" + cardNoClean + "_R.jpg"
				needsFix = true
			}
		}
		
		if needsFix {
			card.LocalImagePath = newPath
			card.ImageURL = "https://mememe-tcg.com/assets/images/card/" + strings.TrimPrefix(newPath, "card_images/")
			db.Save(&card)
			log.Printf("Fixed %s %s %s -> %s", card.CardNo, card.Rarity, card.Name, newPath)
			fixedCount++
		}
	}
	
	log.Printf("Fixed %d card images!", fixedCount)
}