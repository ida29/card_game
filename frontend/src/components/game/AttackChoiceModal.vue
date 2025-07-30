<template>
  <teleport to="body">
    <transition name="modal">
      <div v-if="show" class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
        <div class="bg-gray-800 rounded-lg p-6 max-w-md w-full mx-4">
          <h3 class="text-xl font-bold text-white mb-4">アクションを選択</h3>
          
          <div class="space-y-3">
            <!-- Attack Option -->
            <button
              @click="handleAttack"
              class="w-full bg-red-600 hover:bg-red-700 text-white font-bold py-3 px-4 rounded transition-colors flex items-center justify-center gap-2"
            >
              <span class="text-2xl">⚔️</span>
              <span>アタックする</span>
            </button>
            
            <!-- Use Effect Option -->
            <button
              @click="handleEffect"
              class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-4 rounded transition-colors flex items-center justify-center gap-2"
              :disabled="!hasAttackEffect"
              :class="{ 'opacity-50 cursor-not-allowed': !hasAttackEffect }"
            >
              <span class="text-2xl">✨</span>
              <span>効果を使用する</span>
            </button>
            
            <!-- Effect Description -->
            <div v-if="hasAttackEffect" class="text-sm text-gray-300 text-center">
              {{ effectDescription }}
            </div>
            
            <!-- Cancel -->
            <button
              @click="handleCancel"
              class="w-full bg-gray-600 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded transition-colors"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { DeckCard } from '@/types'

const props = defineProps<{
  show: boolean
  attackerCard: DeckCard | null
}>()

const emit = defineEmits<{
  'choose-attack': []
  'choose-effect': []
  'cancel': []
}>()

const hasAttackEffect = computed(() => {
  if (!props.attackerCard?.card) return false
  const cardNo = props.attackerCard.card.card_no
  
  // List of cards with attack effects
  const attackEffectCards = [
    'F-006', 'F-006 (P)', // ヒヤケラトプス
    'F-008', 'F-008 (P)', // ボーイ
    'F-011', 'F-011 (P)', // ポチ
    'F-020', 'F-020 (P)', // マルカニ
    'F-022', 'F-022 (P)', // ジョニー
    'F-044', 'F-044 (P)', // うっきー
    'F-055', 'F-055 (P)', // Ko2
    'F-102' // くらげ坊(変身)
  ]
  
  return attackEffectCards.includes(cardNo)
})

const effectDescription = computed(() => {
  if (!props.attackerCard?.card) return ''
  const cardNo = props.attackerCard.card.card_no
  
  const effectDescriptions: Record<string, string> = {
    'F-006': 'アタック時：カードを1枚引く',
    'F-006 (P)': 'アタック時：カードを1枚引く',
    'F-008': 'アタック時：パワー3000以下を破壊',
    'F-008 (P)': 'アタック時：パワー3000以下を破壊',
    'F-011': 'アタック時：カードを1枚引く',
    'F-011 (P)': 'アタック時：カードを1枚引く',
    'F-020': 'アタック時：デッキトップを確認',
    'F-020 (P)': 'アタック時：デッキトップを確認',
    'F-022': 'アタック時：デッキトップを破棄',
    'F-022 (P)': 'アタック時：デッキトップを破棄',
    'F-044': 'アタック時：負のエネルギーを破棄',
    'F-044 (P)': 'アタック時：負のエネルギーを破棄',
    'F-055': 'アタック時：相手のふれんどをレスト',
    'F-055 (P)': 'アタック時：相手のふれんどをレスト',
    'F-102': 'アタック時：負のエネルギーでアクティブ'
  }
  
  return effectDescriptions[cardNo] || '効果を発動'
})

const handleAttack = () => {
  emit('choose-attack')
}

const handleEffect = () => {
  emit('choose-effect')
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped>
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s;
}

.modal-enter-from, .modal-leave-to {
  opacity: 0;
}
</style>