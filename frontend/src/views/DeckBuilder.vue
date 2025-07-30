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
          v-for="deck in displayDecks"
          :key="deck.ID"
          class="bg-gray-800 rounded-lg p-4 cursor-pointer hover:bg-gray-700 transition-colors relative"
          :class="{ 'ring-2 ring-blue-500': deckStore.currentDeck?.ID === deck.ID }"
          @click="selectDeck(deck)"
        >
          <!-- Deck Thumbnail -->
          <div class="mb-3 flex justify-center">
            <div class="relative w-20 h-28">
              <!-- Stacked cards effect -->
              <div class="absolute inset-0 bg-gray-700 rounded transform rotate-3 translate-x-1"></div>
              <div class="absolute inset-0 bg-gray-600 rounded transform -rotate-2 -translate-x-1"></div>
              <!-- Main card on top -->
              <div class="absolute inset-0">
                <img 
                  v-if="getMainCardImage(deck)"
                  :src="getMainCardImage(deck)"
                  :alt="deck.name"
                  class="w-full h-full object-cover rounded-lg shadow-lg"
                />
                <div v-else class="w-full h-full bg-gray-500 rounded-lg flex items-center justify-center flex-col p-2">
                  <span class="text-gray-300 text-2xl mb-1">?</span>
                  <span class="text-gray-400 text-[8px] text-center">カードを追加</span>
                </div>
              </div>
            </div>
          </div>
          <h3 class="text-white font-bold text-center">{{ deck.name }}</h3>
          <p class="text-gray-400 text-sm text-center">
            {{ getDeckCardCount(deck) }}/50 カード
          </p>
        </div>
      </div>
    </div>
    
    <!-- Current Deck Info -->
    <div v-if="deckStore.currentDeck" class="mb-8 bg-black bg-opacity-30 rounded-lg p-6">
      <div class="flex justify-between items-center mb-4">
        <div class="flex items-center gap-4">
          <!-- Main Card Display -->
          <div class="relative">
            <div class="relative w-24 h-32">
              <!-- Stacked cards effect -->
              <div class="absolute inset-0 bg-gray-700 rounded transform rotate-3 translate-x-1"></div>
              <div class="absolute inset-0 bg-gray-600 rounded transform -rotate-2 -translate-x-1"></div>
              <!-- Main card -->
              <div class="absolute inset-0 cursor-pointer group" @click="openMainCardSelector">
                <img 
                  v-if="getCurrentMainCardImage()"
                  :src="getCurrentMainCardImage()"
                  :alt="deckStore.currentDeck.name"
                  class="w-full h-full object-cover rounded-lg shadow-lg"
                />
                <div v-else class="w-full h-full bg-gray-500 rounded-lg flex items-center justify-center flex-col p-2">
                  <span class="text-gray-300 text-3xl mb-1">?</span>
                  <span class="text-gray-400 text-[9px] text-center">クリックで設定</span>
                </div>
                <!-- Hover overlay -->
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 rounded-lg transition-all flex items-center justify-center">
                  <span class="text-white text-xs opacity-0 group-hover:opacity-100">クリックで変更</span>
                </div>
              </div>
            </div>
            <p class="text-gray-400 text-xs mt-1 text-center">メインカード</p>
          </div>
          
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
              class="w-full aspect-[2/3] object-cover rounded-lg mb-1"
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
              class="w-full aspect-[2/3] object-cover rounded-lg mb-1"
            />
            <div v-else class="w-full aspect-[2/3] bg-gray-700 rounded-lg mb-1 flex items-center justify-center">
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
    
    <!-- Main Card Selector Modal -->
    <div v-if="showMainCardSelector && deckStore.currentDeck" class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-6 max-w-4xl max-h-[80vh] overflow-auto">
        <h3 class="text-xl font-bold text-white mb-4">メインカードを選択</h3>
        <p class="text-gray-400 mb-4">デッキのサムネイルとして表示するカードを選んでください</p>
        
        <div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
          <div
            v-for="deckCard in deckStore.currentDeck.cards"
            :key="deckCard.card_no"
            class="bg-gray-700 rounded-lg p-2 cursor-pointer hover:bg-gray-600 transition-colors"
            :class="{ 'ring-2 ring-blue-500': tempMainCardNo === deckCard.card_no }"
            @click="tempMainCardNo = deckCard.card_no"
          >
            <img
              v-if="deckCard.card"
              :src="getCardImageUrl(deckCard.card)"
              :alt="deckCard.card.name"
              class="w-full aspect-[2/3] object-cover rounded mb-1"
            />
            <p class="text-white text-xs font-bold truncate">{{ deckCard.card?.name || 'Unknown' }}</p>
          </div>
        </div>
        
        <div class="mt-6 flex justify-end gap-2">
          <button
            @click="cancelMainCardSelection"
            class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition-colors"
          >
            キャンセル
          </button>
          <button
            @click="confirmMainCardSelection"
            :disabled="!tempMainCardNo"
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold py-2 px-4 rounded transition-colors"
          >
            決定
          </button>
        </div>
      </div>
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
const showMainCardSelector = ref(false)
const tempMainCardNo = ref<string | undefined>(undefined)

// Computed property to ensure deck list is reactive
const displayDecks = computed(() => {
  return deckStore.decks
})

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

async function addCardToDeck(card: Card) {
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

function getMainCardImage(deck: Deck): string | null {
  if (!deck.cards || deck.cards.length === 0) {
    return null
  }
  
  // If main_card_no is set, try to find that card
  if (deck.main_card_no) {
    const mainCard = deck.cards.find(dc => dc.card_no === deck.main_card_no)
    if (mainCard?.card) {
      return getCardImageUrl(mainCard.card)
    }
  }
  
  // For existing decks without main_card_no, use the first card as fallback
  // This maintains visual consistency for older decks
  const firstCard = deck.cards[0]
  if (firstCard?.card) {
    return getCardImageUrl(firstCard.card)
  }
  
  return null
}

function getCurrentMainCardImage(): string | null {
  if (!deckStore.currentDeck) return null
  return getMainCardImage(deckStore.currentDeck)
}

function openMainCardSelector() {
  if (!deckStore.currentDeck) return
  tempMainCardNo.value = deckStore.currentDeck.main_card_no
  showMainCardSelector.value = true
}

function cancelMainCardSelection() {
  showMainCardSelector.value = false
  tempMainCardNo.value = undefined
}

async function confirmMainCardSelection() {
  if (!deckStore.currentDeck || !tempMainCardNo.value) return
  
  // Store the previous value in case we need to revert
  const previousMainCard = deckStore.currentDeck.main_card_no
  
  // Update the local value
  deckStore.currentDeck.main_card_no = tempMainCardNo.value
  
  // Close the modal first to provide immediate feedback
  showMainCardSelector.value = false
  
  // Save the deck to persist the main card selection
  try {
    await deckStore.updateDeck(deckStore.currentDeck)
    console.log('Main card saved:', tempMainCardNo.value)
    tempMainCardNo.value = undefined
  } catch (error) {
    console.error('Failed to save main card:', error)
    // Revert to previous value on error
    if (deckStore.currentDeck) {
      deckStore.currentDeck.main_card_no = previousMainCard
    }
    alert('メインカードの保存に失敗しました')
  }
}

onMounted(async () => {
  // Load cards first, then decks (so deck loading can use card data)
  await cardStore.fetchAllCards()
  await deckStore.loadDecks()
})
</script>