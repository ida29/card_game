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
	
	// Fix regular cards that have promo images
	fixes := []struct{
		cardNo string
		rarity string
		correctImage string
	}{
		{"F-013", "R", "card_images/F-013_R.jpg"},
		{"F-023", "C", "card_images/F-023_C.jpg"},
	}
	
	for _, fix := range fixes {
		var card models.Card
		if err := db.Where("card_no = ? AND rarity = ?", fix.cardNo, fix.rarity).First(&card).Error; err == nil {
			if strings.Contains(card.LocalImagePath, "-P") {
				card.LocalImagePath = fix.correctImage
				card.ImageURL = "https://mememe-tcg.com/assets/images/card/" + strings.TrimPrefix(fix.correctImage, "card_images/")
				db.Save(&card)
				log.Printf("Fixed %s %s image to: %s", fix.cardNo, fix.rarity, fix.correctImage)
			}
		}
	}
	
	log.Println("Final image fixes completed!")
}