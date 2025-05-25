#!/bin/bash

echo "=== Testing Enhanced AI Output ==="
echo ""

# Direct API test with both configurations
API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
MODEL="gemini-2.0-flash-exp"

echo "1. Testing with LIMITED tokens (1000)..."
LIMITED=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [{
      "role": "user",
      "parts": [{
        "text": "As a D&D Dungeon Master, describe an epic battle scene between a party of adventurers and a red dragon in its lair. Include detailed descriptions of the environment, the dragon'"'"'s attacks, and the heroes'"'"' actions."
      }]
    }],
    "generationConfig": {
      "temperature": 0.8,
      "maxOutputTokens": 1000
    }
  }' | jq -r '.candidates[0].content.parts[0].text')

echo "Limited response (1000 tokens):"
echo "$LIMITED" | head -5
echo "..."
echo "Length: $(echo -n "$LIMITED" | wc -c) characters"
echo ""

echo "2. Testing with ENHANCED tokens (8192)..."
ENHANCED=$(curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [{
      "role": "user", 
      "parts": [{
        "text": "As a D&D Dungeon Master, describe an epic battle scene between a party of adventurers and a red dragon in its lair. Include detailed descriptions of the environment, the dragon'"'"'s attacks, and the heroes'"'"' actions. Make this a long, immersive narrative."
      }]
    }],
    "generationConfig": {
      "temperature": 0.8,
      "maxOutputTokens": 8192,
      "topK": 40,
      "topP": 0.95
    }
  }' | jq -r '.candidates[0].content.parts[0].text')

echo "Enhanced response (8192 tokens):"
echo "$ENHANCED" | head -5
echo "..."
echo "Length: $(echo -n "$ENHANCED" | wc -c) characters"
echo ""

echo "=== Comparison ==="
echo "Limited (1000 tokens):  $(echo -n "$LIMITED" | wc -c) characters"
echo "Enhanced (8192 tokens): $(echo -n "$ENHANCED" | wc -c) characters"
echo ""

# Check if responses were cut off
if [[ "$LIMITED" == *"..."* ]] || [[ $(echo -n "$LIMITED" | tail -c 20) == *[.!?]* ]]; then
  echo "Limited response appears complete"
else
  echo "Limited response may be truncated"
fi

if [[ "$ENHANCED" == *"..."* ]] || [[ $(echo -n "$ENHANCED" | tail -c 20) == *[.!?]* ]]; then
  echo "Enhanced response appears complete"
else
  echo "Enhanced response may be truncated"
fi