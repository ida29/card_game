<template>
  <div v-if="isSelecting" class="target-selector-overlay">
    <div class="target-selector-modal">
      <h3 class="title">{{ requirements.description }}</h3>
      
      <div class="target-info">
        <p>{{ getSelectionInfo() }}</p>
      </div>
      
      <!-- Battle Area Targets -->
      <div v-if="hasTargetsInBattleArea" class="target-section">
        <h4>バトルエリア</h4>
        <div class="battle-targets">
          <div 
            v-for="target in battleAreaTargets" 
            :key="`${target.type}_${target.id}_${target.location}`"
            class="target-card"
            :class="{ 
              'selected': isSelected(target),
              'selectable': canSelect(target)
            }"
            @click="toggleTarget(target)"
          >
            <img :src="getCardImage(target.id)" :alt="target.id" />
            <div class="target-overlay" v-if="canSelect(target)">
              <span class="checkmark" v-if="isSelected(target)">✓</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Hand Targets -->
      <div v-if="hasTargetsInHand" class="target-section">
        <h4>手札</h4>
        <div class="hand-targets">
          <div 
            v-for="target in handTargets" 
            :key="`${target.type}_${target.id}_${target.location}`"
            class="target-card hand-card"
            :class="{ 
              'selected': isSelected(target),
              'selectable': canSelect(target)
            }"
            @click="toggleTarget(target)"
          >
            <img :src="getCardImage(target.id)" :alt="target.id" />
            <div class="target-overlay" v-if="canSelect(target)">
              <span class="checkmark" v-if="isSelected(target)">✓</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Energy Area Targets -->
      <div v-if="hasTargetsInEnergyArea" class="target-section">
        <h4>エネルギーエリア</h4>
        <div class="energy-targets">
          <div 
            v-for="target in energyAreaTargets" 
            :key="`${target.type}_${target.id}_${target.location}`"
            class="target-card energy-card"
            :class="{ 
              'selected': isSelected(target),
              'selectable': canSelect(target)
            }"
            @click="toggleTarget(target)"
          >
            <img :src="getCardImage(target.id)" :alt="target.id" />
            <div class="target-overlay" v-if="canSelect(target)">
              <span class="checkmark" v-if="isSelected(target)">✓</span>
            </div>
          </div>
        </div>
      </div>
      
      <div class="actions">
        <button 
          class="btn btn-primary" 
          :disabled="!canConfirm"
          @click="confirmSelection"
        >
          決定
        </button>
        <button 
          class="btn btn-secondary" 
          @click="cancelSelection"
          v-if="!requirements.mandatory || selectedTargets.length >= requirements.minTargets"
        >
          キャンセル
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGameStore } from '@/stores/game'
import { useCardStore } from '@/stores/cards'

interface Target {
  type: string
  id: string
  location: string
  data?: any
}

interface TargetRequirements {
  minTargets: number
  maxTargets: number
  targetTypes: string[]
  mandatory: boolean
  description: string
}

interface Props {
  targets: Target[]
  requirements: TargetRequirements
  onConfirm: (targets: Target[]) => void
  onCancel: () => void
}

const props = defineProps<Props>()
const gameStore = useGameStore()
const cardStore = useCardStore()

const isSelecting = ref(true)
const selectedTargets = ref<Target[]>([])

// Computed properties for different target locations
const battleAreaTargets = computed(() => 
  props.targets.filter(t => t.location.includes('battle_area'))
)

const handTargets = computed(() => 
  props.targets.filter(t => t.location === 'hand')
)

const energyAreaTargets = computed(() => 
  props.targets.filter(t => t.location.includes('energy'))
)

const hasTargetsInBattleArea = computed(() => battleAreaTargets.value.length > 0)
const hasTargetsInHand = computed(() => handTargets.value.length > 0)
const hasTargetsInEnergyArea = computed(() => energyAreaTargets.value.length > 0)

const canConfirm = computed(() => {
  const count = selectedTargets.value.length
  return count >= props.requirements.minTargets && 
         (props.requirements.maxTargets === -1 || count <= props.requirements.maxTargets)
})

function getSelectionInfo(): string {
  const min = props.requirements.minTargets
  const max = props.requirements.maxTargets
  const current = selectedTargets.value.length
  
  if (min === max && max === 1) {
    return '1体を選択してください'
  } else if (min === max) {
    return `${min}体を選択してください (${current}/${min})`
  } else if (max === -1) {
    return `${min}体以上選択してください (${current}体選択中)`
  } else {
    return `${min}〜${max}体を選択してください (${current}体選択中)`
  }
}

function isSelected(target: Target): boolean {
  return selectedTargets.value.some(t => 
    t.type === target.type && 
    t.id === target.id && 
    t.location === target.location
  )
}

function canSelect(target: Target): boolean {
  if (isSelected(target)) {
    return true // Can always deselect
  }
  
  // Check if we've reached max targets
  if (props.requirements.maxTargets !== -1 && 
      selectedTargets.value.length >= props.requirements.maxTargets) {
    return false
  }
  
  // Check target type
  if (props.requirements.targetTypes.length > 0 &&
      !props.requirements.targetTypes.includes(target.type)) {
    return false
  }
  
  return true
}

function toggleTarget(target: Target) {
  if (!canSelect(target) && !isSelected(target)) {
    return
  }
  
  if (isSelected(target)) {
    // Deselect
    selectedTargets.value = selectedTargets.value.filter(t => 
      !(t.type === target.type && 
        t.id === target.id && 
        t.location === target.location)
    )
  } else {
    // Select
    selectedTargets.value.push(target)
  }
}

function getCardImage(cardNo: string): string {
  const card = cardStore.getCardByNo(cardNo)
  if (!card) return '/placeholder-card.jpg'
  
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.jpg'
}

function confirmSelection() {
  if (canConfirm.value) {
    props.onConfirm(selectedTargets.value)
    isSelecting.value = false
  }
}

function cancelSelection() {
  props.onCancel()
  isSelecting.value = false
}
</script>

<style scoped>
.target-selector-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.target-selector-modal {
  background: #1a1a1a;
  border-radius: 8px;
  padding: 20px;
  max-width: 800px;
  max-height: 80vh;
  overflow-y: auto;
  color: white;
}

.title {
  font-size: 1.5rem;
  margin-bottom: 10px;
  text-align: center;
}

.target-info {
  text-align: center;
  margin-bottom: 20px;
  color: #888;
}

.target-section {
  margin-bottom: 20px;
}

.target-section h4 {
  margin-bottom: 10px;
  color: #ccc;
}

.battle-targets,
.hand-targets,
.energy-targets {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.target-card {
  position: relative;
  cursor: pointer;
  transition: transform 0.2s;
}

.target-card img {
  width: 100px;
  height: 140px;
  border-radius: 8px;
  object-fit: cover;
}

.hand-card img {
  width: 80px;
  height: 112px;
}

.energy-card img {
  width: 60px;
  height: 84px;
}

.target-card.selectable:hover {
  transform: scale(1.05);
}

.target-card:not(.selectable) {
  opacity: 0.5;
  cursor: not-allowed;
}

.target-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.target-card.selected .target-overlay {
  background: rgba(0, 255, 0, 0.3);
  border: 2px solid #00ff00;
}

.checkmark {
  font-size: 2rem;
  color: #00ff00;
  font-weight: bold;
}

.actions {
  display: flex;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #0056b3;
}

.btn-primary:disabled {
  background: #444;
  cursor: not-allowed;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #545b62;
}
</style>