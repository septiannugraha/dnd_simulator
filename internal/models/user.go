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

