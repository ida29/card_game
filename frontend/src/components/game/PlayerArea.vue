<template>
  <div class="player-area-wrapper">
    <div 
      class="player-area" 
      :class="{ 
        'opponent': isOpponent
      }"
    >
    <!-- Player Info -->
    <div class="player-info flex justify-between items-center mb-2">
      <div class="flex items-center gap-4">
        <h3 class="text-white font-bold px-3 py-1 rounded">
          {{ isOpponent ? '相手' : 'あなた' }}
        </h3>
      </div>
    </div>
    
    <!-- Battle and Energy Areas with Deck and Trash -->
    <div class="game-areas-container mb-4">
      <!-- Left side: Battle and Energy -->
      <div class="left-areas">
        <!-- Battle Area -->
        <div class="battle-zone mb-2" style="max-width: 500px;">
          <h4 class="text-gray-400 text-xs mb-1">バトルエリア</h4>
          <div 
            class="grid grid-cols-5 gap-2 p-3 bg-gray-800/50 rounded-lg min-h-[280px] transition-all"
            :class="{ 
              'bg-gray-700/50': isDraggingOverBattle && !isOpponent && gameStore.currentPhase === 'main',
              'border-2 border-dashed border-gray-500': isDraggingOverBattle && !isOpponent && gameStore.currentPhase === 'main'
            }"
            @dragover.prevent="handleDragOverBattle"
            @dragleave="handleDragLeaveBattle"
            @drop="handleDropBattle"
          >
            <!-- Display friends in grid (max 10 slots: 5 columns x 2 rows) -->
            <div
              v-for="index in 10"
              :key="index"
              class="friend-slot"
            >
              <div
                v-if="player?.friends[index - 1]"
                class="friend-card"
                :class="{ 
                  'can-attack': canAttack(index - 1),
                  'is-attacker': gameStore.selectedAttacker === index - 1 && !isOpponent,
                  'just-played': isJustPlayed(index - 1)
                }"
                @click="handleFriendClick(index - 1)"
                @contextmenu.prevent="handleFriendRightClick(index - 1)"
                @dblclick="showFullScreenCard(player.friends[index - 1].card.card)"
              >
                <GameCard 
                  :card="player.friends[index - 1].card.card" 
                  :size="'small'"
                  :selected="gameStore.selectedAttacker === index - 1 && !isOpponent"
                  :tapped="player.friends[index - 1].tapped"
                  :effective-power="gameStore.getEffectivePower(player.friends[index - 1], isOpponent ? 'opponent' : 'player')"
                />
                <!-- Just Played Indicator -->
                <div 
                  v-if="isJustPlayed(index - 1)" 
                  class="just-played-indicator"
                >
                  <div class="sleep-icon">💤</div>
                  <div class="sleep-text">召喚酔い</div>
                </div>
              </div>
              <div v-else class="empty-slot">
                <div class="w-20 h-28 border-2 border-dashed border-gray-600 rounded-lg"></div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Energy Area (below battle) -->
        <div class="energy-zone" style="max-width: 500px;">
          <h4 class="text-gray-400 text-xs mb-1">エネルギーエリア</h4>
          <div 
            class="flex gap-1 p-2 bg-gray-800/50 rounded-lg min-h-[120px] flex-wrap transition-all"
            :class="{ 
              'bg-gray-700/50': isDraggingOverEnergy && !isOpponent && gameStore.currentPhase === 'energy',
              'border-2 border-dashed border-green-500': isDraggingOverEnergy && !isOpponent && gameStore.currentPhase === 'energy'
            }"
            @dragover.prevent="handleDragOverEnergy"
            @dragleave="handleDragLeaveEnergy"
            @drop="handleDropEnergy"
          >
            <div
              v-for="(energy, index) in player?.energy || []"
              :key="index"
              class="energy-card"
              @click="showFullScreenCard(energy.card.card)"
            >
              <GameCard 
                :card="energy.card.card" 
                :size="'small'"
                :tapped="energy.tapped"
                :show-energy-value="true"
              />
            </div>
          </div>
        </div>
      </div>
      
      <!-- Right side: Deck and Trash -->
      <div class="right-areas">
        <!-- Negative Energy Area (above deck) -->
        <div class="negative-energy-zone mb-1">
          <h4 class="text-gray-400 text-xs mb-1">負のエネルギー</h4>
          <div 
            class="bg-purple-900/30 rounded-lg p-2 border border-purple-600/50 flex items-center justify-center min-h-[60px] cursor-pointer hover:bg-purple-900/40 transition-colors"
            @click="showNegativeEnergyList"
          >
            <div v-if="!player?.negativeEnergy || player.negativeEnergy.length === 0" class="text-purple-400 text-xs">
              空
            </div>
            <div v-else class="relative" :style="{ height: `${Math.max(80, 28 + (player.negativeEnergy.length - 1) * 14)}px`, width: '96px' }">
              <!-- Display all negative energy cards individually -->
              <div 
                v-for="(negativeCard, index) in player.negativeEnergy"
                :key="index"
                class="absolute cursor-pointer transition-all duration-200"
                :class="{ 'opacity-60': !negativeCard.faceUp }"
                :style="{
                  top: `${index * 14}px`,
                  left: '0',
                  zIndex: index
                }"
                @click.stop="handleNegativeEnergyClick(index)"
                @contextmenu.prevent="flipNegativeEnergy(index)"
              >
                <div class="transform rotate-90 origin-center">
                  <GameCard 
                    v-if="negativeCard.faceUp"
                    :card="negativeCard.card.card" 
                    :size="'small'"
                  />
                  <!-- Face-down card back -->
                  <div 
                    v-else
                    class="w-20 h-28 bg-black rounded-lg border-2 border-gray-600 flex items-center justify-center"
                  >
                    <span class="text-gray-400 text-xs font-bold">裏</span>
                  </div>
                </div>
              </div>
              <div class="absolute -bottom-4 right-0 bg-purple-900/70 rounded px-1 text-xs text-purple-200">
                {{ player.negativeEnergy.length }}枚
              </div>
            </div>
          </div>
        </div>
        
        <!-- Deck Area -->
        <div class="deck-zone mb-2">
          <h4 class="text-gray-400 text-xs mb-1">デッキ</h4>
          <div class="bg-gray-800/50 rounded-lg p-2 flex items-center justify-center">
            <div class="relative">
              <!-- Stack effect for deck -->
              <div class="absolute inset-0 bg-gray-700 rounded-lg transform translate-x-1 translate-y-1"></div>
              <div class="absolute inset-0 bg-gray-600 rounded-lg transform translate-x-0.5 translate-y-0.5"></div>
              <div class="relative bg-gradient-to-br from-gray-500 to-gray-700 rounded-lg w-20 h-28 flex items-center justify-center border-2 border-gray-500">
                <div class="text-center">
                  <span class="text-white text-2xl font-bold">{{ player?.deck.length || 0 }}</span>
                  <p class="text-gray-300 text-xs">枚</p>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Trash Area (below deck) -->
        <div class="trash-zone" style="margin-top: 12px;">
          <h4 class="text-gray-400 text-xs mb-1">トラッシュ</h4>
          <div 
            class="bg-gray-800/50 rounded-lg p-2 flex items-center justify-center cursor-pointer hover:bg-gray-800/60 transition-colors"
            @click="showTrashList"
          >
            <div class="relative">
              <!-- Maintain card size even when empty -->
              <div v-if="player?.graveyard.length === 0" class="w-20 h-28 border-2 border-dashed border-gray-600 rounded-lg flex items-center justify-center">
                <span class="text-gray-500 text-sm">空</span>
              </div>
              <div v-else class="relative">
                <!-- Show top card of trash -->
                <div>
                  <GameCard 
                    v-if="player.graveyard[player.graveyard.length - 1]"
                    :card="player.graveyard[player.graveyard.length - 1].card" 
                    :size="'small'"
                  />
                </div>
                <div class="absolute bottom-0 right-0 bg-black/70 rounded px-1 text-xs text-white">
                  {{ player.graveyard.length }}枚
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Field Card Area -->
    <div v-if="player?.field" class="field-zone mb-4">
      <h4 class="text-gray-400 text-xs mb-1">フィールドカード</h4>
      <div 
        class="p-2 bg-purple-900/30 rounded-lg cursor-pointer hover:bg-purple-900/40 transition-colors"
        @click="showFullScreenCard(player.field.card)"
      >
        <GameCard 
          :card="player.field.card" 
          :size="'small'"
        />
      </div>
    </div>
    
    <!-- Hand Container with proper centering -->
    <div class="hand-area-container">
      <!-- Hand Area (only for player) -->
      <div v-if="!isOpponent" class="hand-zone">
        <h4 class="text-gray-400 text-xs mb-1">手札</h4>
        <div class="hand-container relative h-48 flex items-end justify-center">
          <div class="hand-cards relative flex">
            <div
              v-for="(card, index) in player?.hand || []"
              :key="index"
              class="hand-card absolute"
              :class="{ 
                'draggable': canPlayCards, 
                'dragging': draggingCardIndex === index,
                'playable': canPlayCard(card)
              }"
              :style="getCardStyle(index, player?.hand.length || 0)"
              :draggable="canPlayCards"
              @dragstart="handleDragStart($event, index)"
              @dragend="handleDragEnd"
              @click="selectCard(index)"
              @mouseenter="hoveredCard = index"
              @mouseleave="hoveredCard = null"
              @dblclick="showFullScreenCard(card.card)"
            >
              <div class="relative">
                <GameCard 
                  :card="card.card" 
                  :size="'large'"
                  :selected="selectedCardIndex === index"
                />
                <!-- Playable indicator -->
                <div 
                  v-if="canPlayCard(card) && gameStore.currentPhase === 'main'"
                  class="playable-indicator"
                >
                  <div class="glow-effect"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
    </div>
    
    <!-- Full Screen Card Display -->
    <teleport to="body">
      <div 
        v-if="fullScreenCard"
        class="fixed inset-0 bg-black flex items-center justify-center z-50"
        @click="fullScreenCard = null"
      >
        <div class="relative bg-black overflow-hidden p-4" style="width: 90vmin; height: 90vmin; border-radius: 2rem;">
          <img 
            :src="getCardImageUrl(fullScreenCard)"
            :alt="fullScreenCard.name"
            class="absolute inset-0 w-full h-full object-contain"
            style="background-color: black; border-radius: 2rem;"
          />
        </div>
      </div>
    </teleport>
    
    <!-- Card List Modal for Trash -->
    <CardListModal
      :show="showingTrash"
      title="トラッシュ"
      :cards="player?.graveyard || []"
      @close="showingTrash = false"
    />
    
    <!-- Card List Modal for Negative Energy -->
    <CardListModal
      :show="showingNegativeEnergy"
      title="負のエネルギー"
      :cards="player?.negativeEnergy || []"
      :is-negative-energy="true"
      @close="showingNegativeEnergy = false"
    />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGameStore } from '@/stores/game'
import type { PlayerState, Card, DeckCard } from '@/types'
import GameCard from './GameCard.vue'
import CardListModal from './CardListModal.vue'

const props = defineProps<{
  player: PlayerState | undefined
  isOpponent: boolean
  isActive: boolean
}>()

const emit = defineEmits<{
  'play-card': [cardIndex: number, targetZone?: 'friends' | 'energy']
}>()

const gameStore = useGameStore()
const selectedCardIndex = ref<number | null>(null)
const hoveredCard = ref<number | null>(null)
const fullScreenCard = ref<Card | null>(null)
const draggingCardIndex = ref<number | null>(null)
const isDraggingOverBattle = ref(false)
const isDraggingOverEnergy = ref(false)
const showingTrash = ref(false)
const showingNegativeEnergy = ref(false)

const canPlayCards = computed(() => 
  !props.isOpponent && gameStore.canPlayCards
)

const canPlayCard = (deckCard: DeckCard) => {
  if (!props.player || !gameStore.canPlayCards) return false
  
  const card = deckCard.card
  const totalCost = card.cost || 0
  
  // エネルギーフェーズでは全てのカードが使用可能（エネルギーとして）
  if (gameStore.currentPhase === 'energy') return true
  
  // メインフェーズではコストをチェック
  if (gameStore.currentPhase === 'main') {
    // Calculate available energy value
    const availableEnergyValue = props.player.energy
      .filter(e => !e.tapped)
      .reduce((sum, e) => sum + (e.card.card.energy_value || 1), 0)
    
    // Check if we have enough total energy
    if (availableEnergyValue < totalCost) return false
    
    // Check specific color requirements
    const colorRequirements = {
      '赤': card.cost_red || 0,
      '青': card.cost_blue || 0,
      '黄': card.cost_yellow || 0,
      '緑': card.cost_green || 0
    }
    
    // Check each color requirement
    for (const [color, required] of Object.entries(colorRequirements)) {
      if (required > 0) {
        const colorEnergy = props.player.energy
          .filter(e => !e.tapped && e.card.card.color === color)
          .reduce((sum, e) => sum + (e.card.card.energy_value || 1), 0)
        
        if (colorEnergy < required) return false
      }
    }
    
    return true
  }
  
  return false
}

const canAttack = (index: number) => {
  if (!props.player?.friends[index]) return false
  const friend = props.player.friends[index]
  const cardNo = friend.card.card?.card_no
  
  return !props.isOpponent && 
         gameStore.currentPhase === 'main' && 
         gameStore.isPlayerTurn &&
         !friend.tapped && // Can't attack if already tapped
         (friend.playedTurn < gameStore.turnCount || (cardNo && gameStore.canAttackImmediately(cardNo))) // Can attack if played before this turn OR has immediate attack ability
}

const isJustPlayed = (index: number) => {
  if (!props.player?.friends[index]) return false
  const friend = props.player.friends[index]
  return friend.playedTurn === gameStore.turnCount
}

const selectCard = (index: number) => {
  // Only select card for viewing, don't play it
  if (selectedCardIndex.value === index) {
    selectedCardIndex.value = null
  } else {
    selectedCardIndex.value = index
  }
}

const handleFriendClick = (index: number) => {
  const friend = props.player?.friends[index]
  if (!friend) return
  
  if (props.isOpponent) {
    // Opponent's friends cannot be targeted directly
    // Just show the card in full screen
    showFullScreenCard(friend.card.card)
  } else {
    // If this is player's area
    if (canAttack(index)) {
      // If can attack, select for battle
      gameStore.selectBattleTarget(index)
    } else {
      // If can't attack, show the card in full screen
      showFullScreenCard(friend.card.card)
    }
  }
}

const handleFieldClick = () => {
  // Field clicks no longer needed for attacks
}

const handleFriendRightClick = (index: number) => {
  const friend = props.player?.friends[index]
  if (!friend || props.isOpponent) return
  
  // Right-click is disabled for now to avoid confusion
  // Main phase effects are accessed through the action choice modal
}

const handleNegativeEnergyClick = (index: number) => {
  // Do nothing - negative energy cards should not expand
  return
}

const flipNegativeEnergy = (index: number) => {
  if (props.isOpponent) return // Only allow flipping own cards
  if (!props.player?.negativeEnergy[index]) return
  
  const negativeCard = props.player.negativeEnergy[index]
  // Can only flip face-up cards to face-down (to use as energy)
  if (negativeCard.faceUp) {
    gameStore.flipNegativeEnergyCard(index)
  }
}

const showTrashList = () => {
  if (props.player?.graveyard && props.player.graveyard.length > 0) {
    showingTrash.value = true
  }
}

const showNegativeEnergyList = () => {
  if (props.player?.negativeEnergy && props.player.negativeEnergy.length > 0) {
    showingNegativeEnergy.value = true
  }
}

const getCardStyle = (index: number, totalCards: number) => {
  const centerIndex = (totalCards - 1) / 2
  const offset = index - centerIndex
  
  // Fixed spacing regardless of card count for consistent layout
  const baseSpacing = 25 // Fixed spacing between cards
  const maxWidth = 400 // Maximum spread width
  const cardWidth = 80 // Approximate card width
  
  // Calculate spacing - use base spacing unless cards would exceed max width
  const totalWidth = (totalCards - 1) * baseSpacing + cardWidth
  const spacing = totalWidth > maxWidth ? (maxWidth - cardWidth) / (totalCards - 1) : baseSpacing
  
  // Fixed rotation per card position
  const rotation = offset * 5 // 5 degrees per card
  const translateX = offset * spacing
  const translateY = Math.abs(offset) * Math.abs(offset) * 1.5 // Quadratic curve
  
  // Raise hovered card
  const isHovered = hoveredCard.value === index
  const hoverY = isHovered ? -30 : 0
  
  return {
    transform: `translateX(${translateX}px) translateY(${translateY + hoverY}px) rotate(${rotation}deg)`,
    zIndex: isHovered ? 100 : index,
    transition: 'all 0.2s ease'
  }
}

const showFullScreenCard = (card: Card) => {
  fullScreenCard.value = card
}

const getCardImageUrl = (card: Card) => {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.svg'
}

// Drag and Drop handlers
const handleDragStart = (event: DragEvent, index: number) => {
  draggingCardIndex.value = index
  event.dataTransfer!.effectAllowed = 'move'
  event.dataTransfer!.setData('cardIndex', index.toString())
}

const handleDragEnd = () => {
  draggingCardIndex.value = null
  isDraggingOverBattle.value = false
  isDraggingOverEnergy.value = false
}

const handleDragOverBattle = (event: DragEvent) => {
  if (gameStore.currentPhase === 'main' && draggingCardIndex.value !== null) {
    // Check if we can afford the card
    const card = props.player?.hand[draggingCardIndex.value]
    if (card) {
      const totalCost = card.card.cost || 0
      // Calculate total available energy value
      const availableEnergyValue = props.player?.energy
        .filter(e => !e.tapped)
        .reduce((sum, e) => sum + (e.card.card.energy_value || 1), 0) || 0
      
      console.log('handleDragOverBattle:', {
        cardName: card.card.name,
        cardType: card.card.type,
        totalCost,
        availableEnergyValue,
        canAfford: availableEnergyValue >= totalCost
      })
      
      if (availableEnergyValue >= totalCost) {
        event.preventDefault()
        isDraggingOverBattle.value = true
      }
    }
  }
}

const handleDragLeaveBattle = () => {
  isDraggingOverBattle.value = false
}

const handleDropBattle = (event: DragEvent) => {
  event.preventDefault()
  isDraggingOverBattle.value = false
  
  if (gameStore.currentPhase === 'main' && draggingCardIndex.value !== null) {
    // Let playCard determine the zone based on card type
    emit('play-card', draggingCardIndex.value)
    draggingCardIndex.value = null
  }
}

const handleDragOverEnergy = (event: DragEvent) => {
  if (!props.isOpponent && gameStore.currentPhase === 'energy' && draggingCardIndex.value !== null) {
    event.preventDefault()
    isDraggingOverEnergy.value = true
  }
}

const handleDragLeaveEnergy = () => {
  isDraggingOverEnergy.value = false
}

const handleDropEnergy = (event: DragEvent) => {
  event.preventDefault()
  isDraggingOverEnergy.value = false
  
  if (!props.isOpponent && gameStore.currentPhase === 'energy' && draggingCardIndex.value !== null) {
    emit('play-card', draggingCardIndex.value, 'energy')
    draggingCardIndex.value = null
  }
}
</script>

<style scoped>
.hand-card {
  cursor: pointer;
  transform-origin: bottom center;
}

.hand-card.draggable {
  cursor: grab;
}

.hand-card.draggable:active {
  cursor: grabbing;
}

.hand-card.dragging {
  opacity: 0.5;
  cursor: grabbing;
  transform: scale(1.1);
}

.hand-card.playable {
  position: relative;
}

.playable-indicator {
  position: absolute;
  inset: -2px;
  pointer-events: none;
}

.glow-effect {
  position: absolute;
  inset: 0;
  border: 3px solid #10b981;
  border-radius: 0.75rem;
  box-shadow: 
    0 0 15px #10b981,
    0 0 30px #10b981,
    inset 0 0 15px rgba(16, 185, 129, 0.3);
  animation: glow-pulse 2s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.7;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.02);
  }
}

.hand-container {
  perspective: 1000px;
}

.hand-cards {
  width: 400px; /* Reduced width for more compact display */
  height: 200px;
}

/* Remove all rotation for opponent - no longer needed */
.opponent {
  position: relative;
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

.friend-card {
  position: relative;
  transition: transform 0.2s ease;
}

.friend-card:hover {
  transform: translateY(-2px);
}

.friend-card.just-played {
  opacity: 0.9;
}

.just-played-indicator {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 0.5rem;
}

.sleep-icon {
  font-size: 2rem;
  animation: float 2s ease-in-out infinite;
}

.sleep-text {
  color: #fbbf24;
  font-size: 0.75rem;
  font-weight: bold;
  text-shadow: 0 0 4px rgba(0, 0, 0, 0.8);
  margin-top: 0.25rem;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-5px);
  }
}

/* Make cards clickable */
.energy-card {
  cursor: pointer;
}

.field-zone > div {
  cursor: pointer;
}

.game-areas-container {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.left-areas {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.right-areas {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  align-items: center;
  width: 8rem;
  padding-top: calc(1.25rem + 0.75rem + 112px + 0.5rem - 40px - 0.25rem - 1.25rem - 0.25rem - 1.25rem - 0.25rem + 8px); /* Fine-tuned alignment */
}

/* Keep the mirrored layout for opponent without rotation */
.opponent .game-areas-container {
  flex-direction: row-reverse;
}

/* Reverse the order of elements in left areas for opponent */
.opponent .left-areas {
  flex-direction: column-reverse;
}

/* Mirror the right areas padding for opponent */
.opponent .right-areas {
  padding-top: 0;
  padding-bottom: calc(1.25rem + 0.75rem + 112px + 0.5rem - 40px - 0.25rem - 1.25rem - 0.25rem - 1.25rem - 0.25rem + 8px);
}

/* Reverse the order of elements in right areas for opponent */
.opponent .right-areas {
  flex-direction: column-reverse;
}

/* Adjust trash margin for opponent */
.opponent .trash-zone {
  margin-top: 0;
  margin-bottom: 12px;
}

/* Wrapper to handle centering */
.player-area-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
}

/* Ensure proper alignment */
.player-area {
  width: 640px;
  display: flex;
  flex-direction: column;
  align-items: center;
  /* Ensure pixel-perfect alignment */
  transform: translateZ(0);
}

/* Center hand area with field */
.hand-area-container {
  display: flex;
  justify-content: center;
  width: 100%;
}

.hand-zone {
  width: 500px; /* Match battle area width */
}

/* Highlight opponent's field when attacking */
.player-area.attack-target {
  position: relative;
}

.player-area.attack-target::before {
  content: '';
  position: absolute;
  inset: -8px;
  border: 3px solid #ef4444;
  border-radius: 1rem;
  pointer-events: none;
  animation: attack-pulse 1.5s ease-in-out infinite;
  z-index: 1;
}

@keyframes attack-pulse {
  0%, 100% {
    opacity: 0.5;
    box-shadow: 0 0 20px rgba(239, 68, 68, 0.5);
  }
  50% {
    opacity: 1;
    box-shadow: 0 0 40px rgba(239, 68, 68, 0.8);
  }
}
</style>