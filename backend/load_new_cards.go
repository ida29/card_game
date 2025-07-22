package main

import (
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/utils"
)

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	log.Println("Loading new card data from mememe_cards_complete.json...")
	
	// Load all card data
	if err := utils.LoadCompleteCardData(); err != nil {
		log.Fatal("Failed to load card data:", err)
	}

	log.Println("Card data loading completed successfully!")
}