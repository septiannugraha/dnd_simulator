package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"dnd-simulator/internal/database"
	"dnd-simulator/internal/models"
)

type SessionService struct {
	db *database.DB
}

func NewSessionService(db *database.DB) *SessionService {
	return &SessionService{
		db: db,
	}
}

// CreateSession creates a new game session within a campaign
func (s *SessionService) CreateSession(ctx context.Context, req *models.CreateSessionRequest, dmUserID primitive.ObjectID) (*models.GameSession, error) {
	// Verify the campaign exists and user is the DM
	var campaign models.Campaign
	err := s.db.GetCollection("campaigns").FindOne(ctx, bson.M{
		"_id":    req.CampaignID,
		"dm_id": dmUserID,
	}).Decode(&campaign)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("campaign not found or you are not the DM")
		}
		return nil, fmt.Errorf("failed to verify campaign: %w", err)
	}

	// Create the session
	session := &models.GameSession{
		ID:          primitive.NewObjectID(),
		CampaignID:  req.CampaignID,
		Name:        req.Name,
		Description: req.Description,
		Status:      models.SessionStatusPending,
		TurnOrder:   []models.TurnEntry{},
		CurrentTurn: 0,
		Round:       0,
		Players:     []models.SessionPlayer{},
		DMUserID:    dmUserID,
		ChatHistory: []models.SessionChatMsg{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = s.db.GetCollection("sessions").InsertOne(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by ID
func (s *SessionService) GetSession(ctx context.Context, sessionID primitive.ObjectID) (*models.GameSession, error) {
	var session models.GameSession
	err := s.db.GetCollection("sessions").FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

// GetSessionsForCampaign retrieves all sessions for a campaign
func (s *SessionService) GetSessionsForCampaign(ctx context.Context, campaignID primitive.ObjectID) ([]*models.GameSession, error) {
	cursor, err := s.db.GetCollection("sessions").Find(ctx, bson.M{"campaign_id": campaignID})
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions: %w", err)
	}
	defer cursor.Close(ctx)

	var sessions []*models.GameSession
	for cursor.Next(ctx) {
		var session models.GameSession
		if err := cursor.Decode(&session); err != nil {
			return nil, fmt.Errorf("failed to decode session: %w", err)
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// JoinSession adds a player to a session
func (s *SessionService) JoinSession(ctx context.Context, sessionID, userID, characterID primitive.ObjectID) error {
	// Get session
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Check if session is in a joinable state
	if session.Status != models.SessionStatusPending && session.Status != models.SessionStatusActive {
		return errors.New("session is not accepting new players")
	}

	// Verify character exists and belongs to user
	var character models.Character
	err = s.db.GetCollection("characters").FindOne(ctx, bson.M{
		"_id":     characterID,
		"user_id": userID,
	}).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("character not found or does not belong to you")
		}
		return fmt.Errorf("failed to verify character: %w", err)
	}

	// Get user info
	var user models.User
	err = s.db.GetCollection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if player is already in session
	for _, player := range session.Players {
		if player.UserID == userID {
			return errors.New("you are already in this session")
		}
	}

	// Add player to session
	newPlayer := models.SessionPlayer{
		UserID:      userID,
		CharacterID: characterID,
		Username:    user.Username,
		CharName:    character.Name,
		IsConnected: false,
		JoinedAt:    time.Now(),
	}

	_, err = s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{"_id": sessionID},
		bson.M{
			"$push": bson.M{"players": newPlayer},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to join session: %w", err)
	}

	return nil
}

// LeaveSession removes a player from a session
func (s *SessionService) LeaveSession(ctx context.Context, sessionID, userID primitive.ObjectID) error {
	_, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{"_id": sessionID},
		bson.M{
			"$pull": bson.M{"players": bson.M{"user_id": userID}},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to leave session: %w", err)
	}
	return nil
}

// StartSession transitions a session from pending to active
func (s *SessionService) StartSession(ctx context.Context, sessionID, dmUserID primitive.ObjectID) error {
	now := time.Now()
	result, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":       sessionID,
			"dm_user_id": dmUserID,
			"status":    models.SessionStatusPending,
		},
		bson.M{
			"$set": bson.M{
				"status":     models.SessionStatusActive,
				"started_at": &now,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found or you are not the DM, or session is not in pending status")
	}
	return nil
}

// SetInitiative sets or updates a character's initiative in the turn order
func (s *SessionService) SetInitiative(ctx context.Context, sessionID, characterID primitive.ObjectID, initiative int) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Get character info
	var character models.Character
	err = s.db.GetCollection("characters").FindOne(ctx, bson.M{"_id": characterID}).Decode(&character)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// Update or add turn entry
	turnOrder := session.TurnOrder
	found := false
	for i, entry := range turnOrder {
		if entry.CharacterID == characterID {
			turnOrder[i].Initiative = initiative
			turnOrder[i].HasActed = false
			found = true
			break
		}
	}

	if !found {
		newEntry := models.TurnEntry{
			Type:        models.TurnTypePlayer,
			UserID:      character.UserID,
			CharacterID: characterID,
			Initiative:  initiative,
			Name:        character.Name,
			HasActed:    false,
		}
		turnOrder = append(turnOrder, newEntry)
	}

	// Sort turn order by initiative (descending)
	for i := 0; i < len(turnOrder)-1; i++ {
		for j := 0; j < len(turnOrder)-i-1; j++ {
			if turnOrder[j].Initiative < turnOrder[j+1].Initiative {
				turnOrder[j], turnOrder[j+1] = turnOrder[j+1], turnOrder[j]
			}
		}
	}

	// Update session
	_, err = s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{"_id": sessionID},
		bson.M{
			"$set": bson.M{
				"turn_order": turnOrder,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update initiative: %w", err)
	}

	return nil
}

// AdvanceTurn moves to the next turn in the order
func (s *SessionService) AdvanceTurn(ctx context.Context, sessionID, dmUserID primitive.ObjectID, force bool) error {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	// Verify DM permission
	if session.DMUserID != dmUserID {
		return errors.New("only the DM can advance turns")
	}

	if len(session.TurnOrder) == 0 {
		return errors.New("no turn order established")
	}

	// Mark current player as having acted (if not forced)
	if !force && session.CurrentTurn < len(session.TurnOrder) {
		session.TurnOrder[session.CurrentTurn].HasActed = true
	}

	// Advance to next turn
	session.CurrentTurn++
	if session.CurrentTurn >= len(session.TurnOrder) {
		session.CurrentTurn = 0
		session.Round++
		// Reset HasActed for new round
		for i := range session.TurnOrder {
			session.TurnOrder[i].HasActed = false
		}
	}

	// Update session
	_, err = s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{"_id": sessionID},
		bson.M{
			"$set": bson.M{
				"turn_order":   session.TurnOrder,
				"current_turn": session.CurrentTurn,
				"round":        session.Round,
				"updated_at":   time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to advance turn: %w", err)
	}

	return nil
}

// UpdatePlayerConnection updates a player's connection status
func (s *SessionService) UpdatePlayerConnection(ctx context.Context, sessionID, userID primitive.ObjectID, isConnected bool) error {
	_, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":              sessionID,
			"players.user_id": userID,
		},
		bson.M{
			"$set": bson.M{
				"players.$.is_connected": isConnected,
				"updated_at":             time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update player connection: %w", err)
	}
	return nil
}

// AddChatMessage adds a chat message to the session history
func (s *SessionService) AddChatMessage(ctx context.Context, sessionID primitive.ObjectID, msg models.SessionChatMsg) error {
	_, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{"_id": sessionID},
		bson.M{
			"$push": bson.M{"chat_history": msg},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to add chat message: %w", err)
	}
	return nil
}

// UpdateScene updates the current scene description and optional notes
func (s *SessionService) UpdateScene(ctx context.Context, sessionID, dmUserID primitive.ObjectID, scene, notes string) error {
	result, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":        sessionID,
			"dm_user_id": dmUserID,
		},
		bson.M{
			"$set": bson.M{
				"scene":      scene,
				"notes":      notes,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update scene: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found or you are not the DM")
	}
	return nil
}

// EndSession marks a session as completed or cancelled
func (s *SessionService) EndSession(ctx context.Context, sessionID, dmUserID primitive.ObjectID, status models.SessionStatus) error {
	if status != models.SessionStatusCompleted && status != models.SessionStatusCancelled {
		return errors.New("invalid end status")
	}

	now := time.Now()
	result, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":        sessionID,
			"dm_user_id": dmUserID,
			"status": bson.M{
				"$in": []models.SessionStatus{
					models.SessionStatusActive,
					models.SessionStatusPaused,
					models.SessionStatusPending,
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"ended_at":   &now,
				"updated_at": now,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to end session: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found, you are not the DM, or session is already ended")
	}
	return nil
}

// PauseSession pauses an active session
func (s *SessionService) PauseSession(ctx context.Context, sessionID, dmUserID primitive.ObjectID) error {
	result, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":        sessionID,
			"dm_user_id": dmUserID,
			"status":     models.SessionStatusActive,
		},
		bson.M{
			"$set": bson.M{
				"status":     models.SessionStatusPaused,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to pause session: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found, you are not the DM, or session is not active")
	}
	return nil
}

// ResumeSession resumes a paused session
func (s *SessionService) ResumeSession(ctx context.Context, sessionID, dmUserID primitive.ObjectID) error {
	result, err := s.db.GetCollection("sessions").UpdateOne(ctx,
		bson.M{
			"_id":        sessionID,
			"dm_user_id": dmUserID,
			"status":     models.SessionStatusPaused,
		},
		bson.M{
			"$set": bson.M{
				"status":     models.SessionStatusActive,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to resume session: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("session not found, you are not the DM, or session is not paused")
	}
	return nil
}