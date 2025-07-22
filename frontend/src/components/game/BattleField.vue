<template>
  <div class="battle-field h-screen flex flex-col bg-gradient-to-b from-gray-900 to-black">
    <!-- Opponent Area -->
    <div class="opponent-area flex-1 flex flex-col justify-start p-4">
      <PlayerArea
        :player="opponent"
        :is-opponent="true"
        :is-active="!isPlayerTurn"
      />
    </div>
    
    <!-- Center Field -->
    <div class="center-field h-20 bg-gradient-to-r from-transparent via-purple-900/20 to-transparent flex items-center justify-center">
      <div class="text-white text-sm">
        <span v-if="currentPhase === 'battle'">バトルフェーズ</span>
        <span v-else-if="currentPhase === 'main'">メインフェーズ</span>
        <span v-else-if="currentPhase === 'draw'">ドローフェーズ</span>
        <span v-else>{{ currentPhase }}</span>
        <span class="ml-4">ターン {{ turnCount }}</span>
      </div>
    </div>
    
    <!-- Player Area -->
    <div class="player-area flex-1 flex flex-col justify-end p-4">
      <PlayerArea
        :player="player"
        :is-opponent="false"
        :is-active="isPlayerTurn"
        @play-card="handlePlayCard"
      />
    </div>
    
    <!-- Game Controls -->
    <GameControls
      v-if="isPlayerTurn && currentPhase !== 'game_over'"
      :current-phase="currentPhase"
      :can-end-turn="canEndTurn"
      @next-phase="handleNextPhase"
      @end-turn="handleEndTurn"
    />
    
    <!-- Game Over Screen -->
    <div v-if="gameState?.winner" class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
      <div class="bg-gray-800 rounded-lg p-8 text-center max-w-md">
        <h2 class="text-3xl font-bold mb-4" :class="gameState.winner === 'player' ? 'text-green-400' : 'text-red-400'">
          {{ gameState.winner === 'player' ? '勝利！' : '敗北...' }}
        </h2>
        <p class="text-white mb-6">
          {{ gameState.winner === 'player' ? 'おめでとうございます！' : '次は頑張りましょう！' }}
        </p>
        <div class="space-x-4">
          <button
            @click="handleRematch"
            class="btn-primary"
          >
            もう一度対戦
          </button>
          <button
            @click="handleReturnToMenu"
            class="btn-secondary"
          >
            メニューに戻る
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'
import { useRouter } from 'vue-router'
import PlayerArea from './PlayerArea.vue'
import GameControls from './GameControls.vue'

const gameStore = useGameStore()
const router = useRouter()

const gameState = computed(() => gameStore.gameState)
const player = computed(() => gameStore.player)
const opponent = computed(() => gameStore.opponent)
const isPlayerTurn = computed(() => gameStore.isPlayerTurn)
const currentPhase = computed(() => gameStore.currentPhase)
const turnCount = computed(() => gameStore.turnCount)
const canEndTurn = computed(() => currentPhase.value === 'end' || currentPhase.value === 'main')

const handlePlayCard = (cardIndex: number, targetZone: 'friends' | 'energy') => {
  gameStore.playCard('player', cardIndex, targetZone)
}

const handleNextPhase = () => {
  gameStore.nextPhase()
}

const handleEndTurn = () => {
  if (currentPhase.value !== 'end') {
    gameStore.currentPhase = 'end'
  }
  gameStore.endTurn()
}

const handleRematch = () => {
  location.reload() // Simple rematch by reloading
}

const handleReturnToMenu = () => {
  gameStore.resetGame()
  router.push('/game')
}
</script>

<style scoped>
.battle-field {
  background-image: 
    radial-gradient(ellipse at top, rgba(59, 130, 246, 0.1) 0%, transparent 50%),
    radial-gradient(ellipse at bottom, rgba(139, 92, 246, 0.1) 0%, transparent 50%);
}
</style>