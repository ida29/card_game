package handlers

import (
	"log"
	"mememe-tcg/internal/models"
	"mememe-tcg/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeckHandler struct {
	deckService *services.DeckService
}

func NewDeckHandler() *DeckHandler {
	return &DeckHandler{
		deckService: services.NewDeckService(),
	}
}

func (h *DeckHandler) CreateDeck(c *gin.Context) {
	var deck models.Deck
	if err := c.ShouldBindJSON(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For now, use a dummy user ID
	deck.UserID = 1

	if err := h.deckService.CreateDeck(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the created deck with full card data
	createdDeck, err := h.deckService.GetDeck(deck.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created deck"})
		return
	}

	c.JSON(http.StatusCreated, createdDeck)
}

func (h *DeckHandler) GetDeck(c *gin.Context) {
	deckID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	deck, err := h.deckService.GetDeck(uint(deckID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deck not found"})
		return
	}

	c.JSON(http.StatusOK, deck)
}

func (h *DeckHandler) GetUserDecks(c *gin.Context) {
	// For now, use a dummy user ID
	userID := uint(1)

	decks, err := h.deckService.GetUserDecks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure empty arrays instead of null for cards
	for i := range decks {
		if decks[i].Cards == nil {
			decks[i].Cards = []models.DeckCard{}
		}
	}

	c.JSON(http.StatusOK, decks)
}

func (h *DeckHandler) UpdateDeck(c *gin.Context) {
	deckID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	var deck models.Deck
	if err := c.ShouldBindJSON(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the incoming request
	log.Printf("UpdateDeck handler received deck ID %d with %d cards", deckID, len(deck.Cards))

	deck.ID = uint(deckID)
	if err := h.deckService.UpdateDeck(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the updated deck with full card data
	updatedDeck, err := h.deckService.GetDeck(uint(deckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated deck"})
		return
	}

	log.Printf("UpdateDeck handler returning deck with %d cards", len(updatedDeck.Cards))

	c.JSON(http.StatusOK, updatedDeck)
}

func (h *DeckHandler) DeleteDeck(c *gin.Context) {
	deckID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	if err := h.deckService.DeleteDeck(uint(deckID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *DeckHandler) ValidateDeck(c *gin.Context) {
	var cards []models.DeckCard
	if err := c.ShouldBindJSON(&cards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.deckService.ValidateDeckCards(cards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}