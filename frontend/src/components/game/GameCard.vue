<template>
  <div 
    class="game-card"
    :class="[
      sizeClass,
      { 'selected': selected, 'tapped': tapped }
    ]"
  >
    <div v-if="card" class="card-content bg-gradient-to-br from-gray-700 to-gray-800 rounded-lg border-2 p-1 h-full flex flex-col">
      <!-- Card Header -->
      <div class="card-header flex justify-between items-center text-xs mb-1">
        <span class="cost font-bold text-yellow-400">{{ card.cost }}</span>
        <span class="type text-gray-400">{{ card.type }}</span>
      </div>
      
      <!-- Card Image/Name -->
      <div class="card-body flex-1 flex items-center justify-center">
        <p class="text-white text-xs text-center font-medium">{{ card.name }}</p>
      </div>
      
      <!-- Card Footer -->
      <div v-if="card.power" class="card-footer text-center">
        <span class="power text-red-400 font-bold text-sm">{{ card.power }}</span>
      </div>
    </div>
    <div v-else class="card-placeholder bg-gray-700 rounded-lg h-full"></div>
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
}>()

const sizeClass = computed(() => {
  switch (props.size) {
    case 'small':
      return 'w-16 h-20'
    case 'large':
      return 'w-32 h-40'
    default:
      return 'w-20 h-28'
  }
})
</script>

<style scoped>
.game-card {
  transition: all 0.2s;
  cursor: pointer;
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
  border-color: #4b5563;
}

.game-card.selected .card-content {
  border-color: #3b82f6;
}
</style>