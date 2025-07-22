<template>
  <div class="turn-indicator">
    <div class="turn-display" :class="{ 'player-turn': isPlayerTurn, 'opponent-turn': !isPlayerTurn }">
      <div class="turn-content">
        <div class="turn-icon">
          <div v-if="isPlayerTurn" class="player-icon">👤</div>
          <div v-else class="cpu-icon">🤖</div>
        </div>
        <div class="turn-text">
          <div class="turn-label">{{ isPlayerTurn ? 'あなたのターン' : '相手のターン' }}</div>
          <div class="turn-count">ターン {{ turnCount }}</div>
        </div>
      </div>
      
      <!-- CPU思考中インジケーター -->
      <div v-if="!isPlayerTurn && battleMode === 'cpu'" class="cpu-thinking">
        <div class="thinking-dots">
          <div class="dot"></div>
          <div class="dot"></div>
          <div class="dot"></div>
        </div>
        <div class="thinking-text">CPU思考中...</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'

const gameStore = useGameStore()

const isPlayerTurn = computed(() => gameStore.isPlayerTurn)
const turnCount = computed(() => gameStore.turnCount + 1)
const battleMode = computed(() => gameStore.battleMode)
</script>

<style scoped>
.turn-indicator {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 30;
}

.turn-display {
  background: rgba(17, 24, 39, 0.95);
  border: 2px solid transparent;
  border-radius: 12px;
  padding: 16px 20px;
  min-width: 200px;
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
  transition: all 0.4s ease;
}

.turn-display.player-turn {
  border-color: #10b981;
  box-shadow: 0 0 20px rgba(16, 185, 129, 0.3);
}

.turn-display.opponent-turn {
  border-color: #ef4444;
  box-shadow: 0 0 20px rgba(239, 68, 68, 0.3);
}

.turn-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.turn-icon {
  font-size: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.player-turn .player-icon {
  background: rgba(16, 185, 129, 0.2);
  border: 2px solid #10b981;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: player-pulse 2s ease-in-out infinite;
}

.opponent-turn .cpu-icon {
  background: rgba(239, 68, 68, 0.2);
  border: 2px solid #ef4444;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: cpu-pulse 1.5s ease-in-out infinite;
}

.turn-text {
  flex: 1;
}

.turn-label {
  font-size: 1.125rem;
  font-weight: bold;
  margin-bottom: 2px;
}

.player-turn .turn-label {
  color: #10b981;
}

.opponent-turn .turn-label {
  color: #ef4444;
}

.turn-count {
  font-size: 0.875rem;
  color: #9ca3af;
}

.cpu-thinking {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(75, 85, 99, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
}

.thinking-dots {
  display: flex;
  gap: 4px;
}

.dot {
  width: 6px;
  height: 6px;
  background: #6b7280;
  border-radius: 50%;
  animation: thinking-bounce 1.4s ease-in-out infinite both;
}

.dot:nth-child(2) {
  animation-delay: 0.16s;
}

.dot:nth-child(3) {
  animation-delay: 0.32s;
}

.thinking-text {
  font-size: 0.75rem;
  color: #6b7280;
}

@keyframes player-pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 10px rgba(16, 185, 129, 0.3);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
  }
}

@keyframes cpu-pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 10px rgba(239, 68, 68, 0.3);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 0 20px rgba(239, 68, 68, 0.6);
  }
}

@keyframes thinking-bounce {
  0%, 80%, 100% {
    transform: scale(0.8);
    opacity: 0.5;
  }
  40% {
    transform: scale(1.2);
    opacity: 1;
  }
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .turn-indicator {
    top: 10px;
    right: 10px;
  }
  
  .turn-display {
    min-width: 160px;
    padding: 12px 16px;
  }
  
  .turn-icon {
    font-size: 1.5rem;
    width: 40px;
    height: 40px;
  }
  
  .player-turn .player-icon,
  .opponent-turn .cpu-icon {
    width: 40px;
    height: 40px;
  }
}
</style>