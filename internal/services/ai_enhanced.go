package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"dnd-simulator/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnhancedAIService with structured output support
type EnhancedAIService struct {
	apiKey      string
	model       string
	temperature float64
	topK        int
	topP        float64
}

// StructuredAIResponse defines the JSON schema for AI responses
type StructuredAIResponse struct {
	Narrative     string                    `json:"narrative"`
	GameMechanics []StructuredGameMechanic `json:"game_mechanics"`
	NPCActions    []NPCAction              `json:"npc_actions,omitempty"`
	SceneChanges  *SceneChange             `json:"scene_changes,omitempty"`
}

type StructuredGameMechanic struct {
	Type         string            `json:"type"`         // "skill_check", "attack_roll", "saving_throw", "damage", "condition"
	Description  string            `json:"description"`
	Target       string            `json:"target,omitempty"`
	Requirements *RollRequirements `json:"requirements,omitempty"`
	Effects      *MechanicEffects  `json:"effects,omitempty"`
}

type RollRequirements struct {
	RollType      string   `json:"roll_type"`      // "d20", "damage", etc.
	DC            int      `json:"dc,omitempty"`   // Difficulty Class for checks
	Skill         string   `json:"skill,omitempty"` // For skill checks
	SaveType      string   `json:"save_type,omitempty"` // For saving throws
	AttackBonus   int      `json:"attack_bonus,omitempty"`
	DamageDice    string   `json:"damage_dice,omitempty"` // e.g., "2d6+3"
	DamageType    string   `json:"damage_type,omitempty"` // e.g., "fire", "slashing"
	Advantage     bool     `json:"advantage,omitempty"`
	Disadvantage  bool     `json:"disadvantage,omitempty"`
}

type MechanicEffects struct {
	OnSuccess     string   `json:"on_success,omitempty"`
	OnFailure     string   `json:"on_failure,omitempty"`
	Conditions    []string `json:"conditions,omitempty"` // e.g., ["prone", "frightened"]
	Duration      string   `json:"duration,omitempty"`   // e.g., "1 minute", "until end of turn"
}

type NPCAction struct {
	NPCName      string `json:"npc_name"`
	Action       string `json:"action"`
	Target       string `json:"target,omitempty"`
	Initiative   int    `json:"initiative,omitempty"`
}

type SceneChange struct {
	NewLocation  string   `json:"new_location,omitempty"`
	TimeChange   string   `json:"time_change,omitempty"`
	Weather      string   `json:"weather,omitempty"`
	Visibility   string   `json:"visibility,omitempty"`
	Hazards      []string `json:"hazards,omitempty"`
}

func NewEnhancedAIService(apiKey, model string) *EnhancedAIService {
	return &EnhancedAIService{
		apiKey:      apiKey,
		model:       model,
		temperature: 0.8,  // Balanced creativity
		topK:        40,   // Good variety
		topP:        0.95, // High quality sampling
	}
}

func (s *EnhancedAIService) GenerateDMResponse(ctx models.AIContext) (*models.AIResponse, error) {
	prompt := s.buildStructuredPrompt(ctx)
	
	// Define the response schema for structured output
	responseSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"narrative": map[string]interface{}{
				"type": "string",
				"description": "The narrative description of what happens in the scene",
			},
			"game_mechanics": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"skill_check", "attack_roll", "saving_throw", "damage", "condition", "initiative"},
						},
						"description": map[string]interface{}{
							"type": "string",
						},
						"target": map[string]interface{}{
							"type": "string",
						},
						"requirements": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"roll_type": map[string]interface{}{"type": "string"},
								"dc": map[string]interface{}{"type": "integer"},
								"skill": map[string]interface{}{"type": "string"},
								"save_type": map[string]interface{}{"type": "string"},
								"attack_bonus": map[string]interface{}{"type": "integer"},
								"damage_dice": map[string]interface{}{"type": "string"},
								"damage_type": map[string]interface{}{"type": "string"},
								"advantage": map[string]interface{}{"type": "boolean"},
								"disadvantage": map[string]interface{}{"type": "boolean"},
							},
						},
						"effects": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"on_success": map[string]interface{}{"type": "string"},
								"on_failure": map[string]interface{}{"type": "string"},
								"conditions": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{"type": "string"},
								},
								"duration": map[string]interface{}{"type": "string"},
							},
						},
					},
					"required": []string{"type", "description"},
				},
			},
			"npc_actions": map[string]interface{}{
				"type": "array",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"npc_name": map[string]interface{}{"type": "string"},
						"action": map[string]interface{}{"type": "string"},
						"target": map[string]interface{}{"type": "string"},
						"initiative": map[string]interface{}{"type": "integer"},
					},
					"required": []string{"npc_name", "action"},
				},
			},
			"scene_changes": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"new_location": map[string]interface{}{"type": "string"},
					"time_change": map[string]interface{}{"type": "string"},
					"weather": map[string]interface{}{"type": "string"},
					"visibility": map[string]interface{}{"type": "string"},
					"hazards": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{"type": "string"},
					},
				},
			},
		},
		"required": []string{"narrative", "game_mechanics"},
	}

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
			"temperature":       s.temperature,
			"maxOutputTokens":   8192,
			"topK":             s.topK,
			"topP":             s.topP,
			"responseMimeType": "application/json",
			"responseSchema":   responseSchema,
			"stopSequences":    []string{"[END_SCENE]", "[PLAYER_INPUT_REQUIRED]"},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:streamGenerateContent?key=%s",
		s.model, s.apiKey)

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Parse streaming response
	structuredResp, err := s.parseStructuredStreamResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to models.AIResponse
	return s.convertToAIResponse(structuredResp, ctx.SessionID), nil
}

func (s *EnhancedAIService) buildStructuredPrompt(ctx models.AIContext) string {
	var sb strings.Builder

	sb.WriteString("You are an expert Dungeon Master for D&D 5e. Generate a JSON response for the following game situation.\n\n")
	
	// Campaign context
	sb.WriteString(fmt.Sprintf("Campaign: %s\n", ctx.Campaign.Name))
	sb.WriteString(fmt.Sprintf("Setting: %s\n", ctx.Campaign.Description))
	sb.WriteString(fmt.Sprintf("Current Scene: %s\n\n", ctx.CurrentScene))
	
	// Character information
	sb.WriteString("Party Members:\n")
	for _, char := range ctx.Characters {
		sb.WriteString(fmt.Sprintf("- %s: Level %d %s %s, HP: %d/%d, AC: %d\n",
			char.Name, char.Level, char.Race, char.Class,
			char.CurrentHP, char.MaxHP, char.ArmorClass))
	}
	
	// Recent events
	if len(ctx.RecentEvents) > 0 {
		sb.WriteString("\nRecent Events:\n")
		for _, event := range ctx.RecentEvents {
			sb.WriteString(fmt.Sprintf("- %s\n", event.Description))
		}
	}
	
	sb.WriteString("\nJSON Response Requirements:\n")
	sb.WriteString("1. 'narrative': Engaging description of what happens\n")
	sb.WriteString("2. 'game_mechanics': Array of required game mechanics with specific details\n")
	sb.WriteString("3. 'npc_actions': Any NPC actions in combat\n")
	sb.WriteString("4. 'scene_changes': Environmental changes if any\n")
	sb.WriteString("\nFor each game mechanic, specify exact DCs, damage dice, and effects.\n")
	
	// Current action
	sb.WriteString("\n--- CURRENT ACTION ---\n")
	
	// Turn information
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
	
	sb.WriteString("\nGenerate a response that advances the story and provides clear game mechanics.")
	
	return sb.String()
}

func (s *EnhancedAIService) parseStructuredStreamResponse(resp *http.Response) (*StructuredAIResponse, error) {
	// Implementation would parse the SSE stream and extract the JSON response
	// This is simplified for demonstration
	
	// In reality, you'd parse the stream chunks and combine them
	// For now, assuming we get a complete JSON response
	
	var result StructuredAIResponse
	// Parse streaming response logic here...
	
	return &result, nil
}

func (s *EnhancedAIService) convertToAIResponse(structured *StructuredAIResponse, sessionID primitive.ObjectID) *models.AIResponse {
	mechanics := make([]models.GameMechanic, len(structured.GameMechanics))
	
	for i, sm := range structured.GameMechanics {
		mechanic := models.GameMechanic{
			Type:        sm.Type,
			Description: sm.Description,
			Target:      sm.Target,
			Metadata:    make(map[string]interface{}),
		}
		
		// Add requirements to metadata
		if sm.Requirements != nil {
			if sm.Requirements.DC > 0 {
				mechanic.Metadata["dc"] = sm.Requirements.DC
			}
			if sm.Requirements.Skill != "" {
				mechanic.Metadata["skill"] = sm.Requirements.Skill
			}
			if sm.Requirements.DamageDice != "" {
				mechanic.Metadata["damage_dice"] = sm.Requirements.DamageDice
				mechanic.Metadata["damage_type"] = sm.Requirements.DamageType
			}
			if sm.Requirements.Advantage {
				mechanic.Metadata["advantage"] = true
			}
			if sm.Requirements.Disadvantage {
				mechanic.Metadata["disadvantage"] = true
			}
		}
		
		// Add effects to metadata
		if sm.Effects != nil {
			if sm.Effects.OnSuccess != "" {
				mechanic.Metadata["on_success"] = sm.Effects.OnSuccess
			}
			if sm.Effects.OnFailure != "" {
				mechanic.Metadata["on_failure"] = sm.Effects.OnFailure
			}
			if len(sm.Effects.Conditions) > 0 {
				mechanic.Metadata["conditions"] = sm.Effects.Conditions
			}
			if sm.Effects.Duration != "" {
				mechanic.Metadata["duration"] = sm.Effects.Duration
			}
		}
		
		mechanics[i] = mechanic
	}
	
	return &models.AIResponse{
		SessionID:     sessionID,
		Narrative:     structured.Narrative,
		GameMechanics: mechanics,
		Timestamp:     time.Now(),
	}
}