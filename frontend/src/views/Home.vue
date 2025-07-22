<template>
  <div class="text-center">
    <h1 class="text-4xl font-bold text-white mb-8">めめめのくらげ TCG</h1>
    
    <div class="max-w-4xl mx-auto mb-12">
      <p class="text-xl text-gray-300 mb-6">
        「めめめのくらげ」のキャラクターを使った1対1の対戦型トレーディングカードゲーム
      </p>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <router-link
          to="/collection"
          class="bg-gradient-to-r from-blue-600 to-blue-800 p-6 rounded-lg hover:from-blue-700 hover:to-blue-900 transition-all transform hover:scale-105"
        >
          <div class="text-white">
            <h3 class="text-xl font-bold mb-2">カードコレクション</h3>
            <p class="text-blue-200">所有しているカードを確認する</p>
          </div>
        </router-link>
        
        <router-link
          to="/deck-builder"
          class="bg-gradient-to-r from-green-600 to-green-800 p-6 rounded-lg hover:from-green-700 hover:to-green-900 transition-all transform hover:scale-105"
        >
          <div class="text-white">
            <h3 class="text-xl font-bold mb-2">デッキビルダー</h3>
            <p class="text-green-200">50枚のデッキを構築する</p>
          </div>
        </router-link>
        
        <router-link
          to="/game"
          class="bg-gradient-to-r from-purple-600 to-purple-800 p-6 rounded-lg hover:from-purple-700 hover:to-purple-900 transition-all transform hover:scale-105"
        >
          <div class="text-white">
            <h3 class="text-xl font-bold mb-2">対戦</h3>
            <p class="text-purple-200">他のプレイヤーと対戦する</p>
          </div>
        </router-link>
      </div>
    </div>
    
    <div class="bg-black bg-opacity-30 rounded-lg p-8 max-w-2xl mx-auto">
      <h2 class="text-2xl font-bold text-white mb-4">ゲームルール</h2>
      <div class="text-left text-gray-300 space-y-2">
        <p>• デッキは50枚で構築</p>
        <p>• 同じカードは4枚まで</p>
        <p>• ふれんど、サポート、フィールドの3種類のカード</p>
        <p>• 相手の負のエネルギーエリアに7枚置くか、デッキ切れで勝利</p>
        <p>• 5つのフェーズ（スタート、ドロー、エネルギー、メイン、エンド）</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Home page component
</script>