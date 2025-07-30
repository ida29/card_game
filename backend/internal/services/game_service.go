package services

import (
	"fmt"
	"math/rand"
	"github.com/your-org/card_game/backend/internal/game"
	"github.com/your-org/card_game/backend/internal/models"
	"gorm.io/gorm"
)

type GameService struct {
	db           *gorm.DB
	cardService  *CardService
	eventHandler map[string]*game.EventHandler // game_id -> event handler
}

func NewGameService(db *gorm.DB, cardService *CardService) *GameService {
	return &GameService{
		db:           db,
		cardService:  cardService,
		eventHandler: make(map[string]*game.EventHandler),
	}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(player1ID, player2ID uint, deck1ID, deck2ID uint) (*models.Game, error) {
	gameID := generateGameID()
	
	// Initialize game state
	gameState := &models.GameState{
		Player1State: models.PlayerState{
			Deck:           []string{}, // Will be populated from deck
			Hand:           []string{},
			BattleArea:     make(map[string]models.Friend),
			EnergyArea:     []models.EnergyCard{},
			NegativeEnergy: []string{},
			Trash:          []string{},
		},
		Player2State: models.PlayerState{
			Deck:           []string{},
			Hand:           []string{},
			BattleArea:     make(map[string]models.Friend),
			EnergyArea:     []models.EnergyCard{},
			NegativeEnergy: []string{},
			Trash:          []string{},
		},
	}
	
	// Create game
	newGame := &models.Game{
		GameID:       gameID,
		Player1ID:    player1ID,
		Player2ID:    player2ID,
		CurrentTurn:  1,
		CurrentPhase: models.PhaseStart,
		ActivePlayer: 1,
		Status:       models.StatusWaiting,
		GameState:    gameState,
	}
	
	if err := s.db.Create(newGame).Error; err != nil {
		return nil, err
	}
	
	// Create event handler for this game
	s.eventHandler[gameID] = game.NewEventHandler(newGame)
	
	return newGame, nil
}

// PlayCard plays a card from hand
func (s *GameService) PlayCard(gameID string, playerID uint, cardNo string, position string, targets []string) error {
	// Get game
	var gameModel models.Game
	if err := s.db.Where("game_id = ?", gameID).First(&gameModel).Error; err != nil {
		return err
	}
	
	// Get event handler
	handler, exists := s.eventHandler[gameID]
	if !exists {
		handler = game.NewEventHandler(&gameModel)
		s.eventHandler[gameID] = handler
	}
	
	// Determine which player
	var player int
	if gameModel.Player1ID == playerID {
		player = 1
	} else if gameModel.Player2ID == playerID {
		player = 2
	} else {
		return fmt.Errorf("player not in this game")
	}
	
	// Get card details
	card, err := s.cardService.GetCardByNo(cardNo)
	if err != nil {
		return err
	}
	
	// Play the card based on type
	switch card.Type {
	case models.CardTypeFriend:
		if err := s.playFriend(&gameModel, player, cardNo, position); err != nil {
			return err
		}
		
		// Trigger friend played event
		event := game.GameEvent{
			Type:   game.EventFriendPlayed,
			Player: player,
			CardNo: cardNo,
			Phase:  gameModel.CurrentPhase,
		}
		if err := handler.TriggerEvent(event); err != nil {
			return err
		}
		
	case models.CardTypeSupport:
		if err := s.playSupport(&gameModel, player, cardNo, targets); err != nil {
			return err
		}
		
		// Trigger support played event
		event := game.GameEvent{
			Type:   game.EventSupportPlayed,
			Player: player,
			CardNo: cardNo,
			Phase:  gameModel.CurrentPhase,
			Data: map[string]interface{}{
				"targets": targets,
			},
		}
		if err := handler.TriggerEvent(event); err != nil {
			return err
		}
		
	case models.CardTypeField:
		if err := s.playField(&gameModel, player, cardNo); err != nil {
			return err
		}
		
		// Trigger field played event
		event := game.GameEvent{
			Type:   game.EventFieldPlayed,
			Player: player,
			CardNo: cardNo,
			Phase:  gameModel.CurrentPhase,
		}
		if err := handler.TriggerEvent(event); err != nil {
			return err
		}
	}
	
	// Save game state
	return s.db.Save(&gameModel).Error
}

// Attack performs an attack with a friend
func (s *GameService) Attack(gameID string, playerID uint, attackerPos string, targetPos string) error {
	// Get game
	var gameModel models.Game
	if err := s.db.Where("game_id = ?", gameID).First(&gameModel).Error; err != nil {
		return err
	}
	
	// Get event handler
	handler, exists := s.eventHandler[gameID]
	if !exists {
		handler = game.NewEventHandler(&gameModel)
		s.eventHandler[gameID] = handler
	}
	
	// Determine which player
	var player int
	if gameModel.Player1ID == playerID {
		player = 1
	} else if gameModel.Player2ID == playerID {
		player = 2
	} else {
		return fmt.Errorf("player not in this game")
	}
	
	// Get attacker
	var playerState *models.PlayerState
	if player == 1 {
		playerState = &gameModel.GameState.Player1State
	} else {
		playerState = &gameModel.GameState.Player2State
	}
	
	attacker, exists := playerState.BattleArea[attackerPos]
	if !exists {
		return fmt.Errorf("no friend at position %s", attackerPos)
	}
	
	// Trigger attack event
	event := game.GameEvent{
		Type:   game.EventFriendAttacks,
		Player: player,
		CardNo: attacker.CardNo,
		Target: targetPos,
		Phase:  gameModel.CurrentPhase,
	}
	if err := handler.TriggerEvent(event); err != nil {
		return err
	}
	
	// Perform attack logic
	// TODO: Implement attack resolution
	
	// Save game state
	return s.db.Save(&gameModel).Error
}

// Block declares a blocker for an attack
func (s *GameService) Block(gameID string, playerID uint, blockerPos string, attackerPos string) error {
	// Get game
	var gameModel models.Game
	if err := s.db.Where("game_id = ?", gameID).First(&gameModel).Error; err != nil {
		return err
	}
	
	// Get event handler
	handler, exists := s.eventHandler[gameID]
	if !exists {
		handler = game.NewEventHandler(&gameModel)
		s.eventHandler[gameID] = handler
	}
	
	// Determine which player
	var player int
	if gameModel.Player1ID == playerID {
		player = 1
	} else if gameModel.Player2ID == playerID {
		player = 2
	} else {
		return fmt.Errorf("player not in this game")
	}
	
	// Get blocker
	var playerState *models.PlayerState
	if player == 1 {
		playerState = &gameModel.GameState.Player1State
	} else {
		playerState = &gameModel.GameState.Player2State
	}
	
	blocker, exists := playerState.BattleArea[blockerPos]
	if !exists {
		return fmt.Errorf("no friend at position %s", blockerPos)
	}
	
	// Trigger block event
	event := game.GameEvent{
		Type:   game.EventFriendBlocks,
		Player: player,
		CardNo: blocker.CardNo,
		Target: attackerPos,
		Phase:  gameModel.CurrentPhase,
	}
	if err := handler.TriggerEvent(event); err != nil {
		return err
	}
	
	// Perform block logic
	// TODO: Implement block resolution
	
	// Save game state
	return s.db.Save(&gameModel).Error
}

// ChangePhase moves to the next phase
func (s *GameService) ChangePhase(gameID string, playerID uint) error {
	// Get game
	var gameModel models.Game
	if err := s.db.Where("game_id = ?", gameID).First(&gameModel).Error; err != nil {
		return err
	}
	
	// Get event handler
	handler, exists := s.eventHandler[gameID]
	if !exists {
		handler = game.NewEventHandler(&gameModel)
		s.eventHandler[gameID] = handler
	}
	
	// Trigger phase end event
	endEvent := game.GameEvent{
		Type:   game.EventPhaseEnd,
		Player: gameModel.ActivePlayer,
		Phase:  gameModel.CurrentPhase,
	}
	if err := handler.TriggerEvent(endEvent); err != nil {
		return err
	}
	
	// Move to next phase
	switch gameModel.CurrentPhase {
	case models.PhaseStart:
		gameModel.CurrentPhase = models.PhaseDraw
	case models.PhaseDraw:
		gameModel.CurrentPhase = models.PhaseEnergy
	case models.PhaseEnergy:
		gameModel.CurrentPhase = models.PhaseMain
	case models.PhaseMain:
		gameModel.CurrentPhase = models.PhaseEnd
	case models.PhaseEnd:
		// End turn
		gameModel.CurrentPhase = models.PhaseStart
		gameModel.CurrentTurn++
		if gameModel.ActivePlayer == 1 {
			gameModel.ActivePlayer = 2
		} else {
			gameModel.ActivePlayer = 1
		}
		
		// Trigger turn end event
		turnEndEvent := game.GameEvent{
			Type:   game.EventTurnEnd,
			Player: gameModel.ActivePlayer,
			Phase:  gameModel.CurrentPhase,
		}
		handler.TriggerEvent(turnEndEvent)
		
		// Trigger turn start event
		turnStartEvent := game.GameEvent{
			Type:   game.EventTurnStart,
			Player: gameModel.ActivePlayer,
			Phase:  gameModel.CurrentPhase,
		}
		handler.TriggerEvent(turnStartEvent)
	}
	
	// Trigger phase start event
	startEvent := game.GameEvent{
		Type:   game.EventPhaseStart,
		Player: gameModel.ActivePlayer,
		Phase:  gameModel.CurrentPhase,
	}
	if err := handler.TriggerEvent(startEvent); err != nil {
		return err
	}
	
	// Save game state
	return s.db.Save(&gameModel).Error
}

// Helper methods

func (s *GameService) playFriend(game *models.Game, player int, cardNo string, position string) error {
	var playerState *models.PlayerState
	if player == 1 {
		playerState = &game.GameState.Player1State
	} else {
		playerState = &game.GameState.Player2State
	}
	
	// Remove from hand
	for i, card := range playerState.Hand {
		if card == cardNo {
			playerState.Hand = append(playerState.Hand[:i], playerState.Hand[i+1:]...)
			break
		}
	}
	
	// Get card details
	card, err := s.cardService.GetCardByNo(cardNo)
	if err != nil {
		return err
	}
	
	// Add to battle area
	friend := models.Friend{
		CardNo:     cardNo,
		Power:      card.Power,
		IsRest:     false, // Friends enter active by default
		TurnPlayed: game.CurrentTurn,
	}
	
	// Check for effects that modify entry state
	// TODO: Check for "enter rested" effects
	
	playerState.BattleArea[position] = friend
	
	return nil
}

func (s *GameService) playSupport(game *models.Game, player int, cardNo string, targets []string) error {
	var playerState *models.PlayerState
	if player == 1 {
		playerState = &game.GameState.Player1State
	} else {
		playerState = &game.GameState.Player2State
	}
	
	// Remove from hand
	for i, card := range playerState.Hand {
		if card == cardNo {
			playerState.Hand = append(playerState.Hand[:i], playerState.Hand[i+1:]...)
			break
		}
	}
	
	// Support cards go to trash after use
	playerState.Trash = append(playerState.Trash, cardNo)
	
	return nil
}

func (s *GameService) playField(game *models.Game, player int, cardNo string) error {
	var playerState *models.PlayerState
	if player == 1 {
		playerState = &game.GameState.Player1State
	} else {
		playerState = &game.GameState.Player2State
	}
	
	// Remove from hand
	for i, card := range playerState.Hand {
		if card == cardNo {
			playerState.Hand = append(playerState.Hand[:i], playerState.Hand[i+1:]...)
			break
		}
	}
	
	// If there's already a field card, it goes to trash
	if playerState.FieldCard != nil {
		playerState.Trash = append(playerState.Trash, *playerState.FieldCard)
	}
	
	// Set new field card
	playerState.FieldCard = &cardNo
	
	return nil
}

func generateGameID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}