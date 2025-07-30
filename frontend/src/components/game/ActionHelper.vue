<template>
  <transition name="helper">
    <div
      v-if="showHelper && helpMessage"
      class="action-helper fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-40"
    >
      <div class="bg-gradient-to-r from-blue-900 to-purple-900 rounded-lg px-8 py-4 shadow-2xl border-2 border-blue-400">
        <div class="flex items-center gap-4">
          <div class="text-3xl animate-bounce">
            {{ helpIcon }}
          </div>
          <div>
            <p class="text-white font-bold text-lg">{{ helpMessage }}</p>
            <p v-if="helpSubMessage" class="text-sm text-blue-200 mt-1">{{ helpSubMessage }}</p>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useGameStore } from '@/stores/game'

const gameStore = useGameStore()

const showHelper = ref(true)
const helpMessage = ref('')
const helpSubMessage = ref('')
const helpIcon = ref('💡')

const currentPhase = computed(() => gameStore.currentPhase)
const isPlayerTurn = computed(() => gameStore.isPlayerTurn)
const playerHand = computed(() => gameStore.player?.hand || [])
const playerEnergy = computed(() => gameStore.player?.energy || [])
const playerFriends = computed(() => gameStore.player?.friends || [])

// Update help message based on game state
watch([currentPhase, isPlayerTurn, playerHand, playerEnergy], () => {
  if (!isPlayerTurn.value) {
    helpMessage.value = ''
    return
  }

  switch (currentPhase.value) {
    case 'draw':
      helpMessage.value = 'カードを引きました！'
      helpSubMessage.value = 'エネルギーフェーズに進みます'
      helpIcon.value = '🎴'
      break
      
    case 'energy':
      helpMessage.value = 'エネルギーをセットしてください'
      helpSubMessage.value = '手札からカードを掴んでエネルギーエリアにドラッグ＆ドロップ（1枚まで）'
      helpIcon.value = '⚡'
      break
      
    case 'main':
      if (playerHand.value.length === 0) {
        helpMessage.value = '手札がありません'
        helpSubMessage.value = 'バトルフェーズに進みましょう'
        helpIcon.value = '🤔'
      } else if (playerEnergy.value.length === 0) {
        helpMessage.value = 'エネルギーカードをプレイしましょう'
        helpSubMessage.value = '手札のカードをクリックして配置'
        helpIcon.value = '⚡'
      } else {
        helpMessage.value = 'カードをプレイできます'
        helpSubMessage.value = '手札からカードを掴んでバトルエリアにドラッグ＆ドロップ！'
        helpIcon.value = '✨'
      }
      break
      
    case 'battle':
      if (playerFriends.value.length === 0) {
        helpMessage.value = 'ふれんどがいません'
        helpSubMessage.value = 'ターンを終了しましょう'
        helpIcon.value = '😅'
      } else {
        helpMessage.value = 'ふれんどで攻撃！'
        helpSubMessage.value = '緑の枠のふれんどをクリックして攻撃（アタック時効果は自動発動）'
        helpIcon.value = '⚔️'
      }
      break
      
    case 'end':
      helpMessage.value = 'ターン終了中...'
      helpSubMessage.value = ''
      helpIcon.value = '⏳'
      break
  }
})

// Hide helper after a delay
let hideTimeout: number
watch(helpMessage, () => {
  clearTimeout(hideTimeout)
  if (helpMessage.value) {
    showHelper.value = true
    hideTimeout = setTimeout(() => {
      showHelper.value = false
    }, 5000) as unknown as number
  }
})
</script>

<style scoped>
.helper-enter-active,
.helper-leave-active {
  transition: all 0.3s ease;
}

.helper-enter-from {
  transform: translate(-50%, -50%) scale(0.8);
  opacity: 0;
}

.helper-leave-to {
  transform: translate(-50%, -50%) scale(1.2);
  opacity: 0;
}
</style>