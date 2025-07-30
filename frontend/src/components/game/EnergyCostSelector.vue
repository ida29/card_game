<template>
  <teleport to="body">
    <transition name="fade">
      <div
        v-if="show"
        class="fixed inset-0 bg-black/80 flex items-center justify-center z-50"
        @click="handleCancel"
      >
        <div class="cost-selection-container" @click.stop>
          <div class="bg-gray-900 rounded-lg p-6 max-w-6xl w-[90vw] max-h-[90vh] overflow-y-auto">
            <h2 class="text-2xl font-bold text-white mb-4 text-center">
              コストの支払い
            </h2>
            
            <!-- Card Preview and Cost Requirements Side by Side -->
            <div class="flex gap-8 mb-6">
              <!-- Card Preview -->
              <div v-if="cardToPay" class="card-preview flex-shrink-0">
                <div class="text-center">
                  <GameCard
                    :card="cardToPay"
                    size="medium"
                  />
                  <div class="mt-2">
                    <p class="text-white font-bold">{{ cardToPay.name }}</p>
                    <div class="cost-display flex justify-center items-center gap-2 mt-1">
                      <span class="text-gray-400">コスト:</span>
                      <div class="cost-icons flex gap-1">
                        <div v-if="costs.red > 0" class="cost-icon red">{{ costs.red }}</div>
                        <div v-if="costs.blue > 0" class="cost-icon blue">{{ costs.blue }}</div>
                        <div v-if="costs.yellow > 0" class="cost-icon yellow">{{ costs.yellow }}</div>
                        <div v-if="costs.green > 0" class="cost-icon green">{{ costs.green }}</div>
                        <div v-if="costs.colorless > 0" class="cost-icon colorless">{{ costs.colorless }}</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Cost Requirements -->
              <div class="cost-requirements flex-grow">
                <h3 class="text-lg font-bold text-white mb-3">必要なコスト</h3>
                <div class="requirements-grid grid grid-cols-2 gap-3">
                  <div v-if="costs.red > 0" class="requirement-item">
                    <div class="cost-icon red">{{ costs.red }}</div>
                    <span class="text-red-400">赤 × {{ costs.red }}</span>
                    <span class="status" :class="{ 'satisfied': selectedEnergy.red >= costs.red }">
                      ({{ selectedEnergy.red }}/{{ costs.red }})
                    </span>
                  </div>
                  <div v-if="costs.blue > 0" class="requirement-item">
                    <div class="cost-icon blue">{{ costs.blue }}</div>
                    <span class="text-blue-400">青 × {{ costs.blue }}</span>
                    <span class="status" :class="{ 'satisfied': selectedEnergy.blue >= costs.blue }">
                      ({{ selectedEnergy.blue }}/{{ costs.blue }})
                    </span>
                  </div>
                  <div v-if="costs.yellow > 0" class="requirement-item">
                    <div class="cost-icon yellow">{{ costs.yellow }}</div>
                    <span class="text-yellow-400">黄 × {{ costs.yellow }}</span>
                    <span class="status" :class="{ 'satisfied': selectedEnergy.yellow >= costs.yellow }">
                      ({{ selectedEnergy.yellow }}/{{ costs.yellow }})
                    </span>
                  </div>
                  <div v-if="costs.green > 0" class="requirement-item">
                    <div class="cost-icon green">{{ costs.green }}</div>
                    <span class="text-green-400">緑 × {{ costs.green }}</span>
                    <span class="status" :class="{ 'satisfied': selectedEnergy.green >= costs.green }">
                      ({{ selectedEnergy.green }}/{{ costs.green }})
                    </span>
                  </div>
                  <div v-if="costs.colorless > 0" class="requirement-item">
                    <div class="cost-icon colorless">{{ costs.colorless }}</div>
                    <span class="text-gray-400">無色 × {{ costs.colorless }}</span>
                    <span class="status" :class="{ 'satisfied': totalSelectedValue >= totalCost }">
                      ({{ Math.max(0, totalSelectedValue - coloredCost) }}/{{ costs.colorless }})
                    </span>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Available Energy Sources -->
            <div class="energy-sources">
              <h3 class="text-lg font-bold text-white mb-3">使用するエネルギーを選択</h3>
              
              <!-- Regular Energy -->
              <div v-if="availableRegularEnergy.length > 0" class="energy-section mb-4">
                <h4 class="text-md font-semibold text-gray-300 mb-2">エネルギーカード</h4>
                <div class="energy-grid">
                  <div
                    v-for="(energy, index) in availableRegularEnergy"
                    :key="`regular-${index}`"
                    class="energy-option"
                    :class="{ 
                      'selected': energy.selected
                    }"
                    @click="toggleRegularEnergy(index)"
                  >
                    <GameCard
                      :card="energy.card.card"
                      size="small"
                    />
                    <div class="energy-info">
                      <div class="energy-color" :class="getColorClass(energy.card.card.color)">
                        {{ energy.card.card.color }}
                      </div>
                      <div class="energy-value">{{ energy.card.card.energy_value || 1 }}</div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Negative Energy -->
              <div v-if="availableNegativeEnergy.length > 0" class="energy-section mb-4">
                <h4 class="text-md font-semibold text-purple-300 mb-2">負のエネルギーカード（表向き）</h4>
                <div class="energy-grid">
                  <div
                    v-for="(negEnergy, index) in availableNegativeEnergy"
                    :key="`negative-${index}`"
                    class="energy-option negative"
                    :class="{ 'selected': negEnergy.selected }"
                    @click="toggleNegativeEnergy(index)"
                  >
                    <GameCard
                      :card="negEnergy.card.card"
                      size="small"
                    />
                    <div class="energy-info">
                      <div class="energy-color" :class="getColorClass(negEnergy.card.card.color)">
                        {{ negEnergy.card.card.color }}
                      </div>
                      <div class="energy-value">{{ negEnergy.card.card.energy_value || 1 }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex justify-center gap-4 mt-6">
              <button
                @click="handleCancel"
                class="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-bold transition-colors"
              >
                キャンセル
              </button>
              <button
                @click="handleConfirm"
                :disabled="!canConfirm"
                class="px-6 py-3 rounded-lg font-bold transition-colors"
                :class="canConfirm 
                  ? 'bg-green-600 hover:bg-green-500 text-white' 
                  : 'bg-gray-600 text-gray-400 cursor-not-allowed'"
              >
                コストを支払う
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useGameStore } from '@/stores/game'
import type { Card, EnergyCardState, NegativeEnergyCardState } from '@/types'
import GameCard from './GameCard.vue'

const gameStore = useGameStore()

const props = defineProps<{
  show: boolean
  cardToPay: Card | null
  resolve: ((energySelection: any) => void) | null
}>()

const emit = defineEmits<{
  cancel: []
}>()

// Available energy sources
const availableRegularEnergy = ref<(EnergyCardState & { selected: boolean, originalIndex: number })[]>([])
const availableNegativeEnergy = ref<(NegativeEnergyCardState & { selected: boolean, originalIndex: number })[]>([])

// Cost calculation
const costs = computed(() => {
  if (!props.cardToPay) return { red: 0, blue: 0, yellow: 0, green: 0, colorless: 0 }
  
  const card = props.cardToPay
  return {
    red: card.cost_red || 0,
    blue: card.cost_blue || 0,
    yellow: card.cost_yellow || 0,
    green: card.cost_green || 0,
    colorless: Math.max(0, (card.cost || 0) - (card.cost_red || 0) - (card.cost_blue || 0) - (card.cost_yellow || 0) - (card.cost_green || 0))
  }
})

const totalCost = computed(() => {
  return costs.value.red + costs.value.blue + costs.value.yellow + costs.value.green + costs.value.colorless
})

const coloredCost = computed(() => {
  return costs.value.red + costs.value.blue + costs.value.yellow + costs.value.green
})

// Selected energy tracking
const selectedEnergy = computed(() => {
  const result = { red: 0, blue: 0, yellow: 0, green: 0 }
  
  // Count selected regular energy
  availableRegularEnergy.value.forEach(energy => {
    if (energy.selected) {
      const color = energy.card.card.color
      const value = energy.card.card.energy_value || 1
      if (color === '赤') result.red += value
      else if (color === '青') result.blue += value
      else if (color === '黄') result.yellow += value
      else if (color === '緑') result.green += value
    }
  })
  
  // Count selected negative energy
  availableNegativeEnergy.value.forEach(negEnergy => {
    if (negEnergy.selected) {
      const color = negEnergy.card.card.color
      const value = negEnergy.card.card.energy_value || 1
      if (color === '赤') result.red += value
      else if (color === '青') result.blue += value
      else if (color === '黄') result.yellow += value
      else if (color === '緑') result.green += value
    }
  })
  
  return result
})

const totalSelectedValue = computed(() => {
  let total = 0
  availableRegularEnergy.value.forEach(energy => {
    if (energy.selected) total += energy.card.card.energy_value || 1
  })
  availableNegativeEnergy.value.forEach(negEnergy => {
    if (negEnergy.selected) total += negEnergy.card.card.energy_value || 1
  })
  return total
})

const canConfirm = computed(() => {
  // Check if all colored costs are satisfied
  const colorsSatisfied = 
    selectedEnergy.value.red >= costs.value.red &&
    selectedEnergy.value.blue >= costs.value.blue &&
    selectedEnergy.value.yellow >= costs.value.yellow &&
    selectedEnergy.value.green >= costs.value.green
  
  // Check if total cost is exactly satisfied (not more, not less)
  const totalSatisfied = totalSelectedValue.value === totalCost.value
  
  return colorsSatisfied && totalSatisfied
})

// Initialize available energy when component shows
watch(() => props.show, (show) => {
  if (show && gameStore.player) {
    // Reset selections
    availableRegularEnergy.value = gameStore.player.energy
      .filter(e => !e.tapped)
      .map((energy, index) => ({ ...energy, selected: false, originalIndex: index }))
    
    availableNegativeEnergy.value = gameStore.player.negativeEnergy
      .filter(ne => ne.faceUp)
      .map((negEnergy, index) => ({ ...negEnergy, selected: false, originalIndex: index }))
  }
})

const toggleRegularEnergy = (index: number) => {
  const energy = availableRegularEnergy.value[index]
  const energyValue = energy.card.card.energy_value || 1
  const color = energy.card.card.color
  
  if (!energy.selected) {
    // Check if selecting this would exceed the total cost
    if (totalSelectedValue.value + energyValue > totalCost.value) {
      console.log('Cannot select: would exceed total cost')
      return
    }
    
    // Check if selecting this would exceed color requirements
    if (color === '赤' && selectedEnergy.value.red + energyValue > costs.value.red + costs.value.colorless) return
    if (color === '青' && selectedEnergy.value.blue + energyValue > costs.value.blue + costs.value.colorless) return
    if (color === '黄' && selectedEnergy.value.yellow + energyValue > costs.value.yellow + costs.value.colorless) return
    if (color === '緑' && selectedEnergy.value.green + energyValue > costs.value.green + costs.value.colorless) return
  }
  
  energy.selected = !energy.selected
}

const toggleNegativeEnergy = (index: number) => {
  const negEnergy = availableNegativeEnergy.value[index]
  const energyValue = negEnergy.card.card.energy_value || 1
  const color = negEnergy.card.card.color
  
  if (!negEnergy.selected) {
    // Check if selecting this would exceed the total cost
    if (totalSelectedValue.value + energyValue > totalCost.value) {
      console.log('Cannot select: would exceed total cost')
      return
    }
    
    // Check if selecting this would exceed color requirements
    if (color === '赤' && selectedEnergy.value.red + energyValue > costs.value.red + costs.value.colorless) return
    if (color === '青' && selectedEnergy.value.blue + energyValue > costs.value.blue + costs.value.colorless) return
    if (color === '黄' && selectedEnergy.value.yellow + energyValue > costs.value.yellow + costs.value.colorless) return
    if (color === '緑' && selectedEnergy.value.green + energyValue > costs.value.green + costs.value.colorless) return
  }
  
  negEnergy.selected = !negEnergy.selected
}

const getColorClass = (color: string) => {
  switch (color) {
    case '赤': return 'text-red-400'
    case '青': return 'text-blue-400'
    case '黄': return 'text-yellow-400'
    case '緑': return 'text-green-400'
    default: return 'text-gray-400'
  }
}

const handleCancel = () => {
  emit('cancel')
}

const handleConfirm = () => {
  if (!canConfirm.value) return
  
  const selection = {
    regularEnergy: availableRegularEnergy.value
      .filter(e => e.selected)
      .map(e => e.originalIndex),
    negativeEnergy: availableNegativeEnergy.value
      .filter(ne => ne.selected)
      .map(ne => ne.originalIndex)
  }
  
  gameStore.confirmEnergyCostSelection(selection)
}
</script>

<style scoped>
.cost-selection-container {
  max-height: 90vh;
  overflow-y: auto;
}

.cost-icons {
  display: flex;
  align-items: center;
  gap: 4px;
}

.cost-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: bold;
  color: white;
}

.cost-icon.red { background: #ef4444; }
.cost-icon.blue { background: #3b82f6; }
.cost-icon.yellow { background: #eab308; }
.cost-icon.green { background: #10b981; }
.cost-icon.colorless { background: #6b7280; }

.requirements-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.requirement-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: rgba(31, 41, 55, 0.5);
  border-radius: 8px;
}

.status {
  color: #ef4444;
  font-size: 0.875rem;
}

.status.satisfied {
  color: #10b981;
}

.energy-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.energy-grid-larger {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
}

.energy-option {
  position: relative;
  cursor: pointer;
  border: 2px solid transparent;
  border-radius: 8px;
  padding: 8px;
  background: rgba(55, 65, 81, 0.3);
  transition: all 0.2s ease;
}

.energy-option:hover {
  background: rgba(75, 85, 99, 0.5);
}

.energy-option.selected {
  border-color: #10b981;
  background: rgba(16, 185, 129, 0.2);
}

.energy-option.negative {
  border-color: #8b5cf6;
}

.energy-option.negative.selected {
  border-color: #a855f7;
  background: rgba(168, 85, 247, 0.2);
}

.energy-info {
  text-align: center;
  margin-top: 4px;
}

.energy-color {
  font-size: 0.75rem;
  font-weight: bold;
}

.energy-value {
  font-size: 0.625rem;
  color: #9ca3af;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>