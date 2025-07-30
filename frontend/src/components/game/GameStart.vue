<template>
  <div class="game-start fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50">
    <div class="bg-gray-800 rounded-lg p-8 max-w-6xl max-h-[90vh] w-full overflow-auto flex flex-col justify-center">
      <!-- Deck Shuffle Phase -->
      <div v-if="phase === 'shuffle'" class="text-center">
        <h2 class="text-3xl font-bold text-white mb-6 animate-fade-in">ゲーム準備中</h2>
        <div class="max-w-2xl mx-auto">
          <div class="mb-8">
            <!-- Animated card shuffle display -->
            <div class="shuffle-container relative h-64 mb-8">
              <!-- Multiple cards shuffling -->
              <div 
                v-for="i in 12" 
                :key="i"
                class="absolute card-shuffle"
                :class="`card-${i}`"
                :style="{
                  animationDelay: `${(i - 1) * 0.1}s`,
                  left: '50%',
                  top: '50%',
                  transform: 'translate(-50%, -50%)'
                }"
              >
                <div class="w-20 h-28 bg-gradient-to-br from-gray-600 to-gray-800 rounded-lg shadow-2xl border-2 border-gray-500">
                  <div class="w-full h-full flex items-center justify-center">
                    <span class="text-gray-500 text-3xl">?</span>
                  </div>
                </div>
              </div>
              
              <!-- Central glow effect -->
              <div class="absolute inset-0 flex items-center justify-center">
                <div class="w-32 h-32 bg-purple-500 rounded-full opacity-30 animate-pulse-glow blur-xl"></div>
              </div>
            </div>
            
            <p class="text-gray-300 mb-4 animate-fade-in">{{ shuffleMessage }}</p>
            
          </div>
        </div>
      </div>
      
      <!-- Drawing Initial Hand -->
      <div v-else-if="phase === 'drawing'" class="text-center">
        <h2 class="text-3xl font-bold text-white mb-6">初期手札を配っています...</h2>
        <div class="flex justify-center gap-2 mb-8">
          <div
            v-for="i in 5"
            :key="i"
            class="card-draw-animation"
            :style="{ animationDelay: `${i * 0.2}s` }"
          >
            <div class="w-20 h-28 bg-gradient-to-br from-gray-600 to-gray-700 rounded-lg shadow-xl"></div>
          </div>
        </div>
        <p class="text-gray-300">{{ drawMessage }}</p>
      </div>
      
      <!-- Mulligan Phase -->
      <div v-else-if="phase === 'mulligan'" class="text-center">
        <h2 class="text-3xl font-bold text-white mb-4">初期手札の確認</h2>
        <p class="text-gray-300 mb-6">手札を確認して、1回だけ引き直すことができます</p>
        
        <!-- Show Player's Initial Hand -->
        <div class="flex justify-center gap-4 mb-8">
          <div
            v-for="(card, index) in playerHand"
            :key="index"
            class="transform hover:scale-110 transition-all cursor-pointer"
            @click="selectedCard = card"
          >
            <GameCard 
              :card="card.card" 
              :size="'large'"
            />
          </div>
        </div>
        
        <!-- Full Screen Card Display -->
        <div 
          v-if="selectedCard"
          class="fixed inset-0 bg-black bg-opacity-95 flex items-center justify-center z-50 p-8"
          @click="selectedCard = null"
        >
          <img 
            :src="getCardImageUrl(selectedCard.card)"
            :alt="selectedCard.card.name"
            class="max-w-full max-h-full object-contain"
            style="background-color: black; max-height: calc(100vh - 4rem); border-radius: 2rem;"
            @load="onImageLoad"
          />
        </div>
        
        <div class="flex justify-center gap-4">
          <button
            @click="keepHand"
            class="px-6 py-3 bg-green-600 hover:bg-green-700 text-white rounded-lg font-bold transition-colors"
          >
            この手札でOK
          </button>
          <button
            @click="mulligan"
            :disabled="hasMulliganed"
            class="px-6 py-3 bg-orange-600 hover:bg-orange-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-lg font-bold transition-colors"
          >
            {{ hasMulliganed ? '引き直し済み' : '手札を引き直す' }}
          </button>
        </div>
        
        <p v-if="hasMulliganed" class="text-yellow-400 mt-4">
          引き直しは1回のみです
        </p>
      </div>
      
      <!-- Rock Paper Scissors -->
      <div v-else-if="phase === 'janken'" class="text-center">
        <h2 class="text-3xl font-bold text-white mb-6">先攻・後攻を決めよう！</h2>
        <p class="text-gray-300 mb-8">じゃんけんで勝った方が先攻・後攻を選べます</p>
        
        <div class="grid grid-cols-3 gap-4 max-w-md mx-auto mb-8">
          <button
            v-for="choice in choices"
            :key="choice.value"
            @click="!playerChoice && selectChoice(choice.value)"
            class="relative group"
            :class="{ 
              'opacity-50': playerChoice && playerChoice !== choice.value,
              'ring-4 ring-blue-500': playerChoice === choice.value
            }"
          >
            <div class="bg-gray-700 rounded-lg p-6 transition-all group-hover:scale-110">
              <div class="text-6xl mb-2">{{ choice.emoji }}</div>
              <div class="text-white font-bold">{{ choice.name }}</div>
            </div>
          </button>
        </div>
        
        <!-- Result -->
        <div v-if="playerChoice && cpuChoice" class="mt-8">
          <div class="flex justify-center items-center gap-8 mb-4">
            <div class="text-center">
              <p class="text-gray-400 mb-2">あなた</p>
              <div class="text-6xl">{{ getChoiceEmoji(playerChoice) }}</div>
            </div>
            <div class="text-4xl text-yellow-400 animate-pulse">VS</div>
            <div class="text-center">
              <p class="text-gray-400 mb-2">CPU</p>
              <div class="text-6xl">{{ getChoiceEmoji(cpuChoice) }}</div>
            </div>
          </div>
          
          <div class="text-2xl font-bold mt-6" :class="resultColor">
            {{ resultMessage }}
          </div>
        </div>
      </div>
      
      <!-- Turn Order Selection -->
      <div v-else-if="phase === 'turn-order'" class="text-center">
        <h2 class="text-3xl font-bold text-white mb-6">じゃんけんに勝ちました！</h2>
        <p class="text-gray-300 mb-8">先攻・後攻を選んでください</p>
        
        <div class="grid grid-cols-2 gap-8 max-w-2xl mx-auto">
          <button
            @click="selectTurnOrder('first')"
            class="bg-gradient-to-br from-blue-600 to-blue-800 hover:from-blue-700 hover:to-blue-900 rounded-lg p-8 transition-all transform hover:scale-105"
          >
            <div class="text-white">
              <div class="text-6xl mb-4">☀️</div>
              <h3 class="text-2xl font-bold mb-2">先攻</h3>
              <p class="text-sm opacity-90">最初にターンを開始します</p>
            </div>
          </button>
          
          <button
            @click="selectTurnOrder('second')"
            class="bg-gradient-to-br from-purple-600 to-purple-800 hover:from-purple-700 hover:to-purple-900 rounded-lg p-8 transition-all transform hover:scale-105"
          >
            <div class="text-white">
              <div class="text-6xl mb-4">🌙</div>
              <h3 class="text-2xl font-bold mb-2">後攻</h3>
              <p class="text-sm opacity-90">相手の後にターンを開始します</p>
            </div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useGameStore } from '@/stores/game'
import GameCard from './GameCard.vue'
import type { DeckCard, Card } from '@/types'

const emit = defineEmits<{
  'game-ready': []
}>()

const gameStore = useGameStore()

const phase = ref<'shuffle' | 'drawing' | 'mulligan' | 'janken' | 'turn-order'>('shuffle')
const playerChoice = ref<'rock' | 'paper' | 'scissors' | null>(null)
const cpuChoice = ref<'rock' | 'paper' | 'scissors' | null>(null)
const resultMessage = ref('')
const resultColor = ref('')
const shuffleMessage = ref('デッキをシャッフル中...')
const shuffleProgress = ref(0)
const drawMessage = ref('5枚のカードを配布中...')
const hasMulliganed = ref(false)
const selectedCard = ref<DeckCard | null>(null)

const playerHand = computed(() => gameStore.player?.hand || [])

const choices = [
  { value: 'rock' as const, emoji: '✊', name: 'グー' },
  { value: 'scissors' as const, emoji: '✌️', name: 'チョキ' },
  { value: 'paper' as const, emoji: '✋', name: 'パー' }
]

const getChoiceEmoji = (choice: string) => {
  return choices.find(c => c.value === choice)?.emoji || '?'
}

const selectChoice = (choice: 'rock' | 'paper' | 'scissors') => {
  playerChoice.value = choice
  
  // CPU makes choice
  setTimeout(() => {
    const cpuChoices: ('rock' | 'paper' | 'scissors')[] = ['rock', 'paper', 'scissors']
    cpuChoice.value = cpuChoices[Math.floor(Math.random() * 3)]
    
    // Determine winner
    const result = determineWinner(playerChoice.value!, cpuChoice.value)
    
    if (result === 'win') {
      resultMessage.value = 'あなたの勝ち！'
      resultColor.value = 'text-green-400'
    } else if (result === 'lose') {
      resultMessage.value = 'CPUの勝ち！CPUが先攻・後攻を選びます'
      resultColor.value = 'text-red-400'
    } else {
      resultMessage.value = 'あいこ！もう一度！'
      resultColor.value = 'text-yellow-400'
      setTimeout(() => {
        playerChoice.value = null
        cpuChoice.value = null
      }, 1500)
      return
    }
    
    // Move to turn order selection if win/lose, otherwise retry
    setTimeout(() => {
      if (result === 'win') {
        phase.value = 'turn-order'
      } else if (result === 'lose') {
        // CPU chooses randomly
        const cpuChoice = Math.random() < 0.5 ? 'first' : 'second'
        if (cpuChoice === 'first') {
          gameStore.currentPlayer = 'opponent'
          resultMessage.value = 'CPUが先攻を選びました！'
        } else {
          gameStore.currentPlayer = 'player'
          resultMessage.value = 'CPUが後攻を選びました！あなたが先攻です！'
        }
        
        // Start game after showing the message
        setTimeout(() => {
          emit('game-ready')
        }, 2000)
      }
    }, 2000)
  }, 500)
}

const determineWinner = (player: string, cpu: string): 'win' | 'lose' | 'draw' => {
  if (player === cpu) return 'draw'
  
  if (
    (player === 'rock' && cpu === 'scissors') ||
    (player === 'scissors' && cpu === 'paper') ||
    (player === 'paper' && cpu === 'rock')
  ) {
    return 'win'
  }
  
  return 'lose'
}

const startDrawingAnimation = () => {
  let cardsDrawn = 0
  const drawInterval = setInterval(() => {
    cardsDrawn++
    drawMessage.value = `${cardsDrawn}枚目を配布中...`
    
    if (cardsDrawn >= 5) {
      clearInterval(drawInterval)
      drawMessage.value = '準備完了！'
      
      setTimeout(() => {
        phase.value = 'mulligan'
      }, 1000)
    }
  }, 400)
}

const keepHand = () => {
  // Move to janken phase
  phase.value = 'janken'
}

const mulligan = () => {
  if (hasMulliganed.value) return
  
  hasMulliganed.value = true
  
  // Return cards to deck and shuffle
  gameStore.mulliganHand('player')
  
  // Show redraw animation
  const currentPhase = phase.value
  phase.value = 'drawing'
  drawMessage.value = '手札を引き直しています...'
  
  setTimeout(() => {
    // Only return to mulligan if we haven't moved to janken
    if (phase.value === 'drawing') {
      phase.value = 'mulligan'
    }
  }, 2000)
}

const getCardImageUrl = (card: Card) => {
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.svg'
}

const onImageLoad = (e: Event) => {
  const img = e.target as HTMLImageElement
  // Force black background on the image element itself
  img.style.backgroundColor = 'black'
}

const selectTurnOrder = (order: 'first' | 'second') => {
  if (order === 'first') {
    gameStore.currentPlayer = 'player'
  } else {
    gameStore.currentPlayer = 'opponent'
  }
  
  console.log('Turn order selected:', order, 'currentPlayer:', gameStore.currentPlayer)
  emit('game-ready')
}

const startShuffleAnimation = () => {
  shuffleProgress.value = 0
  shuffleMessage.value = 'デッキをシャッフル中...'
  
  // Wait for animation to complete (3s + some buffer)
  setTimeout(() => {
    shuffleMessage.value = 'シャッフル完了！'
    
    // Move to drawing phase
    setTimeout(() => {
      phase.value = 'drawing'
      startDrawingAnimation()
    }, 1000)
  }, 3500)
}

onMounted(() => {
  // Check if game is properly initialized
  if (!gameStore.gameState) {
    console.error('GameStart: No game state found')
    emit('game-ready') // Skip to prevent getting stuck
    return
  }
  
  // Reset game state
  gameStore.turnCount = 0
  
  // Start shuffle animation
  startShuffleAnimation()
})
</script>

<style scoped>
/* Card shuffle animation */
@keyframes shuffleCard {
  0% {
    transform: translate(-50%, -50%) rotate(0deg) translateX(0);
    opacity: 0.8;
  }
  25% {
    transform: translate(-50%, -50%) rotate(180deg) translateX(150px) scale(1.1);
    opacity: 1;
  }
  50% {
    transform: translate(-50%, -50%) rotate(360deg) translateX(0) translateY(-100px) scale(0.9);
    opacity: 0.8;
  }
  75% {
    transform: translate(-50%, -50%) rotate(540deg) translateX(-150px) scale(1.05);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) rotate(720deg) translateX(0);
    opacity: 0.8;
  }
}

.card-shuffle {
  animation: shuffleCard 3s ease-in-out;
  animation-fill-mode: forwards;
}

/* Different paths for each card */
.card-1 { animation-duration: 2.8s; }
.card-2 { animation-duration: 3.2s; animation-direction: reverse; }
.card-3 { animation-duration: 2.6s; }
.card-4 { animation-duration: 3.4s; animation-direction: reverse; }
.card-5 { animation-duration: 2.9s; }
.card-6 { animation-duration: 3.1s; animation-direction: reverse; }
.card-7 { animation-duration: 2.7s; }
.card-8 { animation-duration: 3.3s; animation-direction: reverse; }
.card-9 { animation-duration: 2.5s; }
.card-10 { animation-duration: 3.5s; animation-direction: reverse; }
.card-11 { animation-duration: 2.8s; }
.card-12 { animation-duration: 3.2s; animation-direction: reverse; }

/* Fade in animation */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fadeIn 0.6s ease-out;
}

/* Pulse glow animation */
@keyframes pulseGlow {
  0%, 100% {
    transform: scale(1);
    opacity: 0.3;
  }
  50% {
    transform: scale(1.5);
    opacity: 0.5;
  }
}

.animate-pulse-glow {
  animation: pulseGlow 2s ease-in-out infinite;
}

/* Shimmer effect */
@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(200%);
  }
}

.shimmer {
  animation: shimmer 2s linear infinite;
}

/* Card draw animation (existing) */
@keyframes cardDraw {
  0% {
    transform: translateY(-100px) rotate(-20deg);
    opacity: 0;
  }
  50% {
    transform: translateY(0) rotate(10deg);
    opacity: 1;
  }
  100% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
}

.card-draw-animation {
  animation: cardDraw 0.6s ease-out forwards;
  opacity: 0;
}

/* Force black background on all images */
:deep(img) {
  background-color: black !important;
  image-rendering: -webkit-optimize-contrast;
  image-rendering: crisp-edges;
}
</style>