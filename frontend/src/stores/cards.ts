import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Card, CardType, CardColor } from '@/types'
import { cardService } from '@/services/api'

export const useCardStore = defineStore('cards', () => {
  const cards = ref<Card[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const friendCards = computed(() => 
    cards.value.filter(card => card.type === 'ふれんど')
  )

  const supportCards = computed(() => 
    cards.value.filter(card => card.type === 'サポート')
  )

  const fieldCards = computed(() => 
    cards.value.filter(card => card.type === 'フィールド')
  )

  const getCardsByColor = (color: CardColor) => 
    cards.value.filter(card => card.color === color)

  const getCardsByRarity = (rarity: string) => 
    cards.value.filter(card => card.rarity === rarity)

  async function fetchAllCards() {
    loading.value = true
    error.value = null
    try {
      cards.value = await cardService.getAllCards()
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch cards'
    } finally {
      loading.value = false
    }
  }

  async function searchCards(query: string) {
    loading.value = true
    error.value = null
    try {
      return await cardService.searchCards(query)
    } catch (err: any) {
      error.value = err.message || 'Failed to search cards'
      return []
    } finally {
      loading.value = false
    }
  }

  function getCardByNumber(cardNo: string): Card | undefined {
    return cards.value.find(card => card.card_no === cardNo)
  }

  return {
    cards,
    loading,
    error,
    friendCards,
    supportCards,
    fieldCards,
    getCardsByColor,
    getCardsByRarity,
    fetchAllCards,
    searchCards,
    getCardByNumber,
  }
})