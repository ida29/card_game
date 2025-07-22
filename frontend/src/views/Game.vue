<template>
  <div>
    <h1 class="text-3xl font-bold text-white mb-8">対戦</h1>
    
    <div v-if="!gameMode" class="bg-black bg-opacity-30 rounded-lg p-8">
      <h2 class="text-2xl font-bold text-white mb-8 text-center">対戦モードを選択</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl mx-auto">
        <button
          @click="selectMode('pvp')"
          class="bg-gradient-to-br from-blue-600 to-blue-800 hover:from-blue-700 hover:to-blue-900 rounded-lg p-8 transition-all transform hover:scale-105"
        >
          <div class="text-white">
            <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <h3 class="text-xl font-bold mb-2">誰かと対戦</h3>
            <p class="text-sm opacity-90">他のプレイヤーとオンラインで対戦します</p>
          </div>
        </button>
        
        <button
          @click="selectMode('cpu')"
          class="bg-gradient-to-br from-purple-600 to-purple-800 hover:from-purple-700 hover:to-purple-900 rounded-lg p-8 transition-all transform hover:scale-105"
        >
          <div class="text-white">
            <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <h3 class="text-xl font-bold mb-2">CPUと対戦</h3>
            <p class="text-sm opacity-90">コンピューターと練習対戦をします</p>
          </div>
        </button>
      </div>
      
      <div class="mt-8 text-center">
        <router-link
          to="/"
          class="text-gray-400 hover:text-white transition-colors"
        >
          ← ホームに戻る
        </router-link>
      </div>
    </div>
    
    <div v-else-if="gameMode === 'cpu-difficulty'" class="bg-black bg-opacity-30 rounded-lg p-8">
      <h2 class="text-2xl font-bold text-white mb-8 text-center">CPU難易度を選択</h2>
      
      <div class="space-y-4 max-w-md mx-auto">
        <button
          @click="selectDifficulty('easy')"
          class="w-full bg-gradient-to-r from-green-600 to-green-700 hover:from-green-700 hover:to-green-800 rounded-lg p-4 text-white transition-all"
        >
          <h3 class="text-lg font-bold">かんたん</h3>
          <p class="text-sm opacity-90">初心者向け - CPUは基本的な動きをします</p>
        </button>
        
        <button
          @click="selectDifficulty('normal')"
          class="w-full bg-gradient-to-r from-yellow-600 to-yellow-700 hover:from-yellow-700 hover:to-yellow-800 rounded-lg p-4 text-white transition-all"
        >
          <h3 class="text-lg font-bold">ふつう</h3>
          <p class="text-sm opacity-90">標準的な強さ - バランスの良い対戦相手</p>
        </button>
        
        <button
          @click="selectDifficulty('hard')"
          class="w-full bg-gradient-to-r from-red-600 to-red-700 hover:from-red-700 hover:to-red-800 rounded-lg p-4 text-white transition-all"
        >
          <h3 class="text-lg font-bold">むずかしい</h3>
          <p class="text-sm opacity-90">上級者向け - CPUは最適な戦略を使います</p>
        </button>
      </div>
      
      <div class="mt-6 text-center">
        <button
          @click="gameMode = null"
          class="text-gray-400 hover:text-white transition-colors"
        >
          ← モード選択に戻る
        </button>
      </div>
    </div>
    
    <div v-else-if="gameMode === 'deck-selection'" class="bg-black bg-opacity-30 rounded-lg p-8">
      <h2 class="text-2xl font-bold text-white mb-6 text-center">デッキを選択</h2>
      
      <div v-if="availableDecks.length === 0" class="text-center">
        <p class="text-gray-300 mb-4">対戦可能なデッキがありません</p>
        <router-link
          to="/deck-builder"
          class="btn-primary"
        >
          デッキを作成する
        </router-link>
      </div>
      
      <div v-else class="space-y-4 max-w-2xl mx-auto">
        <div
          v-for="deck in availableDecks"
          :key="deck.id"
          @click="selectDeck(deck)"
          class="bg-gray-800 hover:bg-gray-700 rounded-lg p-4 cursor-pointer transition-colors"
        >
          <div class="flex justify-between items-center">
            <div>
              <h3 class="text-lg font-bold text-white">{{ deck.name }}</h3>
              <p class="text-sm text-gray-400">{{ getDeckCardCount(deck) }}枚</p>
            </div>
            <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </div>
      </div>
      
      <div class="mt-6 text-center">
        <button
          @click="gameMode = null"
          class="text-gray-400 hover:text-white transition-colors"
        >
          ← モード選択に戻る
        </button>
      </div>
    </div>
    
    <div v-else-if="gameMode === 'battle'">
      <GameStart v-if="!gameReady" @game-ready="handleGameReady" />
      <BattleField v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDeckStore } from '@/stores/decks'
import { useGameStore } from '@/stores/game'
import { useCardStore } from '@/stores/cards'
import type { Deck, DeckCard } from '@/types'
import BattleField from '@/components/game/BattleField.vue'
import GameStart from '@/components/game/GameStart.vue'

const deckStore = useDeckStore()
const gameStore = useGameStore()
const cardStore = useCardStore()

const gameMode = ref<null | 'pvp' | 'cpu' | 'cpu-difficulty' | 'deck-selection' | 'battle'>(null)
const battleMode = ref<null | 'pvp' | 'cpu'>(null)
const selectedDeck = ref<Deck | null>(null)
const gameReady = ref(false)

const getDeckCardCount = (deck: Deck) => {
  return deck.cards?.reduce((sum, card) => sum + (card.quantity || 1), 0) || 0
}

const availableDecks = computed(() => {
  console.log('All decks:', deckStore.decks)
  deckStore.decks.forEach((deck, index) => {
    console.log(`Deck ${index}: ${deck.name}, cards:`, getDeckCardCount(deck))
  })
  
  // Count total cards considering quantity
  return deckStore.decks.filter(deck => {
    const totalCards = getDeckCardCount(deck)
    console.log(`Deck ${deck.name}: ${totalCards} total cards`)
    return totalCards >= 40
  })
})

onMounted(async () => {
  await deckStore.loadDecks()
  // Load all cards for CPU deck generation
  await cardStore.fetchAllCards()
})

const selectMode = (mode: 'pvp' | 'cpu') => {
  battleMode.value = mode
  if (mode === 'cpu') {
    gameMode.value = 'cpu-difficulty'
  } else {
    gameMode.value = 'deck-selection'
  }
}

const selectDifficulty = (difficulty: 'easy' | 'normal' | 'hard') => {
  console.log('selectDifficulty called:', { difficulty, battleMode: battleMode.value })
  gameStore.setCPUDifficulty(difficulty)
  gameMode.value = 'deck-selection'
  console.log('After selectDifficulty:', { gameMode: gameMode.value, battleMode: battleMode.value })
}

const selectDeck = (deck: Deck) => {
  console.log('selectDeck called:', { deck, battleMode: battleMode.value })
  selectedDeck.value = deck
  gameReady.value = false // Reset game ready state
  gameMode.value = 'battle'
  console.log('About to call startBattle with:', { selectedDeck: !!selectedDeck.value, battleMode: battleMode.value, gameReady: gameReady.value })
  startBattle()
}

const startBattle = () => {
  console.log('startBattle called:', { selectedDeck: selectedDeck.value, battleMode: battleMode.value })
  if (!selectedDeck.value || !battleMode.value) {
    console.error('Cannot start battle - missing requirements:', { 
      hasSelectedDeck: !!selectedDeck.value, 
      hasBattleMode: !!battleMode.value 
    })
    return
  }
  
  // Expand deck cards based on quantity
  const expandedDeck: DeckCard[] = []
  selectedDeck.value.cards.forEach(deckCard => {
    for (let i = 0; i < (deckCard.quantity || 1); i++) {
      expandedDeck.push({
        ...deckCard,
        quantity: 1
      })
    }
  })
  
  console.log('Expanded deck size:', expandedDeck.length)
  
  // Initialize game but don't start phases yet
  gameStore.initializeGame(battleMode.value, expandedDeck, false)
}

const handleGameReady = () => {
  console.log('handleGameReady called - about to start game phases')
  // Start the actual game phases after janken
  gameStore.startGame()
  gameReady.value = true
  console.log('Game ready set to true, should now show BattleField')
}
</script>