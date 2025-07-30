<template>
  <teleport to="body">
    <transition name="modal">
      <div v-if="show" class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
        <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
          <h3 class="text-xl font-bold text-white mb-4">アクションを選択</h3>
          
          <div class="space-y-3">
            <!-- Main Phase Effect Option -->
            <button
              v-if="hasMainPhaseEffect"
              @click="handleMainEffect"
              class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-4 rounded transition-colors flex items-center justify-center gap-2"
              :disabled="!canUseEffect"
              :class="{ 'opacity-50 cursor-not-allowed': !canUseEffect }"
            >
              <span class="text-2xl">💪</span>
              <span>{{ mainEffectDescription }}</span>
            </button>
            
            <!-- Energy requirement notice -->
            <div v-if="hasMainPhaseEffect && !canUseEffect" class="text-sm text-red-400 text-center">
              エネルギーが不足しています
            </div>
            
            <!-- Cancel -->
            <button
              @click="handleCancel"
              class="w-full bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition-colors"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FriendCardState } from '@/types'
import { useGameStore } from '@/stores/game'

const props = defineProps<{
  show: boolean
  friendState: FriendCardState | null
}>()

const emit = defineEmits<{
  'use-main-effect': []
  'cancel': []
}>()

const gameStore = useGameStore()

const hasMainPhaseEffect = computed(() => {
  if (!props.friendState?.card?.card) return false
  const cardNo = props.friendState.card.card.card_no
  
  // List of cards with main phase effects
  const mainPhaseEffectCards = [
    'F-002', // なみだぶくろん
  ]
  
  return mainPhaseEffectCards.includes(cardNo)
})

const mainEffectDescription = computed(() => {
  if (!props.friendState?.card?.card) return ''
  const cardNo = props.friendState.card.card.card_no
  
  const effectDescriptions: Record<string, string> = {
    'F-002': '効果を使用する（コスト①）：パワー+1000（このターン）',
  }
  
  return effectDescriptions[cardNo] || '効果を発動'
})

const canUseEffect = computed(() => {
  if (!gameStore.player) return false
  // Check if player has at least one untapped energy
  const availableEnergy = gameStore.player.energy.filter(e => !e.tapped)
  return availableEnergy.length > 0
})

const handleMainEffect = () => {
  emit('use-main-effect')
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s;
}

.modal-enter-from, .modal-leave-to {
  opacity: 0;
}
</style>