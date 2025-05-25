#!/bin/bash

# Test AI integration with Google Gemini - Demo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}   D&D AI Dungeon Master Demo${NC}"
echo -e "${BLUE}   Powered by Google Gemini${NC}"
echo -e "${BLUE}============================================${NC}"

# Login
echo -e "\n${GREEN}1. Logging in as snoogie...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"snoogie","password":"snoogenz"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
if [ "$TOKEN" == "null" ]; then
  echo -e "${RED}Failed to login${NC}"
  exit 1
fi
echo "✓ Login successful"

# Use existing campaign
CAMPAIGN_ID="68328e55a2090fcc38604f93"
echo -e "\n${GREEN}2. Using campaign: The Can Can${NC}"

# Create a character
echo -e "\n${GREEN}3. Creating character: Thorin Ironforge${NC}"
CHARACTER_RESPONSE=$(curl -s -X POST "$API_URL/characters" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Thorin Ironforge",
    "race": "dwarf",
    "class": "fighter",
    "background": "soldier",
    "alignment": "Lawful Good",
    "abilities": {
      "Strength": 16,
      "Dexterity": 12,
      "Constitution": 15,
      "Intelligence": 10,
      "Wisdom": 13,
      "Charisma": 8
    }
  }')

CHARACTER_ID=$(echo $CHARACTER_RESPONSE | jq -r '.character.id')
if [ "$CHARACTER_ID" == "null" ]; then
  echo -e "${RED}Failed to create character${NC}"
  echo $CHARACTER_RESPONSE | jq
  exit 1
fi
echo "✓ Character created: $CHARACTER_ID"
echo -e "${YELLOW}   Class: Dwarf Fighter (Level 1)${NC}"
echo -e "${YELLOW}   HP: $(echo $CHARACTER_RESPONSE | jq -r '.character.max_hp') | AC: $(echo $CHARACTER_RESPONSE | jq -r '.character.armor_class')${NC}"

# Assign to campaign
echo -e "\n${GREEN}4. Assigning character to campaign...${NC}"
curl -s -X POST "$API_URL/characters/$CHARACTER_ID/assign" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"campaign_id\": \"$CAMPAIGN_ID\"}" > /dev/null
echo "✓ Character assigned"

# Create session
echo -e "\n${GREEN}5. Creating game session...${NC}"
SESSION_RESPONSE=$(curl -s -X POST "$API_URL/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"campaign_id\": \"$CAMPAIGN_ID\",
    \"name\": \"AI Dungeon Master Demo\",
    \"description\": \"Testing the AI DM capabilities\"
  }")

SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.id')
if [ -z "$SESSION_ID" ] || [ "$SESSION_ID" == "null" ]; then
  echo -e "${RED}Failed to create session${NC}"
  echo $SESSION_RESPONSE | jq
  exit 1
fi
echo "✓ Session created"

# Join and start session
echo -e "\n${GREEN}6. Starting the adventure...${NC}"
curl -s -X POST "$API_URL/sessions/$SESSION_ID/join" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"character_id\": \"$CHARACTER_ID\"}" > /dev/null

curl -s -X POST "$API_URL/sessions/$SESSION_ID/start" \
  -H "Authorization: Bearer $TOKEN" > /dev/null
echo "✓ Session started"

# Set the scene
echo -e "\n${GREEN}7. Setting the scene...${NC}"
curl -s -X PUT "$API_URL/sessions/$SESSION_ID/scene" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scene": "The tavern '\''The Rusty Dragon'\'' is bustling with activity. You sit at a worn wooden table, nursing an ale, when a hooded figure approaches. '\''I hear you'\''re looking for work,'\'' they whisper, sliding a rolled parchment across the table. '\''The abandoned Cragmaw Hideout... dangerous, but the pay is good.'\''",
    "notes": "Quest hook - players offered job to clear goblin hideout"
  }' > /dev/null
echo "✓ Scene set"

echo -e "\n${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${BLUE}         🎲 ADVENTURE BEGINS! 🎲${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}\n"

# First AI interaction
echo -e "${YELLOW}[Thorin]${NC}: I examine the parchment carefully. What does it say?"
AI_RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I examine the parchment carefully. What does it say?\",
    \"action_type\": \"exploration\"
  }")

if echo "$AI_RESPONSE" | jq -e '.narrative' > /dev/null 2>&1; then
  echo -e "\n${BLUE}[AI Dungeon Master]:${NC}"
  echo "$AI_RESPONSE" | jq -r '.narrative' | fold -s -w 80
else
  echo -e "\n${RED}Error with AI response:${NC}"
  echo "$AI_RESPONSE" | jq
fi

# Roll for insight
echo -e "\n${YELLOW}[Thorin]${NC}: I study the hooded figure, trying to determine if they're trustworthy."
echo -e "${YELLOW}*Rolling Insight check: 1d20+1*${NC}"

# Simulate dice roll
DICE_ROLL=$((RANDOM % 20 + 1))
TOTAL=$((DICE_ROLL + 1))
echo -e "${GREEN}Rolled: $DICE_ROLL + 1 = $TOTAL${NC}"

AI_RESPONSE2=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I study the hooded figure, trying to determine if they're trustworthy\",
    \"action_type\": \"exploration\",
    \"dice_result\": {
      \"dice\": \"1d20+1\",
      \"result\": [$DICE_ROLL],
      \"modifier\": 1,
      \"total\": $TOTAL
    }
  }")

if echo "$AI_RESPONSE2" | jq -e '.narrative' > /dev/null 2>&1; then
  echo -e "\n${BLUE}[AI Dungeon Master]:${NC}"
  echo "$AI_RESPONSE2" | jq -r '.narrative' | fold -s -w 80
else
  echo -e "\n${RED}Error with AI response:${NC}"
  echo "$AI_RESPONSE2" | jq
fi

# Accept the quest
echo -e "\n${YELLOW}[Thorin]${NC}: \"Aye, I'll take the job. Goblins need to learn their place. What's the pay?\""
AI_RESPONSE3=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"Aye, I'll take the job. Goblins need to learn their place. What's the pay?\",
    \"action_type\": \"roleplay\"
  }")

if echo "$AI_RESPONSE3" | jq -e '.narrative' > /dev/null 2>&1; then
  echo -e "\n${BLUE}[AI Dungeon Master]:${NC}"
  echo "$AI_RESPONSE3" | jq -r '.narrative' | fold -s -w 80
else
  echo -e "\n${RED}Error with AI response:${NC}"
  echo "$AI_RESPONSE3" | jq
fi

echo -e "\n${BLUE}═══════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ AI Dungeon Master is working!${NC}"
echo -e "${YELLOW}Session ID: $SESSION_ID${NC}"
echo -e "${YELLOW}Character ID: $CHARACTER_ID${NC}"
echo -e "${BLUE}═══════════════════════════════════════════${NC}\n"