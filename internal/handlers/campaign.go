package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/models"
	"dnd-simulator/internal/services"
)

type CampaignHandler struct {
	campaignService *services.CampaignService
}

func NewCampaignHandler(campaignService *services.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		campaignService: campaignService,
	}
}

type CreateCampaignRequest struct {
	Name        string                   `json:"name" binding:"required,min=3,max=100"`
	Description string                   `json:"description" binding:"max=1000"`
	WorldInfo   string                   `json:"world_info" binding:"max=5000"`
	Settings    models.CampaignSettings  `json:"settings"`
}

type UpdateCampaignRequest struct {
	Name        string                   `json:"name" binding:"required,min=3,max=100"`
	Description string                   `json:"description" binding:"max=1000"`
	WorldInfo   string                   `json:"world_info" binding:"max=5000"`
	Settings    models.CampaignSettings  `json:"settings"`
}

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set default settings if not provided
	if req.Settings.MaxPlayers == 0 {
		req.Settings.MaxPlayers = 6
	}

	campaign, err := h.campaignService.CreateCampaign(
		userID.(primitive.ObjectID),
		req.Name,
		req.Description,
		req.WorldInfo,
		req.Settings,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"campaign": campaign})
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	campaigns, err := h.campaignService.GetUserCampaigns(userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(campaignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user has access to this campaign
	hasAccess := campaign.DMID == userID.(primitive.ObjectID)
	if !hasAccess {
		for _, playerID := range campaign.PlayerIDs {
			if playerID == userID.(primitive.ObjectID) {
				hasAccess = true
				break
			}
		}
	}

	if !hasAccess && !campaign.Settings.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaign": campaign})
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	var req UpdateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	campaign, err := h.campaignService.UpdateCampaign(
		campaignID,
		userID.(primitive.ObjectID),
		req.Name,
		req.Description,
		req.WorldInfo,
		req.Settings,
	)
	if err != nil {
		if err.Error() == "only the DM can update this campaign" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaign": campaign})
}

func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.campaignService.DeleteCampaign(campaignID, userID.(primitive.ObjectID))
	if err != nil {
		if err.Error() == "only the DM can delete this campaign" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}

func (h *CampaignHandler) JoinCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.campaignService.JoinCampaign(campaignID, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined campaign"})
}

func (h *CampaignHandler) LeaveCampaign(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.campaignService.LeaveCampaign(campaignID, userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully left campaign"})
}

func (h *CampaignHandler) GetPublicCampaigns(c *gin.Context) {
	campaigns, err := h.campaignService.GetPublicCampaigns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get public campaigns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}