<template>
  <div>
    <h1 class="text-3xl font-bold text-white mb-8">カードコレクション</h1>
    
    <!-- Search and Filters -->
    <div class="mb-8 bg-black bg-opacity-30 rounded-lg p-6">
      <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
        <div>
          <label class="block text-white text-sm font-bold mb-2">検索</label>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="カード名で検索..."
            class="w-full px-3 py-2 bg-gray-800 text-white rounded border border-gray-600 focus:border-blue-500"
          />
        </div>
        
        <div>
          <label class="block text-white text-sm font-bold mb-2">タイプ</label>
          <select
            v-model="selectedType"
            class="w-full px-3 py-2 bg-gray-800 text-white rounded border border-gray-600 focus:border-blue-500"
          >
            <option value="">すべて</option>
            <option value="ふれんど">ふれんど</option>
            <option value="サポート">サポート</option>
            <option value="フィールド">フィールド</option>
          </select>
        </div>
        
        <div>
          <label class="block text-white text-sm font-bold mb-2">色</label>
          <select
            v-model="selectedColor"
            class="w-full px-3 py-2 bg-gray-800 text-white rounded border border-gray-600 focus:border-blue-500"
          >
            <option value="">すべて</option>
            <option value="赤">赤</option>
            <option value="青">青</option>
            <option value="黄">黄</option>
            <option value="緑">緑</option>
          </select>
        </div>
        
        <div>
          <label class="block text-white text-sm font-bold mb-2">レアリティ</label>
          <select
            v-model="selectedRarity"
            class="w-full px-3 py-2 bg-gray-800 text-white rounded border border-gray-600 focus:border-blue-500"
          >
            <option value="">すべて</option>
            <option value="C">C</option>
            <option value="U">U</option>
            <option value="R">R</option>
            <option value="SR">SR</option>
            <option value="SEC">SEC</option>
          </select>
        </div>
        
        <div>
          <label class="block text-white text-sm font-bold mb-2">カードタイプ</label>
          <select
            v-model="selectedVariant"
            class="w-full px-3 py-2 bg-gray-800 text-white rounded border border-gray-600 focus:border-blue-500"
          >
            <option value="">すべて</option>
            <option value="normal">通常カードのみ</option>
            <option value="promo">プロモカード (P)</option>
            <option value="parallel">パラレルカード (-P)</option>
          </select>
        </div>
      </div>
    </div>
    
    <!-- Loading State -->
    <div v-if="cardStore.loading" class="text-center text-white">
      <p>カードを読み込み中...</p>
    </div>
    
    <!-- Error State -->
    <div v-else-if="cardStore.error" class="text-center text-red-400">
      <p>エラー: {{ cardStore.error }}</p>
      <button
        @click="cardStore.fetchAllCards"
        class="btn-primary mt-4"
      >
        再試行
      </button>
    </div>
    
    <!-- Cards Grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
      <CardComponent
        v-for="card in filteredCards"
        :key="card.card_no"
        :card="card"
        @card-click="openCardModal"
      />
    </div>
    
    <!-- Card Modal -->
    <div
      v-if="selectedCard"
      class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50 p-4"
      @click="closeCardModal"
    >
      <div
        class="bg-gray-900 rounded-lg p-6 max-w-4xl mx-auto max-h-[90vh] overflow-y-auto"
        @click.stop
      >
        <div class="flex justify-between items-start mb-4">
          <h2 class="text-2xl font-bold text-white">{{ selectedCard.name }}</h2>
          <button
            @click="closeCardModal"
            class="text-gray-400 hover:text-white text-2xl font-bold"
          >
            ×
          </button>
        </div>
        
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="flex items-center justify-center">
            <img
              :src="getCardImageUrl(selectedCard)"
              :alt="selectedCard.name"
              class="max-w-full h-auto rounded-lg shadow-2xl"
              style="max-height: 70vh;"
            />
          </div>
          
          <div class="space-y-3 text-white">
            <div class="bg-gray-800 rounded-lg p-4">
              <h3 class="text-lg font-bold mb-3">カード情報</h3>
              <div class="space-y-2">
                <p><strong>カード番号:</strong> {{ selectedCard.card_no }}</p>
                <p><strong>タイプ:</strong> {{ selectedCard.type }}</p>
                <p><strong>色:</strong> {{ selectedCard.color }}</p>
                <p><strong>コスト:</strong> {{ selectedCard.cost }}
                  <span v-if="selectedCard.cost_red > 0" class="text-red-400 ml-2">(赤{{ selectedCard.cost_red }})</span>
                  <span v-if="selectedCard.cost_blue > 0" class="text-blue-400 ml-2">(青{{ selectedCard.cost_blue }})</span>
                  <span v-if="selectedCard.cost_yellow > 0" class="text-yellow-400 ml-2">(黄{{ selectedCard.cost_yellow }})</span>
                  <span v-if="selectedCard.cost_green > 0" class="text-green-400 ml-2">(緑{{ selectedCard.cost_green }})</span>
                  <span v-if="selectedCard.cost_colorless > 0" class="text-gray-400 ml-2">(無{{ selectedCard.cost_colorless }})</span>
                </p>
                <p v-if="selectedCard.power"><strong>パワー:</strong> {{ selectedCard.power }}</p>
                <p><strong>レアリティ:</strong> {{ selectedCard.rarity }}</p>
              </div>
            </div>
            
            <div v-if="selectedCard.effect" class="bg-gray-800 rounded-lg p-4">
              <h3 class="text-lg font-bold mb-2">効果</h3>
              <p class="text-gray-300">{{ selectedCard.effect }}</p>
            </div>
            
            <div v-if="selectedCard.flavor_text" class="bg-gray-800 rounded-lg p-4">
              <h3 class="text-lg font-bold mb-2">フレーバーテキスト</h3>
              <p class="text-gray-400 italic">{{ selectedCard.flavor_text }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useCardStore } from '@/stores/cards'
import CardComponent from '@/components/CardComponent.vue'
import type { Card, CardType, CardColor, CardRarity } from '@/types'

const cardStore = useCardStore()

const searchQuery = ref('')
const selectedType = ref<CardType | ''>('')
const selectedColor = ref<CardColor | ''>('')
const selectedRarity = ref<CardRarity | ''>('')
const selectedVariant = ref<'' | 'normal' | 'promo' | 'parallel'>('')
const selectedCard = ref<Card | null>(null)

const filteredCards = computed(() => {
  let filtered = cardStore.cards

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(card =>
      card.name.toLowerCase().includes(query) ||
      card.effect?.toLowerCase().includes(query)
    )
  }

  if (selectedType.value) {
    filtered = filtered.filter(card => card.type === selectedType.value)
  }

  if (selectedColor.value) {
    filtered = filtered.filter(card => card.color === selectedColor.value)
  }

  if (selectedRarity.value) {
    filtered = filtered.filter(card => card.rarity === selectedRarity.value)
  }

  if (selectedVariant.value) {
    if (selectedVariant.value === 'normal') {
      filtered = filtered.filter(card => !card.card_no.endsWith('-P') && !card.card_no.includes('(P)') && !card.rarity.endsWith('-P'))
    } else if (selectedVariant.value === 'promo') {
      filtered = filtered.filter(card => card.card_no.includes('(P)'))
    } else if (selectedVariant.value === 'parallel') {
      filtered = filtered.filter(card => card.card_no.endsWith('-P') || card.rarity.endsWith('-P'))
    }
  }

  return filtered
})

function openCardModal(card: Card) {
  selectedCard.value = card
}

function closeCardModal() {
  selectedCard.value = null
}

function getCardImageUrl(card: Card) {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.jpg'
}

onMounted(() => {
  cardStore.fetchAllCards()
})
</script>