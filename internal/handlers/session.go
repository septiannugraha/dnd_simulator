package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"dnd-simulator/internal/models"
	"dnd-simulator/internal/services"
)

type SessionHandler struct {
	sessionService  *services.SessionService
	campaignService *services.CampaignService
}

func NewSessionHandler(sessionService *services.SessionService, campaignService *services.CampaignService) *SessionHandler {
	return &SessionHandler{
		sessionService:  sessionService,
		campaignService: campaignService,
	}
}

// CreateSession creates a new game session within a campaign
// POST /api/sessions
func (h *SessionHandler) CreateSession(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjID := userID.(primitive.ObjectID)
	session, err := h.sessionService.CreateSession(c.Request.Context(), &req, userObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// GetSession retrieves a specific session
// GET /api/sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetCampaignSessions retrieves all sessions for a campaign
// GET /api/campaigns/:id/sessions
func (h *SessionHandler) GetCampaignSessions(c *gin.Context) {
	campaignID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	sessions, err := h.sessionService.GetSessionsForCampaign(c.Request.Context(), campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

// JoinSession adds a player to a session
// POST /api/sessions/:id/join
func (h *SessionHandler) JoinSession(c *gin.Context) {
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

	var req models.JoinSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.JoinSession(c.Request.Context(), sessionID, userObjID, req.CharacterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined session"})
}

// LeaveSession removes a player from a session
// POST /api/sessions/:id/leave
func (h *SessionHandler) LeaveSession(c *gin.Context) {
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

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.LeaveSession(c.Request.Context(), sessionID, userObjID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left session"})
}

// StartSession transitions a session from pending to active
// POST /api/sessions/:id/start
func (h *SessionHandler) StartSession(c *gin.Context) {
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

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.StartSession(c.Request.Context(), sessionID, userObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session started successfully"})
}

// SetInitiative sets or updates a character's initiative
// POST /api/sessions/:id/initiative
func (h *SessionHandler) SetInitiative(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	var req models.SetInitiativeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.sessionService.SetInitiative(c.Request.Context(), sessionID, req.CharacterID, req.Initiative)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Initiative set successfully"})
}

// AdvanceTurn moves to the next turn in the order
// POST /api/sessions/:id/turn/advance
func (h *SessionHandler) AdvanceTurn(c *gin.Context) {
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

	var req models.AdvanceTurnRequest
	c.ShouldBindJSON(&req) // Optional body

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.AdvanceTurn(c.Request.Context(), sessionID, userObjID, req.Force)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Turn advanced successfully"})
}

// GetSessionStatus returns current session state for WebSocket clients
// GET /api/sessions/:id/status
func (h *SessionHandler) GetSessionStatus(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	session, err := h.sessionService.GetSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return just the turn/game state info
	response := gin.H{
		"session_id":    session.ID,
		"status":        session.Status,
		"current_turn":  session.CurrentTurn,
		"round":         session.Round,
		"turn_order":    session.TurnOrder,
		"scene":         session.Scene,
		"player_count":  len(session.Players),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateScene allows DM to update the current scene description
// PUT /api/sessions/:id/scene
func (h *SessionHandler) UpdateScene(c *gin.Context) {
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

	var req struct {
		Scene string `json:"scene" binding:"required"`
		Notes string `json:"notes,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.UpdateScene(c.Request.Context(), sessionID, userObjID, req.Scene, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scene updated successfully"})
}

// EndSession marks a session as completed or cancelled
// POST /api/sessions/:id/end
func (h *SessionHandler) EndSession(c *gin.Context) {
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

	var req struct {
		Status string `json:"status" binding:"required,oneof=completed cancelled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.EndSession(c.Request.Context(), sessionID, userObjID, models.SessionStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session ended successfully"})
}

// PauseSession pauses an active session
// POST /api/sessions/:id/pause
func (h *SessionHandler) PauseSession(c *gin.Context) {
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

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.PauseSession(c.Request.Context(), sessionID, userObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session paused successfully"})
}

// ResumeSession resumes a paused session
// POST /api/sessions/:id/resume
func (h *SessionHandler) ResumeSession(c *gin.Context) {
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

	userObjID := userID.(primitive.ObjectID)
	err = h.sessionService.ResumeSession(c.Request.Context(), sessionID, userObjID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session resumed successfully"})
}