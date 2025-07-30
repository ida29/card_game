<template>
  <div class="effect-display" v-if="hasActiveEffects">
    <!-- Effect Stack Display -->
    <div class="effect-stack" v-if="effectStack.length > 0">
      <h3>効果スタック</h3>
      <div class="stack-items">
        <div 
          v-for="(effect, index) in effectStack" 
          :key="index"
          class="stack-item"
          :class="{ 'resolving': index === 0 }"
        >
          <div class="effect-source">
            <img :src="getCardImage(effect.source)" :alt="effect.source" />
          </div>
          <div class="effect-info">
            <p class="effect-name">{{ effect.source }}</p>
            <p class="effect-description">{{ effect.description }}</p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Persistent Effects Display -->
    <div class="persistent-effects" v-if="persistentEffects.length > 0">
      <h3>常在効果</h3>
      <div class="effect-list">
        <div 
          v-for="effect in persistentEffects" 
          :key="effect.cardNo"
          class="persistent-effect"
          :class="{ 'active': effect.isActive }"
        >
          <img :src="getCardImage(effect.cardNo)" :alt="effect.cardNo" />
          <div class="effect-tooltip">
            <p class="card-name">{{ effect.cardName }}</p>
            <p class="effect-text">{{ effect.description }}</p>
            <p class="effect-status" v-if="!effect.isActive">
              条件未達成
            </p>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Effect Animation Overlay -->
    <transition name="effect-animation">
      <div v-if="currentAnimation" class="effect-animation-overlay">
        <div class="animation-content">
          <div 
            v-if="currentAnimation.type === 'destroy'" 
            class="effect-icon destroy-effect"
          >
            💥
          </div>
          <div 
            v-if="currentAnimation.type === 'power-up'" 
            class="effect-icon powerup-effect"
          >
            ⚡
          </div>
          <div 
            v-if="currentAnimation.type === 'draw'" 
            class="effect-icon draw-effect"
          >
            🎴
          </div>
          <p class="animation-text">{{ currentAnimation.text }}</p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useGameStore } from '@/stores/game'
import { useCardStore } from '@/stores/cards'

interface StackEffect {
  source: string
  description: string
  player: number
}

interface PersistentEffect {
  cardNo: string
  cardName: string
  description: string
  isActive: boolean
}

interface EffectAnimation {
  type: string
  text: string
}

const gameStore = useGameStore()
const cardStore = useCardStore()

const effectStack = ref<StackEffect[]>([])
const persistentEffects = ref<PersistentEffect[]>([])
const currentAnimation = ref<EffectAnimation | null>(null)

const hasActiveEffects = computed(() => 
  effectStack.value.length > 0 || persistentEffects.value.length > 0
)

// Watch for effect updates from the game store
watch(() => gameStore.effectUpdates, (updates) => {
  if (updates) {
    // Update effect stack
    if (updates.stack) {
      effectStack.value = updates.stack
    }
    
    // Update persistent effects
    if (updates.persistent) {
      updatePersistentEffects(updates.persistent)
    }
    
    // Show animation if needed
    if (updates.animation) {
      showAnimation(updates.animation)
    }
  }
}, { deep: true })

function updatePersistentEffects(effects: any[]) {
  persistentEffects.value = effects.map(effect => {
    const card = cardStore.getCardByNo(effect.cardNo)
    return {
      cardNo: effect.cardNo,
      cardName: card?.name || effect.cardNo,
      description: effect.description,
      isActive: effect.isActive
    }
  })
}

function showAnimation(animation: EffectAnimation) {
  currentAnimation.value = animation
  setTimeout(() => {
    currentAnimation.value = null
  }, 2000)
}

function getCardImage(cardNo: string): string {
  const card = cardStore.getCardByNo(cardNo)
  if (!card) return '/placeholder-card.jpg'
  
  if (card.local_image_path) {
    return `/api/v1/images/${card.local_image_path.replace('card_images/', '')}`
  }
  return card.image_url || '/placeholder-card.jpg'
}

// Public methods that can be called from parent
defineExpose({
  addToStack(effect: StackEffect) {
    effectStack.value.unshift(effect)
  },
  
  resolveTopEffect() {
    if (effectStack.value.length > 0) {
      effectStack.value.shift()
    }
  },
  
  clearStack() {
    effectStack.value = []
  },
  
  showEffect(type: string, text: string) {
    showAnimation({ type, text })
  }
})
</script>

<style scoped>
.effect-display {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  max-width: 300px;
}

.effect-stack,
.persistent-effects {
  background: rgba(0, 0, 0, 0.8);
  border: 2px solid #444;
  border-radius: 8px;
  padding: 10px;
  margin-bottom: 10px;
}

.effect-stack h3,
.persistent-effects h3 {
  color: #fff;
  margin: 0 0 10px;
  font-size: 14px;
}

.stack-item {
  display: flex;
  align-items: center;
  padding: 5px;
  margin-bottom: 5px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  transition: all 0.3s;
}

.stack-item.resolving {
  background: rgba(255, 255, 0, 0.2);
  border: 1px solid yellow;
}

.effect-source img {
  width: 40px;
  height: 56px;
  object-fit: cover;
  border-radius: 4px;
}

.effect-info {
  margin-left: 10px;
  flex: 1;
}

.effect-name {
  color: #fff;
  font-size: 12px;
  font-weight: bold;
  margin: 0;
}

.effect-description {
  color: #ccc;
  font-size: 11px;
  margin: 2px 0 0;
}

.persistent-effect {
  display: inline-block;
  position: relative;
  margin: 5px;
}

.persistent-effect img {
  width: 50px;
  height: 70px;
  object-fit: cover;
  border-radius: 4px;
  opacity: 0.5;
  transition: opacity 0.3s;
}

.persistent-effect.active img {
  opacity: 1;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
}

.effect-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid #666;
  border-radius: 4px;
  padding: 8px;
  min-width: 150px;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.3s;
}

.persistent-effect:hover .effect-tooltip {
  opacity: 1;
}

.card-name {
  color: #fff;
  font-size: 12px;
  font-weight: bold;
  margin: 0 0 4px;
}

.effect-text {
  color: #ccc;
  font-size: 11px;
  margin: 0;
}

.effect-status {
  color: #f66;
  font-size: 10px;
  margin: 4px 0 0;
}

/* Effect Animation Overlay */
.effect-animation-overlay {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 2000;
}

.animation-content {
  text-align: center;
}

.effect-icon {
  font-size: 100px;
  animation: effect-pulse 0.5s ease-out;
}

.destroy-effect {
  animation: destroy-shake 0.5s;
}

.powerup-effect {
  animation: powerup-glow 1s;
}

.draw-effect {
  animation: draw-spin 1s;
}

.animation-text {
  color: #fff;
  font-size: 24px;
  font-weight: bold;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);
  margin-top: 20px;
}

/* Animations */
@keyframes effect-pulse {
  0% {
    transform: scale(0.5);
    opacity: 0;
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

@keyframes destroy-shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10px); }
  75% { transform: translateX(10px); }
}

@keyframes powerup-glow {
  0%, 100% {
    filter: brightness(1);
  }
  50% {
    filter: brightness(2) drop-shadow(0 0 20px yellow);
  }
}

@keyframes draw-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Vue transitions */
.effect-animation-enter-active,
.effect-animation-leave-active {
  transition: opacity 0.3s;
}

.effect-animation-enter-from,
.effect-animation-leave-to {
  opacity: 0;
}

.effect-stack,
.persistent-effects {
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid #444;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 10px;
}

.effect-stack h3,
.persistent-effects h3 {
  color: #fff;
  font-size: 1.1rem;
  margin-bottom: 10px;
  text-align: center;
}

.stack-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.stack-item {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 8px;
  transition: all 0.3s;
}

.stack-item.resolving {
  background: rgba(255, 200, 0, 0.2);
  border: 1px solid #ffc800;
}

.effect-source img {
  width: 40px;
  height: 56px;
  border-radius: 4px;
  object-fit: cover;
}

.effect-info {
  margin-left: 10px;
  flex: 1;
}

.effect-name {
  color: #fff;
  font-size: 0.9rem;
  font-weight: bold;
  margin: 0;
}

.effect-description {
  color: #ccc;
  font-size: 0.8rem;
  margin: 4px 0 0 0;
}

.effect-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.persistent-effect {
  position: relative;
  cursor: pointer;
}

.persistent-effect img {
  width: 50px;
  height: 70px;
  border-radius: 4px;
  object-fit: cover;
  transition: all 0.2s;
}

.persistent-effect.active img {
  box-shadow: 0 0 10px rgba(0, 255, 0, 0.5);
}

.persistent-effect:not(.active) img {
  opacity: 0.5;
}

.effect-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.95);
  border: 1px solid #444;
  border-radius: 4px;
  padding: 10px;
  width: 200px;
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.2s;
  z-index: 10;
}

.persistent-effect:hover .effect-tooltip {
  opacity: 1;
}

.card-name {
  color: #fff;
  font-weight: bold;
  margin: 0 0 5px 0;
}

.effect-text {
  color: #ccc;
  font-size: 0.8rem;
  margin: 0;
}

.effect-status {
  color: #ff6666;
  font-size: 0.8rem;
  margin: 5px 0 0 0;
}

.effect-animation-overlay {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 1000;
}

.animation-content {
  text-align: center;
}

.animation-content img {
  width: 200px;
  height: 200px;
}

.animation-text {
  color: #fff;
  font-size: 1.5rem;
  font-weight: bold;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.8);
  margin-top: 10px;
}

.effect-animation-enter-active,
.effect-animation-leave-active {
  transition: opacity 0.3s;
}

.effect-animation-enter-from,
.effect-animation-leave-to {
  opacity: 0;
}
</style>