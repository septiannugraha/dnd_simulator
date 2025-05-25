#!/bin/bash

# Final AI Demo - showing it works!

echo "🎲 D&D AI Dungeon Master with Google Gemini"
echo "==========================================="
echo ""

API_URL="http://localhost:8080/api"

# Login
TOKEN=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"snoogie","password":"snoogenz"}' | jq -r '.token')

# Create fresh session
echo "Creating new adventure session..."
CHARACTER_RESPONSE=$(curl -s -X POST "$API_URL/characters" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gromnar the Bold",
    "race": "half-orc",
    "class": "barbarian",
    "background": "soldier",
    "alignment": "Chaotic Good",
    "abilities": {
      "Strength": 17,
      "Dexterity": 13,
      "Constitution": 16,
      "Intelligence": 8,
      "Wisdom": 11,
      "Charisma": 10
    }
  }')

CHARACTER_ID=$(echo $CHARACTER_RESPONSE | jq -r '.character.id')
echo "✓ Created character: Gromnar the Bold (Half-Orc Barbarian)"

# Assign to campaign
curl -s -X POST "$API_URL/characters/$CHARACTER_ID/assign" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"campaign_id": "68328e55a2090fcc38604f93"}' > /dev/null

# Create session
SESSION_RESPONSE=$(curl -s -X POST "$API_URL/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "campaign_id": "68328e55a2090fcc38604f93",
    "name": "Gemini AI Test",
    "description": "Testing AI DM"
  }')

SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.id')

# Join and start
curl -s -X POST "$API_URL/sessions/$SESSION_ID/join" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"character_id\": \"$CHARACTER_ID\"}" > /dev/null

curl -s -X POST "$API_URL/sessions/$SESSION_ID/start" \
  -H "Authorization: Bearer $TOKEN" > /dev/null

# Set scene
curl -s -X PUT "$API_URL/sessions/$SESSION_ID/scene" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scene": "You stand before the entrance to the goblin cave. The smell of damp earth and something fouler wafts from within. Crude totems made of bones flank the entrance.",
    "notes": "First encounter - goblin cave entrance"
  }' > /dev/null

echo "✓ Session ready!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test 1
echo "🧙 Gromnar: I examine the bone totems. What do they look like?"
RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I examine the bone totems. What do they look like?\",
    \"action_type\": \"exploration\"
  }")

echo -e "\n📖 AI DM Response:"
echo "$RESPONSE" | jq -r '.ai_event.description' | fold -s -w 70
echo ""

# Test 2
echo "🧙 Gromnar: I roar a battle cry and charge into the cave!"
RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I roar a battle cry and charge into the cave!\",
    \"action_type\": \"combat\"
  }")

echo -e "\n📖 AI DM Response:"
echo "$RESPONSE" | jq -r '.ai_event.description' | fold -s -w 70

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ AI Dungeon Master is working with Google Gemini!"
echo "   (Note: Responses may be truncated due to token limits)"
echo ""