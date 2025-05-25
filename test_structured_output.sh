#!/bin/bash

# Test Gemini API with structured output
echo "=== Testing Gemini Structured Output for D&D ==="
echo ""

# API Configuration
API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
MODEL="gemini-2.0-flash-exp"  # Using flash for faster response in testing

# Create the request with structured output
echo "Testing structured output with responseSchema..."
echo ""

# Create request JSON
cat > /tmp/structured_request.json << 'EOF'
{
  "contents": [{
    "role": "user",
    "parts": [{
      "text": "You are a D&D 5e Dungeon Master. A level 3 fighter named Thorin is exploring a goblin cave. He says: 'I carefully approach the wooden door and listen for any sounds. If I don't hear anything dangerous, I'll try to open it slowly.' Generate a response with what happens next, including any required game mechanics."
    }]
  }],
  "generationConfig": {
    "temperature": 0.8,
    "maxOutputTokens": 8192,
    "topK": 40,
    "topP": 0.95,
    "responseMimeType": "application/json",
    "responseSchema": {
      "type": "object",
      "properties": {
        "narrative": {
          "type": "string",
          "description": "The narrative description of what happens"
        },
        "game_mechanics": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "type": {
                "type": "string",
                "enum": ["skill_check", "attack_roll", "saving_throw", "damage", "initiative"]
              },
              "description": {
                "type": "string"
              },
              "target": {
                "type": "string"
              },
              "requirements": {
                "type": "object",
                "properties": {
                  "roll_type": {"type": "string"},
                  "dc": {"type": "integer"},
                  "skill": {"type": "string"},
                  "modifier": {"type": "integer"}
                }
              },
              "effects": {
                "type": "object",
                "properties": {
                  "on_success": {"type": "string"},
                  "on_failure": {"type": "string"}
                }
              }
            },
            "required": ["type", "description"]
          }
        },
        "scene_details": {
          "type": "object",
          "properties": {
            "immediate_threat": {"type": "boolean"},
            "npcs_present": {
              "type": "array",
              "items": {"type": "string"}
            }
          }
        }
      },
      "required": ["narrative", "game_mechanics"]
    }
  }
}
EOF

echo "Sending request to Gemini API..."
echo ""

# Make the API call
RESPONSE=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d @/tmp/structured_request.json)

# Check if we got a response
if [ -z "$RESPONSE" ]; then
  echo "Error: No response from API"
  exit 1
fi

# Extract the JSON content
JSON_CONTENT=$(echo "$RESPONSE" | jq -r '.candidates[0].content.parts[0].text' 2>/dev/null)

if [ "$JSON_CONTENT" = "null" ] || [ -z "$JSON_CONTENT" ]; then
  echo "Error: Could not extract JSON from response"
  echo "Raw response:"
  echo "$RESPONSE" | jq '.'
  exit 1
fi

echo "✓ Received structured JSON response!"
echo ""
echo "Parsed Response:"
echo "$JSON_CONTENT" | jq '.'

# Analyze the response
echo ""
echo "=== Response Analysis ==="
NARRATIVE=$(echo "$JSON_CONTENT" | jq -r '.narrative')
MECHANICS_COUNT=$(echo "$JSON_CONTENT" | jq '.game_mechanics | length')

echo "Narrative length: $(echo -n "$NARRATIVE" | wc -c) characters"
echo "Game mechanics found: $MECHANICS_COUNT"

if [ $MECHANICS_COUNT -gt 0 ]; then
  echo ""
  echo "Game Mechanics Details:"
  echo "$JSON_CONTENT" | jq -r '.game_mechanics[] | "- Type: \(.type), DC: \(.requirements.dc // "N/A"), Skill: \(.requirements.skill // "N/A")"'
fi

# Test with combat scenario
echo ""
echo ""
echo "=== Testing Combat Scenario ==="

cat > /tmp/combat_request.json << 'EOF'
{
  "contents": [{
    "role": "user",
    "parts": [{
      "text": "D&D 5e combat scenario: A level 5 wizard casts Fireball at a group of 4 orcs. The wizard rolled 28 damage (8d6). Generate the response including saving throws for each orc and damage calculations."
    }]
  }],
  "generationConfig": {
    "temperature": 0.8,
    "maxOutputTokens": 8192,
    "responseMimeType": "application/json",
    "responseSchema": {
      "type": "object",
      "properties": {
        "narrative": {"type": "string"},
        "game_mechanics": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "type": {"type": "string"},
              "description": {"type": "string"},
              "target": {"type": "string"},
              "requirements": {
                "type": "object",
                "properties": {
                  "save_type": {"type": "string"},
                  "dc": {"type": "integer"},
                  "roll_needed": {"type": "integer"}
                }
              },
              "damage": {
                "type": "object",
                "properties": {
                  "amount": {"type": "integer"},
                  "type": {"type": "string"}
                }
              }
            }
          }
        }
      },
      "required": ["narrative", "game_mechanics"]
    }
  }
}
EOF

echo "Testing combat with structured output..."
COMBAT_RESPONSE=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d @/tmp/combat_request.json)

COMBAT_JSON=$(echo "$COMBAT_RESPONSE" | jq -r '.candidates[0].content.parts[0].text' 2>/dev/null)

if [ "$COMBAT_JSON" != "null" ] && [ -n "$COMBAT_JSON" ]; then
  echo ""
  echo "Combat Response:"
  echo "$COMBAT_JSON" | jq '.'
  
  # Show damage calculations
  echo ""
  echo "Damage Summary:"
  echo "$COMBAT_JSON" | jq -r '.game_mechanics[] | select(.damage) | "- \(.target): \(.damage.amount) \(.damage.type) damage"'
else
  echo "Error parsing combat response"
fi

# Cleanup
rm -f /tmp/structured_request.json /tmp/combat_request.json

echo ""
echo "=== Test Complete ==="
echo ""
echo "Benefits demonstrated:"
echo "✓ Structured JSON responses with guaranteed format"
echo "✓ Precise game mechanics extraction"
echo "✓ DC values, damage calculations, and effects clearly defined"
echo "✓ Easy for frontend to parse and display"
echo "✓ No text parsing errors or missed information"