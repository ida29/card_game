<template>
  <div class="mulligan-card h-full">
    <div class="card-content rounded-lg border-2 border-gray-600 h-full overflow-hidden relative bg-gray-900">
      <!-- Card Image Background -->
      <div 
        v-if="!isCPUCard && cardImageUrl"
        class="absolute inset-0 bg-cover bg-center opacity-30"
        :style="{ backgroundImage: `url(${cardImageUrl})` }"
      />
      
      <!-- Card Info -->
      <div class="relative h-full flex flex-col p-4">
        <!-- Card Header -->
        <div class="flex justify-between items-center mb-3">
          <span class="cost text-2xl font-bold text-yellow-400 bg-black/80 rounded-lg px-3 py-2">
            コスト: {{ card.cost }}
          </span>
          <span class="type text-lg text-white bg-black/80 rounded-lg px-3 py-2">
            {{ card.type }}
          </span>
        </div>
        
        <!-- Card Name -->
        <div class="mb-3 text-center">
          <h3 class="text-white text-2xl font-bold bg-black/80 rounded-lg px-4 py-2 inline-block">
            {{ card.name }}
          </h3>
        </div>
        
        <!-- Card Image Section -->
        <div class="flex-1 mb-4 rounded-lg overflow-hidden bg-black/50 max-h-32">
          <img 
            v-if="cardImageUrl"
            :src="cardImageUrl" 
            :alt="card.name"
            class="w-full h-full object-cover rounded-lg"
            @error="onImageError"
          />
          <div v-else class="w-full h-full flex items-center justify-center">
            <span class="text-gray-500 text-5xl">📋</span>
          </div>
        </div>
        
        <!-- Card Details -->
        <div class="space-y-2">
          <!-- Color and Rarity -->
          <div class="flex justify-between items-center">
            <span 
              class="text-sm px-2 py-1 rounded font-semibold"
              :class="colorClasses[card.color]"
            >
              {{ card.color }}
            </span>
            <span 
              class="text-sm px-2 py-1 rounded font-semibold"
              :class="rarityClasses[card.rarity]"
            >
              {{ card.rarity }}
            </span>
          </div>
          
          <!-- Power -->
          <div v-if="card.power" class="text-center">
            <span class="power text-2xl font-bold text-red-400 bg-black/70 rounded px-4 py-2 inline-block">
              パワー: {{ card.power }}
            </span>
          </div>
          
          <!-- Effect Text -->
          <div v-if="card.effect" class="bg-black/90 rounded-lg p-4 max-h-32 overflow-y-auto border border-gray-600">
            <p class="text-base text-white leading-relaxed whitespace-pre-wrap">{{ card.effect }}</p>
          </div>
          
          <!-- Flavor Text -->
          <div v-if="card.flavor_text && !card.effect" class="bg-black/60 rounded p-2">
            <p class="text-xs text-gray-300 italic">{{ card.flavor_text }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '@/types'

const props = defineProps<{
  card: Card
}>()

const cardImageUrl = computed(() => {
  if (props.card.local_image_path) {
    return `/api/v1/images/${props.card.local_image_path.replace('card_images/', '')}`
  }
  return props.card.image_url || ''
})

const isCPUCard = computed(() => {
  return props.card.card_no?.startsWith('CPU-') || false
})

const colorClasses: Record<string, string> = {
  '赤': 'bg-red-600 text-white',
  '青': 'bg-blue-600 text-white',
  '黄': 'bg-yellow-600 text-black',
  '緑': 'bg-green-600 text-white',
  '紫': 'bg-purple-600 text-white'
}

const rarityClasses: Record<string, string> = {
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

const onImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}
</script>

<style scoped>
.mulligan-card {
  transition: all 0.3s ease;
}

.mulligan-card:hover {
  filter: brightness(1.1);
}
</style>