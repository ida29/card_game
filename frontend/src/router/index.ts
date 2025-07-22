import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Collection from '@/views/Collection.vue'
import DeckBuilder from '@/views/DeckBuilder.vue'
import Game from '@/views/Game.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/collection',
      name: 'collection',
      component: Collection
    },
    {
      path: '/deck-builder',
      name: 'deck-builder',
      component: DeckBuilder
    },
    {
      path: '/game',
      name: 'game',
      component: Game
    }
  ]
})

export default router