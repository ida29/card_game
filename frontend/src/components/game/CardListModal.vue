<template>
  <teleport to="body">
    <transition name="fade">
      <div
        v-if="show"
        class="fixed inset-0 bg-black/90 flex items-center justify-center z-50 p-4"
        @click="$emit('close')"
      >
        <div 
          class="bg-gray-900 rounded-lg p-6 max-w-6xl max-h-[90vh] overflow-hidden flex flex-col"
          @click.stop
        >
          <!-- Header -->
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-2xl font-bold text-white">{{ title }}</h2>
            <button
              @click="$emit('close')"
              class="text-gray-400 hover:text-white transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <!-- Card count -->
          <div class="text-gray-400 mb-4">
            合計: {{ cards.length }}枚
          </div>
          
          <!-- Card grid -->
          <div class="flex-1 overflow-y-auto">
            <div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
              <div
                v-for="(card, index) in cards"
                :key="index"
                class="cursor-pointer transform transition-transform hover:scale-105"
                @click="handleCardClick(card, index)"
              >
                <GameCard
                  v-if="!isNegativeEnergy || card.faceUp"
                  :card="isNegativeEnergy ? card.card.card : card.card"
                  size="medium"
                />
                <!-- Face-down negative energy card -->
                <div 
                  v-else
                  class="w-full aspect-[5/7] bg-black rounded-lg border-2 border-gray-600 flex items-center justify-center"
                >
                  <span class="text-gray-400 text-sm font-bold">裏</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Full screen card display -->
        <div 
          v-if="selectedCard"
          class="fixed inset-0 bg-black bg-opacity-95 flex items-center justify-center z-60 p-8"
          @click="selectedCard = null"
        >
          <img 
            :src="getCardImageUrl(isNegativeEnergy ? selectedCard.card.card : selectedCard.card)"
            :alt="isNegativeEnergy ? selectedCard.card.card.name : selectedCard.card.name"
            class="max-w-full max-h-full object-contain"
            style="background-color: black; max-height: calc(100vh - 4rem); border-radius: 2rem;"
          />
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { DeckCard, NegativeEnergyCardState } from '@/types'
import GameCard from './GameCard.vue'

const props = defineProps<{
  show: boolean
  title: string
  cards: DeckCard[] | NegativeEnergyCardState[]
  isNegativeEnergy?: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const selectedCard = ref<any>(null)

const handleCardClick = (card: any, index: number) => {
  // Don't allow expansion for face-down negative energy cards
  if (props.isNegativeEnergy && !card.faceUp) {
    return
  }
  
  // For trash and face-up negative energy cards, don't allow expansion
  // Just return without setting selectedCard
  return
}

const getCardImageUrl = (card: any) => {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.svg'
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>