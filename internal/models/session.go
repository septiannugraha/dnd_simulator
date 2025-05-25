package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Game Session model for structured gameplay
type GameSession struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CampaignID  primitive.ObjectID `bson:"campaign_id" json:"campaign_id" binding:"required"`
	Name        string             `bson:"name" json:"name" binding:"required,min=2,max=100"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	
	// Session State
	Status      SessionStatus      `bson:"status" json:"status"`
	Scene       string             `bson:"scene,omitempty" json:"scene,omitempty"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty"`
	
	// Turn Management
	TurnOrder   []TurnEntry        `bson:"turn_order" json:"turn_order"`
	CurrentTurn int                `bson:"current_turn" json:"current_turn"`
	Round       int                `bson:"round" json:"round"`
	
	// Participants
	Players     []SessionPlayer    `bson:"players" json:"players"`
	DMUserID    primitive.ObjectID `bson:"dm_user_id" json:"dm_user_id"`
	
	// Session Data
	ChatHistory []SessionChatMsg   `bson:"chat_history,omitempty" json:"chat_history,omitempty"`
	
	// Metadata
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	StartedAt   *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`
	EndedAt     *time.Time         `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
}

type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"   // Created but not started
	SessionStatusActive    SessionStatus = "active"    // Currently running
	SessionStatusPaused    SessionStatus = "paused"    // Temporarily paused
	SessionStatusCompleted SessionStatus = "completed" // Finished
	SessionStatusCancelled SessionStatus = "cancelled" // Cancelled
)

type SessionPlayer struct {
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CharacterID primitive.ObjectID `bson:"character_id" json:"character_id"`
	Username    string             `bson:"username" json:"username"`
	CharName    string             `bson:"character_name" json:"character_name"`
	IsConnected bool               `bson:"is_connected" json:"is_connected"`
	JoinedAt    time.Time          `bson:"joined_at" json:"joined_at"`
}

type TurnEntry struct {
	Type         TurnType           `bson:"type" json:"type"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	CharacterID  primitive.ObjectID `bson:"character_id,omitempty" json:"character_id,omitempty"`
	Initiative   int                `bson:"initiative" json:"initiative"`
	Name         string             `bson:"name" json:"name"`
	HasActed     bool               `bson:"has_acted" json:"has_acted"`
}

type TurnType string

const (
	TurnTypePlayer TurnType = "player"
	TurnTypeNPC    TurnType = "npc"
	TurnTypeEvent  TurnType = "event"
)

// Session-specific chat message (for persistence)
type SessionChatMsg struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Message   string             `bson:"message" json:"message"`
	Type      string             `bson:"type" json:"type"` // "ic", "ooc", "roll", "system"
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

// Request/Response DTOs for session operations
type CreateSessionRequest struct {
	CampaignID  primitive.ObjectID `json:"campaign_id" binding:"required"`
	Name        string             `json:"name" binding:"required,min=2,max=100"`
	Description string             `json:"description,omitempty"`
}

type JoinSessionRequest struct {
	CharacterID primitive.ObjectID `json:"character_id" binding:"required"`
}

type SetInitiativeRequest struct {
	CharacterID primitive.ObjectID `json:"character_id" binding:"required"`
	Initiative  int                `json:"initiative" binding:"required,min=1,max=30"`
}

type AdvanceTurnRequest struct {
	Force bool `json:"force,omitempty"` // Force advance even if player hasn't acted
}

type SessionResponse struct {
	GameSession
	Campaign *Campaign `json:"campaign,omitempty"`
}