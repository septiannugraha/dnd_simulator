package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	Email     string             `bson:"email" json:"email" binding:"required,email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Character struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	CampaignID   primitive.ObjectID `bson:"campaign_id,omitempty" json:"campaign_id,omitempty"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Race         string             `bson:"race" json:"race" binding:"required"`
	Class        string             `bson:"class" json:"class" binding:"required"`
	Level        int                `bson:"level" json:"level"`
	HitPoints    int                `bson:"hit_points" json:"hit_points"`
	ArmorClass   int                `bson:"armor_class" json:"armor_class"`
	Attributes   Attributes         `bson:"attributes" json:"attributes"`
	Equipment    []string           `bson:"equipment" json:"equipment"`
	Spells       []string           `bson:"spells" json:"spells"`
	Background   string             `bson:"background" json:"background"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type Attributes struct {
	Strength     int `bson:"strength" json:"strength"`
	Dexterity    int `bson:"dexterity" json:"dexterity"`
	Constitution int `bson:"constitution" json:"constitution"`
	Intelligence int `bson:"intelligence" json:"intelligence"`
	Wisdom       int `bson:"wisdom" json:"wisdom"`
	Charisma     int `bson:"charisma" json:"charisma"`
}

type Campaign struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name" binding:"required"`
	Description string               `bson:"description" json:"description"`
	DMID        primitive.ObjectID   `bson:"dm_id" json:"dm_id"`
	PlayerIDs   []primitive.ObjectID `bson:"player_ids" json:"player_ids"`
	WorldInfo   string               `bson:"world_info" json:"world_info"`
	Settings    CampaignSettings     `bson:"settings" json:"settings"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updated_at"`
}

type CampaignSettings struct {
	IsPublic     bool `bson:"is_public" json:"is_public"`
	MaxPlayers   int  `bson:"max_players" json:"max_players"`
	AllowGuests  bool `bson:"allow_guests" json:"allow_guests"`
}

type GameSession struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	CampaignID   primitive.ObjectID   `bson:"campaign_id" json:"campaign_id"`
	Name         string               `bson:"name" json:"name"`
	Status       string               `bson:"status" json:"status"` // "waiting", "active", "paused", "ended"
	CurrentTurn  int                  `bson:"current_turn" json:"current_turn"`
	TurnOrder    []primitive.ObjectID `bson:"turn_order" json:"turn_order"`
	PlayerIDs    []primitive.ObjectID `bson:"player_ids" json:"player_ids"`
	CurrentScene string               `bson:"current_scene" json:"current_scene"`
	SessionLog   []SessionEvent       `bson:"session_log" json:"session_log"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
}

type SessionEvent struct {
	Type      string    `bson:"type" json:"type"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Content   string    `bson:"content" json:"content"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}