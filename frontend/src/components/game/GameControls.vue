<template>
  <div class="game-controls fixed right-6 top-1/2 transform -translate-y-1/2 flex flex-col gap-3">
    <button
      v-if="showNextPhase"
      @click="$emit('next-phase')"
      class="control-button next-phase"
    >
      <div class="button-content">
        <div class="button-icon">→</div>
        <div class="button-text">
          <div class="button-label">次のフェーズへ</div>
          <div class="button-hint">{{ getPhaseHint() }}</div>
        </div>
      </div>
    </button>
    
    <button
      v-if="!showNextPhase && canEndTurn"
      @click="$emit('end-turn')"
      class="control-button end-turn"
    >
      <div class="button-content">
        <div class="button-icon">⏸</div>
        <div class="button-text">
          <div class="button-label">ターン終了</div>
          <div class="button-hint">相手のターンへ</div>
        </div>
      </div>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPhase: string
  canEndTurn: boolean
  energyPlayedThisTurn: { player: boolean, opponent: boolean }
}>()

const emit = defineEmits<{
  'next-phase': []
  'end-turn': []
}>()

const showNextPhase = computed(() => {
  // Show "Next Phase" for draw phase
  if (props.currentPhase === 'draw') {
    return true
  }
  
  // For energy phase, only show if energy has been played
  if (props.currentPhase === 'energy') {
    return props.energyPlayedThisTurn?.player === true
  }
  
  return false
})

const getPhaseHint = () => {
  switch (props.currentPhase) {
    case 'draw':
      return 'エネルギーフェーズへ'
    case 'energy':
      return 'メインフェーズへ'
    case 'main':
      return 'エンドフェーズへ'
    default:
      return ''
  }
}
</script>

<style scoped>
.game-controls {
  z-index: 30;
}

.control-button {
  position: relative;
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.9) 0%, rgba(139, 92, 246, 0.9) 100%);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 1rem;
  padding: 1rem 1.5rem;
  color: white;
  font-weight: bold;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  box-shadow: 
    0 4px 20px rgba(59, 130, 246, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  min-width: 200px;
  overflow: hidden;
}

.control-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s ease;
}

.control-button:hover::before {
  left: 100%;
}

.control-button:hover {
  transform: translateX(-5px) scale(1.05);
  box-shadow: 
    0 6px 30px rgba(59, 130, 246, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.control-button:active {
  transform: translateX(-3px) scale(1.02);
}

.control-button.end-turn {
  background: linear-gradient(135deg, rgba(239, 68, 68, 0.9) 0%, rgba(236, 72, 153, 0.9) 100%);
  box-shadow: 
    0 4px 20px rgba(239, 68, 68, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.control-button.end-turn:hover {
  box-shadow: 
    0 6px 30px rgba(239, 68, 68, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.button-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.button-icon {
  font-size: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
}

.button-text {
  flex: 1;
  text-align: left;
}

.button-label {
  font-size: 1rem;
  font-weight: bold;
  margin-bottom: 0.125rem;
}

.button-hint {
  font-size: 0.75rem;
  opacity: 0.8;
  font-weight: normal;
}
</style>