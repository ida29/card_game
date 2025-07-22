<template>
  <div class="player-area" :class="{ 'opponent': isOpponent }">
    <!-- Player Info -->
    <div class="player-info flex justify-between items-center mb-2">
      <div class="flex items-center gap-4">
        <h3 class="text-white font-bold">
          {{ isOpponent ? '相手' : 'あなた' }}
        </h3>
        <div class="flex items-center gap-2">
          <span 
            class="text-red-400 px-2 py-1 rounded transition-all"
            :class="{
              'cursor-pointer hover:bg-red-900/50': canTargetPlayer,
              'ring-2 ring-red-500 animate-pulse': canTargetPlayer && gameStore.availableTargets.includes(-1)
            }"
            @click="handlePlayerClick"
          >
            ❤️ {{ player?.life || 0 }}
          </span>
          <span class="text-blue-400">🃏 {{ player?.deck.length || 0 }}</span>
        </div>
      </div>
      <div v-if="isActive" class="text-green-400 text-sm">
        アクティブ
      </div>
    </div>
    
    <!-- Battle Area -->
    <div class="battle-zone mb-4">
      <h4 class="text-gray-400 text-xs mb-1">バトルエリア</h4>
      <div class="flex gap-2 p-2 bg-gray-800/50 rounded-lg min-h-[120px]">
        <div
          v-for="(friend, index) in player?.friends || []"
          :key="index"
          class="friend-card"
          :class="{ 
            'can-attack': canAttack(index),
            'is-attacker': gameStore.selectedAttacker === index && !isOpponent,
            'is-target': gameStore.availableTargets.includes(index) && isOpponent
          }"
          @click="handleFriendClick(index)"
        >
          <GameCard 
            :card="friend.card" 
            :size="'small'"
            :selected="gameStore.selectedAttacker === index && !isOpponent"
          />
        </div>
        <div
          v-for="empty in emptyFriendSlots"
          :key="`empty-${empty}`"
          class="empty-slot"
        >
          <div class="w-16 h-20 border-2 border-dashed border-gray-600 rounded-lg"></div>
        </div>
      </div>
    </div>
    
    <!-- Energy Area -->
    <div class="energy-zone mb-4">
      <h4 class="text-gray-400 text-xs mb-1">エネルギーエリア</h4>
      <div class="flex gap-1 p-2 bg-gray-800/50 rounded-lg min-h-[60px] flex-wrap">
        <div
          v-for="(energy, index) in player?.energy || []"
          :key="index"
          class="energy-card"
        >
          <div class="w-10 h-12 bg-gradient-to-br from-yellow-600 to-orange-600 rounded border-2 border-yellow-400"></div>
        </div>
      </div>
    </div>
    
    <!-- Hand Area (only for player) -->
    <div v-if="!isOpponent" class="hand-zone">
      <h4 class="text-gray-400 text-xs mb-1">手札</h4>
      <div class="flex gap-2 p-2 bg-gray-800/50 rounded-lg min-h-[140px] overflow-x-auto">
        <div
          v-for="(card, index) in player?.hand || []"
          :key="index"
          class="hand-card"
          :class="{ 'draggable': canPlayCards }"
          @click="selectCard(index)"
        >
          <GameCard 
            :card="card.card" 
            :size="'medium'"
            :selected="selectedCardIndex === index"
          />
        </div>
      </div>
    </div>
    
    <!-- Opponent's hand (hidden) -->
    <div v-else class="hand-zone">
      <h4 class="text-gray-400 text-xs mb-1">手札</h4>
      <div class="flex gap-2 p-2 bg-gray-800/50 rounded-lg min-h-[100px]">
        <div
          v-for="n in player?.hand.length || 0"
          :key="n"
          class="hidden-card"
        >
          <div class="w-16 h-24 bg-gradient-to-br from-gray-700 to-gray-800 rounded-lg border-2 border-gray-600 flex items-center justify-center">
            <span class="text-gray-500 text-2xl">?</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGameStore } from '@/stores/game'
import type { PlayerState } from '@/types'
import GameCard from './GameCard.vue'

const props = defineProps<{
  player: PlayerState | undefined
  isOpponent: boolean
  isActive: boolean
}>()

const emit = defineEmits<{
  'play-card': [cardIndex: number, targetZone: 'friends' | 'energy']
}>()

const gameStore = useGameStore()
const selectedCardIndex = ref<number | null>(null)

const canPlayCards = computed(() => 
  !props.isOpponent && gameStore.canPlayCards
)

const emptyFriendSlots = computed(() => {
  const currentFriends = props.player?.friends.length || 0
  return Math.max(0, 5 - currentFriends)
})

const canAttack = (index: number) => {
  return !props.isOpponent && 
         gameStore.currentPhase === 'battle' && 
         gameStore.isPlayerTurn &&
         props.player?.friends[index]
}

const canTargetPlayer = computed(() => {
  return props.isOpponent &&
         gameStore.selectedAttacker !== null &&
         gameStore.availableTargets.includes(-1)
})

const selectCard = (index: number) => {
  if (!canPlayCards.value) return
  
  if (selectedCardIndex.value === index) {
    // Play the card
    const hasEnergy = (props.player?.energy.length || 0) > 0
    const targetZone = hasEnergy ? 'friends' : 'energy'
    emit('play-card', index, targetZone)
    selectedCardIndex.value = null
  } else {
    selectedCardIndex.value = index
  }
}

const handleFriendClick = (index: number) => {
  if (props.isOpponent) {
    // If this is opponent's area and we're selecting a target
    if (gameStore.selectedAttacker !== null && gameStore.availableTargets.includes(index)) {
      gameStore.executeBattle(index)
    }
  } else {
    // If this is player's area and it's battle phase
    if (canAttack(index)) {
      gameStore.selectBattleTarget(index)
    }
  }
}

const handlePlayerClick = () => {
  if (canTargetPlayer.value) {
    gameStore.executeBattle('player')
  }
}
</script>

<style scoped>
.hand-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.hand-card:hover {
  transform: translateY(-10px);
}

.hand-card.draggable {
  cursor: grab;
}

.opponent {
  transform: rotate(180deg);
}

.opponent > * {
  transform: rotate(180deg);
}

.friend-card.can-attack {
  cursor: pointer;
  position: relative;
}

.friend-card.can-attack::after {
  content: '';
  position: absolute;
  inset: -4px;
  border: 2px solid #10b981;
  border-radius: 0.5rem;
  opacity: 0.5;
  animation: pulse 2s infinite;
}

.friend-card.is-attacker::after {
  border-color: #3b82f6;
  opacity: 1;
}

.friend-card.is-target {
  cursor: pointer;
}

.friend-card.is-target::after {
  content: '';
  position: absolute;
  inset: -4px;
  border: 2px solid #ef4444;
  border-radius: 0.5rem;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 0.5;
  }
  50% {
    transform: scale(1.05);
    opacity: 1;
  }
}
</style>