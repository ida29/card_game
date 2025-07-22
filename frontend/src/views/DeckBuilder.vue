<template>
  <div>
    <h1 class="text-3xl font-bold text-white mb-8">デッキビルダー</h1>
    
    <!-- Deck Selection -->
    <div class="mb-8 bg-black bg-opacity-30 rounded-lg p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold text-white">デッキ選択</h2>
        <button
          @click="createNewDeck"
          class="btn-primary"
        >
          新しいデッキ
        </button>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
        <div
          v-for="deck in deckStore.decks"
          :key="deck.ID"
          class="bg-gray-800 rounded-lg p-4 cursor-pointer hover:bg-gray-700 transition-colors"
          :class="{ 'ring-2 ring-blue-500': deckStore.currentDeck?.ID === deck.ID }"
          @click="selectDeck(deck)"
        >
          <h3 class="text-white font-bold">{{ deck.name }}</h3>
          <p class="text-gray-400 text-sm">
            {{ getDeckCardCount(deck) }}/50 カード
          </p>
        </div>
      </div>
    </div>
    
    <!-- Current Deck Info -->
    <div v-if="deckStore.currentDeck" class="mb-8 bg-black bg-opacity-30 rounded-lg p-6">
      <div class="flex justify-between items-center mb-4">
        <div>
          <h2 class="text-xl font-bold text-white">{{ deckStore.currentDeck.name }}</h2>
          <p class="text-gray-400">
            {{ deckStore.totalCards }}/50 カード
            <span
              class="ml-2 px-2 py-1 rounded text-xs"
              :class="deckStore.isValidDeck ? 'bg-green-600' : 'bg-red-600'"
            >
              {{ deckStore.isValidDeck ? '有効' : '無効' }}
            </span>
          </p>
        </div>
        
        <div class="space-x-2">
          <button
            @click="saveDeck"
            :disabled="!deckStore.isValidDeck"
            class="btn-primary disabled:bg-gray-600 disabled:cursor-not-allowed"
          >
            保存
          </button>
          <button
            @click="deleteDeck"
            class="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition-colors"
          >
            削除
          </button>
        </div>
      </div>
    </div>
    
    <!-- Deck Editing Section - Only show when a deck is selected -->
    <div v-if="deckStore.currentDeck" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- Card Collection (Left Side) -->
      <div class="bg-black bg-opacity-30 rounded-lg p-6">
        <h3 class="text-xl font-bold text-white mb-4">所持カード一覧</h3>
        
        <!-- Filters -->
        <div class="mb-4 grid grid-cols-2 gap-2">
          <select
            v-model="selectedType"
            class="px-3 py-2 bg-gray-800 text-white rounded border border-gray-600"
          >
            <option value="">すべてのタイプ</option>
            <option value="ふれんど">ふれんど</option>
            <option value="サポート">サポート</option>
            <option value="フィールド">フィールド</option>
          </select>
          
          <select
            v-model="selectedColor"
            class="px-3 py-2 bg-gray-800 text-white rounded border border-gray-600"
          >
            <option value="">すべての色</option>
            <option value="赤">赤</option>
            <option value="青">青</option>
            <option value="黄">黄</option>
            <option value="緑">緑</option>
          </select>
        </div>
        
        <div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-2 max-h-[600px] overflow-y-auto">
          <div
            v-for="card in filteredAvailableCards"
            :key="card.card_no"
            class="bg-gray-800 rounded-lg p-1 cursor-pointer hover:bg-gray-700 transition-colors relative group"
            @click="addCardToDeck(card)"
          >
            <img
              :src="getCardImageUrl(card)"
              :alt="card.name"
              class="w-full aspect-[2/3] object-cover rounded mb-1"
            />
            <p class="text-white text-[10px] font-bold truncate">{{ card.name }}</p>
            <p class="text-gray-400 text-[10px]">{{ card.type }}</p>
            <div class="absolute inset-0 bg-green-600 bg-opacity-0 group-hover:bg-opacity-20 rounded-lg transition-all flex items-center justify-center">
              <span class="text-white text-xl font-bold opacity-0 group-hover:opacity-100 transition-opacity">+</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Current Deck (Right Side) -->
      <div class="bg-black bg-opacity-30 rounded-lg p-6">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-xl font-bold text-white">デッキカード一覧</h3>
          <button
            @click="clearDeck"
            class="bg-red-600 hover:bg-red-700 text-white font-bold py-1 px-3 rounded text-sm transition-colors"
          >
            クリア
          </button>
        </div>
        
        <div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-2 max-h-[600px] overflow-y-auto">
          <div
            v-for="deckCard in deckStore.currentDeck.cards"
            :key="deckCard.card_no"
            class="bg-gray-800 rounded-lg p-1 cursor-pointer hover:bg-gray-700 transition-colors relative group"
            @click="removeCardFromDeck(deckCard.card_no)"
          >
            <img
              v-if="deckCard.card"
              :src="getCardImageUrl(deckCard.card)"
              :alt="deckCard.card.name"
              class="w-full aspect-[2/3] object-cover rounded mb-1"
            />
            <div v-else class="w-full aspect-[2/3] bg-gray-700 rounded mb-1 flex items-center justify-center">
              <span class="text-gray-500 text-[10px]">No Image</span>
            </div>
            <p class="text-white text-[10px] font-bold truncate">{{ deckCard.card?.name || 'Unknown' }}</p>
            <p class="text-gray-400 text-[10px]">{{ deckCard.card?.type || 'Unknown' }}</p>
            <div class="absolute top-1 right-1 bg-blue-600 text-white rounded-full w-5 h-5 flex items-center justify-center text-[10px] font-bold">
              {{ deckCard.quantity }}
            </div>
            <div class="absolute inset-0 bg-red-600 bg-opacity-0 group-hover:bg-opacity-20 rounded-lg transition-all flex items-center justify-center">
              <span class="text-white text-xl font-bold opacity-0 group-hover:opacity-100 transition-opacity">−</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- No Deck Selected Message -->
    <div v-else class="bg-black bg-opacity-30 rounded-lg p-12 text-center">
      <p class="text-gray-400 text-xl">デッキを選択してください</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCardStore } from '@/stores/cards'
import { useDeckStore } from '@/stores/decks'
import type { Card, Deck, CardType, CardColor } from '@/types'

const cardStore = useCardStore()
const deckStore = useDeckStore()

const selectedType = ref<CardType | ''>('')
const selectedColor = ref<CardColor | ''>('')

const filteredAvailableCards = computed(() => {
  let filtered = cardStore.cards

  if (selectedType.value) {
    filtered = filtered.filter(card => card.type === selectedType.value)
  }

  if (selectedColor.value) {
    filtered = filtered.filter(card => card.color === selectedColor.value)
  }

  return filtered
})

function getDeckCardCount(deck: Deck): number {
  return deck.cards.reduce((sum, card) => sum + card.quantity, 0)
}

async function selectDeck(deck: Deck) {
  await deckStore.fetchDeck(deck.ID)
}

async function createNewDeck() {
  const name = prompt('デッキ名を入力してください:')
  if (name) {
    try {
      await deckStore.createDeck(name)
    } catch (error) {
      alert('デッキの作成に失敗しました')
    }
  }
}

async function saveDeck() {
  if (deckStore.currentDeck) {
    try {
      await deckStore.updateDeck(deckStore.currentDeck)
      alert('デッキを保存しました')
    } catch (error) {
      alert('デッキの保存に失敗しました')
    }
  }
}

async function deleteDeck() {
  if (deckStore.currentDeck && confirm('このデッキを削除しますか？')) {
    try {
      await deckStore.deleteDeck(deckStore.currentDeck.ID)
    } catch (error) {
      alert('デッキの削除に失敗しました')
    }
  }
}

function addCardToDeck(card: Card) {
  const result = deckStore.addCardToDeck(card.card_no, card)
  if (!result.success) {
    switch (result.reason) {
      case 'max-copies':
        alert('同じカードはデッキに４枚までしか入れられません。')
        break
      case 'deck-full':
        alert('デッキは50枚までしか入れられません。')
        break
      case 'no-deck':
        alert('デッキが選択されていません。')
        break
      case 'card-not-found':
        alert('カードが見つかりません。')
        break
    }
  }
}

function removeCardFromDeck(cardNo: string) {
  deckStore.removeCardFromDeck(cardNo)
}

function clearDeck() {
  console.log('Clear deck button clicked')
  if (confirm('デッキのすべてのカードを削除しますか？')) {
    console.log('User confirmed clear deck')
    console.log('Current deck before clear:', deckStore.currentDeck)
    console.log('Cards before clear:', deckStore.currentDeck?.cards?.length)
    deckStore.clearCurrentDeck()
    console.log('Cards after clear:', deckStore.currentDeck?.cards?.length)
  } else {
    console.log('User cancelled clear deck')
  }
}

function getCardImageUrl(card: Card) {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.jpg'
}

onMounted(async () => {
  await Promise.all([
    cardStore.fetchAllCards(),
    deckStore.fetchUserDecks()
  ])
})
</script>