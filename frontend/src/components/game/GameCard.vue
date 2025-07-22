<template>
  <div 
    class="game-card"
    :class="[
      sizeClass,
      { 'selected': selected, 'tapped': tapped }
    ]"
  >
    <div v-if="card" class="card-content rounded-lg border-2 border-gray-700 h-full overflow-hidden relative bg-gray-900">
      <!-- Card Image Background -->
      <div 
        v-if="hasValidImage"
        class="absolute inset-0"
      >
        <img 
          :src="cardImageUrl"
          :alt="card.name"
          class="absolute inset-0 w-full h-full object-cover rounded-lg"
          @error="onImageError"
        />
        <div class="absolute inset-0 bg-gradient-to-b from-transparent via-transparent to-black/80"></div>
      </div>
      
      <!-- Placeholder Card Background -->
      <div 
        v-else
        class="absolute inset-0 bg-gradient-to-br rounded-lg"
        :class="placeholderGradient"
      >
        <div class="absolute inset-0 bg-black/20 rounded-lg"></div>
      </div>
      
      <!-- Card Info Overlay -->
      <div class="relative h-full flex flex-col p-1">
        <!-- Card Header -->
        <div class="card-header flex justify-between items-center text-xs mb-auto">
          <span class="cost font-bold text-yellow-400 bg-black/50 rounded px-1">{{ card.cost }}</span>
          <span class="type text-white bg-black/50 rounded px-1 text-[10px]">{{ card.type }}</span>
        </div>
        
        <!-- Energy Value Indicator (for energy cards) -->
        <div v-if="showEnergyValue" class="absolute top-1 left-1">
          <span class="energy-value font-bold text-green-400 bg-black/70 rounded-full px-2 py-1 text-xs">
            ⚡{{ card.energy_value || 1 }}
          </span>
        </div>
        
        <!-- Card Name -->
        <div class="card-body mt-auto">
          <p class="text-white text-xs text-center font-bold bg-black/60 rounded px-1 py-0.5">{{ card.name }}</p>
        </div>
        
        <!-- Card Footer -->
        <div v-if="card.power" class="card-footer text-center mt-1">
          <span class="power text-red-400 font-bold text-sm bg-black/60 rounded px-2">{{ card.power }}</span>
        </div>
      </div>
    </div>
    <div v-else class="card-placeholder bg-gray-800 rounded-lg h-full border-2 border-gray-700"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Card } from '@/types'

const props = defineProps<{
  card: Card | undefined
  size?: 'small' | 'medium' | 'large'
  selected?: boolean
  tapped?: boolean
  showEnergyValue?: boolean
}>()

const sizeClass = computed(() => {
  switch (props.size) {
    case 'small':
      return 'w-20 h-28'
    case 'large':
      return 'w-28 h-40'
    default:
      return 'w-24 h-32'
  }
})

const cardImageUrl = computed(() => {
  if (!props.card) return ''
  
  if (props.card.local_image_path && props.card.local_image_path !== 'card_images/placeholder.jpg') {
    return `/api/v1/images/${props.card.local_image_path.replace('card_images/', '')}`
  }
  return props.card.image_url || ''
})

const hasValidImage = computed(() => {
  const url = cardImageUrl.value
  return url && url !== '' && !url.includes('placeholder')
})

const placeholderGradient = computed(() => {
  if (!props.card) return ''
  
  const colorGradients: Record<string, string> = {
    '赤': 'from-red-700 to-red-900',
    '青': 'from-blue-700 to-blue-900',
    '黄': 'from-yellow-700 to-yellow-900',
    '緑': 'from-green-700 to-green-900',
    '紫': 'from-purple-700 to-purple-900'
  }
  
  return colorGradients[props.card.color] || 'from-gray-700 to-gray-900'
})

const onImageError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}
</script>

<style scoped>
.game-card {
  transition: all 0.2s;
  cursor: pointer;
  background-color: black;
  border-radius: 0.5rem;
}

.game-card:hover {
  transform: scale(1.05);
  z-index: 10;
}

.game-card.selected {
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.8);
  transform: scale(1.1);
}

.game-card.tapped {
  transform: rotate(90deg);
}

.card-content {
  background-clip: padding-box;
}

.game-card.selected .card-content {
  border-color: #3b82f6;
}
</style>