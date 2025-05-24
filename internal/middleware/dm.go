package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/services"
)

func DMMiddleware(campaignService *services.CampaignService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get campaign ID from URL parameter
		campaignIDStr := c.Param("id")
		campaignID, err := primitive.ObjectIDFromHex(campaignIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
			c.Abort()
			return
		}

		// Get user ID from auth middleware
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get campaign and check if user is DM
		campaign, err := campaignService.GetCampaignByID(campaignID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
			c.Abort()
			return
		}

		if campaign.DMID != userID.(primitive.ObjectID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the DM can perform this action"})
			c.Abort()
			return
		}

		// Store campaign in context for handler use
		c.Set("campaign", campaign)
		c.Next()
	}
}