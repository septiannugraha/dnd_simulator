package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/models"
	"dnd-simulator/internal/services"
	"dnd-simulator/internal/websocket"
)

type WebSocketHandler struct {
	hub         *websocket.Hub
	diceService *services.DiceService
}

func NewWebSocketHandler(hub *websocket.Hub, diceService *services.DiceService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		diceService: diceService,
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	h.hub.HandleWebSocket(c)
}

// SendChatMessage sends a chat message to a session via REST API
func (h *WebSocketHandler) SendChatMessage(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	var req struct {
		Content     string             `json:"content" binding:"required"`
		CharacterID primitive.ObjectID `json:"character_id,omitempty"`
		IsIC        bool               `json:"is_ic"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageType := models.MessageTypeChatOOC
	if req.IsIC {
		messageType = models.MessageTypeChatIC
	}

	message := models.WSMessage{
		Type:      messageType,
		Timestamp: time.Now(),
		UserID:    userID.(primitive.ObjectID),
		Username:  username.(string),
		SessionID: sessionID,
		Data: map[string]interface{}{
			"content":      req.Content,
			"character_id": req.CharacterID,
			"is_ic":        req.IsIC,
		},
	}

	h.hub.BroadcastToSession(sessionID, message)
	c.JSON(http.StatusOK, gin.H{"message": "Chat message sent"})
}

// RollDice rolls dice and broadcasts the result to session
func (h *WebSocketHandler) RollDice(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	var req struct {
		Dice        string             `json:"dice" binding:"required"`
		Purpose     string             `json:"purpose"`
		CharacterID primitive.ObjectID `json:"character_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Roll the dice
	diceResult, err := h.diceService.ParseAndRoll(req.Dice, req.Purpose)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	diceResult.CharacterID = req.CharacterID

	// Add special messages for critical hits/fails
	var specialMessage string
	if req.Dice == "1d20" || req.Dice == "1d20+0" {
		if diceResult.Result[0] == 20 {
			specialMessage = h.diceService.GetCriticalHitMessage()
		} else if diceResult.Result[0] == 1 {
			specialMessage = h.diceService.GetCriticalFailMessage()
		}
	}

	message := models.WSMessage{
		Type:      models.MessageTypeDiceResult,
		Timestamp: time.Now(),
		UserID:    userID.(primitive.ObjectID),
		Username:  username.(string),
		SessionID: sessionID,
		Data: map[string]interface{}{
			"dice":            diceResult.Dice,
			"result":          diceResult.Result,
			"total":           diceResult.Total,
			"modifier":        diceResult.Modifier,
			"purpose":         diceResult.Purpose,
			"character_id":    diceResult.CharacterID,
			"special_message": specialMessage,
		},
	}

	h.hub.BroadcastToSession(sessionID, message)
	c.JSON(http.StatusOK, gin.H{
		"message": "Dice rolled",
		"result":  diceResult,
	})
}

// UpdateCharacter broadcasts character update to session
func (h *WebSocketHandler) UpdateCharacter(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	var req struct {
		CharacterID primitive.ObjectID `json:"character_id" binding:"required"`
		Field       string             `json:"field" binding:"required"`
		Value       interface{}        `json:"value" binding:"required"`
		OldValue    interface{}        `json:"old_value,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := models.WSMessage{
		Type:      models.MessageTypeCharacterUpdate,
		Timestamp: time.Now(),
		UserID:    userID.(primitive.ObjectID),
		Username:  username.(string),
		SessionID: sessionID,
		Data: map[string]interface{}{
			"character_id": req.CharacterID,
			"field":        req.Field,
			"value":        req.Value,
			"old_value":    req.OldValue,
		},
	}

	h.hub.BroadcastToSession(sessionID, message)
	c.JSON(http.StatusOK, gin.H{"message": "Character update broadcasted"})
}

// GetSessionStatus returns information about active connections in a session
func (h *WebSocketHandler) GetSessionStatus(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	clientCount := h.hub.GetSessionClients(sessionID)
	
	c.JSON(http.StatusOK, gin.H{
		"session_id":     sessionID,
		"active_clients": clientCount,
		"status":         "active",
	})
}

// RollQuickDice provides common D&D dice rolls
func (h *WebSocketHandler) RollQuickDice(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := primitive.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	}

	diceType := c.Param("dice")
	purpose := c.Query("purpose")
	if purpose == "" {
		purpose = "Quick roll"
	}

	// Roll the dice
	diceResult, err := h.diceService.RollStandardDice(diceType, purpose)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := models.WSMessage{
		Type:      models.MessageTypeDiceResult,
		Timestamp: time.Now(),
		UserID:    userID.(primitive.ObjectID),
		Username:  username.(string),
		SessionID: sessionID,
		Data: map[string]interface{}{
			"dice":     diceResult.Dice,
			"result":   diceResult.Result,
			"total":    diceResult.Total,
			"modifier": diceResult.Modifier,
			"purpose":  diceResult.Purpose,
		},
	}

	h.hub.BroadcastToSession(sessionID, message)
	c.JSON(http.StatusOK, gin.H{
		"message": "Dice rolled",
		"result":  diceResult,
	})
}