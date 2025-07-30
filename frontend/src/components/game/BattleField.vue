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
    
    <!-- Drawn Card Display -->
    <transition name="drawn-card-fade">
      <div 
        v-if="gameStore.drawnCardDisplay.show && gameStore.drawnCardDisplay.card"
        class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50"
        style="pointer-events: none;"
      >
        <div class="drawn-card-container">
          <div class="draw-effect-text">DRAW!</div>
          <img 
            :src="getCardImageUrl(gameStore.drawnCardDisplay.card.card)"
            :alt="gameStore.drawnCardDisplay.card.card.name"
            class="drawn-card-image"
          />
          <div class="card-name-label">{{ gameStore.drawnCardDisplay.card.card.name }}</div>
        </div>
      </div>
    </transition>
    
    <!-- Action Choice Modal -->
    <ActionChoiceModal
      :show="gameStore.actionChoice.show"
      :friend-state="gameStore.actionChoice.friendState"
      :friend-index="gameStore.actionChoice.friendIndex"
      @choose-attack="gameStore.executeActionChoice('attack')"
      @choose-main-effect="gameStore.executeActionChoice('main-effect')"
      @cancel="gameStore.cancelActionChoice"
    />
    
    <!-- Main Phase Action Modal -->
    <MainPhaseActionModal
      :show="gameStore.mainPhaseAction.show"
      :friend-state="gameStore.mainPhaseAction.friendState"
      @use-main-effect="gameStore.executeMainPhaseAction"
      @cancel="gameStore.cancelMainPhaseAction"
    />
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
import ActionChoiceModal from './ActionChoiceModal.vue'
import MainPhaseActionModal from './MainPhaseActionModal.vue'

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
const energyPlayedThisTurn = computed(() => gameStore.energyPlayedThisTurn || { player: false, opponent: false })

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
  // Navigate first, then reset to avoid the initialization screen
  router.push('/game')
  // Reset after navigation to clean up state
  setTimeout(() => {
    gameStore.resetGame()
  }, 100)
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

const getCardImageUrl = (card: any) => {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.svg'
}
</script>

<style scoped>
.battle-field {
  background-image: 
    radial-gradient(ellipse at top, rgba(59, 130, 246, 0.05) 0%, transparent 50%),
    radial-gradient(ellipse at bottom, rgba(139, 92, 246, 0.05) 0%, transparent 50%);
}

/* Drawn Card Display */
.drawn-card-container {
  position: relative;
  animation: draw-card-appear 0.5s ease-out;
}

.draw-effect-text {
  position: absolute;
  top: -60px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 48px;
  font-weight: bold;
  color: #fbbf24;
  text-shadow: 
    0 0 20px rgba(251, 191, 36, 0.8),
    0 0 40px rgba(251, 191, 36, 0.6),
    2px 2px 4px rgba(0, 0, 0, 0.8);
  animation: draw-text-pulse 1s ease-in-out;
}

.drawn-card-image {
  width: 300px;
  height: 420px;
  object-fit: cover;
  border-radius: 16px;
  box-shadow: 
    0 0 30px rgba(255, 255, 255, 0.5),
    0 10px 40px rgba(0, 0, 0, 0.8);
  animation: card-glow 2s ease-in-out;
}

.card-name-label {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.9);
  padding: 8px 20px;
  border-radius: 8px;
  color: white;
  font-weight: bold;
  font-size: 18px;
  white-space: nowrap;
}

/* Animations */
@keyframes draw-card-appear {
  from {
    transform: scale(0) rotate(180deg);
    opacity: 0;
  }
  to {
    transform: scale(1) rotate(0deg);
    opacity: 1;
  }
}

@keyframes draw-text-pulse {
  0%, 100% {
    transform: translateX(-50%) scale(1);
  }
  50% {
    transform: translateX(-50%) scale(1.2);
  }
}

@keyframes card-glow {
  0%, 100% {
    filter: brightness(1);
  }
  50% {
    filter: brightness(1.2) drop-shadow(0 0 20px rgba(255, 255, 255, 0.8));
  }
}

/* Transition */
.drawn-card-fade-enter-active,
.drawn-card-fade-leave-active {
  transition: opacity 0.3s ease;
}

.drawn-card-fade-enter-from,
.drawn-card-fade-leave-to {
  opacity: 0;
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