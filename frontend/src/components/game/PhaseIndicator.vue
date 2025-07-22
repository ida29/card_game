<template>
  <div class="phase-indicator">
    <div class="phase-container">
      <!-- Current Phase Display -->
      <div class="current-phase">
        <div class="phase-icon">{{ currentPhase.icon }}</div>
        <div class="phase-name">{{ currentPhase.name }}</div>
        <div class="turn-info">ターン {{ turnCount }} - {{ isPlayerTurn ? 'あなた' : '相手' }}</div>
      </div>
      
      <!-- Phase Progress -->
      <div class="phase-progress">
        <div
          v-for="(phase, index) in phases"
          :key="phase.id"
          class="phase-step"
          :class="{
            'active': currentPhaseIndex === index,
            'completed': currentPhaseIndex > index,
            'upcoming': currentPhaseIndex < index
          }"
        >
          <div class="step-dot">
            <span v-if="currentPhaseIndex > index">✓</span>
            <span v-else>{{ index + 1 }}</span>
          </div>
          <div class="step-name">{{ phase.shortName }}</div>
          <div v-if="index < phases.length - 1" class="step-line"></div>
        </div>
      </div>
      
      <!-- Phase Description -->
      <div class="phase-description">
        {{ currentPhase.description }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useGameStore } from '@/stores/game'

const gameStore = useGameStore()

const phases = [
  {
    id: 'start',
    name: 'スタートフェーズ',
    shortName: 'スタート',
    description: 'カードをアクティブにします',
    icon: '☀️'
  },
  {
    id: 'draw',
    name: 'ドローフェーズ',
    shortName: 'ドロー',
    description: 'カードを1枚引きます',
    icon: '🎴'
  },
  {
    id: 'energy',
    name: 'エネルギーフェーズ',
    shortName: 'エネルギー',
    description: 'エネルギーを1枚セットできます',
    icon: '⚡'
  },
  {
    id: 'main',
    name: 'メインフェーズ',
    shortName: 'メイン',
    description: 'カードをプレイし、攻撃できます',
    icon: '🃏'
  },
  {
    id: 'end',
    name: 'エンドフェーズ',
    shortName: 'エンド',
    description: 'ターンを終了します',
    icon: '🌙'
  }
]

const currentPhaseIndex = computed(() => {
  return phases.findIndex(p => p.id === gameStore.currentPhase)
})

const currentPhase = computed(() => {
  return phases[currentPhaseIndex.value] || phases[0]
})

const isPlayerTurn = computed(() => gameStore.isPlayerTurn)
const turnCount = computed(() => gameStore.turnCount + 1)
</script>

<style scoped>
.phase-indicator {
  position: fixed;
  left: 20px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 40;
}

.phase-container {
  background: rgba(17, 24, 39, 0.95);
  border: 2px solid rgba(55, 65, 81, 0.8);
  border-radius: 12px;
  padding: 20px;
  min-width: 200px;
  backdrop-filter: blur(10px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
}

.current-phase {
  text-align: center;
  margin-bottom: 20px;
}

.phase-icon {
  font-size: 2.5rem;
  margin-bottom: 8px;
  animation: pulse 2s infinite;
}

.phase-name {
  font-size: 1.25rem;
  font-weight: bold;
  color: #fbbf24;
  margin-bottom: 4px;
}

.turn-info {
  font-size: 0.75rem;
  color: #9ca3af;
}

.phase-progress {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.phase-step {
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
}

.step-dot {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: bold;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.phase-step.completed .step-dot {
  background: #10b981;
  color: white;
}

.phase-step.active .step-dot {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  color: white;
  transform: scale(1.1);
  box-shadow: 0 0 15px rgba(59, 130, 246, 0.5);
}

.phase-step.upcoming .step-dot {
  background: #374151;
  color: #9ca3af;
}

.step-name {
  font-size: 0.875rem;
  font-weight: 500;
  flex: 1;
  transition: color 0.3s ease;
}

.phase-step.active .step-name {
  color: #fbbf24;
  font-weight: bold;
}

.phase-step.completed .step-name {
  color: #10b981;
}

.phase-step.upcoming .step-name {
  color: #6b7280;
}

.step-line {
  position: absolute;
  left: 16px;
  top: 32px;
  width: 2px;
  height: 12px;
  background: #374151;
  transition: background 0.3s ease;
}

.phase-step.completed .step-line {
  background: #10b981;
}

.phase-description {
  font-size: 0.875rem;
  color: #d1d5db;
  text-align: center;
  padding: 12px;
  background: rgba(31, 41, 55, 0.5);
  border-radius: 8px;
  border: 1px solid rgba(75, 85, 99, 0.3);
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    opacity: 1;
  }
  50% {
    transform: scale(1.05);
    opacity: 0.8;
  }
}

/* レスポンシブ対応 */
@media (max-width: 768px) {
  .phase-indicator {
    left: 10px;
    transform: translateY(-50%) scale(0.9);
  }
  
  .phase-container {
    min-width: 180px;
    padding: 16px;
  }
}
</style>