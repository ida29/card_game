<template>
  <div id="app" class="min-h-screen bg-gradient-to-br from-blue-900 via-purple-900 to-indigo-900">
    <nav class="bg-black bg-opacity-50 backdrop-blur-sm">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-white text-xl font-bold">
              めめめのくらげ TCG
            </router-link>
          </div>
          <div class="flex space-x-4">
            <router-link
              to="/collection"
              class="text-white hover:text-blue-300 transition-colors"
            >
              カードコレクション
            </router-link>
            <router-link
              to="/deck-builder"
              class="text-white hover:text-blue-300 transition-colors"
            >
              デッキビルダー
            </router-link>
            <router-link
              to="/game"
              class="text-white hover:text-blue-300 transition-colors"
            >
              対戦
            </router-link>
          </div>
        </div>
      </div>
    </nav>
    
    <main class="container mx-auto px-4 py-8">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
// App component setup
</script>

<style scoped>
/* App-specific styles */
</style>