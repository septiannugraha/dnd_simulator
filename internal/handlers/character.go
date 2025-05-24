package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/data"
	"dnd-simulator/internal/services"
)

type CharacterHandler struct {
	characterService *services.CharacterService
}

func NewCharacterHandler(characterService *services.CharacterService) *CharacterHandler {
	return &CharacterHandler{
		characterService: characterService,
	}
}

func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req services.CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	character, err := h.characterService.CreateCharacter(userID.(primitive.ObjectID), req)
	if err != nil {
		if err.Error() == "invalid race" || err.Error() == "invalid class" || err.Error() == "invalid background" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"character": character})
}

func (h *CharacterHandler) GetCharacters(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	characters, err := h.characterService.GetUserCharacters(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get characters"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"characters": characters})
}

func (h *CharacterHandler) GetCharacter(c *gin.Context) {
	characterIDStr := c.Param("id")
	characterID, err := primitive.ObjectIDFromHex(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	character, err := h.characterService.GetCharacterByID(characterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user owns this character
	if character.UserID != userID.(primitive.ObjectID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}

func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
	characterIDStr := c.Param("id")
	characterID, err := primitive.ObjectIDFromHex(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character, err := h.characterService.UpdateCharacter(characterID, userID.(primitive.ObjectID), updates)
	if err != nil {
		if err.Error() == "you can only update your own characters" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update character"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"character": character})
}

func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
	characterIDStr := c.Param("id")
	characterID, err := primitive.ObjectIDFromHex(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.characterService.DeleteCharacter(characterID, userID.(primitive.ObjectID))
	if err != nil {
		if err.Error() == "you can only delete your own characters" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete character"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character deleted successfully"})
}

func (h *CharacterHandler) AssignToCampaign(c *gin.Context) {
	characterIDStr := c.Param("id")
	characterID, err := primitive.ObjectIDFromHex(characterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	var req struct {
		CampaignID string `json:"campaign_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaignID, err := primitive.ObjectIDFromHex(req.CampaignID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.characterService.AssignToCampaign(characterID, campaignID, userID.(primitive.ObjectID))
	if err != nil {
		if err.Error() == "you can only assign your own characters" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign character to campaign"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character assigned to campaign successfully"})
}

// Helper endpoints to get available races, classes, backgrounds
func (h *CharacterHandler) GetRaces(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"races": data.Races})
}

func (h *CharacterHandler) GetClasses(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"classes": data.Classes})
}

func (h *CharacterHandler) GetBackgrounds(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"backgrounds": data.Backgrounds})
}