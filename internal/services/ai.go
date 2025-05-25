package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dnd-simulator/internal/config"
	"dnd-simulator/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AIService handles AI-powered DM responses
type AIService struct {
	apiKey       string
	apiEndpoint  string
	model        string
	httpClient   *http.Client
	maxTokens    int
	temperature  float64
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		apiKey:      cfg.OpenAIKey, // Will be added to config
		apiEndpoint: "https://api.openai.com/v1/chat/completions",
		model:       "gpt-4-turbo-preview",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxTokens:   1000,
		temperature: 0.8,
	}
}

// AIContext contains all the context needed for AI to generate appropriate responses
type AIContext struct {
	Campaign      *models.Campaign     `json:"campaign"`
	CurrentScene  string               `json:"current_scene"`
	Characters    []models.Character   `json:"characters"`
	RecentEvents  []models.GameEvent   `json:"recent_events"`
	PlayerAction  models.PlayerAction  `json:"player_action"`
	TurnOrder     []models.TurnEntry   `json:"turn_order"`
	CurrentTurn   int                  `json:"current_turn"`
}

// GenerateResponse generates an AI DM response based on the game context and player action
func (s *AIService) GenerateResponse(ctx context.Context, aiContext *AIContext) (*models.AIResponse, error) {
	// Build the system prompt with D&D rules and campaign context
	systemPrompt := s.buildSystemPrompt(aiContext)
	
	// Build the user prompt with the current action
	userPrompt := s.buildUserPrompt(aiContext)
	
	// Prepare the API request
	requestBody := map[string]interface{}{
		"model": s.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"max_tokens":  s.maxTokens,
		"temperature": s.temperature,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}
	
	// Parse response
	var apiResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(apiResp.Choices) == 0 {
		return nil, errors.New("no response from AI")
	}
	
	// Extract narrative and any game mechanics
	narrative := apiResp.Choices[0].Message.Content
	mechanics := s.extractGameMechanics(narrative)
	
	return &models.AIResponse{
		ID:            primitive.NewObjectID(),
		Narrative:     narrative,
		GameMechanics: mechanics,
		Timestamp:     time.Now(),
		TokensUsed:    apiResp.Usage.TotalTokens,
	}, nil
}

// buildSystemPrompt creates the system prompt with D&D rules and campaign context
func (s *AIService) buildSystemPrompt(ctx *AIContext) string {
	var sb strings.Builder
	
	sb.WriteString("You are an expert Dungeon Master for a D&D 5e campaign. ")
	sb.WriteString("Your responses should be immersive, descriptive, and follow D&D 5e rules.\n\n")
	
	// Campaign context
	sb.WriteString(fmt.Sprintf("Campaign: %s\n", ctx.Campaign.Name))
	sb.WriteString(fmt.Sprintf("Setting: %s\n", ctx.Campaign.Description))
	if ctx.Campaign.WorldInfo != "" {
		sb.WriteString(fmt.Sprintf("World Info: %s\n", ctx.Campaign.WorldInfo))
	}
	if ctx.CurrentScene != "" {
		sb.WriteString(fmt.Sprintf("Current Scene: %s\n", ctx.CurrentScene))
	}
	
	// Character information
	sb.WriteString("\nActive Characters:\n")
	for _, char := range ctx.Characters {
		sb.WriteString(fmt.Sprintf("- %s: Level %d %s %s, HP: %d/%d, AC: %d\n",
			char.Name, char.Level, char.Race, char.Class,
			char.CurrentHP, char.MaxHP, char.ArmorClass))
	}
	
	// Recent events for continuity
	if len(ctx.RecentEvents) > 0 {
		sb.WriteString("\nRecent Events:\n")
		for _, event := range ctx.RecentEvents {
			sb.WriteString(fmt.Sprintf("- %s\n", event.Description))
		}
	}
	
	sb.WriteString("\nRules:\n")
	sb.WriteString("1. Stay true to D&D 5e mechanics\n")
	sb.WriteString("2. Be descriptive but concise\n")
	sb.WriteString("3. When combat or skill checks are needed, specify the type of roll required\n")
	sb.WriteString("4. Maintain campaign continuity and character consistency\n")
	sb.WriteString("5. Create engaging narratives that give players meaningful choices\n")
	
	return sb.String()
}

// buildUserPrompt creates the user prompt with the current player action
func (s *AIService) buildUserPrompt(ctx *AIContext) string {
	var sb strings.Builder
	
	// Current turn information
	if len(ctx.TurnOrder) > 0 && ctx.CurrentTurn < len(ctx.TurnOrder) {
		currentChar := ctx.TurnOrder[ctx.CurrentTurn]
		sb.WriteString(fmt.Sprintf("It is %s's turn.\n\n", currentChar.Name))
	}
	
	// Player action
	sb.WriteString(fmt.Sprintf("Player: %s\n", ctx.PlayerAction.CharacterName))
	sb.WriteString(fmt.Sprintf("Action: %s\n", ctx.PlayerAction.Action))
	
	if ctx.PlayerAction.Target != "" {
		sb.WriteString(fmt.Sprintf("Target: %s\n", ctx.PlayerAction.Target))
	}
	
	if ctx.PlayerAction.DiceRoll != nil {
		sb.WriteString(fmt.Sprintf("Dice Roll: %s = %d\n", 
			ctx.PlayerAction.DiceRoll.Dice, ctx.PlayerAction.DiceRoll.Total))
	}
	
	sb.WriteString("\nRespond as the DM, describing the outcome and what happens next.")
	
	return sb.String()
}

// extractGameMechanics parses the AI response for game mechanics instructions
func (s *AIService) extractGameMechanics(narrative string) []models.GameMechanic {
	var mechanics []models.GameMechanic
	
	// Look for common patterns like "roll a DC 15 Perception check"
	// This is a simplified version - could be enhanced with regex
	if strings.Contains(strings.ToLower(narrative), "roll") {
		// Extract roll requirements
		mechanic := models.GameMechanic{
			Type:        "dice_roll",
			Description: "Roll required",
			// Would parse DC, skill type, etc. from narrative
		}
		mechanics = append(mechanics, mechanic)
	}
	
	return mechanics
}

// OpenAI API response structures
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}