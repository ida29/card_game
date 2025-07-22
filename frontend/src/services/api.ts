import axios from 'axios'
import type { Card, Deck, DeckCard } from '@/types'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

export const cardService = {
  async getAllCards(promo?: boolean): Promise<Card[]> {
    const params = promo !== undefined ? { promo: promo.toString() } : {}
    const response = await api.get('/cards', { params })
    return response.data
  },

  async getCardByNumber(cardNo: string): Promise<Card> {
    const response = await api.get(`/cards/${cardNo}`)
    return response.data
  },

  async getCardsByType(type: string): Promise<Card[]> {
    const response = await api.get('/cards/type', { params: { type } })
    return response.data
  },

  async getCardsByColor(color: string): Promise<Card[]> {
    const response = await api.get('/cards/color', { params: { color } })
    return response.data
  },

  async searchCards(query: string): Promise<Card[]> {
    const response = await api.get('/cards/search', { params: { q: query } })
    return response.data
  },

  async getPromoCards(): Promise<Card[]> {
    const response = await api.get('/cards/promo')
    return response.data
  },
}

export const deckService = {
  async createDeck(deck: Partial<Deck>): Promise<Deck> {
    const response = await api.post('/decks', deck)
    return response.data
  },

  async getUserDecks(): Promise<Deck[]> {
    const response = await api.get('/decks')
    return response.data
  },

  async getDeck(id: number): Promise<Deck> {
    const response = await api.get(`/decks/${id}`)
    return response.data
  },

  async updateDeck(id: number, deck: Partial<Deck>): Promise<Deck> {
    const response = await api.put(`/decks/${id}`, deck)
    return response.data
  },

  async deleteDeck(id: number): Promise<void> {
    await api.delete(`/decks/${id}`)
  },

  async validateDeck(cards: DeckCard[]): Promise<{ valid: boolean; error?: string }> {
    try {
      const response = await api.post('/decks/validate', cards)
      return response.data
    } catch (error: any) {
      return { valid: false, error: error.response?.data?.error || 'Validation failed' }
    }
  },
}

export default api