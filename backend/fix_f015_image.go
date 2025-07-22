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
	
	// Fix F-015 (non-promo) image
	var card models.Card
	if err := db.Where("card_no = ? AND rarity = ?", "F-015", "R").First(&card).Error; err == nil {
		card.LocalImagePath = "card_images/F-015_R.jpg"
		card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-015_R.jpg"
		db.Save(&card)
		log.Printf("Fixed F-015 R image path")
	}
	
	// Also fix F-016 SR
	if err := db.Where("card_no = ? AND rarity = ?", "F-016", "SR").First(&card).Error; err == nil {
		card.LocalImagePath = "card_images/F-016_SR.jpg"
		card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-016_SR.jpg"
		db.Save(&card)
		log.Printf("Fixed F-016 SR image path")
	}
	
	// Fix F-064 SR
	if err := db.Where("card_no = ? AND rarity = ?", "F-064", "SR").First(&card).Error; err == nil {
		card.LocalImagePath = "card_images/F-064_SR.jpg"
		card.ImageURL = "https://mememe-tcg.com/assets/images/card/F-064_SR.jpg"
		db.Save(&card)
		log.Printf("Fixed F-064 SR image path")
	}
	
	log.Println("Image fixes completed!")
}