package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	temperature  float64
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.Config) *AIService {
	return &AIService{
		apiKey:      cfg.GeminiAPIKey,
		apiEndpoint: "https://generativelanguage.googleapis.com/v1beta/models",
		model:       "gemini-2.5-pro-preview-05-06",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
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
	// Build the combined prompt
	prompt := s.buildCombinedPrompt(aiContext)
	
	// Prepare the API request for Gemini
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{
						"text": prompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": s.temperature,
			"maxOutputTokens": 8192,
			"responseMimeType": "text/plain",
			"topK": 40,
			"topP": 0.95,
			"stopSequences": []string{"[END_SCENE]", "[AWAIT_PLAYER_ACTION]"},
		},
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request - using streamGenerateContent endpoint
	url := fmt.Sprintf("%s/%s:streamGenerateContent?key=%s&alt=sse", s.apiEndpoint, s.model, s.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorBody)
		return nil, fmt.Errorf("API returned status code %d: %v", resp.StatusCode, errorBody)
	}
	
	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Parse SSE data
	var fullText strings.Builder
	var totalTokens int
	
	// Split by "data: " to get individual SSE messages
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")
			jsonData = strings.TrimSpace(jsonData)
			if jsonData == "" || jsonData == "[DONE]" {
				continue
			}
			
			var chunk GeminiStreamResponse
			if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
				continue
			}
			
			// Extract text from candidates
			if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
				text := chunk.Candidates[0].Content.Parts[0].Text
				fullText.WriteString(text)
			}
			
			// Update token count
			if chunk.UsageMetadata.TotalTokenCount > 0 {
				totalTokens = chunk.UsageMetadata.TotalTokenCount
			}
		}
	}
	
	narrative := strings.TrimSpace(fullText.String())
	if narrative == "" {
		return nil, errors.New("empty response from AI")
	}
	
	// Extract game mechanics from the narrative
	mechanics := s.extractGameMechanics(narrative)
	
	return &models.AIResponse{
		ID:            primitive.NewObjectID(),
		Narrative:     narrative,
		GameMechanics: mechanics,
		Timestamp:     time.Now(),
		TokensUsed:    totalTokens,
	}, nil
}

// buildCombinedPrompt creates a single prompt with all context for Gemini
func (s *AIService) buildCombinedPrompt(ctx *AIContext) string {
	var sb strings.Builder
	
	// System context
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
	
	// Current action
	sb.WriteString("\n--- CURRENT ACTION ---\n")
	
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
	lowerNarrative := strings.ToLower(narrative)
	
	// Check for dice roll requirements
	if strings.Contains(lowerNarrative, "roll") || strings.Contains(lowerNarrative, "make a") {
		// Look for DC mentions
		if strings.Contains(lowerNarrative, "dc") {
			mechanic := models.GameMechanic{
				Type:        "dice_roll",
				Description: "Skill check required",
			}
			mechanics = append(mechanics, mechanic)
		}
		
		// Look for attack rolls
		if strings.Contains(lowerNarrative, "attack") || strings.Contains(lowerNarrative, "hit") {
			mechanic := models.GameMechanic{
				Type:        "attack_roll",
				Description: "Attack roll required",
			}
			mechanics = append(mechanics, mechanic)
		}
		
		// Look for saving throws
		if strings.Contains(lowerNarrative, "saving throw") || strings.Contains(lowerNarrative, "save") {
			mechanic := models.GameMechanic{
				Type:        "saving_throw",
				Description: "Saving throw required",
			}
			mechanics = append(mechanics, mechanic)
		}
	}
	
	// Check for damage mentions
	if strings.Contains(lowerNarrative, "damage") || strings.Contains(lowerNarrative, "takes") {
		mechanic := models.GameMechanic{
			Type:        "damage",
			Description: "Damage calculation needed",
		}
		mechanics = append(mechanics, mechanic)
	}
	
	return mechanics
}

// GeminiStreamResponse represents a single chunk in the streaming response
type GeminiStreamResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason,omitempty"`
		Index        int    `json:"index"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}