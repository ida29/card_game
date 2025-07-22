package main

import (
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

	db := database.GetDB()
	
	// Clear existing cards to avoid conflicts
	log.Println("Clearing existing cards...")
	if err := db.Exec("DELETE FROM cards").Error; err != nil {
		log.Printf("Warning: Failed to clear cards: %v", err)
	}
	
	log.Println("Reloading all card data...")
	
	// Load all card data
	if err := utils.LoadAllCardData(); err != nil {
		log.Fatal("Failed to load card data:", err)
	}

	// Count cards
	var total int64
	db.Model(&models.Card{}).Count(&total)
	log.Printf("Total cards loaded: %d", total)
	
	// Count promo cards
	var promoCount int64
	db.Model(&models.Card{}).Where("card_no LIKE '%(P)%' OR is_promo = ?", true).Count(&promoCount)
	log.Printf("Promo cards: %d", promoCount)
	
	// Count parallel cards
	var parallelCount int64
	db.Model(&models.Card{}).Where("card_no LIKE '%-P'").Count(&parallelCount)
	log.Printf("Parallel cards: %d", parallelCount)
	
	log.Println("Card data reload completed successfully!")
}