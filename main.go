package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"dnd-simulator/internal/auth"
	"dnd-simulator/internal/config"
	"dnd-simulator/internal/database"
	"dnd-simulator/internal/handlers"
	"dnd-simulator/internal/middleware"
	"dnd-simulator/internal/services"
	"dnd-simulator/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userService := services.NewUserService(db)
	campaignService := services.NewCampaignService(db)
	characterService := services.NewCharacterService(db)
	diceService := services.NewDiceService()

	// Initialize WebSocket hub and start it
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, jwtService)
	campaignHandler := handlers.NewCampaignHandler(campaignService)
	characterHandler := handlers.NewCharacterHandler(characterService)
	wsHandler := handlers.NewWebSocketHandler(hub, diceService)

	// Setup router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "D&D Simulator API is running",
		})
	})

	// API routes
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.Me)
		}

		// Campaign routes
		campaigns := api.Group("/campaigns", middleware.AuthMiddleware(jwtService))
		{
			campaigns.POST("", campaignHandler.CreateCampaign)                    // Create campaign
			campaigns.GET("", campaignHandler.GetCampaigns)                       // List user's campaigns
			campaigns.GET("/public", campaignHandler.GetPublicCampaigns)          // List public campaigns
			campaigns.GET("/:id", campaignHandler.GetCampaign)                    // Get campaign details
			campaigns.PUT("/:id", campaignHandler.UpdateCampaign)                 // Update campaign (DM only)
			campaigns.DELETE("/:id", campaignHandler.DeleteCampaign)              // Delete campaign (DM only)
			campaigns.POST("/:id/join", campaignHandler.JoinCampaign)             // Join campaign
			campaigns.POST("/:id/leave", campaignHandler.LeaveCampaign)           // Leave campaign
		}

		// Character routes
		characters := api.Group("/characters", middleware.AuthMiddleware(jwtService))
		{
			characters.POST("", characterHandler.CreateCharacter)                 // Create character
			characters.GET("", characterHandler.GetCharacters)                    // List user's characters
			characters.GET("/:id", characterHandler.GetCharacter)                 // Get character details
			characters.PUT("/:id", characterHandler.UpdateCharacter)              // Update character
			characters.DELETE("/:id", characterHandler.DeleteCharacter)           // Delete character
			characters.POST("/:id/assign", characterHandler.AssignToCampaign)     // Assign to campaign
		}

		// D&D Data routes (for character creation)
		dnd := api.Group("/dnd")
		{
			dnd.GET("/races", characterHandler.GetRaces)                          // Get available races
			dnd.GET("/classes", characterHandler.GetClasses)                      // Get available classes
			dnd.GET("/backgrounds", characterHandler.GetBackgrounds)              // Get available backgrounds
		}

		// WebSocket and real-time routes
		sessions := api.Group("/sessions", middleware.AuthMiddleware(jwtService))
		{
			// WebSocket connection endpoint
			sessions.GET("/:id/ws", wsHandler.HandleWebSocket)                    // WebSocket connection
			
			// REST API for real-time features
			sessions.POST("/:id/chat", wsHandler.SendChatMessage)                 // Send chat message
			sessions.POST("/:id/dice", wsHandler.RollDice)                        // Roll custom dice
			sessions.POST("/:id/dice/:dice", wsHandler.RollQuickDice)             // Quick dice roll (d20, d6, etc.)
			sessions.POST("/:id/character-update", wsHandler.UpdateCharacter)     // Broadcast character update
			sessions.GET("/:id/status", wsHandler.GetSessionStatus)               // Get session status
		}
	}

	log.Printf("Starting D&D Simulator server on :%s", cfg.Port)
	r.Run(":" + cfg.Port)
}