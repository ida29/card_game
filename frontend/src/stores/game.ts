import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GameState, PlayerState, Card, DeckCard } from '@/types'

export type BattleMode = 'pvp' | 'cpu'
export type GamePhase = 'setup' | 'draw' | 'main' | 'battle' | 'end' | 'game_over'
export type CPUDifficulty = 'easy' | 'normal' | 'hard'

interface BattleAction {
  attacker: number
  target: number | 'player'
  damage: number
}

export const useGameStore = defineStore('game', () => {
  const battleMode = ref<BattleMode | null>(null)
  const cpuDifficulty = ref<CPUDifficulty>('normal')
  const gameState = ref<GameState | null>(null)
  const currentPhase = ref<GamePhase>('setup')
  const currentPlayer = ref<'player' | 'opponent'>('player')
  const turnCount = ref(0)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const battleActions = ref<BattleAction[]>([])
  const selectedAttacker = ref<number | null>(null)
  const availableTargets = ref<number[]>([])
  const cpuThinkingTime = ref(1000)

  const player = computed(() => gameState.value?.players.player)
  const opponent = computed(() => gameState.value?.players.opponent)
  const isPlayerTurn = computed(() => currentPlayer.value === 'player')
  const canPlayCards = computed(() => 
    isPlayerTurn.value && currentPhase.value === 'main'
  )

  function initializeGame(mode: BattleMode, playerDeck: DeckCard[]) {
    battleMode.value = mode
    turnCount.value = 0
    currentPhase.value = 'setup'
    
    const playerState: PlayerState = {
      deck: shuffleDeck([...playerDeck]),
      hand: [],
      friends: [],
      energy: [],
      graveyard: [],
      life: 6,
      supports: []
    }

    const opponentState: PlayerState = {
      deck: mode === 'cpu' ? generateCPUDeck() : [],
      hand: [],
      friends: [],
      energy: [],
      graveyard: [],
      life: 6,
      supports: []
    }

    gameState.value = {
      players: {
        player: playerState,
        opponent: opponentState
      },
      currentPlayer: 'player' as const,
      turn: 1,
      winner: null
    }

    // Draw initial hands
    drawCards('player', 5)
    drawCards('opponent', 5)
    
    currentPhase.value = 'draw'
  }

  function drawCards(player: 'player' | 'opponent', count: number) {
    if (!gameState.value) return
    
    const playerState = gameState.value.players[player]
    for (let i = 0; i < count; i++) {
      if (playerState.deck.length > 0) {
        const card = playerState.deck.shift()!
        playerState.hand.push(card)
      }
    }
  }

  function playCard(player: 'player' | 'opponent', cardIndex: number, targetZone: 'friends' | 'energy') {
    if (!gameState.value || !canPlayCards.value) return false
    
    const playerState = gameState.value.players[player]
    if (cardIndex < 0 || cardIndex >= playerState.hand.length) return false
    
    const card = playerState.hand[cardIndex]
    
    if (targetZone === 'friends') {
      if (playerState.friends.length >= 5) return false
      playerState.friends.push(card)
    } else if (targetZone === 'energy') {
      if (playerState.energy.length >= 10) return false
      playerState.energy.push(card)
    }
    
    playerState.hand.splice(cardIndex, 1)
    return true
  }

  function nextPhase() {
    if (!gameState.value) return
    
    switch (currentPhase.value) {
      case 'draw':
        currentPhase.value = 'main'
        break
      case 'main':
        currentPhase.value = 'battle'
        break
      case 'battle':
        currentPhase.value = 'end'
        break
      case 'end':
        endTurn()
        break
    }
  }

  function endTurn() {
    if (!gameState.value) return
    
    currentPlayer.value = currentPlayer.value === 'player' ? 'opponent' : 'player'
    gameState.value.currentPlayer = currentPlayer.value
    
    if (currentPlayer.value === 'player') {
      turnCount.value++
      gameState.value.turn++
    }
    
    currentPhase.value = 'draw'
    drawCards(currentPlayer.value, 1)
    
    if (battleMode.value === 'cpu' && currentPlayer.value === 'opponent') {
      executeCPUTurn()
    }
  }

  function executeCPUTurn() {
    setTimeout(() => {
      if (!gameState.value) return
      
      const cpu = gameState.value.players.opponent
      const player = gameState.value.players.player
      
      // CPU AI Strategy
      analyzeBoardState()
      
      // Phase 1: Play energy cards (prioritize if we have less than 3)
      if (cpu.energy.length < 3) {
        const energyCards = cpu.hand
          .map((card, index) => ({ card: card.card, index }))
          .filter(item => item.card?.type === 'エネルギー')
        
        // Play up to 1 energy card per turn
        if (energyCards.length > 0) {
          const cardToPlay = energyCards[0]
          playCard('opponent', cardToPlay.index, 'energy')
        }
      }
      
      // Phase 2: Play friend cards based on strategy
      const availableEnergy = cpu.energy.length
      const friendsInHand = cpu.hand
        .map((card, index) => ({ card: card.card, index }))
        .filter(item => item.card?.type === 'ふれんど')
        .sort((a, b) => {
          // Sort by power (descending) and cost (ascending)
          const powerDiff = (b.card?.power || 0) - (a.card?.power || 0)
          const costDiff = (a.card?.cost || 0) - (b.card?.cost || 0)
          return powerDiff + costDiff * 100
        })
      
      // Play friends based on available energy and board state
      let energyUsed = 0
      for (const friendCard of friendsInHand) {
        if (!friendCard.card) continue
        
        const cost = friendCard.card.cost || 1
        if (energyUsed + cost <= availableEnergy && cpu.friends.length < 5) {
          // Strategic decision: Play stronger cards when player has more friends
          const shouldPlay = 
            player.friends.length > cpu.friends.length || // We're behind
            (friendCard.card.power || 0) >= 2000 || // Strong card
            cpu.friends.length < 2 // Need board presence
          
          if (shouldPlay) {
            if (playCard('opponent', friendCard.index, 'friends')) {
              energyUsed += cost
            }
          }
        }
      }
      
      // Phase 3: Battle phase decisions
      currentPhase.value = 'battle'
      
      // CPU battle strategy based on difficulty
      if (cpu.friends.length > 0) {
        executeCPUBattlePhase()
      }
      
      // End turn after a delay
      setTimeout(() => {
        currentPhase.value = 'end'
        nextPhase()
      }, cpuThinkingTime.value)
    }, cpuThinkingTime.value)
  }
  
  function analyzeBoardState() {
    if (!gameState.value) return
    
    const cpu = gameState.value.players.opponent
    const player = gameState.value.players.player
    
    // Analyze board state for strategic decisions
    const boardAnalysis = {
      cpuFriends: cpu.friends.length,
      playerFriends: player.friends.length,
      cpuEnergy: cpu.energy.length,
      playerEnergy: player.energy.length,
      cpuLife: cpu.life,
      playerLife: player.life,
      handSize: cpu.hand.length,
      deckSize: cpu.deck.length
    }
    
    console.log('CPU Board Analysis:', boardAnalysis)
    return boardAnalysis
  }
  
  function executeCPUBattlePhase() {
    if (!gameState.value) return
    
    const cpu = gameState.value.players.opponent
    const player = gameState.value.players.player
    
    // Get all CPU friends that can attack
    const cpuAttackers = cpu.friends.map((friend, index) => ({
      index,
      card: friend.card,
      power: friend.card?.power || 0
    })).filter(attacker => attacker.power > 0)
    
    // Different strategies based on difficulty
    switch (cpuDifficulty.value) {
      case 'easy':
        executeEasyCPUBattle(cpuAttackers, player)
        break
      case 'normal':
        executeNormalCPUBattle(cpuAttackers, player)
        break
      case 'hard':
        executeHardCPUBattle(cpuAttackers, player)
        break
    }
  }
  
  function executeEasyCPUBattle(attackers: any[], player: PlayerState) {
    // Easy: Random attacks, sometimes skips attacking
    attackers.forEach((attacker, i) => {
      if (Math.random() > 0.3) { // 70% chance to attack
        setTimeout(() => {
          if (player.friends.length > 0) {
            // Attack random friend
            const targetIndex = Math.floor(Math.random() * player.friends.length)
            performBattle('opponent', attacker.index, targetIndex)
          } else {
            // Attack player directly
            performBattle('opponent', attacker.index, 'player')
          }
        }, i * 500)
      }
    })
  }
  
  function executeNormalCPUBattle(attackers: any[], player: PlayerState) {
    // Normal: Prioritize weaker targets, always attacks
    const playerTargets = player.friends.map((friend, index) => ({
      index,
      power: friend.card?.power || 0
    })).sort((a, b) => a.power - b.power)
    
    attackers.forEach((attacker, i) => {
      setTimeout(() => {
        if (playerTargets.length > 0) {
          // Find best target (weakest that can be defeated)
          const killableTargets = playerTargets.filter(target => 
            attacker.power >= target.power
          )
          
          if (killableTargets.length > 0) {
            performBattle('opponent', attacker.index, killableTargets[0].index)
          } else {
            // Attack strongest if can't kill any
            performBattle('opponent', attacker.index, playerTargets[playerTargets.length - 1].index)
          }
        } else {
          // Attack player directly
          performBattle('opponent', attacker.index, 'player')
        }
      }, i * 600)
    })
  }
  
  function executeHardCPUBattle(attackers: any[], player: PlayerState) {
    // Hard: Optimal targeting, considers life totals, maximizes damage
    const playerTargets = player.friends.map((friend, index) => ({
      index,
      power: friend.card?.power || 0,
      threat: friend.card?.power || 0 // Could be enhanced with abilities
    }))
    
    // Sort attackers by power (strongest first for hard mode)
    const sortedAttackers = [...attackers].sort((a, b) => b.power - a.power)
    
    sortedAttackers.forEach((attacker, i) => {
      setTimeout(() => {
        // If player is low on life and we can attack directly, prioritize that
        if (player.life <= 2 && playerTargets.length === 0) {
          performBattle('opponent', attacker.index, 'player')
          return
        }
        
        if (playerTargets.length > 0) {
          // Find optimal target
          const optimalTargets = playerTargets
            .filter(target => attacker.power >= target.power)
            .sort((a, b) => b.threat - a.threat)
          
          if (optimalTargets.length > 0) {
            // Kill highest threat that we can defeat
            performBattle('opponent', attacker.index, optimalTargets[0].index)
          } else {
            // If can't kill anything, attack highest threat to damage it
            const highestThreat = playerTargets.sort((a, b) => b.threat - a.threat)[0]
            performBattle('opponent', attacker.index, highestThreat.index)
          }
        } else {
          // Attack player directly with optimal timing
          performBattle('opponent', attacker.index, 'player')
        }
      }, i * 400) // Faster execution for hard mode
    })
  }

  function shuffleDeck(deck: DeckCard[]): DeckCard[] {
    const shuffled = [...deck]
    for (let i = shuffled.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]]
    }
    return shuffled
  }

  function generateCPUDeck(): DeckCard[] {
    // Generate a balanced CPU deck with proper card distribution
    const deck: DeckCard[] = []
    let cardId = 0
    
    // Energy cards (15 cards - balanced colors)
    const energyColors = ['赤', '青', '黄', '緑', '紫']
    for (let i = 0; i < 15; i++) {
      deck.push({
        ID: cardId++,
        deck_id: 0,
        card_no: `CPU-E-${i}`,
        quantity: 1,
        card: {
          ID: cardId,
          card_no: `CPU-E-${i}`,
          name: `${energyColors[i % 5]}エネルギー`,
          type: 'エネルギー',
          color: energyColors[i % 5],
          rarity: 'C',
          cost: 0,
          power: 0,
          description: 'CPUのエネルギーカード',
          image_url: '',
          artist: 'CPU',
          set_name: 'CPU Deck'
        }
      })
    }
    
    // Low cost friends (10 cards)
    for (let i = 0; i < 10; i++) {
      const color = energyColors[i % 5]
      deck.push({
        ID: cardId++,
        deck_id: 0,
        card_no: `CPU-F-L-${i}`,
        quantity: 1,
        card: {
          ID: cardId,
          card_no: `CPU-F-L-${i}`,
          name: `${color}の見習い`,
          type: 'ふれんど',
          color: color,
          rarity: 'C',
          cost: 1,
          power: 1000 + (i * 100),
          description: '低コストのCPUふれんど',
          image_url: '',
          artist: 'CPU',
          set_name: 'CPU Deck'
        }
      })
    }
    
    // Medium cost friends (10 cards)
    for (let i = 0; i < 10; i++) {
      const color = energyColors[i % 5]
      deck.push({
        ID: cardId++,
        deck_id: 0,
        card_no: `CPU-F-M-${i}`,
        quantity: 1,
        card: {
          ID: cardId,
          card_no: `CPU-F-M-${i}`,
          name: `${color}の戦士`,
          type: 'ふれんど',
          color: color,
          rarity: 'U',
          cost: 2,
          power: 2000 + (i * 150),
          description: '中コストのCPUふれんど',
          image_url: '',
          artist: 'CPU',
          set_name: 'CPU Deck'
        }
      })
    }
    
    // High cost friends (5 cards)
    for (let i = 0; i < 5; i++) {
      const color = energyColors[i % 5]
      deck.push({
        ID: cardId++,
        deck_id: 0,
        card_no: `CPU-F-H-${i}`,
        quantity: 1,
        card: {
          ID: cardId,
          card_no: `CPU-F-H-${i}`,
          name: `${color}の英雄`,
          type: 'ふれんど',
          color: color,
          rarity: 'R',
          cost: 3,
          power: 3000 + (i * 200),
          description: '高コストのCPUふれんど',
          image_url: '',
          artist: 'CPU',
          set_name: 'CPU Deck'
        }
      })
    }
    
    return shuffleDeck(deck)
  }

  function resetGame() {
    gameState.value = null
    battleMode.value = null
    currentPhase.value = 'setup'
    currentPlayer.value = 'player'
    turnCount.value = 0
    error.value = null
  }

  function performBattle(attacker: 'player' | 'opponent', attackerIndex: number, targetIndex: number | 'player') {
    if (!gameState.value) return false
    
    const attackingPlayer = gameState.value.players[attacker]
    const defendingPlayer = gameState.value.players[attacker === 'player' ? 'opponent' : 'player']
    
    if (!attackingPlayer.friends[attackerIndex]) return false
    
    const attackerCard = attackingPlayer.friends[attackerIndex].card
    if (!attackerCard || !attackerCard.power) return false
    
    if (targetIndex === 'player') {
      // Direct attack on player
      defendingPlayer.life = Math.max(0, defendingPlayer.life - 1)
      
      if (defendingPlayer.life === 0) {
        currentPhase.value = 'game_over'
        gameState.value.winner = attacker
      }
    } else {
      // Attack on friend
      const targetCard = defendingPlayer.friends[targetIndex]?.card
      if (!targetCard || !targetCard.power) return false
      
      // Simple battle: compare power
      if (attackerCard.power >= targetCard.power) {
        // Target is defeated
        const defeated = defendingPlayer.friends.splice(targetIndex, 1)[0]
        defendingPlayer.graveyard.push(defeated)
      }
      
      if (targetCard.power >= attackerCard.power) {
        // Attacker is defeated
        const defeated = attackingPlayer.friends.splice(attackerIndex, 1)[0]
        attackingPlayer.graveyard.push(defeated)
      }
    }
    
    return true
  }
  
  function selectBattleTarget(attackerIndex: number) {
    if (!gameState.value || !isPlayerTurn.value || currentPhase.value !== 'battle') return
    
    selectedAttacker.value = attackerIndex
    
    // Calculate available targets
    const opponentFriends = gameState.value.players.opponent.friends
    availableTargets.value = opponentFriends.map((_, index) => index)
    
    // Can attack player directly if no friends
    if (opponentFriends.length === 0) {
      availableTargets.value = [-1] // -1 represents direct player attack
    }
  }
  
  function executeBattle(targetIndex: number | 'player') {
    if (selectedAttacker.value === null) return
    
    performBattle('player', selectedAttacker.value, targetIndex)
    selectedAttacker.value = null
    availableTargets.value = []
  }
  
  function setCPUDifficulty(difficulty: CPUDifficulty) {
    cpuDifficulty.value = difficulty
    cpuThinkingTime.value = difficulty === 'easy' ? 500 : difficulty === 'normal' ? 1000 : 1500
  }

  return {
    battleMode,
    cpuDifficulty,
    gameState,
    currentPhase,
    currentPlayer,
    turnCount,
    loading,
    error,
    battleActions,
    selectedAttacker,
    availableTargets,
    player,
    opponent,
    isPlayerTurn,
    canPlayCards,
    initializeGame,
    drawCards,
    playCard,
    nextPhase,
    endTurn,
    resetGame,
    performBattle,
    selectBattleTarget,
    executeBattle,
    setCPUDifficulty
  }
})