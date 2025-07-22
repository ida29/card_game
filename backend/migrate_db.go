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
	
	// First, drop the old unique index
	if err := db.Exec("DROP INDEX IF EXISTS idx_cards_card_no").Error; err != nil {
		log.Printf("Warning: Failed to drop old index: %v", err)
	}
	
	// Add is_promo column if it doesn't exist
	if err := db.Exec("ALTER TABLE cards ADD COLUMN is_promo BOOLEAN DEFAULT FALSE").Error; err != nil {
		log.Printf("Warning: Failed to add is_promo column (may already exist): %v", err)
	}
	
	// Run auto migration to update schema
	if err := db.AutoMigrate(&models.Card{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	// Create new composite unique index
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_card_no_rarity ON cards(card_no, rarity)").Error; err != nil {
		log.Printf("Warning: Failed to create composite index: %v", err)
	}
	
	log.Println("Database migration completed successfully!")
}