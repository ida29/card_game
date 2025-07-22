<template>
  <teleport to="body">
    <transition name="fade">
      <div 
        v-if="blockingDecision.show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
      >
        <div class="bg-gray-800 rounded-lg p-6 max-w-4xl w-full mx-4">
          <h2 class="text-2xl font-bold text-white mb-4 text-center">
            相手が攻撃してきました！
          </h2>
          
          <!-- Attacker Display -->
          <div class="flex justify-center mb-6">
            <div class="text-center">
              <p class="text-gray-400 mb-2">攻撃カード</p>
              <GameCard 
                v-if="blockingDecision.attacker.card"
                :card="blockingDecision.attacker.card.card" 
                :size="'large'"
              />
            </div>
          </div>
          
          <p class="text-gray-300 text-center mb-6">
            ブロックするカードを選択するか、「ブロックしない」を選んでください
          </p>
          
          <!-- Available Blockers -->
          <div class="mb-6">
            <p class="text-gray-400 mb-2">ブロック可能なカード:</p>
            <div class="flex flex-wrap justify-center gap-4">
              <div
                v-for="index in blockingDecision.availableBlockers"
                :key="index"
                class="cursor-pointer transform hover:scale-105 transition-all"
                @click="selectBlocker(index)"
              >
                <div class="relative">
                  <GameCard 
                    v-if="player?.friends[index]"
                    :card="player.friends[index].card.card" 
                    :size="'medium'"
                  />
                  <div class="absolute inset-0 border-2 border-blue-500 rounded-lg pointer-events-none"></div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- No Block Button -->
          <div class="flex justify-center">
            <button
              @click="selectBlocker(null)"
              class="px-8 py-3 bg-red-600 hover:bg-red-700 text-white rounded-lg font-bold transition-colors text-lg"
            >
              ブロックしない（ダイレクトアタックを受ける）
            </button>
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

const blockingDecision = computed(() => gameStore.blockingDecision)
const player = computed(() => gameStore.player)

const selectBlocker = (index: number | null) => {
  gameStore.resolveBlocking(index)
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