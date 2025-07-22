package handlers

import (
	"mememe-tcg/internal/models"
	"mememe-tcg/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardService *services.CardService
}

func NewCardHandler() *CardHandler {
	return &CardHandler{
		cardService: services.NewCardService(),
	}
}

func (h *CardHandler) GetAllCards(c *gin.Context) {
	// Check if promo filter is applied
	promoFilter := c.Query("promo")
	
	var cards []models.Card
	var err error
	
	if promoFilter == "true" {
		cards, err = h.cardService.GetPromoCards()
	} else if promoFilter == "false" {
		cards, err = h.cardService.GetNonPromoCards()
	} else {
		cards, err = h.cardService.GetAllCards()
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) GetCardByNumber(c *gin.Context) {
	cardNo := c.Param("cardNo")
	
	card, err := h.cardService.GetCardByNumber(cardNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) GetCardsByType(c *gin.Context) {
	cardType := models.CardType(c.Query("type"))
	
	cards, err := h.cardService.GetCardsByType(cardType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) GetCardsByColor(c *gin.Context) {
	color := models.CardColor(c.Query("color"))
	
	cards, err := h.cardService.GetCardsByColor(color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) SearchCards(c *gin.Context) {
	query := c.Query("q")
	
	cards, err := h.cardService.SearchCards(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) GetPromoCards(c *gin.Context) {
	cards, err := h.cardService.GetPromoCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}