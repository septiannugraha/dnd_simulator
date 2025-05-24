package services

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/database"
	"dnd-simulator/internal/models"
)

type CampaignService struct {
	db *database.DB
}

func NewCampaignService(db *database.DB) *CampaignService {
	return &CampaignService{db: db}
}

func (s *CampaignService) CreateCampaign(dmID primitive.ObjectID, name, description, worldInfo string, settings models.CampaignSettings) (*models.Campaign, error) {
	collection := s.db.GetCollection("campaigns")

	campaign := &models.Campaign{
		Name:        name,
		Description: description,
		DMID:        dmID,
		PlayerIDs:   []primitive.ObjectID{},
		WorldInfo:   worldInfo,
		Settings:    settings,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(context.Background(), campaign)
	if err != nil {
		return nil, err
	}

	campaign.ID = result.InsertedID.(primitive.ObjectID)
	return campaign, nil
}

func (s *CampaignService) GetCampaignByID(campaignID primitive.ObjectID) (*models.Campaign, error) {
	collection := s.db.GetCollection("campaigns")

	var campaign models.Campaign
	err := collection.FindOne(context.Background(), bson.M{"_id": campaignID}).Decode(&campaign)
	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (s *CampaignService) GetUserCampaigns(userID primitive.ObjectID) ([]models.Campaign, error) {
	collection := s.db.GetCollection("campaigns")

	// Find campaigns where user is either DM or player
	filter := bson.M{
		"$or": []bson.M{
			{"dm_id": userID},
			{"player_ids": bson.M{"$in": []primitive.ObjectID{userID}}},
		},
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campaigns []models.Campaign
	if err = cursor.All(context.Background(), &campaigns); err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (s *CampaignService) UpdateCampaign(campaignID, dmID primitive.ObjectID, name, description, worldInfo string, settings models.CampaignSettings) (*models.Campaign, error) {
	collection := s.db.GetCollection("campaigns")

	// First check if user is the DM
	var campaign models.Campaign
	err := collection.FindOne(context.Background(), bson.M{"_id": campaignID}).Decode(&campaign)
	if err != nil {
		return nil, err
	}

	if campaign.DMID != dmID {
		return nil, errors.New("only the DM can update this campaign")
	}

	update := bson.M{
		"$set": bson.M{
			"name":        name,
			"description": description,
			"world_info":  worldInfo,
			"settings":    settings,
			"updated_at":  time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": campaignID}, update)
	if err != nil {
		return nil, err
	}

	return s.GetCampaignByID(campaignID)
}

func (s *CampaignService) DeleteCampaign(campaignID, dmID primitive.ObjectID) error {
	collection := s.db.GetCollection("campaigns")

	// First check if user is the DM
	var campaign models.Campaign
	err := collection.FindOne(context.Background(), bson.M{"_id": campaignID}).Decode(&campaign)
	if err != nil {
		return err
	}

	if campaign.DMID != dmID {
		return errors.New("only the DM can delete this campaign")
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": campaignID})
	return err
}

func (s *CampaignService) JoinCampaign(campaignID, userID primitive.ObjectID) error {
	collection := s.db.GetCollection("campaigns")

	// Check if campaign exists and get current player count
	var campaign models.Campaign
	err := collection.FindOne(context.Background(), bson.M{"_id": campaignID}).Decode(&campaign)
	if err != nil {
		return err
	}

	// Check if user is already a player or the DM
	if campaign.DMID == userID {
		return errors.New("DM cannot join their own campaign as a player")
	}

	for _, playerID := range campaign.PlayerIDs {
		if playerID == userID {
			return errors.New("user is already a player in this campaign")
		}
	}

	// Check max players limit
	if len(campaign.PlayerIDs) >= campaign.Settings.MaxPlayers {
		return errors.New("campaign is full")
	}

	// Add user to player list
	update := bson.M{
		"$push": bson.M{"player_ids": userID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": campaignID}, update)
	return err
}

func (s *CampaignService) LeaveCampaign(campaignID, userID primitive.ObjectID) error {
	collection := s.db.GetCollection("campaigns")

	// Remove user from player list
	update := bson.M{
		"$pull": bson.M{"player_ids": userID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": campaignID}, update)
	return err
}

func (s *CampaignService) GetPublicCampaigns() ([]models.Campaign, error) {
	collection := s.db.GetCollection("campaigns")

	filter := bson.M{"settings.is_public": true}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var campaigns []models.Campaign
	if err = cursor.All(context.Background(), &campaigns); err != nil {
		return nil, err
	}

	return campaigns, nil
}