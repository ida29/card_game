<template>
  <div class="game-controls fixed bottom-4 right-4 flex gap-2">
    <button
      v-if="showNextPhase"
      @click="$emit('next-phase')"
      class="btn-primary text-sm"
    >
      次のフェーズへ
    </button>
    
    <button
      v-if="canEndTurn"
      @click="$emit('end-turn')"
      class="btn-secondary text-sm"
    >
      ターン終了
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPhase: string
  canEndTurn: boolean
}>()

const emit = defineEmits<{
  'next-phase': []
  'end-turn': []
}>()

const showNextPhase = computed(() => {
  return props.currentPhase === 'draw' || props.currentPhase === 'main' || props.currentPhase === 'battle'
})
</script>

<style scoped>
.game-controls {
  z-index: 100;
}

button {
  padding: 0.5rem 1rem;
  font-weight: bold;
  transition: all 0.2s;
}
</style>