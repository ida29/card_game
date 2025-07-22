<template>
  <teleport to="body">
    <transition name="fade">
      <div
        v-if="counterDecision.show"
        class="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
      >
        <div class="counter-selection-container">
          <div class="bg-gray-900 rounded-lg p-6 max-w-4xl">
            <h2 class="text-2xl font-bold text-white mb-4 text-center">
              カウンタータイミング
            </h2>
            
            <!-- Battle Preview -->
            <div class="battle-preview mb-6">
              <div class="flex items-center justify-center gap-8">
                <!-- Attacker -->
                <div class="text-center">
                  <p class="text-gray-400 mb-2">アタッカー</p>
                  <GameCard
                    v-if="counterDecision.attacker.card"
                    :card="counterDecision.attacker.card.card"
                    size="medium"
                  />
                  <p class="text-red-400 font-bold mt-2">
                    パワー: {{ counterDecision.attacker.card?.card.power || 0 }}
                  </p>
                </div>
                
                <div class="text-4xl text-white">VS</div>
                
                <!-- Blocker -->
                <div class="text-center">
                  <p class="text-gray-400 mb-2">ブロッカー</p>
                  <GameCard
                    v-if="counterDecision.blocker?.card"
                    :card="counterDecision.blocker.card.card"
                    size="medium"
                  />
                  <p class="text-blue-400 font-bold mt-2">
                    パワー: {{ counterDecision.blocker?.card.card.power || 0 }}
                  </p>
                </div>
              </div>
            </div>
            
            <!-- Counter Cards -->
            <div class="mb-4">
              <p class="text-gray-300 mb-3">
                【カウンター】カードを使用できます（1枚まで）
              </p>
              
              <div v-if="availableCounters.length > 0" class="flex gap-3 flex-wrap justify-center">
                <div
                  v-for="(card, index) in availableCounters"
                  :key="index"
                  class="counter-card-option"
                  @click="selectCounter(index)"
                >
                  <GameCard
                    :card="card.card"
                    size="medium"
                  />
                </div>
              </div>
              
              <div v-else class="text-center text-gray-500 py-8">
                使用可能な【カウンター】カードがありません
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex justify-center gap-4">
              <button
                @click="skipCounter"
                class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-bold transition-colors"
              >
                カウンターを使用しない
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'
import GameCard from './GameCard.vue'

const gameStore = useGameStore()

const counterDecision = computed(() => gameStore.counterDecision)

// Get available counter cards from hand
const availableCounters = computed(() => {
  if (!gameStore.player?.hand) return []
  
  return gameStore.player.hand.filter(card => {
    // Check if card has counter ability
    // This should check for "【カウンター】" in card text or a counter flag
    return card.card.abilities?.includes('カウンター') || 
           card.card.description?.includes('【カウンター】')
  })
})

const selectCounter = (index: number) => {
  const card = availableCounters.value[index]
  gameStore.resolveCounter(card)
}

const skipCounter = () => {
  gameStore.resolveCounter(null)
}
</script>

<style scoped>
.counter-selection-container {
  max-height: 90vh;
  overflow-y: auto;
}

.counter-card-option {
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.counter-card-option:hover {
  transform: translateY(-5px);
}

.counter-card-option::after {
  content: '';
  position: absolute;
  inset: -4px;
  border: 2px solid transparent;
  border-radius: 0.5rem;
  transition: all 0.2s ease;
}

.counter-card-option:hover::after {
  border-color: #fbbf24;
  box-shadow: 0 0 20px rgba(251, 191, 36, 0.5);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>