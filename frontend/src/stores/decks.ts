import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Deck, DeckCard, Card } from '@/types'
import { deckService } from '@/services/api'
import { useCardStore } from './cards'

export const useDeckStore = defineStore('decks', () => {
  const decks = ref<Deck[]>([])
  const currentDeck = ref<Deck | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const totalCards = computed(() => {
    if (!currentDeck.value) return 0
    return currentDeck.value.cards.reduce((sum, card) => sum + card.quantity, 0)
  })

  const isValidDeck = computed(() => {
    if (!currentDeck.value) return false
    const cardCounts = new Map<string, number>()
    let totalCount = 0

    for (const deckCard of currentDeck.value.cards) {
      const currentCount = cardCounts.get(deckCard.card_no) || 0
      cardCounts.set(deckCard.card_no, currentCount + deckCard.quantity)
      totalCount += deckCard.quantity

      if (cardCounts.get(deckCard.card_no)! > 4) {
        return false
      }
    }

    return totalCount === 50
  })

  async function fetchUserDecks() {
    loading.value = true
    error.value = null
    try {
      decks.value = await deckService.getUserDecks()
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch decks'
    } finally {
      loading.value = false
    }
  }

  async function fetchDeck(id: number) {
    loading.value = true
    error.value = null
    try {
      currentDeck.value = await deckService.getDeck(id)
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch deck'
    } finally {
      loading.value = false
    }
  }

  async function createDeck(name: string) {
    loading.value = true
    error.value = null
    try {
      const newDeck = await deckService.createDeck({
        name,
        user_id: 1, // Dummy user ID
        cards: [],
        is_active: false,
      })
      decks.value.push(newDeck)
      currentDeck.value = newDeck
      return newDeck
    } catch (err: any) {
      error.value = err.message || 'Failed to create deck'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateDeck(deck: Deck) {
    loading.value = true
    error.value = null
    try {
      console.log('Updating deck with cards:', deck.cards.length)
      deck.cards.forEach((card, idx) => {
        console.log(`Card ${idx}: ${card.card_no} x${card.quantity}`)
      })
      
      const updatedDeck = await deckService.updateDeck(deck.ID, deck)
      console.log('Received updated deck with cards:', updatedDeck.cards.length)
      
      const index = decks.value.findIndex(d => d.ID === deck.ID)
      if (index !== -1) {
        decks.value[index] = updatedDeck
      }
      if (currentDeck.value?.ID === deck.ID) {
        currentDeck.value = updatedDeck
      }
      return updatedDeck
    } catch (err: any) {
      error.value = err.message || 'Failed to update deck'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteDeck(id: number) {
    loading.value = true
    error.value = null
    try {
      await deckService.deleteDeck(id)
      decks.value = decks.value.filter(d => d.ID !== id)
      if (currentDeck.value?.ID === id) {
        currentDeck.value = null
      }
    } catch (err: any) {
      error.value = err.message || 'Failed to delete deck'
      throw err
    } finally {
      loading.value = false
    }
  }

  function setCurrentDeck(deck: Deck | null) {
    currentDeck.value = deck
  }

  function addCardToDeck(cardNo: string, card?: Card) {
    if (!currentDeck.value) return { success: false, reason: 'no-deck' }

    const existingCard = currentDeck.value.cards.find(c => c.card_no === cardNo)
    
    // Check if adding one more card would exceed 50 total
    const currentTotal = currentDeck.value.cards.reduce((sum, card) => sum + card.quantity, 0)
    if (currentTotal >= 50) {
      return { success: false, reason: 'deck-full' }
    }
    
    if (existingCard) {
      if (existingCard.quantity >= 4) {
        return { success: false, reason: 'max-copies' }
      }
      existingCard.quantity++
      // Ensure card data is present
      if (!existingCard.card && card) {
        existingCard.card = card
      }
    } else {
      // Get card from the card store if not provided
      const cardStore = useCardStore()
      const actualCard = card || cardStore.cards.find(c => c.card_no === cardNo)
      
      if (!actualCard) return { success: false, reason: 'card-not-found' }
      
      currentDeck.value.cards.push({
        ID: 0,
        deck_id: currentDeck.value.ID,
        card_no: cardNo,
        quantity: 1,
        card: actualCard,
      })
    }

    return { success: true, reason: null }
  }

  function removeCardFromDeck(cardNo: string) {
    if (!currentDeck.value) return

    const cardIndex = currentDeck.value.cards.findIndex(c => c.card_no === cardNo)
    if (cardIndex === -1) return

    const deckCard = currentDeck.value.cards[cardIndex]
    if (deckCard.quantity > 1) {
      deckCard.quantity--
    } else {
      currentDeck.value.cards.splice(cardIndex, 1)
    }
  }

  function clearCurrentDeck() {
    console.log('clearCurrentDeck called', currentDeck.value)
    if (!currentDeck.value) {
      console.log('No current deck to clear')
      return
    }
    console.log('Clearing deck with', currentDeck.value.cards.length, 'cards')
    
    // Clear the cards array
    currentDeck.value.cards = []
    
    // Also update the deck in the decks array to maintain consistency
    const deckIndex = decks.value.findIndex(d => d.ID === currentDeck.value!.ID)
    if (deckIndex !== -1) {
      decks.value[deckIndex] = { ...currentDeck.value }
    }
    
    console.log('Deck cleared, remaining cards:', currentDeck.value.cards.length)
  }

  return {
    decks,
    currentDeck,
    loading,
    error,
    totalCards,
    isValidDeck,
    fetchUserDecks,
    fetchDeck,
    createDeck,
    updateDeck,
    deleteDeck,
    setCurrentDeck,
    addCardToDeck,
    removeCardFromDeck,
    clearCurrentDeck,
  }
})