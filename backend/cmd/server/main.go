package main

import (
	"log"
	"mememe-tcg/internal/database"
	"mememe-tcg/internal/handlers"
	"mememe-tcg/internal/utils"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Load all card data from unified sources
	if err := utils.LoadAllCardData(); err != nil {
		log.Println("Warning: Failed to load card data:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Static files for card images
	r.Static("/api/v1/images", "./data/card_images")

	// Initialize handlers
	cardHandler := handlers.NewCardHandler()
	deckHandler := handlers.NewDeckHandler()

	// API routes
	api := r.Group("/api/v1")
	{
		// Card routes
		cards := api.Group("/cards")
		{
			cards.GET("", cardHandler.GetAllCards)
			cards.GET("/promo", cardHandler.GetPromoCards)
			cards.GET("/:cardNo", cardHandler.GetCardByNumber)
			cards.GET("/type", cardHandler.GetCardsByType)
			cards.GET("/color", cardHandler.GetCardsByColor)
			cards.GET("/search", cardHandler.SearchCards)
		}

		// Deck routes
		decks := api.Group("/decks")
		{
			decks.POST("", deckHandler.CreateDeck)
			decks.GET("", deckHandler.GetUserDecks)
			decks.GET("/:id", deckHandler.GetDeck)
			decks.PUT("/:id", deckHandler.UpdateDeck)
			decks.DELETE("/:id", deckHandler.DeleteDeck)
			decks.POST("/validate", deckHandler.ValidateDeck)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}