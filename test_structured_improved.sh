#!/bin/bash

# Improved structured output test with better schema
echo "=== Improved Structured Output Test ==="
echo ""

API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
MODEL="gemini-2.0-flash-exp"

# Test 1: Skill Check Scenario
echo "1. Testing Skill Check Scenario..."
cat > /tmp/skill_check.json << 'EOF'
{
  "contents": [{
    "role": "user", 
    "parts": [{
      "text": "D&D 5e: A rogue is trying to pick a lock on a treasure chest. Generate the scene and mechanics."
    }]
  }],
  "generationConfig": {
    "responseMimeType": "application/json",
    "responseSchema": {
      "type": "object",
      "properties": {
        "narrative": {"type": "string"},
        "mechanics": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "type": {"type": "string"},
              "target": {"type": "string"},
              "dc": {"type": "integer"},
              "skill": {"type": "string"},
              "onSuccess": {"type": "string"},
              "onFailure": {"type": "string"}
            }
          }
        }
      }
    }
  }
}
EOF

RESPONSE1=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d @/tmp/skill_check.json)

echo "Response:"
echo "$RESPONSE1" | jq -r '.candidates[0].content.parts[0].text' | jq '.'

# Test 2: Combat with Multiple Targets
echo ""
echo "2. Testing Combat with Multiple Targets..."
cat > /tmp/combat.json << 'EOF'
{
  "contents": [{
    "role": "user",
    "parts": [{
      "text": "D&D 5e combat: A wizard casts Lightning Bolt hitting 3 goblins in a line. Generate mechanics for each."
    }]
  }],
  "generationConfig": {
    "responseMimeType": "application/json",
    "responseSchema": {
      "type": "object",
      "properties": {
        "narrative": {"type": "string"},
        "targets": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "name": {"type": "string"},
              "saveDC": {"type": "integer"},
              "saveType": {"type": "string"},
              "baseDamage": {"type": "string"},
              "damageType": {"type": "string"}
            }
          }
        }
      }
    }
  }
}
EOF

RESPONSE2=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d @/tmp/combat.json)

echo "Response:"
echo "$RESPONSE2" | jq -r '.candidates[0].content.parts[0].text' | jq '.'

# Test 3: Complex Scene with Multiple Mechanics
echo ""
echo "3. Testing Complex Scene..."
cat > /tmp/complex.json << 'EOF'
{
  "contents": [{
    "role": "user",
    "parts": [{
      "text": "D&D 5e: The party enters a trapped hallway. There's a pressure plate that triggers poison darts. A character steps on it."
    }]
  }],
  "generationConfig": {
    "responseMimeType": "application/json",
    "responseSchema": {
      "type": "object",
      "properties": {
        "narrative": {"type": "string"},
        "trapDetection": {
          "type": "object",
          "properties": {
            "passivePerceptionDC": {"type": "integer"},
            "activeCheckDC": {"type": "integer"}
          }
        },
        "trapEffects": {
          "type": "object", 
          "properties": {
            "attackBonus": {"type": "integer"},
            "damage": {"type": "string"},
            "damageType": {"type": "string"},
            "additionalEffect": {"type": "string"}
          }
        }
      }
    }
  }
}
EOF

RESPONSE3=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d @/tmp/complex.json)

echo "Response:"
echo "$RESPONSE3" | jq -r '.candidates[0].content.parts[0].text' | jq '.'

# Cleanup
rm -f /tmp/skill_check.json /tmp/combat.json /tmp/complex.json

echo ""
echo "=== Summary ==="
echo "Structured output provides:"
echo "✓ Consistent JSON format"
echo "✓ Type-safe responses"
echo "✓ No parsing required"
echo "✓ Direct integration with game mechanics"