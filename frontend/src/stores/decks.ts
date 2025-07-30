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
      const fetchedDecks = await deckService.getUserDecks()
      
      // Ensure card data is populated for API-fetched decks
      const cardStore = useCardStore()
      if (cardStore.cards.length === 0) {
        await cardStore.fetchAllCards()
      }
      
      for (const deck of fetchedDecks) {
        if (deck.cards) {
          for (const deckCard of deck.cards) {
            if (!deckCard.card || !deckCard.card.name) {
              const fullCard = cardStore.cards.find(c => c.card_no === deckCard.card_no)
              if (fullCard) {
                deckCard.card = fullCard
              }
            }
          }
        }
      }
      
      decks.value = fetchedDecks
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
      // Don't automatically modify main_card_no - let the user explicitly set it
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch deck'
      // Try to load from local storage as fallback
      const localDeck = decks.value.find(d => d.ID === id)
      if (localDeck) {
        currentDeck.value = JSON.parse(JSON.stringify(localDeck)) // Deep copy
        
        // Ensure card data is populated
        const cardStore = useCardStore()
        if (currentDeck.value && currentDeck.value.cards) {
          for (const deckCard of currentDeck.value.cards) {
            if (!deckCard.card || !deckCard.card.name) {
              const fullCard = cardStore.cards.find(c => c.card_no === deckCard.card_no)
              if (fullCard) {
                deckCard.card = fullCard
              }
            }
          }
        }
        
        // Don't automatically modify main_card_no - let the user explicitly set it
      }
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
      // Save to localStorage
      localStorage.setItem('mememe_decks', JSON.stringify(decks.value))
      return newDeck
    } catch (err: any) {
      console.log('API failed, creating local deck:', err)
      // Fallback: Create deck locally
      const newDeck: Deck = {
        ID: Date.now(), // Use timestamp as ID
        name,
        user_id: 1,
        cards: [],
        is_active: false,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString()
      }
      decks.value.push(newDeck)
      currentDeck.value = newDeck
      // Save to localStorage
      localStorage.setItem('mememe_decks', JSON.stringify(decks.value))
      return newDeck
    } finally {
      loading.value = false
    }
  }

  async function updateDeck(deck: Deck) {
    loading.value = true
    error.value = null
    try {
      console.log('Updating deck with cards:', deck.cards.length)
      console.log('Main card no:', deck.main_card_no)
      deck.cards.forEach((card, idx) => {
        console.log(`Card ${idx}: ${card.card_no} x${card.quantity}`)
      })
      
      const updatedDeck = await deckService.updateDeck(deck.ID, deck)
      console.log('Received updated deck with cards:', updatedDeck.cards.length)
      console.log('Received main card no:', updatedDeck.main_card_no)
      
      const index = decks.value.findIndex(d => d.ID === deck.ID)
      if (index !== -1) {
        // Preserve the main_card_no from the input deck if it exists
        if (deck.main_card_no && !updatedDeck.main_card_no) {
          updatedDeck.main_card_no = deck.main_card_no
        }
        decks.value[index] = updatedDeck
      }
      if (currentDeck.value?.ID === deck.ID) {
        // Preserve the main_card_no from the input deck if it exists
        if (deck.main_card_no && !updatedDeck.main_card_no) {
          updatedDeck.main_card_no = deck.main_card_no
        }
        currentDeck.value = updatedDeck
      }
      // Save to localStorage
      localStorage.setItem('mememe_decks', JSON.stringify(decks.value))
      return updatedDeck
    } catch (err: any) {
      console.log('API update failed, updating locally:', err)
      // Fallback: Update deck locally
      const index = decks.value.findIndex(d => d.ID === deck.ID)
      if (index !== -1) {
        deck.updated_at = new Date().toISOString()
        decks.value[index] = { ...deck }
        if (currentDeck.value?.ID === deck.ID) {
          currentDeck.value = { ...deck }
        }
        // Save to localStorage with main_card_no
        console.log('Saving to localStorage with main_card_no:', deck.main_card_no)
        localStorage.setItem('mememe_decks', JSON.stringify(decks.value))
        return deck
      }
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

  async function loadDecks() {
    try {
      // Try to load from localStorage first
      const savedDecks = localStorage.getItem('mememe_decks')
      if (savedDecks) {
        const parsedDecks = JSON.parse(savedDecks)
        console.log('Loaded decks from localStorage:', parsedDecks)
        
        // Populate card data for each deck card
        const cardStore = useCardStore()
        // Ensure cards are loaded
        if (cardStore.cards.length === 0) {
          await cardStore.fetchAllCards()
        }
        
        // Populate card data for each deck
        for (const deck of parsedDecks) {
          if (deck.cards) {
            for (const deckCard of deck.cards) {
              if (!deckCard.card || !deckCard.card.name) {
                // Find the full card data
                const fullCard = cardStore.cards.find(c => c.card_no === deckCard.card_no)
                if (fullCard) {
                  deckCard.card = fullCard
                }
              }
            }
          }
        }
        
        decks.value = parsedDecks
      }
      
      // Also try to fetch from server
      try {
        await fetchUserDecks()
        console.log('Fetched decks from server:', decks.value)
      } catch (fetchError) {
        console.log('Could not fetch from server, using local decks only:', fetchError)
      }
    } catch (error) {
      console.error('Error loading decks:', error)
      // Initialize with empty array if all else fails
      decks.value = []
    }
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
    loadDecks,
  }
})