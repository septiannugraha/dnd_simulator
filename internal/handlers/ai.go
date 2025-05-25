package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"dnd-simulator/internal/models"
	"dnd-simulator/internal/services"
)

type AIHandler struct {
	aiService       *services.AIService
	sessionService  *services.SessionService
	characterService *services.CharacterService
	campaignService *services.CampaignService
	eventService    *services.EventService
}

func NewAIHandler(
	aiService *services.AIService,
	sessionService *services.SessionService,
	characterService *services.CharacterService,
	campaignService *services.CampaignService,
	eventService *services.EventService,
) *AIHandler {
	return &AIHandler{
		aiService:       aiService,
		sessionService:  sessionService,
		characterService: characterService,
		campaignService: campaignService,
		eventService:    eventService,
	}
}

// ProcessPlayerAction handles a player action and generates an AI DM response
// POST /api/sessions/:id/action
func (h *AIHandler) ProcessPlayerAction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// Parse player action
	var req struct {
		CharacterID primitive.ObjectID `json:"character_id" binding:"required"`
		Action      string             `json:"action" binding:"required"`
		ActionType  string             `json:"action_type" binding:"required,oneof=combat roleplay exploration"`
		Target      string             `json:"target,omitempty"`
		DiceRoll    *models.DiceRoll   `json:"dice_roll,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get session
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Verify session is active
	if session.Status != models.SessionStatusActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not active"})
		return
	}

	// Get character
	character, err := h.characterService.GetCharacterByID(req.CharacterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Verify character belongs to user
	if character.UserID != userID.(primitive.ObjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Character does not belong to you"})
		return
	}

	// Get campaign
	campaign, err := h.campaignService.GetCampaignByID(session.CampaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaign"})
		return
	}

	// Get all characters in session
	characterIDs := make([]primitive.ObjectID, 0, len(session.Players))
	for _, player := range session.Players {
		characterIDs = append(characterIDs, player.CharacterID)
	}
	
	characters, err := h.characterService.GetCharactersByIDs(c.Request.Context(), characterIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get characters"})
		return
	}

	// Get recent events from the event service
	recentEvents, err := h.eventService.GetRecentEvents(c.Request.Context(), sessionID, 10)
	if err != nil {
		// Log error but continue - recent events are helpful but not critical
		recentEvents = []models.GameEvent{}
	}

	// Build AI context
	aiContext := &services.AIContext{
		Campaign:     campaign,
		CurrentScene: session.Scene,
		Characters:   characters,
		RecentEvents: recentEvents,
		PlayerAction: models.PlayerAction{
			CharacterID:   req.CharacterID,
			CharacterName: character.Name,
			Action:        req.Action,
			ActionType:    req.ActionType,
			Target:        req.Target,
			DiceRoll:      req.DiceRoll,
			Timestamp:     time.Now(),
		},
		TurnOrder:   session.TurnOrder,
		CurrentTurn: session.CurrentTurn,
	}

	// Generate AI response
	aiResponse, err := h.aiService.GenerateResponse(c.Request.Context(), aiContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate AI response"})
		return
	}

	// Set session ID
	aiResponse.SessionID = sessionID

	// Store the player action event
	playerAction := models.PlayerAction{
		CharacterID:   req.CharacterID,
		CharacterName: character.Name,
		Action:        req.Action,
		ActionType:    req.ActionType,
		Target:        req.Target,
		DiceRoll:      req.DiceRoll,
		Timestamp:     time.Now(),
	}
	
	gameEvent, err := h.eventService.StorePlayerAction(c.Request.Context(), sessionID, &playerAction)
	if err != nil {
		// Log error but continue - we still want to return the AI response
		gameEvent = &models.GameEvent{
			ID:          primitive.NewObjectID(),
			SessionID:   sessionID,
			Type:        "player_action",
			Description: req.Action,
			Timestamp:   time.Now(),
			ActorID:     req.CharacterID,
		}
	}

	// Store AI response event
	aiEvent, err := h.eventService.StoreAIResponse(c.Request.Context(), sessionID, aiResponse)
	if err != nil {
		// Log error but continue
		aiEvent = &models.GameEvent{
			ID:          primitive.NewObjectID(),
			SessionID:   sessionID,
			Type:        "ai_response",
			Description: aiResponse.Narrative,
			Timestamp:   aiResponse.Timestamp,
		}
	}

	// Broadcast to WebSocket clients
	wsMessage := models.WSMessage{
		Type:      models.MessageTypeAIResponse,
		Timestamp: time.Now(),
		UserID:    userID.(primitive.ObjectID),
		SessionID: sessionID,
		Data: map[string]interface{}{
			"character_name": character.Name,
			"action":         req.Action,
			"ai_response":    aiResponse,
			"game_event":     gameEvent,
			"ai_event":       aiEvent,
		},
	}

	// Send through WebSocket hub (would need to inject hub)
	// h.hub.BroadcastToSession(sessionID, wsMessage)

	c.JSON(http.StatusOK, gin.H{
		"response":    aiResponse,
		"game_event":  gameEvent,
		"ai_event":    aiEvent,
		"ws_message":  wsMessage,
	})
}

// GetNarrativeHistory returns the narrative history for a session
// GET /api/sessions/:id/narrative
func (h *AIHandler) GetNarrativeHistory(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	// Get the session to verify it exists
	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	// Get all narrative events for the session
	events, err := h.eventService.GetSessionNarrative(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve narrative history"})
		return
	}

	// Get event count for statistics
	eventCount, err := h.eventService.CountSessionEvents(c.Request.Context(), sessionID)
	if err != nil {
		eventCount = int64(len(events))
	}
	
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"session_name": session.Name,
		"narrative": events,
		"event_count": eventCount,
		"current_scene": session.Scene,
	})
}

// GetEventsByType returns events of a specific type for a session
// GET /api/sessions/:id/events/:type
func (h *AIHandler) GetEventsByType(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	eventType := c.Param("type")
	validTypes := []string{"player_action", "ai_response", "dice_roll", "combat", "narrative"}
	isValid := false
	for _, vt := range validTypes {
		if eventType == vt {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event type"})
		return
	}

	events, err := h.eventService.GetEventsByType(c.Request.Context(), sessionID, eventType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"event_type": eventType,
		"events": events,
		"count": len(events),
	})
}