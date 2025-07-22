<template>
  <div 
    class="card relative overflow-hidden transform transition-all duration-300 hover:scale-105 cursor-pointer"
    @click="$emit('card-click', card)"
  >
    <div class="aspect-[2/3] overflow-hidden rounded-t-lg">
      <img 
        :src="cardImageUrl"
        :alt="card.name"
        class="w-full h-full object-cover"
        @error="onImageError"
      />
    </div>
    
    <div class="p-3">
      <div class="flex justify-between items-start mb-2">
        <h3 class="text-sm font-bold text-white truncate flex-1">{{ card.name }}</h3>
        <span 
          class="text-xs px-2 py-1 rounded ml-2"
          :class="rarityColors[card.rarity]"
        >
          {{ displayRarity }}
        </span>
      </div>
      
      <div class="flex justify-between items-center mb-2">
        <span 
          class="text-xs px-2 py-1 rounded"
          :class="typeColors[card.type]"
        >
          {{ card.type }}
        </span>
        <span 
          class="text-xs px-2 py-1 rounded"
          :class="colorClasses[card.color]"
        >
          {{ card.color }}
        </span>
      </div>
      
      <div class="flex justify-between items-center text-xs text-gray-300">
        <div class="flex items-center gap-1">
          <span>コスト: {{ card.cost }}</span>
          <span v-if="card.cost_red > 0" class="text-red-500">(赤{{ card.cost_red }})</span>
          <span v-if="card.cost_blue > 0" class="text-blue-500">(青{{ card.cost_blue }})</span>
          <span v-if="card.cost_yellow > 0" class="text-yellow-500">(黄{{ card.cost_yellow }})</span>
          <span v-if="card.cost_green > 0" class="text-green-500">(緑{{ card.cost_green }})</span>
          <span v-if="card.cost_colorless > 0" class="text-gray-400">(無{{ card.cost_colorless }})</span>
        </div>
        <span v-if="card.power">パワー: {{ card.power }}</span>
      </div>
      
      <div v-if="card.effect" class="mt-2">
        <p class="text-xs text-gray-400 line-clamp-2">{{ card.effect }}</p>
      </div>
    </div>
    
    <div v-if="showQuantity && quantity" class="absolute top-2 right-2 bg-blue-600 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center font-bold">
      {{ quantity }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '@/types'

interface Props {
  card: Card
  quantity?: number
  showQuantity?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  quantity: 0,
  showQuantity: false,
})

defineEmits<{
  'card-click': [card: Card]
}>()

const cardImageUrl = computed(() => {
  if (props.card.local_image_path) {
    return `/api/v1/images/${props.card.local_image_path.replace('card_images/', '')}`
  }
  return props.card.image_url || '/placeholder-card.svg'
})

const isParallel = computed(() => {
  return props.card.card_no.endsWith('-P') || props.card.rarity.endsWith('-P')
})

const isPromo = computed(() => {
  return props.card.card_no.includes('(P)') || props.card.is_promo
})

const displayRarity = computed(() => {
  // Just return the rarity as-is
  // Promo cards are identified by (P) in card number, not rarity
  // Parallel cards have -P in their rarity (e.g., SR-P)
  return props.card.rarity
})

const rarityColors = {
  'C': 'bg-gray-600 text-white',
  'U': 'bg-green-600 text-white',
  'R': 'bg-blue-600 text-white',
  'SR': 'bg-purple-600 text-white',
  'SEC': 'bg-yellow-600 text-black',
  'C-P': 'bg-gradient-to-r from-gray-600 to-purple-600 text-white',
  'U-P': 'bg-gradient-to-r from-green-600 to-purple-600 text-white',
  'R-P': 'bg-gradient-to-r from-blue-600 to-purple-600 text-white',
  'SR-P': 'bg-gradient-to-r from-purple-600 to-pink-600 text-white',
  'SEC-P': 'bg-gradient-to-r from-yellow-600 to-purple-600 text-black',
}

const typeColors = {
  'ふれんど': 'bg-orange-600 text-white',
  'サポート': 'bg-green-600 text-white',
  'フィールド': 'bg-blue-600 text-white',
}

const colorClasses = {
  '赤': 'bg-red-600 text-white',
  '青': 'bg-blue-600 text-white',
  '黄': 'bg-yellow-600 text-black',
  '緑': 'bg-green-600 text-white',
}

function onImageError(event: Event) {
  const img = event.target as HTMLImageElement
  img.src = '/placeholder-card.svg'
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>