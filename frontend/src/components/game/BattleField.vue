<template>
  <div class="battle-field h-screen flex flex-col bg-gradient-to-b from-gray-900/95 to-gray-900/95">
    <!-- Loading State -->
    <div v-if="!gameState" class="flex items-center justify-center h-full">
      <div class="text-white text-xl">ゲームを初期化中...</div>
    </div>
    
    <!-- Game Content -->
    <template v-else>
      <!-- Opponent's Hand at the very top -->
      <div class="opponent-hand-top py-2">
        <div class="flex justify-center">
          <div class="hand-zone" style="width: 500px;">
            <div class="hand-container relative h-24 flex items-center justify-center">
              <div class="hand-cards relative flex">
                <div
                  v-for="(_, index) in opponent?.hand || []"
                  :key="index"
                  class="hidden-card absolute"
                  :style="getOpponentCardStyle(index, opponent?.hand.length || 0)"
                >
                  <div class="w-16 h-24 bg-gradient-to-br from-gray-700 to-gray-800 border-2 border-gray-600 flex items-center justify-center shadow-lg" style="border-radius: 0.5rem;">
                    <span class="text-gray-500 text-2xl">?</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Phase Indicator -->
      <PhaseIndicator />
      
      <!-- Turn Indicator -->
      <TurnIndicator />
      
      <!-- Action Helper -->
      <ActionHelper />
      
      <!-- Opponent Area -->
      <div class="opponent-area flex-1 flex flex-col justify-start p-4">
        <div class="flex justify-center">
          <PlayerArea
            :player="opponent"
            :is-opponent="true"
            :is-active="!isPlayerTurn"
          />
        </div>
      </div>
    
    <!-- Center Field -->
    <div class="center-field h-20 bg-gradient-to-r from-transparent via-purple-900/20 to-transparent flex items-center justify-center">
      <!-- Battle Effects Area -->
      <transition name="battle-effect">
        <div v-if="showBattleEffect" class="battle-effect-container">
          <div class="text-6xl animate-pulse">⚔️</div>
        </div>
      </transition>
    </div>
    
    <!-- Player Area -->
    <div class="player-area flex-1 flex flex-col justify-end p-4">
      <div class="flex justify-center">
        <PlayerArea
          :player="player"
          :is-opponent="false"
          :is-active="isPlayerTurn"
          @play-card="handlePlayCard"
        />
      </div>
    </div>
    
    <!-- Game Controls -->
    <GameControls
      v-if="isPlayerTurn && currentPhase !== 'game_over'"
      :current-phase="currentPhase"
      :can-end-turn="canEndTurn"
      :energy-played-this-turn="energyPlayedThisTurn"
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
    
    <!-- Battle Animation -->
    <BattleAnimation
      :show="gameStore.battleAnimation.show"
      :attacker-card="gameStore.battleAnimation.attacker?.card"
      :defender-card="gameStore.battleAnimation.defender?.card"
      :attacker-defeated="gameStore.battleAnimation.attackerDefeated"
      :defender-defeated="gameStore.battleAnimation.defenderDefeated"
      @animation-complete="gameStore.hideBattleAnimation"
    />
    
    <!-- Blocking Decision -->
    <BlockingDecision />
    
    <!-- Counter Selection -->
    <CounterSelection />
    
    <!-- Energy Cost Selector -->
    <EnergyCostSelector
      v-if="energyCostSelection && gameStore"
      :show="energyCostSelection.show"
      :card-to-pay="energyCostSelection.cardToPay"
      :resolve="energyCostSelection.resolve"
      @cancel="gameStore?.cancelEnergyCostSelection"
    />
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useGameStore } from '@/stores/game'
import { useRouter } from 'vue-router'
import PlayerArea from './PlayerArea.vue'
import GameControls from './GameControls.vue'
import PhaseIndicator from './PhaseIndicator.vue'
import TurnIndicator from './TurnIndicator.vue'
import ActionHelper from './ActionHelper.vue'
import BattleAnimation from './BattleAnimation.vue'
import BlockingDecision from './BlockingDecision.vue'
import CounterSelection from './CounterSelection.vue'
import EnergyCostSelector from './EnergyCostSelector.vue'

const gameStore = useGameStore()
const router = useRouter()

const gameState = computed(() => gameStore.gameState)
const player = computed(() => gameStore.player)
const opponent = computed(() => gameStore.opponent)
const isPlayerTurn = computed(() => gameStore.isPlayerTurn)
const currentPhase = computed(() => gameStore.currentPhase)
const turnCount = computed(() => gameStore.turnCount)
const canEndTurn = computed(() => currentPhase.value === 'end' || currentPhase.value === 'main')
const showBattleEffect = ref(false)
const energyPlayedThisTurn = computed(() => gameStore.energyPlayedThisTurn)

// Safe access to energyCostSelection
const energyCostSelection = computed(() => {
  console.log('gameStore:', gameStore)
  console.log('gameStore.energyCostSelection:', gameStore?.energyCostSelection)
  
  if (!gameStore || !gameStore.energyCostSelection) {
    return {
      show: false,
      cardToPay: null,
      resolve: null
    }
  }
  
  return gameStore.energyCostSelection
})

const handlePlayCard = async (cardIndex: number, targetZone: 'friends' | 'energy') => {
  await gameStore.playCard('player', cardIndex, targetZone)
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

const getOpponentCardStyle = (index: number, totalCards: number) => {
  const centerIndex = (totalCards - 1) / 2
  const offset = index - centerIndex
  const spread = Math.min(totalCards * 8, 60) // Smaller spread for opponent
  const rotation = (offset * spread) / totalCards
  const translateX = offset * 15 // Smaller spacing
  
  return {
    transform: `translateX(${translateX}px) rotate(${rotation}deg)`,
    zIndex: index
  }
}
</script>

<style scoped>
.battle-field {
  background-image: 
    radial-gradient(ellipse at top, rgba(59, 130, 246, 0.05) 0%, transparent 50%),
    radial-gradient(ellipse at bottom, rgba(139, 92, 246, 0.05) 0%, transparent 50%);
}

.opponent-hand-top .hand-cards {
  position: relative;
  height: 100%;
  width: 400px;
}

.opponent-hand-top .hidden-card {
  transform-origin: center center;
}
</style>