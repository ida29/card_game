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
	
	// Fix F-068 name
	var card models.Card
	if err := db.Where("card_no = ?", "F-068").First(&card).Error; err == nil {
		card.Name = "デコーレーション" // Use correct dash character
		db.Save(&card)
		log.Printf("Fixed F-068 name to: %s", card.Name)
	} else {
		log.Printf("Could not find F-068: %v", err)
	}
	
	log.Println("F-068 name fix completed!")
}