package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dnd-simulator/internal/database"
	"dnd-simulator/internal/models"
)

// EventService handles game event tracking and retrieval
type EventService struct {
	db *database.DB
}

// NewEventService creates a new event service instance
func NewEventService(db *database.DB) *EventService {
	return &EventService{
		db: db,
	}
}

// StoreEvent stores a game event in the database
func (s *EventService) StoreEvent(ctx context.Context, event *models.GameEvent) error {
	if event.ID.IsZero() {
		event.ID = primitive.NewObjectID()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	_, err := s.db.GetCollection("game_events").InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	return nil
}

// StorePlayerAction stores a player action as a game event
func (s *EventService) StorePlayerAction(ctx context.Context, sessionID primitive.ObjectID, action *models.PlayerAction) (*models.GameEvent, error) {
	event := &models.GameEvent{
		ID:          primitive.NewObjectID(),
		SessionID:   sessionID,
		Type:        "player_action",
		Description: action.Action,
		Timestamp:   action.Timestamp,
		ActorID:     action.CharacterID,
		Data: map[string]interface{}{
			"character_name": action.CharacterName,
			"action_type":    action.ActionType,
			"target":         action.Target,
			"dice_roll":      action.DiceRoll,
		},
	}

	err := s.StoreEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// StoreAIResponse stores an AI response as a game event
func (s *EventService) StoreAIResponse(ctx context.Context, sessionID primitive.ObjectID, response *models.AIResponse) (*models.GameEvent, error) {
	event := &models.GameEvent{
		ID:          primitive.NewObjectID(),
		SessionID:   sessionID,
		Type:        "ai_response",
		Description: response.Narrative,
		Timestamp:   response.Timestamp,
		Data: map[string]interface{}{
			"game_mechanics": response.GameMechanics,
			"tokens_used":    response.TokensUsed,
		},
	}

	err := s.StoreEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetRecentEvents retrieves recent events for a session
func (s *EventService) GetRecentEvents(ctx context.Context, sessionID primitive.ObjectID, limit int) ([]models.GameEvent, error) {
	if limit <= 0 {
		limit = 10 // Default to last 10 events
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.db.GetCollection("game_events").Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.GameEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}

	return events, nil
}

// GetSessionNarrative retrieves all narrative events for a session
func (s *EventService) GetSessionNarrative(ctx context.Context, sessionID primitive.ObjectID) ([]models.GameEvent, error) {
	filter := bson.M{
		"session_id": sessionID,
		"type": bson.M{
			"$in": []string{"player_action", "ai_response", "narrative"},
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := s.db.GetCollection("game_events").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find narrative events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.GameEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// GetEventsByType retrieves events of a specific type for a session
func (s *EventService) GetEventsByType(ctx context.Context, sessionID primitive.ObjectID, eventType string) ([]models.GameEvent, error) {
	filter := bson.M{
		"session_id": sessionID,
		"type":       eventType,
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := s.db.GetCollection("game_events").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.GameEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

// DeleteSessionEvents deletes all events for a session (cleanup)
func (s *EventService) DeleteSessionEvents(ctx context.Context, sessionID primitive.ObjectID) error {
	_, err := s.db.GetCollection("game_events").DeleteMany(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return fmt.Errorf("failed to delete session events: %w", err)
	}
	return nil
}

// CountSessionEvents counts the total events for a session
func (s *EventService) CountSessionEvents(ctx context.Context, sessionID primitive.ObjectID) (int64, error) {
	count, err := s.db.GetCollection("game_events").CountDocuments(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}
	return count, nil
}

// GetEventsSince retrieves events that occurred after a specific timestamp
func (s *EventService) GetEventsSince(ctx context.Context, sessionID primitive.ObjectID, since time.Time) ([]models.GameEvent, error) {
	filter := bson.M{
		"session_id": sessionID,
		"timestamp": bson.M{
			"$gt": since,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := s.db.GetCollection("game_events").Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find events: %w", err)
	}
	defer cursor.Close(ctx)

	var events []models.GameEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}