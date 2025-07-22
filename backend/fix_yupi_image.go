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
	
	// Fix F-023 (non-promo) Yupi image
	var card models.Card
	if err := db.Where("card_no = ? AND rarity = ?", "F-023", "C").First(&card).Error; err == nil {
		card.LocalImagePath = "card_images/F-023_C.jpg"
		card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-023_C.jpg"
		db.Save(&card)
		log.Printf("Fixed F-023 C ユピ image path: %s", card.LocalImagePath)
	} else {
		log.Printf("Could not find F-023 C: %v", err)
	}
	
	// Verify F-023 (P) has correct image
	if err := db.Where("card_no = ? AND rarity = ?", "F-023 (P)", "C").First(&card).Error; err == nil {
		log.Printf("F-023 (P) C ユピ current image: %s", card.LocalImagePath)
		// Make sure it has the promo image
		if card.LocalImagePath != "card_images/F-023-P_C_Blue.jpg" {
			card.LocalImagePath = "card_images/F-023-P_C_Blue.jpg"
			card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-023-P_C_Blue.jpg"
			db.Save(&card)
			log.Printf("Fixed F-023 (P) C ユピ promo image path")
		}
	}
	
	log.Println("Yupi image fixes completed!")
}