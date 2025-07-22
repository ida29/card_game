package main

import (
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
)

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := database.GetDB()
	
	// Fix F-016 SR (regular, not promo)
	var card models.Card
	if err := db.Where("card_no = ? AND rarity = ? AND is_promo = ?", "F-016", "SR", false).First(&card).Error; err == nil {
		card.LocalImagePath = "card_images/F-016_SR.jpg"
		card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-016_SR.jpg"
		db.Save(&card)
		log.Printf("Fixed F-016 SR (regular) image path to: %s", card.LocalImagePath)
	} else {
		log.Printf("Could not find F-016 SR (regular): %v", err)
		
		// Try without is_promo filter
		if err := db.Where("card_no = ? AND rarity = ?", "F-016", "SR").First(&card).Error; err == nil {
			if card.LocalImagePath != "card_images/F-016_SR.jpg" {
				card.LocalImagePath = "card_images/F-016_SR.jpg"
				card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-016_SR.jpg"
				db.Save(&card)
				log.Printf("Fixed F-016 SR image path to: %s", card.LocalImagePath)
			}
		}
	}
	
	log.Println("F-016 final fix completed!")
}