package main

import (
	"fmt"
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/models"
	"mememe-tcg/internal/utils"
)

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Load all card data
	if err := utils.LoadAllCardData(); err != nil {
		log.Fatal("Failed to load card data:", err)
	}

	// Check for F-064-P
	db := database.GetDB()
	var card models.Card
	if err := db.Where("card_no = ?", "F-064-P").First(&card).Error; err != nil {
		fmt.Println("F-064-P not found:", err)
	} else {
		fmt.Printf("Found F-064-P: %s (%s)\n", card.Name, card.Rarity)
	}

	// List all -P cards
	var parallelCards []models.Card
	db.Where("card_no LIKE ?", "%-P").Find(&parallelCards)
	fmt.Printf("\nTotal -P cards found: %d\n", len(parallelCards))
	for _, c := range parallelCards {
		fmt.Printf("- %s: %s\n", c.CardNo, c.Name)
	}
}