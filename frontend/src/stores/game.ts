import { defineStore } from 'pinia'
import { ref, computed, nextTick } from 'vue'
import type { Card, DeckCard, CardType, CardColor, CardRarity } from '@/types'
import { useCardStore } from '@/stores/cards'

export type BattleMode = 'pvp' | 'cpu'
export type GamePhase = 'setup' | 'start' | 'draw' | 'energy' | 'main' | 'end' | 'game_over'
export type CPUDifficulty = 'easy' | 'normal' | 'hard'

interface BattleAction {
  attacker: number
  target: number | 'player'
  damage: number
}

interface CardState {
  card: DeckCard
  tapped: boolean
  justPlayed: boolean
  turnPlayed: number
}

interface EnergyCardState {
  card: DeckCard
  tapped: boolean
}

interface FriendCardState {
  card: DeckCard
  tapped: boolean
  playedTurn: number // Track which turn this friend was played
}

interface NegativeEnergyCardState {
  card: DeckCard
  faceUp: boolean // True = face up (can be inspected), False = face down (used as energy)
}

interface PlayerState {
  deck: DeckCard[]
  hand: DeckCard[]
  friends: FriendCardState[]
  energy: EnergyCardState[]
  graveyard: DeckCard[]
  negativeEnergy: NegativeEnergyCardState[]
  life: number
  supports: DeckCard[]
  field?: DeckCard
}

interface GameState {
  players: {
    player: PlayerState
    opponent: PlayerState
  }
  currentPlayer: 'player' | 'opponent'
  turn: number
  winner: 'player' | 'opponent' | null
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
  const cpuThinkingTime = ref(2500) // Increased base thinking time
  const energyPlayedThisTurn = ref<{player: boolean, opponent: boolean}>({player: false, opponent: false})
  
  // Blocking decision state
  const blockingDecision = ref<{
    show: boolean
    attacker: { player: 'player' | 'opponent', index: number, card: DeckCard | null }
    availableBlockers: number[]
    resolve: ((blocker: number | null) => void) | null
  }>({
    show: false,
    attacker: { player: 'player', index: -1, card: null },
    availableBlockers: [],
    resolve: null
  })
  
  // Battle animation state
  const battleAnimation = ref<{
    show: boolean
    attacker: DeckCard | null
    defender: DeckCard | null
    attackerDefeated: boolean
    defenderDefeated: boolean
  }>({
    show: false,
    attacker: null,
    defender: null,
    attackerDefeated: false,
    defenderDefeated: false
  })
  
  // Counter decision state
  const counterDecision = ref<{
    show: boolean
    attacker: { player: 'player' | 'opponent', index: number, card: DeckCard | null }
    blocker: { player: 'player' | 'opponent', index: number, card: DeckCard | null } | null
    resolve: ((counter: DeckCard | null) => void) | null
  }>({
    show: false,
    attacker: { player: 'player', index: -1, card: null },
    blocker: null,
    resolve: null
  })

  // Energy cost selection state
  const energyCostSelection = ref<{
    show: boolean
    cardToPay: Card | null
    resolve: ((selection: any) => void) | null
  }>({
    show: false,
    cardToPay: null,
    resolve: null
  })

  const player = computed(() => gameState.value?.players.player)
  const opponent = computed(() => gameState.value?.players.opponent)
  const isPlayerTurn = computed(() => currentPlayer.value === 'player')
  const canPlayCards = computed(() => 
    isPlayerTurn.value && (currentPhase.value === 'main' || currentPhase.value === 'energy')
  )

  function initializeGame(mode: BattleMode, playerDeck: DeckCard[], startPhases: boolean = true) {
    console.log('Initializing game with mode:', mode)
    console.log('Player deck size:', playerDeck.length)
    battleMode.value = mode
    turnCount.value = 0
    currentPhase.value = 'setup'
    
    const playerState: PlayerState = {
      deck: shuffleDeck([...playerDeck]),
      hand: [],
      friends: [],
      energy: [],
      graveyard: [],
      negativeEnergy: [],
      life: 0, // Not used anymore - life is represented by negative energy
      supports: []
    }
    
    console.log('Player state deck size after shuffle:', playerState.deck.length)

    const opponentState: PlayerState = {
      deck: mode === 'cpu' ? generateCPUDeck() : [],
      hand: [],
      friends: [],
      energy: [],
      graveyard: [],
      negativeEnergy: [],
      life: 0, // Not used anymore - life is represented by negative energy
      supports: []
    }

    gameState.value = {
      players: {
        player: playerState,
        opponent: opponentState
      },
      currentPlayer: currentPlayer.value,
      turn: 1,
      winner: null
    }

    console.log('Game state initialized')
    
    // Draw initial hands for both players
    drawCards('player', 5)
    drawCards('opponent', 5)
    
    if (startPhases) {
      startGame()
    }
  }
  
  function startGame() {
    if (!gameState.value) return
    
    // Ensure currentPlayer is set in gameState
    if (gameState.value) {
      gameState.value.currentPlayer = currentPlayer.value
    }
    
    // Reset energy played flags at game start
    energyPlayedThisTurn.value = {player: false, opponent: false}
    
    // Start with start phase for the first player
    currentPhase.value = 'start'
    console.log('Game starting, starting player:', currentPlayer.value, 'phase:', currentPhase.value, 'isPlayerTurn:', isPlayerTurn.value)
    
    // Immediately schedule the phase progression
    console.log('Scheduling start phase progression...')
    
    // First timeout: untap cards
    window.setTimeout(() => {
      console.log('Start phase timeout 1: Untapping cards for:', currentPlayer.value)
      untapAllCards(currentPlayer.value)
      
      // Second timeout: progress to next phase
      window.setTimeout(() => {
        console.log('Start phase timeout 2: Progressing to next phase')
        nextPhase()
      }, 1000)
    }, 1000)
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
  
  function drawToEnergy(player: 'player' | 'opponent') {
    if (!gameState.value) return
    
    const playerState = gameState.value.players[player]
    if (playerState.deck.length > 0) {
      const card = playerState.deck.shift()!
      playerState.energy.push(card)
    }
  }
  
  function untapAllCards(player: 'player' | 'opponent') {
    if (!gameState.value) return
    
    const playerState = gameState.value.players[player]
    
    // Untap all energy cards
    playerState.energy.forEach(energy => {
      energy.tapped = false
    })
    
    // Untap all friend cards
    playerState.friends.forEach(friend => {
      friend.tapped = false
    })
  }
  
  function checkVictoryConditions() {
    if (!gameState.value) return
    
    const player = gameState.value.players.player
    const opponent = gameState.value.players.opponent
    
    // Check deck out condition - lose if deck is 0 at start of turn
    if (player.deck.length === 0 && currentPlayer.value === 'player') {
      currentPhase.value = 'game_over'
      gameState.value.winner = 'opponent'
      return
    }
    
    if (opponent.deck.length === 0 && currentPlayer.value === 'opponent') {
      currentPhase.value = 'game_over'
      gameState.value.winner = 'player'
      return
    }
    
    // Check negative energy condition - lose if negative energy reaches 6
    if (player.negativeEnergy.length >= 7) {
      currentPhase.value = 'game_over'
      gameState.value.winner = 'opponent'
      return
    }
    
    if (opponent.negativeEnergy.length >= 7) {
      currentPhase.value = 'game_over'
      gameState.value.winner = 'player'
      return
    }
  }

  async function playCard(player: 'player' | 'opponent', cardIndex: number, targetZone?: 'friends' | 'energy') {
    if (!gameState.value) return false
    
    // Check if it's the player's turn
    if (player !== currentPlayer.value) {
      console.log(`Cannot play card: it's not ${player}'s turn (current: ${currentPlayer.value})`)
      return false
    }
    
    // Only check canPlayCards for human player
    if (player === 'player' && !canPlayCards.value) return false
    
    const playerState = gameState.value.players[player]
    if (cardIndex < 0 || cardIndex >= playerState.hand.length) return false
    
    const card = playerState.hand[cardIndex]
    const cardData = card.card
    
    // Determine target zone based on card type if not specified
    if (!targetZone && currentPhase.value === 'main') {
      switch (cardData.type) {
        case 'ふれんど':
          targetZone = 'friends'
          break
        case 'サポート':
          return playSupport(player, cardIndex)
        case 'フィールド':
          return playField(player, cardIndex)
      }
    }
    
    if (targetZone === 'friends') {
      if (playerState.friends.length >= 10) return false
      
      // Check and pay cost
      if (!(await payCost(player, cardData))) {
        return false
      }
      
      playerState.friends.push({
        card: card,
        tapped: false,
        playedTurn: turnCount.value
      })
    } else if (targetZone === 'energy') {
      // Debug logging for energy setting
      console.log('Attempting to set energy:', {
        player,
        currentPlayer: currentPlayer.value,
        isPlayerTurn: player === currentPlayer.value,
        currentPhase: currentPhase.value,
        energyPlayedThisTurn: energyPlayedThisTurn.value,
        turnCount: turnCount.value,
        cardName: cardData.name
      })
      
      // Check if energy was already played this turn
      if (currentPhase.value === 'energy' && energyPlayedThisTurn.value[player]) {
        console.log(`${player} already played energy this turn`)
        return false
      }
      if (playerState.energy.length >= 10) return false
      playerState.energy.push({
        card: card,
        tapped: false
      })
      
      // Mark that energy was played this turn
      if (currentPhase.value === 'energy') {
        energyPlayedThisTurn.value[player] = true
        console.log(`Energy successfully set for ${player}, flag updated:`, energyPlayedThisTurn.value)
        
        // Auto-progress to main phase after playing energy
        setTimeout(() => {
          nextPhase()
        }, 500)
      }
    }
    
    playerState.hand.splice(cardIndex, 1)
    return true
  }
  
  function playSupport(player: 'player' | 'opponent', cardIndex: number): boolean {
    if (!gameState.value) return false
    
    const playerState = gameState.value.players[player]
    const card = playerState.hand[cardIndex]
    const cardData = card.card
    
    // Check if it's main phase
    if (currentPhase.value !== 'main') {
      console.log('Support cards can only be played during main phase')
      return false
    }
    
    // Check cost and pay
    if (!payCost(player, cardData)) {
      return false
    }
    
    // TODO: Execute support card effect
    console.log(`Playing support card: ${cardData.name}`)
    console.log(`Effect: ${cardData.effect}`)
    
    // Move to trash after effect
    playerState.graveyard.push(card)
    playerState.hand.splice(cardIndex, 1)
    
    return true
  }
  
  function playField(player: 'player' | 'opponent', cardIndex: number): boolean {
    if (!gameState.value) return false
    
    const playerState = gameState.value.players[player]
    const card = playerState.hand[cardIndex]
    const cardData = card.card
    
    // Check if it's main phase
    if (currentPhase.value !== 'main') {
      console.log('Field cards can only be played during main phase')
      return false
    }
    
    // Check cost and pay
    if (!payCost(player, cardData)) {
      return false
    }
    
    // If there's already a field card, move it to trash
    if (playerState.field) {
      playerState.graveyard.push(playerState.field)
    }
    
    // Place new field card
    playerState.field = card
    playerState.hand.splice(cardIndex, 1)
    
    console.log(`Playing field card: ${cardData.name}`)
    
    return true
  }
  
  async function payCost(player: 'player' | 'opponent', cardData: Card): Promise<boolean> {
    if (!gameState.value) return false
    
    const playerState = gameState.value.players[player]
    const totalCost = cardData.cost || 0
    
    if (totalCost === 0) return true
    
    // Check if we can pay the cost first
    if (!canPayCost(player, cardData)) {
      console.log(`Cannot pay cost for ${cardData.name}`)
      return false
    }
    
    // For human players, show selection UI
    if (player === 'player') {
      try {
        const selection = await getEnergyCostSelection(cardData)
        if (!selection) return false // Cancelled
        
        // Apply the selected energy consumption
        applyEnergySelection(selection)
        return true
      } catch (error) {
        console.log('Energy selection cancelled')
        return false
      }
    } else {
      // For CPU, use automatic selection
      return payCostAutomatically(player, cardData)
    }
  }
  
  function canPayCost(player: 'player' | 'opponent', cardData: Card): boolean {
    if (!gameState.value) return false
    
    const playerState = gameState.value.players[player]
    const totalCost = cardData.cost || 0
    
    if (totalCost === 0) return true
    
    // Get available energy sources
    const untappedEnergy = playerState.energy.filter(e => !e.tapped)
    const faceUpNegativeEnergy = playerState.negativeEnergy.filter(ne => ne.faceUp)
    
    // Calculate color requirements
    const colorRequirements = {
      red: cardData.cost_red || 0,
      blue: cardData.cost_blue || 0,
      yellow: cardData.cost_yellow || 0,
      green: cardData.cost_green || 0
    }
    
    // Check each color requirement
    for (const [colorKey, required] of Object.entries(colorRequirements)) {
      if (required > 0) {
        const colorName = getColorName(colorKey)
        const availableOfColor = [
          ...untappedEnergy.filter(e => e.card.card.color === colorName),
          ...faceUpNegativeEnergy.filter(ne => ne.card.card.color === colorName)
        ]
        const totalOfColor = availableOfColor.reduce((sum, source) => {
          return sum + (source.card.card.energy_value || 1)
        }, 0)
        
        if (totalOfColor < required) {
          return false
        }
      }
    }
    
    // Check total cost
    const totalAvailable = [
      ...untappedEnergy,
      ...faceUpNegativeEnergy
    ].reduce((sum, source) => sum + (source.card.card.energy_value || 1), 0)
    
    return totalAvailable >= totalCost
  }
  
  function getColorName(colorKey: string): string {
    switch (colorKey) {
      case 'red': return '赤'
      case 'blue': return '青' 
      case 'yellow': return '黄'
      case 'green': return '緑'
      default: return ''
    }
  }
  
  function applyEnergySelection(selection: any) {
    if (!gameState.value) return
    
    const playerState = gameState.value.players.player
    
    // Apply regular energy selections
    selection.regularEnergy.forEach((index: number) => {
      if (playerState.energy[index]) {
        playerState.energy[index].tapped = true
      }
    })
    
    // Apply negative energy selections
    selection.negativeEnergy.forEach((index: number) => {
      if (playerState.negativeEnergy[index]) {
        playerState.negativeEnergy[index].faceUp = false
      }
    })
  }
  
  function payCostAutomatically(player: 'player' | 'opponent', cardData: Card): boolean {
    // Use the existing automatic payment logic for CPU
    if (!gameState.value) return false
    
    const playerState = gameState.value.players[player]
    const totalCost = cardData.cost || 0
    
    // Get available energy sources
    const untappedEnergy = playerState.energy.filter(e => !e.tapped)
    const faceUpNegativeEnergy = playerState.negativeEnergy.filter(ne => ne.faceUp)
    
    const getEnergyValue = (energy: EnergyCardState) => energy.card.card.energy_value || 1
    const getNegativeEnergyValue = (negEnergy: NegativeEnergyCardState) => negEnergy.card.card.energy_value || 1
    
    // Check color requirements
    const colorRequirements = {
      '赤': cardData.cost_red || 0,
      '青': cardData.cost_blue || 0,
      '黄': cardData.cost_yellow || 0,
      '緑': cardData.cost_green || 0
    }
    
    let costRemaining = totalCost
    
    // Pay required colors first
    for (const [color, required] of Object.entries(colorRequirements)) {
      if (required > 0) {
        let colorCostRemaining = required
        const regularEnergy = untappedEnergy.filter(e => e.card.card.color === color)
        const negativeEnergy = faceUpNegativeEnergy.filter(ne => ne.card.card.color === color)
        const colorSources = [...regularEnergy, ...negativeEnergy]
        
        for (const source of colorSources) {
          if (colorCostRemaining <= 0) break
          
          if ('tapped' in source) {
            source.tapped = true
            const value = getEnergyValue(source)
            colorCostRemaining -= value
            costRemaining -= value
          } else {
            source.faceUp = false
            const value = getNegativeEnergyValue(source)
            colorCostRemaining -= value
            costRemaining -= value
          }
        }
      }
    }
    
    // Pay remaining colorless cost
    if (costRemaining > 0) {
      const remainingEnergy = playerState.energy.filter(e => !e.tapped)
      const remainingNegative = playerState.negativeEnergy.filter(ne => ne.faceUp)
      const allRemaining = [...remainingEnergy, ...remainingNegative]
      
      for (const source of allRemaining) {
        if (costRemaining <= 0) break
        
        if ('tapped' in source) {
          source.tapped = true
          costRemaining -= getEnergyValue(source)
        } else {
          source.faceUp = false
          costRemaining -= getNegativeEnergyValue(source)
        }
      }
    }
    
    return true
  }

  function nextPhase() {
    if (!gameState.value) return
    
    console.log('nextPhase called, current phase:', currentPhase.value)
    
    switch (currentPhase.value) {
      case 'start':
        currentPhase.value = 'draw'
        // Handle draw phase automatically
        setTimeout(() => {
          // First player doesn't draw on their first turn
          const isFirstPlayerFirstTurn = turnCount.value === 0 && currentPlayer.value === 'player'
          if (!isFirstPlayerFirstTurn) {
            drawCards(currentPlayer.value, 1)
          }
          
          // Move to energy phase
          currentPhase.value = 'energy'
          
          // For CPU, automatically progress to main phase
          if (currentPlayer.value === 'opponent' && battleMode.value === 'cpu') {
            setTimeout(() => {
              // Move to main phase
              currentPhase.value = 'main'
              executeCPUTurn()
            }, 1000)
          }
          // For player, wait for manual energy play
        }, 1000)
        break
      case 'draw':
        currentPhase.value = 'energy'
        console.log('Moved to energy phase, energyPlayedThisTurn status:', energyPlayedThisTurn.value)
        console.log('Current player:', currentPlayer.value, 'Turn count:', turnCount.value)
        // For CPU, handle energy phase automatically
        if (currentPlayer.value === 'opponent' && battleMode.value === 'cpu') {
          setTimeout(() => {
            drawToEnergy(currentPlayer.value)
            currentPhase.value = 'main'
          }, 2000) // Increased energy phase delay
        }
        // For player, wait for manual energy play
        break
      case 'energy':
        currentPhase.value = 'main'
        break
      case 'main':
        currentPhase.value = 'end'
        endTurn()
        break
      case 'end':
        endTurn()
        break
    }
  }

  function endTurn() {
    if (!gameState.value) return
    
    // End phase - discard down to 7 cards
    const playerState = gameState.value.players[currentPlayer.value]
    if (playerState.hand.length > 7) {
      // For now, discard from the end
      const discardCount = playerState.hand.length - 7
      for (let i = 0; i < discardCount; i++) {
        const discarded = playerState.hand.pop()
        if (discarded) {
          playerState.graveyard.push(discarded)
        }
      }
    }
    
    currentPlayer.value = currentPlayer.value === 'player' ? 'opponent' : 'player'
    gameState.value.currentPlayer = currentPlayer.value
    
    if (currentPlayer.value === 'player') {
      turnCount.value++
      gameState.value.turn++
    }
    
    // Start new turn with start phase
    currentPhase.value = 'start'
    
    // Reset turn-based flags
    energyPlayedThisTurn.value = {player: false, opponent: false}
    console.log('Turn started, energyPlayedThisTurn reset:', energyPlayedThisTurn.value)
    
    // Check victory conditions at start of turn
    checkVictoryConditions()
    
    // Progress through phases
    setTimeout(() => {
      // Untap all cards
      untapAllCards(currentPlayer.value)
      
      // Always auto-progress through start phase
      setTimeout(() => {
        nextPhase()
      }, 1000)
    }, 1000)
  }

  function executeCPUTurn() {
    // Only execute if it's actually CPU's turn
    if (currentPlayer.value !== 'opponent' || battleMode.value !== 'cpu') {
      console.log('executeCPUTurn called but not CPU turn, returning')
      return
    }
    
    const turnTimeout = setTimeout(() => {
      console.error('CPU turn timeout - forcing end turn')
      currentPhase.value = 'end'
      nextPhase()
    }, 10000) // 10 second timeout
    
    // Wait for main phase before making decisions
    const waitForMainPhase = setInterval(() => {
      if (currentPhase.value === 'main') {
        clearInterval(waitForMainPhase)
        
        setTimeout(async () => {
          try {
            if (!gameState.value) {
              clearTimeout(turnTimeout)
              return
            }
            
            const cpu = gameState.value.players.opponent
            const player = gameState.value.players.player
            
            console.log('CPU turn starting in main phase...')
            
            // CPU AI Strategy
            analyzeBoardState()
            
            // Phase 1: Play energy cards (prioritize if we have less than 3)
            if (cpu.energy.length < 3 && !energyPlayedThisTurn.value.opponent) {
              // Play any card as energy, prefer low-cost cards
              // Find the lowest cost card in hand
              let lowestCostIndex = -1
              let lowestCost = Infinity
              
              for (let i = 0; i < cpu.hand.length; i++) {
                const cost = cpu.hand[i].card.cost || 0
                if (cost < lowestCost) {
                  lowestCost = cost
                  lowestCostIndex = i
                }
              }
              
              if (lowestCostIndex >= 0) {
                await playCard('opponent', lowestCostIndex, 'energy')
              }
            }
        
        // Phase 2: Play friend cards based on strategy
        const availableEnergy = cpu.energy.length
        let energyUsed = 0
        
        // Play friends from the end of hand to avoid index issues
        for (let i = cpu.hand.length - 1; i >= 0; i--) {
          const card = cpu.hand[i].card
          if (!card || card.type !== 'ふれんど') continue
          
          const cost = card.cost || 1
          if (energyUsed + cost <= availableEnergy && cpu.friends.length < 10) {
            // Strategic decision: Play stronger cards when player has more friends
            const shouldPlay = 
              player.friends.length > cpu.friends.length || // We're behind
              (card.power || 0) >= 2000 || // Strong card
              cpu.friends.length < 2 // Need board presence
            
            if (shouldPlay) {
              // Add delay between card plays for visibility
              await new Promise(resolve => setTimeout(resolve, 800))
              if (await playCard('opponent', i, 'friends')) {
                energyUsed += cost
              }
            }
          }
        }
        
            // Phase 3: Battle decisions
            // CPU battle strategy based on difficulty
            if (cpu.friends.length > 0) {
              await executeCPUBattlePhase()
            }
            
            // End turn after a delay
            setTimeout(() => {
              clearTimeout(turnTimeout)
              currentPhase.value = 'end'
              nextPhase()
            }, cpuThinkingTime.value)
          } catch (error) {
            console.error('Error in CPU turn:', error)
            clearTimeout(turnTimeout)
            currentPhase.value = 'end'
            nextPhase()
          }
        }, cpuThinkingTime.value)
      }
    }, 100)
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
  
  function analyzeBattleSituation(cpuAttackers: any[], player: PlayerState) {
    if (!gameState.value) return { shouldAttack: true, reason: 'No game state' }
    
    const cpu = gameState.value.players.opponent
    
    // Get player's untapped friends that can block
    const playerBlockers = player.friends
      .map((friend, index) => ({ 
        index, 
        power: friend.card.card?.power || 0, 
        tapped: friend.tapped 
      }))
      .filter(blocker => !blocker.tapped)
    
    // Count available attackers vs blockers
    const attackerCount = cpuAttackers.length
    const blockerCount = playerBlockers.length
    
    // Check if player has stronger friends than CPU's strongest attacker
    const cpuMaxPower = Math.max(...cpuAttackers.map(a => a.power), 0)
    const playerMaxPower = Math.max(...playerBlockers.map(b => b.power), 0)
    
    // Strategic conditions to avoid attacking
    const isOutnumbered = attackerCount < blockerCount
    const isOutpowered = playerMaxPower > cpuMaxPower
    const shouldAvoidAttack = isOutnumbered && isOutpowered
    
    console.log('CPU Battle Analysis:', {
      attackerCount,
      blockerCount,
      cpuMaxPower,
      playerMaxPower,
      isOutnumbered,
      isOutpowered,
      shouldAvoidAttack
    })
    
    return {
      shouldAttack: !shouldAvoidAttack,
      reason: shouldAvoidAttack 
        ? `Avoiding attack: outnumbered (${attackerCount} vs ${blockerCount}) and outpowered (${cpuMaxPower} vs ${playerMaxPower})`
        : 'Safe to attack',
      attackerCount,
      blockerCount,
      cpuMaxPower,
      playerMaxPower,
      isOutnumbered,
      isOutpowered
    }
  }
  
  async function executeCPUBattlePhase() {
    if (!gameState.value) return
    
    const cpu = gameState.value.players.opponent
    const player = gameState.value.players.player
    
    // Get all CPU friends that can attack (not tapped and not played this turn)
    const cpuAttackers = cpu.friends.map((friend, index) => ({
      index,
      card: friend.card.card,
      power: friend.card.card?.power || 0,
      tapped: friend.tapped,
      playedTurn: friend.playedTurn
    })).filter(attacker => 
      attacker.power > 0 && 
      !attacker.tapped && 
      attacker.playedTurn < turnCount.value // Can't attack on the turn it was played
    )
    
    // Perform strategic analysis
    const strategicAnalysis = analyzeBattleSituation(cpuAttackers, player)
    
    // Different strategies based on difficulty
    switch (cpuDifficulty.value) {
      case 'easy':
        await executeEasyCPUBattle(cpuAttackers, player, strategicAnalysis)
        break
      case 'normal':
        await executeNormalCPUBattle(cpuAttackers, player, strategicAnalysis)
        break
      case 'hard':
        await executeHardCPUBattle(cpuAttackers, player, strategicAnalysis)
        break
    }
  }
  
  async function executeEasyCPUBattle(attackers: any[], player: PlayerState, analysis: any) {
    // Easy: Random attacks, sometimes skips attacking
    // Easy mode ignores strategic analysis 50% of the time
    const shouldIgnoreStrategy = Math.random() > 0.5
    
    if (!analysis.shouldAttack && !shouldIgnoreStrategy) {
      console.log('CPU (easy) decides not to attack:', analysis.reason)
      return
    }
    
    for (let i = 0; i < attackers.length; i++) {
      const attacker = attackers[i]
      if (Math.random() > 0.3) { // 70% chance to attack
        await new Promise(resolve => setTimeout(resolve, 1200))
        
        // Friends can only attack the player directly
        await performBattle('opponent', attacker.index, 'player')
      }
    }
  }
  
  async function executeNormalCPUBattle(attackers: any[], player: PlayerState, analysis: any) {
    // Normal: Strategic direct attacks, respects tactical analysis
    if (!analysis.shouldAttack) {
      console.log('CPU (normal) decides not to attack:', analysis.reason)
      return
    }
    
    // If safe to attack, proceed with all attackers
    for (let i = 0; i < attackers.length; i++) {
      const attacker = attackers[i]
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      // Friends can only attack the player directly
      // The player will get a chance to block
      await performBattle('opponent', attacker.index, 'player')
    }
  }
  
  async function executeHardCPUBattle(attackers: any[], player: PlayerState, analysis: any) {
    // Hard: Strategic direct attacks with optimal ordering and advanced analysis
    if (!analysis.shouldAttack) {
      console.log('CPU (hard) decides not to attack:', analysis.reason)
      return
    }
    
    // Additional hard mode analysis: selective attacking
    if (analysis.isOutnumbered && !analysis.isOutpowered) {
      // If only outnumbered but not outpowered, attack with strongest only
      console.log('CPU (hard) using selective attack strategy (outnumbered but stronger)')
      const strongestAttacker = attackers.reduce((prev, curr) => 
        prev.power > curr.power ? prev : curr
      )
      await new Promise(resolve => setTimeout(resolve, 400))
      await performBattle('opponent', strongestAttacker.index, 'player')
      return
    }
    
    // Sort attackers by power (strongest first for hard mode)
    const sortedAttackers = [...attackers].sort((a, b) => b.power - a.power)
    
    for (let i = 0; i < sortedAttackers.length; i++) {
      const attacker = sortedAttackers[i]
      await new Promise(resolve => setTimeout(resolve, 1000)) // Adjusted for better visibility
      
      // Friends can only attack the player directly
      // In hard mode, CPU attacks with strongest friends first
      // to pressure the player into difficult blocking decisions
      await performBattle('opponent', attacker.index, 'player')
    }
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
    const cardStore = useCardStore()
    const deck: DeckCard[] = []
    let deckCardId = 0
    
    // Get all available cards
    const allCards = cardStore.cards
    
    // If no cards are loaded, generate a basic deck with proper structure
    if (allCards.length === 0) {
      console.warn('No cards loaded, generating basic CPU deck')
      return generateBasicCPUDeck()
    }
    
    // Separate cards by type and cost
    const friendCards = allCards.filter(c => c.type === 'ふれんど' && c.power && c.power > 0)
    const lowCostFriends = friendCards.filter(c => c.cost <= 1)
    const midCostFriends = friendCards.filter(c => c.cost === 2)
    const highCostFriends = friendCards.filter(c => c.cost >= 3)
    const supportCards = allCards.filter(c => c.type === 'サポート')
    const fieldCards = allCards.filter(c => c.type === 'フィールド')
    
    // Build a balanced deck
    // 15 energy cards (use low-cost friends as energy)
    const energyCards = [...lowCostFriends, ...midCostFriends].slice(0, 15)
    energyCards.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // 15 low-cost friends
    const selectedLowCost = selectRandomCards(lowCostFriends, 15)
    selectedLowCost.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // 10 mid-cost friends
    const selectedMidCost = selectRandomCards(midCostFriends, 10)
    selectedMidCost.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // 5 high-cost friends
    const selectedHighCost = selectRandomCards(highCostFriends, 5)
    selectedHighCost.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // 3 support cards
    const selectedSupport = selectRandomCards(supportCards, 3)
    selectedSupport.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // 2 field cards
    const selectedField = selectRandomCards(fieldCards, 2)
    selectedField.forEach(card => {
      deck.push({
        ID: deckCardId++,
        deck_id: 0,
        card_no: card.card_no,
        quantity: 1,
        card: { ...card }
      })
    })
    
    // Ensure we have exactly 50 cards
    while (deck.length < 50) {
      const randomCard = friendCards[Math.floor(Math.random() * friendCards.length)]
      if (randomCard) {
        deck.push({
          ID: deckCardId++,
          deck_id: 0,
          card_no: randomCard.card_no,
          quantity: 1,
          card: { ...randomCard }
        })
      }
    }
    
    return shuffleDeck(deck.slice(0, 50))
  }
  
  function selectRandomCards(cards: Card[], count: number): Card[] {
    const shuffled = [...cards].sort(() => Math.random() - 0.5)
    return shuffled.slice(0, Math.min(count, shuffled.length))
  }
  
  function generateBasicCPUDeck(): DeckCard[] {
    // Fallback deck with proper structure but placeholder cards
    const deck: DeckCard[] = []
    let cardId = 0
    
    // Create 50 basic cards with proper image paths
    for (let i = 0; i < 50; i++) {
      const cost = i < 20 ? 1 : i < 35 ? 2 : 3
      const power = cost * 1000
      const color = ['赤', '青', '黄', '緑'][i % 4] as CardColor
      
      deck.push({
        ID: cardId++,
        deck_id: 0,
        card_no: `BASIC-${i}`,
        quantity: 1,
        card: {
          ID: cardId,
          card_no: `BASIC-${i}`,
          name: `基本カード${i}`,
          type: 'ふれんど' as CardType,
          color: color,
          rarity: 'C' as CardRarity,
          cost: cost,
          cost_red: color === '赤' ? cost : 0,
          cost_blue: color === '青' ? cost : 0,
          cost_yellow: color === '黄' ? cost : 0,
          cost_green: color === '緑' ? cost : 0,
          cost_colorless: 0,
          power: power,
          effect: '',
          flavor_text: '',
          image_url: '',
          local_image_path: 'card_images/placeholder.jpg',
          energy_icons: [],
          energy_value: 1,
          is_counter: false,
          is_main_counter: false,
          is_promo: false
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

  async function performBattle(attacker: 'player' | 'opponent', attackerIndex: number, targetIndex: number | 'player') {
    if (!gameState.value) return false
    
    const attackingPlayer = gameState.value.players[attacker]
    const defendingPlayer = gameState.value.players[attacker === 'player' ? 'opponent' : 'player']
    
    if (!attackingPlayer.friends[attackerIndex]) return false
    
    const attackerFriend = attackingPlayer.friends[attackerIndex]
    const attackerCard = attackerFriend.card.card
    if (!attackerCard || !attackerCard.power) return false
    
    // Tap the attacking card
    attackerFriend.tapped = true
    
    if (targetIndex === 'player') {
      // Direct attack on player - check if defender has untapped friends to block
      const defender = attacker === 'player' ? 'opponent' : 'player'
      const untappedFriends = defendingPlayer.friends
        .map((f, i) => ({ friend: f, index: i }))
        .filter(({ friend }) => !friend.tapped)
      
      if (untappedFriends.length > 0 && defender === 'player') {
        // Ask player if they want to block
        const blockerIndex = await getBlockingDecision(attacker, attackerIndex, attackerFriend.card)
        
        if (blockerIndex !== null) {
          // Player chose to block - redirect attack to the blocker
          targetIndex = blockerIndex
          
          // Tap the blocking friend
          defendingPlayer.friends[blockerIndex].tapped = true
        }
      } else if (untappedFriends.length > 0 && defender === 'opponent') {
        // CPU blocking decision - check if beneficial to block
        const blockerIndex = getCPUBlockingDecision(attackerFriend.card, untappedFriends)
        
        if (blockerIndex !== null) {
          // CPU chose to block - redirect attack to the blocker
          console.log('CPU blocks with friend at index:', blockerIndex)
          targetIndex = blockerIndex
          
          // Tap the blocking friend
          defendingPlayer.friends[blockerIndex].tapped = true
        }
      }
    }
    
    if (targetIndex === 'player') {
      // Direct attack on player - show animation
      battleAnimation.value = {
        show: true,
        attacker: attackerFriend.card,
        defender: null, // null indicates player target
        attackerDefeated: false,
        defenderDefeated: true
      }
      
      // Wait for animation to complete
      await new Promise(resolve => {
        const checkAnimation = setInterval(() => {
          if (!battleAnimation.value.show) {
            clearInterval(checkAnimation)
            resolve(true)
          }
        }, 100)
      })
      
      // Apply damage after animation
      if (defendingPlayer.deck.length > 0) {
        const topCard = defendingPlayer.deck.shift()!
        // Ensure reactivity by creating a new array - cards start face up
        defendingPlayer.negativeEnergy = [...defendingPlayer.negativeEnergy, { card: topCard, faceUp: true }]
        
        // Check victory condition
        if (defendingPlayer.negativeEnergy.length >= 7) {
          currentPhase.value = 'game_over'
          gameState.value.winner = attacker
        }
      }
    } else {
      // Attack on friend
      const targetFriend = defendingPlayer.friends[targetIndex]
      if (!targetFriend) return false
      
      const targetCard = targetFriend.card.card
      if (!targetCard || !targetCard.power) return false
      
      // Counter timing - defender can use counter cards
      const defender = attacker === 'player' ? 'opponent' : 'player'
      if (defender === 'player') {
        const counter = await getCounterDecision(
          attacker, attackerIndex, attackerFriend.card,
          defender, targetIndex, targetFriend.card
        )
        
        if (counter) {
          // Apply counter effects
          await applyCounterEffects(counter, attackerFriend, targetFriend)
        }
      }
      
      // Determine battle outcome after counter effects
      const attackerPower = attackerFriend.card.card.power || 0
      const defenderPower = targetFriend.card.card.power || 0
      const attackerDefeated = defenderPower >= attackerPower
      const defenderDefeated = attackerPower >= defenderPower
      
      // Show battle animation
      battleAnimation.value = {
        show: true,
        attacker: attackerFriend.card,
        defender: targetFriend.card,
        attackerDefeated,
        defenderDefeated
      }
      
      // Wait for animation to complete
      await new Promise(resolve => {
        const checkAnimation = setInterval(() => {
          if (!battleAnimation.value.show) {
            clearInterval(checkAnimation)
            resolve(true)
          }
        }, 100)
      })
      
      // Apply battle results after animation
      if (defenderDefeated) {
        // Target is defeated - goes to trash
        const defeated = defendingPlayer.friends.splice(targetIndex, 1)[0]
        defendingPlayer.graveyard.push(defeated.card)
      }
      
      if (attackerDefeated) {
        // Attacker is defeated - goes to trash
        const defeated = attackingPlayer.friends.splice(attackerIndex, 1)[0]
        attackingPlayer.graveyard.push(defeated.card)
      }
    }
    
    return true
  }
  
  function hideBattleAnimation() {
    battleAnimation.value.show = false
  }
  
  function selectBattleTarget(attackerIndex: number) {
    if (!gameState.value || !isPlayerTurn.value || currentPhase.value !== 'main') return
    
    selectedAttacker.value = attackerIndex
    
    // In this game, friends can only attack the player directly
    // The opponent can choose to block with their friends
    availableTargets.value = [-1] // -1 represents direct player attack
  }
  
  async function executeBattle(targetIndex: number | 'player') {
    if (selectedAttacker.value === null) return
    
    await performBattle('player', selectedAttacker.value, targetIndex)
    selectedAttacker.value = null
    availableTargets.value = []
  }
  
  function setCPUDifficulty(difficulty: CPUDifficulty) {
    cpuDifficulty.value = difficulty
    cpuThinkingTime.value = difficulty === 'easy' ? 1500 : difficulty === 'normal' ? 2500 : 3500
  }

  function mulliganHand(player: 'player' | 'opponent') {
    if (!gameState.value) return
    
    const playerState = gameState.value.players[player]
    
    // Return all cards to deck
    playerState.deck.push(...playerState.hand)
    playerState.hand = []
    
    // Shuffle deck
    playerState.deck = shuffleDeck(playerState.deck)
    
    // Draw new hand
    drawCards(player, 5)
  }

  async function getBlockingDecision(attacker: 'player' | 'opponent', attackerIndex: number, attackerCard: DeckCard): Promise<number | null> {
    if (!gameState.value) return null
    
    const defendingPlayer = gameState.value.players[attacker === 'player' ? 'opponent' : 'player']
    const availableBlockers = defendingPlayer.friends
      .map((f, i) => ({ friend: f, index: i }))
      .filter(({ friend }) => !friend.tapped)
      .map(({ index }) => index)
    
    if (availableBlockers.length === 0) return null
    
    return new Promise<number | null>((resolve) => {
      blockingDecision.value = {
        show: true,
        attacker: { player: attacker, index: attackerIndex, card: attackerCard },
        availableBlockers,
        resolve
      }
    })
  }
  
  async function getCounterDecision(
    attacker: 'player' | 'opponent', 
    attackerIndex: number, 
    attackerCard: DeckCard,
    defender: 'player' | 'opponent', 
    blockerIndex: number, 
    blockerCard: DeckCard
  ): Promise<DeckCard | null> {
    if (!gameState.value) return null
    
    return new Promise<DeckCard | null>((resolve) => {
      counterDecision.value = {
        show: true,
        attacker: { player: attacker, index: attackerIndex, card: attackerCard },
        blocker: { player: defender, index: blockerIndex, card: blockerCard },
        resolve
      }
    })
  }
  
  async function applyCounterEffects(counter: DeckCard, attacker: FriendCardState, defender: FriendCardState): Promise<void> {
    if (!gameState.value) return
    
    // Remove counter card from hand
    const playerState = gameState.value.players.player
    const handIndex = playerState.hand.findIndex(h => h.ID === counter.ID)
    if (handIndex !== -1) {
      playerState.hand.splice(handIndex, 1)
      // Counter cards go to trash after use
      playerState.graveyard.push(counter)
    }
    
    // Apply counter effects based on card
    // This is a simple implementation - in practice, you'd check the specific counter effects
    if (counter.card.description?.includes('パワー+')) {
      // Example: Add power to defending friend
      const powerBoost = 1000 // Could parse this from card description
      // Note: This modifies the actual card object temporarily
      if (defender.card.card.power) {
        defender.card.card.power += powerBoost
      }
    }
  }
  
  function getCPUBlockingDecision(
    attackerCard: DeckCard, 
    untappedFriends: { friend: FriendCardState, index: number }[]
  ): number | null {
    const attackerPower = attackerCard.card.power || 0
    
    console.log('CPU blocking decision - Attacker:', attackerCard.card.name, 'Power:', attackerPower)
    console.log('Available blockers:', untappedFriends.map(f => ({
      name: f.friend.card.card.name,
      power: f.friend.card.card.power,
      index: f.index
    })))
    
    // Analyze all possible blocking options
    const blockingOptions = untappedFriends.map(({ friend, index }) => {
      const blockerPower = friend.card.card.power || 0
      
      // Calculate outcome
      const attackerDefeated = blockerPower >= attackerPower
      const blockerDefeated = attackerPower >= blockerPower
      const isMutualDestruction = attackerDefeated && blockerDefeated
      
      return {
        index,
        blockerPower,
        blockerName: friend.card.card.name,
        attackerDefeated,
        blockerDefeated,
        isMutualDestruction,
        powerDifference: blockerPower - attackerPower
      }
    })
    
    console.log('Blocking analysis:', blockingOptions)
    
    // Strategy based on CPU difficulty
    switch (cpuDifficulty.value) {
      case 'easy':
        // Easy: 40% chance to block, sometimes seeks mutual destruction
        if (Math.random() > 0.6) {
          // 30% chance to specifically look for mutual destruction
          if (Math.random() > 0.7) {
            const mutualOptions = blockingOptions.filter(opt => opt.isMutualDestruction)
            if (mutualOptions.length > 0) {
              console.log('CPU (easy) randomly chooses mutual destruction')
              return mutualOptions[Math.floor(Math.random() * mutualOptions.length)].index
            }
          }
          
          const destructionOptions = blockingOptions.filter(opt => opt.attackerDefeated)
          if (destructionOptions.length > 0 && Math.random() > 0.5) {
            console.log('CPU (easy) randomly chooses to destroy attacker')
            return destructionOptions[Math.floor(Math.random() * destructionOptions.length)].index
          }
          console.log('CPU (easy) blocks randomly')
          return blockingOptions[Math.floor(Math.random() * blockingOptions.length)].index
        }
        return null
        
      case 'normal':
        // Normal: Actively seek mutual destruction opportunities
        
        // 1. First, check for any mutual destruction opportunities (prioritize these)
        const mutualDestructionOptions = blockingOptions.filter(opt => opt.isMutualDestruction)
        console.log('Mutual destruction options:', mutualDestructionOptions)
        if (mutualDestructionOptions.length > 0) {
          console.log('CPU actively chooses mutual destruction block')
          // Choose the mutual destruction with highest power for maximum impact
          mutualDestructionOptions.sort((a, b) => b.blockerPower - a.blockerPower)
          console.log('Selected blocker for mutual destruction:', mutualDestructionOptions[0])
          return mutualDestructionOptions[0].index
        }
        
        // 2. Then look for opportunities to destroy the attacker
        const destructionOptions = blockingOptions.filter(opt => opt.attackerDefeated)
        if (destructionOptions.length > 0) {
          const favorableDestruction = destructionOptions.filter(opt => !opt.blockerDefeated)
          if (favorableDestruction.length > 0) {
            console.log('CPU chooses favorable destruction block')
            return favorableDestruction[0].index
          }
          
          // Even if blocker dies, destroying attacker is often worth it
          console.log('CPU chooses destruction block (trading for attacker)')
          return destructionOptions[0].index
        }
        
        // 3. If can't destroy attacker, 60% chance to block anyway to prevent damage
        if (Math.random() > 0.4) {
          console.log('CPU blocks to prevent damage')
          return blockingOptions[0].index
        }
        return null
        
      case 'hard':
        // Hard: Aggressively seek mutual destruction and optimal trades
        
        // 1. First check for mutual destruction - in hard mode, always take it
        const hardMutualOptions = blockingOptions.filter(opt => opt.isMutualDestruction)
        if (hardMutualOptions.length > 0) {
          console.log('CPU (hard) aggressively chooses mutual destruction')
          // Choose the mutual destruction with highest power for maximum board impact
          hardMutualOptions.sort((a, b) => b.blockerPower - a.blockerPower)
          return hardMutualOptions[0].index
        }
        
        // 2. Then look for any way to destroy the attacker
        const hardDestructionOptions = blockingOptions.filter(opt => opt.attackerDefeated)
        if (hardDestructionOptions.length > 0) {
          // Among destruction options, prefer favorable trades
          const favorableHardOptions = hardDestructionOptions.filter(opt => !opt.blockerDefeated)
          if (favorableHardOptions.length > 0) {
            console.log('CPU (hard) chooses optimal destruction block')
            // Choose the blocker with lowest power that still wins (resource efficiency)
            favorableHardOptions.sort((a, b) => a.blockerPower - b.blockerPower)
            return favorableHardOptions[0].index
          }
          
          // Even unfavorable trades are worth it in hard mode
          console.log('CPU (hard) trades to destroy attacker')
          // Choose the weakest blocker for the trade
          hardDestructionOptions.sort((a, b) => a.blockerPower - b.blockerPower)
          return hardDestructionOptions[0].index
        }
        
        // 2. If can't destroy attacker, still block frequently to prevent damage
        if (Math.random() > 0.2) { // 80% chance to block on hard
          console.log('CPU (hard) blocks to prevent damage')
          return blockingOptions[0].index
        }
        return null
        
      default:
        return null
    }
  }
  
  function resolveBlocking(blockerIndex: number | null) {
    if (blockingDecision.value.resolve) {
      blockingDecision.value.resolve(blockerIndex)
      blockingDecision.value = {
        show: false,
        attacker: { player: 'player', index: -1, card: null },
        availableBlockers: [],
        resolve: null
      }
    }
  }
  
  function flipNegativeEnergyCard(index: number) {
    if (!gameState.value) return
    
    const playerState = gameState.value.players.player
    if (index < 0 || index >= playerState.negativeEnergy.length) return
    
    const negativeCard = playerState.negativeEnergy[index]
    if (negativeCard.faceUp) {
      // Flip face-up card to face-down (to use as energy payment)
      negativeCard.faceUp = false
    }
    // Once face-down, cards cannot be flipped back up
  }
  
  function resolveCounter(counter: DeckCard | null) {
    if (counterDecision.value.resolve) {
      counterDecision.value.resolve(counter)
      // Reset counter decision state
      counterDecision.value = {
        show: false,
        attacker: { player: 'player', index: -1, card: null },
        blocker: null,
        resolve: null
      }
    }
  }
  
  async function getEnergyCostSelection(cardToPay: Card): Promise<any> {
    return new Promise((resolve) => {
      energyCostSelection.value = {
        show: true,
        cardToPay,
        resolve
      }
    })
  }
  
  function cancelEnergyCostSelection() {
    energyCostSelection.value = {
      show: false,
      cardToPay: null,
      resolve: null
    }
  }
  
  function confirmEnergyCostSelection(selection: any) {
    if (energyCostSelection.value.resolve) {
      energyCostSelection.value.resolve(selection)
      energyCostSelection.value = {
        show: false,
        cardToPay: null,
        resolve: null
      }
    }
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
    battleAnimation,
    blockingDecision,
    counterDecision,
    energyCostSelection,
    energyPlayedThisTurn,
    initializeGame,
    startGame,
    drawCards,
    playCard,
    nextPhase,
    endTurn,
    resetGame,
    performBattle,
    selectBattleTarget,
    executeBattle,
    setCPUDifficulty,
    mulliganHand,
    hideBattleAnimation,
    resolveBlocking,
    flipNegativeEnergyCard,
    resolveCounter,
    getEnergyCostSelection,
    cancelEnergyCostSelection,
    confirmEnergyCostSelection
  }
})