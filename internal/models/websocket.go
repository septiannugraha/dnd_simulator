package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WebSocket message types
const (
	// Connection management
	MessageTypeJoinSession    = "join_session"
	MessageTypeLeaveSession   = "leave_session"
	MessageTypePlayerJoined   = "player_joined"
	MessageTypePlayerLeft     = "player_left"
	
	// Chat messages
	MessageTypeChat           = "chat"
	MessageTypeChatIC         = "chat_ic"  // In-character
	MessageTypeChatOOC        = "chat_ooc" // Out-of-character
	
	// Dice rolling
	MessageTypeDiceRoll       = "dice_roll"
	MessageTypeDiceResult     = "dice_result"
	
	// Character updates
	MessageTypeCharacterUpdate = "character_update"
	MessageTypeCharacterSync   = "character_sync"
	
	// Game state
	MessageTypeGameState      = "game_state"
	MessageTypeTurnOrder      = "turn_order"
	MessageTypeTurnAdvance    = "turn_advance"
	
	// System messages
	MessageTypeError          = "error"
	MessageTypeSuccess        = "success"
	MessageTypeNotification   = "notification"
)

// WebSocket message structure
type WSMessage struct {
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    primitive.ObjectID     `json:"user_id,omitempty"`
	Username  string                 `json:"username,omitempty"`
	SessionID primitive.ObjectID     `json:"session_id,omitempty"`
}

// Specific message data structures
type ChatMessage struct {
	Content     string             `json:"content"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty"`
	IsIC        bool               `json:"is_ic"` // In-character vs out-of-character
}

type DiceRoll struct {
	Dice        string             `json:"dice"`        // e.g., "1d20+5", "3d6"
	Result      []int              `json:"result"`      // Individual die results
	Total       int                `json:"total"`       // Final total
	Modifier    int                `json:"modifier"`    // Applied modifier
	Purpose     string             `json:"purpose"`     // e.g., "Attack roll", "Saving throw"
	CharacterID primitive.ObjectID `json:"character_id,omitempty"`
}

type CharacterUpdate struct {
	CharacterID primitive.ObjectID     `json:"character_id"`
	Field       string                 `json:"field"`
	Value       interface{}            `json:"value"`
	OldValue    interface{}            `json:"old_value,omitempty"`
}

type GameState struct {
	SessionID     primitive.ObjectID   `json:"session_id"`
	CurrentTurn   int                  `json:"current_turn"`
	TurnOrder     []primitive.ObjectID `json:"turn_order"`
	Status        string               `json:"status"`
	CurrentScene  string               `json:"current_scene"`
}

type PlayerJoinedData struct {
	UserID      primitive.ObjectID `json:"user_id"`
	Username    string             `json:"username"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty"`
	Character   *Character         `json:"character,omitempty"`
}

type PlayerLeftData struct {
	UserID   primitive.ObjectID `json:"user_id"`
	Username string             `json:"username"`
}

// Connection represents an active WebSocket connection
type Connection struct {
	UserID      primitive.ObjectID `json:"user_id"`
	Username    string             `json:"username"`
	SessionID   primitive.ObjectID `json:"session_id"`
	CharacterID primitive.ObjectID `json:"character_id,omitempty"`
	ConnectedAt time.Time          `json:"connected_at"`
}