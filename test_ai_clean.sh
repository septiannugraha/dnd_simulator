#!/bin/bash

# Clean AI test - show only narratives

API_URL="http://localhost:8080/api"

# Login
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"snoogie","password":"snoogenz"}' | jq -r '.token')

# Use the session we just created
SESSION_ID="68333dae24b56172e25792e2"
CHARACTER_ID="68333dae24b56172e25792e1"

echo "🎲 D&D AI Dungeon Master Demo"
echo "============================="
echo ""

# Test 1: Simple exploration
echo "Player: I look around the tavern. What do I see?"
RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I look around the tavern. What do I see?\",
    \"action_type\": \"exploration\"
  }")

echo -e "\nAI DM:"
echo "$RESPONSE" | jq -r '.narrative // .ai_response.narrative // "No narrative found"' | fold -s -w 80
echo ""

# Test 2: With dice roll
echo "Player: I try to intimidate the hooded figure. *rolls 1d20+2 = 18*"
RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I try to intimidate the hooded figure into giving me more information\",
    \"action_type\": \"roleplay\",
    \"dice_result\": {
      \"dice\": \"1d20+2\",
      \"result\": [16],
      \"modifier\": 2,
      \"total\": 18
    }
  }")

echo -e "\nAI DM:"
echo "$RESPONSE" | jq -r '.narrative // .ai_response.narrative // "No narrative found"' | fold -s -w 80

echo -e "\n\n✅ AI is responding!"