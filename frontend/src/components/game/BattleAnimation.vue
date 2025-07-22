<template>
  <teleport to="body">
    <transition name="battle-fade">
      <div 
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80"
        @click="handleSkip"
      >
        <div class="battle-animation-container">
          <!-- Battle Lightning Effect -->
          <div class="battle-effects">
            <div class="lightning-effect"></div>
            <div class="impact-zone"></div>
          </div>
          
          <!-- Attacker Card -->
          <div 
            class="battle-card attacker"
            :class="{ 
              'shake': battleStarted, 
              'defeated': attackerDefeated,
              'mutual-destruction': isMutualDestruction && battleStarted
            }"
          >
            <GameCard 
              :card="attackerCard" 
              size="large"
            />
            <div class="power-display">
              <span class="power-value">{{ attackerCard?.power || 0 }}</span>
            </div>
            
            <!-- Attacker Destruction Effect -->
            <transition name="destruction">
              <div v-if="attackerDefeated && attackerCard" class="destruction-effect">
                <!-- Glass shatter particles -->
                <div class="glass-particle" v-for="i in 24" :key="`attacker-glass-${i}`"></div>
                <!-- Crack lines -->
                <div class="crack-line" v-for="i in 8" :key="`attacker-crack-${i}`"></div>
                <!-- Impact shockwave -->
                <div class="shockwave"></div>
                <!-- Flash effect -->
                <div class="flash-effect"></div>
              </div>
            </transition>
          </div>
          
          <!-- VS Text -->
          <div class="vs-text" :class="{ 
            'active': battleStarted,
            'mutual-destruction': isMutualDestruction && battleStarted
          }">
            <span v-if="!isMutualDestruction">VS</span>
            <span v-else class="mutual-text">相打ち!</span>
          </div>
          
          <!-- Defender Card or Player Silhouette -->
          <div 
            class="battle-card defender"
            :class="{ 
              'shake': battleStarted, 
              'defeated': defenderDefeated,
              'mutual-destruction': isMutualDestruction && battleStarted
            }"
          >
            <!-- Show card if defending, otherwise show player silhouette -->
            <template v-if="defenderCard">
              <GameCard 
                :card="defenderCard" 
                size="large"
              />
              <div class="power-display">
                <span class="power-value">{{ defenderCard?.power || 0 }}</span>
              </div>
            </template>
            
            <!-- Player Silhouette for direct attack -->
            <template v-else>
              <div class="player-silhouette">
                <svg width="112" height="160" viewBox="0 0 112 160" class="silhouette-svg">
                  <path d="M56 20 C40 20 28 32 28 48 C28 58 32 66 38 70 L30 80 C20 85 10 95 10 110 L10 140 C10 150 18 160 30 160 L82 160 C94 160 102 150 102 140 L102 110 C102 95 92 85 82 80 L74 70 C80 66 84 58 84 48 C84 32 72 20 56 20 Z" 
                    fill="url(#playerGradient)" 
                    stroke="#ef4444" 
                    stroke-width="2"
                  />
                  <defs>
                    <linearGradient id="playerGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                      <stop offset="0%" style="stop-color:#1f2937;stop-opacity:1" />
                      <stop offset="100%" style="stop-color:#111827;stop-opacity:1" />
                    </linearGradient>
                  </defs>
                </svg>
                <div class="direct-attack-text">
                  DIRECT ATTACK!
                </div>
              </div>
            </template>
            
            <!-- Destruction Effect -->
            <transition name="destruction">
              <div v-if="defenderDefeated && defenderCard" class="destruction-effect">
                <!-- Glass shatter particles -->
                <div class="glass-particle" v-for="i in 24" :key="`glass-${i}`"></div>
                <!-- Crack lines -->
                <div class="crack-line" v-for="i in 8" :key="`crack-${i}`"></div>
                <!-- Impact shockwave -->
                <div class="shockwave"></div>
                <!-- Flash effect -->
                <div class="flash-effect"></div>
              </div>
            </transition>
            
            <!-- Damage Effect for direct attack -->
            <transition name="damage">
              <div v-if="defenderDefeated && !defenderCard" class="damage-effect">
                <div class="damage-number">-1</div>
              </div>
            </transition>
          </div>
        </div>
        
        <!-- Skip Hint -->
        <div class="skip-hint">
          クリックでスキップ
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed, onUnmounted } from 'vue'
import type { Card } from '@/types'
import GameCard from './GameCard.vue'

const props = defineProps<{
  show: boolean
  attackerCard: Card | undefined
  defenderCard: Card | undefined
  attackerDefeated: boolean
  defenderDefeated: boolean
}>()

const emit = defineEmits<{
  'animation-complete': []
}>()

const battleStarted = ref(false)

// Check if it's mutual destruction (both defeated)
const isMutualDestruction = computed(() => 
  props.attackerDefeated && props.defenderDefeated && props.defenderCard
)

let animationTimeout: number

watch(() => props.show, (newVal) => {
  if (newVal) {
    battleStarted.value = false
    
    // Start battle animation after a short delay
    setTimeout(() => {
      battleStarted.value = true
      
      // Complete animation after showing the result
      animationTimeout = setTimeout(() => {
        emit('animation-complete')
      }, 2500) as unknown as number
    }, 500)
  } else {
    battleStarted.value = false
    clearTimeout(animationTimeout)
  }
})

const handleSkip = () => {
  clearTimeout(animationTimeout)
  emit('animation-complete')
}

onUnmounted(() => {
  clearTimeout(animationTimeout)
})
</script>

<style scoped>
.battle-animation-container {
  position: relative;
  display: flex;
  align-items: center;
  gap: 4rem;
}

.battle-card {
  position: relative;
  transform: scale(1.5);
  transition: all 0.3s ease;
}

.battle-card.attacker {
  animation: slide-in-left 0.5s ease-out;
}

.battle-card.defender {
  animation: slide-in-right 0.5s ease-out;
}

.battle-card.shake {
  animation: shake 0.3s ease-in-out;
}

.battle-card.mutual-destruction {
  animation: mutual-destruction-shake 0.8s ease-in-out;
}

.battle-card.defeated {
  opacity: 0.7;
  transform: scale(1.3) translateY(20px);
}

.power-display {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(255, 0, 0, 0.9);
  color: white;
  font-size: 2rem;
  font-weight: bold;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  box-shadow: 0 4px 12px rgba(255, 0, 0, 0.5);
}

.vs-text {
  font-size: 4rem;
  font-weight: bold;
  color: #fff;
  text-shadow: 0 0 20px rgba(255, 255, 255, 0.8);
  opacity: 0;
  transform: scale(0);
  transition: all 0.5s ease;
}

.vs-text.active {
  opacity: 1;
  transform: scale(1);
  animation: pulse 1s ease-in-out infinite;
}

.vs-text.mutual-destruction {
  color: #ff6b35;
  text-shadow: 0 0 30px rgba(255, 107, 53, 0.9);
  animation: mutual-pulse 0.5s ease-in-out infinite;
}

.mutual-text {
  font-size: 0.8em;
  background: linear-gradient(45deg, #ff6b35, #f7931e, #ff6b35);
  background-size: 200% 200%;
  animation: gradient-shift 1s ease-in-out infinite;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.battle-effects {
  position: absolute;
  inset: -200px;
  pointer-events: none;
}

.lightning-effect {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 800px;
  height: 4px;
  background: linear-gradient(90deg, transparent, #fff, transparent);
  opacity: 0;
  animation: lightning 0.3s ease-out 0.8s;
}

.impact-zone {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.8), transparent);
  opacity: 0;
  animation: impact 0.5s ease-out 1s;
}

.destruction-effect {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.glass-particle {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 15px;
  height: 20px;
  background: linear-gradient(135deg, rgba(255,255,255,0.9), rgba(200,200,255,0.7));
  transform-origin: center;
  animation: glass-shatter 1.2s ease-out forwards;
  clip-path: polygon(20% 0%, 80% 10%, 100% 50%, 75% 100%, 25% 90%, 0% 40%);
}

/* Create random angles and distances for each particle */
.glass-particle:nth-child(1) { --angle: 15deg; --distance: 150px; --rotation: 720deg; animation-delay: 0s; }
.glass-particle:nth-child(2) { --angle: 30deg; --distance: 180px; --rotation: -540deg; animation-delay: 0.02s; }
.glass-particle:nth-child(3) { --angle: 45deg; --distance: 120px; --rotation: 900deg; animation-delay: 0.01s; }
.glass-particle:nth-child(4) { --angle: 60deg; --distance: 200px; --rotation: -720deg; animation-delay: 0.03s; }
.glass-particle:nth-child(5) { --angle: 75deg; --distance: 140px; --rotation: 630deg; animation-delay: 0s; }
.glass-particle:nth-child(6) { --angle: 90deg; --distance: 170px; --rotation: -810deg; animation-delay: 0.02s; }
.glass-particle:nth-child(7) { --angle: 105deg; --distance: 190px; --rotation: 540deg; animation-delay: 0.01s; }
.glass-particle:nth-child(8) { --angle: 120deg; --distance: 130px; --rotation: -900deg; animation-delay: 0.03s; }
.glass-particle:nth-child(9) { --angle: 135deg; --distance: 160px; --rotation: 720deg; animation-delay: 0s; }
.glass-particle:nth-child(10) { --angle: 150deg; --distance: 180px; --rotation: -630deg; animation-delay: 0.02s; }
.glass-particle:nth-child(11) { --angle: 165deg; --distance: 140px; --rotation: 810deg; animation-delay: 0.01s; }
.glass-particle:nth-child(12) { --angle: 180deg; --distance: 200px; --rotation: -540deg; animation-delay: 0.03s; }
.glass-particle:nth-child(13) { --angle: 195deg; --distance: 150px; --rotation: 900deg; animation-delay: 0s; }
.glass-particle:nth-child(14) { --angle: 210deg; --distance: 170px; --rotation: -720deg; animation-delay: 0.02s; }
.glass-particle:nth-child(15) { --angle: 225deg; --distance: 130px; --rotation: 630deg; animation-delay: 0.01s; }
.glass-particle:nth-child(16) { --angle: 240deg; --distance: 190px; --rotation: -810deg; animation-delay: 0.03s; }
.glass-particle:nth-child(17) { --angle: 255deg; --distance: 160px; --rotation: 540deg; animation-delay: 0s; }
.glass-particle:nth-child(18) { --angle: 270deg; --distance: 140px; --rotation: -900deg; animation-delay: 0.02s; }
.glass-particle:nth-child(19) { --angle: 285deg; --distance: 180px; --rotation: 720deg; animation-delay: 0.01s; }
.glass-particle:nth-child(20) { --angle: 300deg; --distance: 150px; --rotation: -630deg; animation-delay: 0.03s; }
.glass-particle:nth-child(21) { --angle: 315deg; --distance: 170px; --rotation: 810deg; animation-delay: 0s; }
.glass-particle:nth-child(22) { --angle: 330deg; --distance: 130px; --rotation: -540deg; animation-delay: 0.02s; }
.glass-particle:nth-child(23) { --angle: 345deg; --distance: 200px; --rotation: 900deg; animation-delay: 0.01s; }
.glass-particle:nth-child(24) { --angle: 360deg; --distance: 160px; --rotation: -720deg; animation-delay: 0.03s; }

.crack-line {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 2px;
  height: 100px;
  background: linear-gradient(to bottom, transparent, rgba(255,255,255,0.8), transparent);
  transform-origin: center;
  animation: crack-appear 0.3s ease-out forwards;
}

.crack-line:nth-child(1) { --crack-angle: 0deg; }
.crack-line:nth-child(2) { --crack-angle: 45deg; }
.crack-line:nth-child(3) { --crack-angle: 90deg; }
.crack-line:nth-child(4) { --crack-angle: 135deg; }
.crack-line:nth-child(5) { --crack-angle: 22.5deg; }
.crack-line:nth-child(6) { --crack-angle: 67.5deg; }
.crack-line:nth-child(7) { --crack-angle: 112.5deg; }
.crack-line:nth-child(8) { --crack-angle: 157.5deg; }

.shockwave {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100px;
  height: 100px;
  border: 3px solid rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  transform: translate(-50%, -50%);
  animation: shockwave-expand 0.6s ease-out forwards;
}

.flash-effect {
  position: absolute;
  inset: -20%;
  background: radial-gradient(circle, rgba(255,255,255,0.9), transparent);
  animation: flash-bright 0.3s ease-out forwards;
}

.skip-hint {
  position: absolute;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
}

@keyframes slide-in-left {
  from {
    transform: translateX(-100px) scale(1.5);
    opacity: 0;
  }
  to {
    transform: translateX(0) scale(1.5);
    opacity: 1;
  }
}

@keyframes slide-in-right {
  from {
    transform: translateX(100px) scale(1.5);
    opacity: 0;
  }
  to {
    transform: translateX(0) scale(1.5);
    opacity: 1;
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0) scale(1.5); }
  25% { transform: translateX(-5px) scale(1.5); }
  75% { transform: translateX(5px) scale(1.5); }
}

@keyframes mutual-destruction-shake {
  0% { transform: translate(0, 0) scale(1.5); }
  5% { transform: translate(-8px, -3px) scale(1.5); }
  10% { transform: translate(8px, 3px) scale(1.5); }
  15% { transform: translate(-6px, 4px) scale(1.5); }
  20% { transform: translate(6px, -4px) scale(1.5); }
  25% { transform: translate(-10px, -2px) scale(1.5); }
  30% { transform: translate(10px, 2px) scale(1.5); }
  35% { transform: translate(-7px, 5px) scale(1.5); }
  40% { transform: translate(7px, -5px) scale(1.5); }
  45% { transform: translate(-12px, -1px) scale(1.5); }
  50% { transform: translate(12px, 1px) scale(1.5); }
  55% { transform: translate(-9px, 6px) scale(1.5); }
  60% { transform: translate(9px, -6px) scale(1.5); }
  65% { transform: translate(-5px, -7px) scale(1.5); }
  70% { transform: translate(5px, 7px) scale(1.5); }
  75% { transform: translate(-3px, 3px) scale(1.5); }
  80% { transform: translate(3px, -3px) scale(1.5); }
  85% { transform: translate(-2px, -2px) scale(1.5); }
  90% { transform: translate(2px, 2px) scale(1.5); }
  95% { transform: translate(-1px, 1px) scale(1.5); }
  100% { transform: translate(0, 0) scale(1.5); }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.2); }
}

@keyframes mutual-pulse {
  0%, 100% { 
    transform: scale(1);
    text-shadow: 0 0 30px rgba(255, 107, 53, 0.9);
  }
  50% { 
    transform: scale(1.3);
    text-shadow: 0 0 50px rgba(255, 107, 53, 1);
  }
}

@keyframes gradient-shift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes lightning {
  0% { opacity: 0; width: 0; }
  50% { opacity: 1; width: 800px; }
  100% { opacity: 0; width: 800px; }
}

@keyframes impact {
  0% { 
    opacity: 0; 
    transform: translate(-50%, -50%) scale(0);
  }
  50% { 
    opacity: 1; 
    transform: translate(-50%, -50%) scale(1);
  }
  100% { 
    opacity: 0; 
    transform: translate(-50%, -50%) scale(2);
  }
}

@keyframes glass-shatter {
  0% {
    transform: translate(-50%, -50%) scale(1) rotate(0deg);
    opacity: 1;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translate(
      calc(-50% + var(--distance) * cos(var(--angle))),
      calc(-50% + var(--distance) * sin(var(--angle)))
    ) scale(0.3) rotate(var(--rotation));
    opacity: 0;
  }
}

@keyframes crack-appear {
  0% {
    transform: translate(-50%, -50%) rotate(var(--crack-angle)) scaleY(0);
    opacity: 0;
  }
  50% {
    transform: translate(-50%, -50%) rotate(var(--crack-angle)) scaleY(1);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -50%) rotate(var(--crack-angle)) scaleY(1);
    opacity: 0.3;
  }
}

@keyframes shockwave-expand {
  0% {
    transform: translate(-50%, -50%) scale(0.5);
    opacity: 1;
    border-width: 8px;
  }
  100% {
    transform: translate(-50%, -50%) scale(3);
    opacity: 0;
    border-width: 1px;
  }
}

@keyframes flash-bright {
  0% {
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}

.battle-fade-enter-active,
.battle-fade-leave-active {
  transition: opacity 0.3s ease;
}

.battle-fade-enter-from,
.battle-fade-leave-to {
  opacity: 0;
}

.destruction-enter-active {
  transition: all 0.3s ease;
}

.destruction-enter-from {
  opacity: 0;
  transform: scale(0);
}

.player-silhouette {
  position: relative;
  width: 112px;
  height: 160px;
  transform: scale(1.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.silhouette-svg {
  filter: drop-shadow(0 0 10px rgba(239, 68, 68, 0.5));
}

.direct-attack-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #ef4444;
  font-size: 0.875rem;
  font-weight: bold;
  text-shadow: 0 0 10px rgba(239, 68, 68, 0.8);
  white-space: nowrap;
  animation: flash 0.5s ease-in-out infinite alternate;
}

@keyframes flash {
  from { opacity: 0.8; }
  to { opacity: 1; }
}

.damage-effect {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
}

.damage-number {
  font-size: 4rem;
  font-weight: bold;
  color: #ef4444;
  text-shadow: 
    0 0 20px rgba(239, 68, 68, 0.8),
    0 0 40px rgba(239, 68, 68, 0.6);
  animation: damage-float 1s ease-out forwards;
}

@keyframes damage-float {
  0% {
    transform: translate(-50%, -50%) scale(0.5);
    opacity: 0;
  }
  50% {
    transform: translate(-50%, -70%) scale(1.5);
    opacity: 1;
  }
  100% {
    transform: translate(-50%, -100%) scale(1);
    opacity: 0;
  }
}

.damage-enter-active {
  transition: all 0.3s ease;
}

.damage-enter-from {
  opacity: 0;
}
</style>