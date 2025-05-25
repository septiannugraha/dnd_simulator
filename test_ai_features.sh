#!/bin/bash

# Test AI DM features

echo "Testing AI DM Features"
echo "====================="

# Login as player
echo -e "\n1. Login as player..."
PLAYER_LOGIN=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "player1", "password": "password123"}')

PLAYER_TOKEN=$(echo $PLAYER_LOGIN | jq -r '.token')
echo "Player Token obtained"

# Create a character with proper abilities
echo -e "\n2. Creating a character..."
CHARACTER_CREATE=$(curl -s -X POST "http://localhost:8080/api/characters" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Aragorn Strider",
    "race": "human",
    "class": "ranger",
    "background": "soldier",
    "abilities": {
      "strength": 15,
      "dexterity": 16,
      "constitution": 14,
      "intelligence": 12,
      "wisdom": 15,
      "charisma": 10
    },
    "alignment": "Chaotic Good"
  }')

echo $CHARACTER_CREATE | jq

CHARACTER_ID=$(echo $CHARACTER_CREATE | jq -r '.id')
echo "Character ID: $CHARACTER_ID"

if [ "$CHARACTER_ID" != "null" ]; then
  # Join an existing session (using the one we created in previous test)
  SESSION_ID="6832129da02c259d35d1d2e7"
  
  echo -e "\n3. Joining session..."
  curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/join" \
    -H "Authorization: Bearer $PLAYER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"character_id\": \"$CHARACTER_ID\"}" | jq

  echo -e "\n4. Testing AI DM - Player action (will fail without OpenAI key)..."
  AI_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/action" \
    -H "Authorization: Bearer $PLAYER_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"character_id\": \"$CHARACTER_ID\",
      \"action\": \"I cautiously approach the cave entrance and listen for any sounds from within.\",
      \"action_type\": \"exploration\"
    }")
  
  echo $AI_RESPONSE | jq

  echo -e "\n5. Get narrative history to see if action was recorded..."
  curl -s -X GET "http://localhost:8080/api/sessions/$SESSION_ID/narrative" \
    -H "Authorization: Bearer $PLAYER_TOKEN" | jq

  echo -e "\n6. Get player_action events..."
  curl -s -X GET "http://localhost:8080/api/sessions/$SESSION_ID/events/player_action" \
    -H "Authorization: Bearer $PLAYER_TOKEN" | jq

  echo -e "\n7. Testing dice roll..."
  DICE_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/dice" \
    -H "Authorization: Bearer $PLAYER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "dice": "1d20+5",
      "purpose": "Perception check",
      "character_id": "'$CHARACTER_ID'"
    }')
  
  echo $DICE_RESPONSE | jq

  echo -e "\n8. Get WebSocket connection status..."
  curl -s -X GET "http://localhost:8080/api/sessions/$SESSION_ID/ws/status" \
    -H "Authorization: Bearer $PLAYER_TOKEN" | jq
fi