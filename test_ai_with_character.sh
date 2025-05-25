#!/bin/bash

# Test AI integration with Google Gemini - Full workflow

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api"

echo "Testing AI DM Integration with Google Gemini - Full Workflow"
echo "==========================================================="

# First, get auth token
echo -e "\n${GREEN}1. Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"snoogie","password":"snoogenz"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
if [ "$TOKEN" == "null" ]; then
  echo -e "${RED}Failed to login${NC}"
  echo $LOGIN_RESPONSE
  exit 1
fi
echo "Token obtained successfully"

# Get campaign ID
CAMPAIGN_ID="6834c84fe74f674dd17e5ff0"

# Create a character for testing
echo -e "\n${GREEN}2. Creating test character...${NC}"
CHARACTER_RESPONSE=$(curl -s -X POST "$API_URL/characters" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Thorin Ironforge",
    "race": "Dwarf",
    "class": "Fighter",
    "level": 3,
    "background": "Soldier",
    "alignment": "Lawful Good",
    "stats": {
      "strength": 16,
      "dexterity": 12,
      "constitution": 15,
      "intelligence": 10,
      "wisdom": 13,
      "charisma": 8
    }
  }')

CHARACTER_ID=$(echo $CHARACTER_RESPONSE | jq -r '.character.id')
if [ "$CHARACTER_ID" == "null" ]; then
  echo -e "${RED}Failed to create character${NC}"
  echo $CHARACTER_RESPONSE
else
  echo "Character created: $CHARACTER_ID"
  echo $CHARACTER_RESPONSE | jq '.character'
fi

# Assign character to campaign
echo -e "\n${GREEN}3. Assigning character to campaign...${NC}"
ASSIGN_RESPONSE=$(curl -s -X POST "$API_URL/characters/$CHARACTER_ID/assign" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"campaign_id\": \"$CAMPAIGN_ID\"}")

echo $ASSIGN_RESPONSE | jq

# Create a new session
echo -e "\n${GREEN}4. Creating new game session...${NC}"
SESSION_RESPONSE=$(curl -s -X POST "$API_URL/sessions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"campaign_id\": \"$CAMPAIGN_ID\",
    \"name\": \"AI Test Session\",
    \"description\": \"Testing AI DM functionality\"
  }")

SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.session.id')
if [ "$SESSION_ID" == "null" ]; then
  echo -e "${RED}Failed to create session${NC}"
  echo $SESSION_RESPONSE
  exit 1
fi
echo "Session created: $SESSION_ID"

# Join session with character
echo -e "\n${GREEN}5. Joining session with character...${NC}"
JOIN_RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/join" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"character_id\": \"$CHARACTER_ID\"}")

echo $JOIN_RESPONSE | jq

# Start the session
echo -e "\n${GREEN}6. Starting session...${NC}"
START_RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/start" \
  -H "Authorization: Bearer $TOKEN")

echo $START_RESPONSE | jq

# Set the scene
echo -e "\n${GREEN}7. Setting the scene...${NC}"
SCENE_RESPONSE=$(curl -s -X PUT "$API_URL/sessions/$SESSION_ID/scene" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scene": "You stand at the entrance of the Lost Mine of Phandelver. The cave mouth yawns before you, dark and foreboding. You can hear the faint sound of dripping water echoing from within.",
    "notes": "Players are about to enter the goblin cave"
  }')

echo $SCENE_RESPONSE | jq

# Test AI action - Exploration
echo -e "\n${BLUE}8. Testing AI DM response - Exploration...${NC}"
AI_RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I light a torch and cautiously enter the cave, looking for any signs of danger\",
    \"action_type\": \"exploration\"
  }")

echo -e "${BLUE}AI Response:${NC}"
echo $AI_RESPONSE | jq '.'

# Test AI action with dice roll - Perception check
echo -e "\n${BLUE}9. Testing AI DM response - Perception check...${NC}"
AI_RESPONSE_PERCEPTION=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I carefully examine the cave entrance for traps\",
    \"action_type\": \"skill_check\",
    \"dice_result\": {
      \"dice\": \"1d20+1\",
      \"result\": [14],
      \"modifier\": 1,
      \"total\": 15
    }
  }")

echo -e "${BLUE}AI Response:${NC}"
echo $AI_RESPONSE_PERCEPTION | jq '.'

# Test AI action - Combat
echo -e "\n${BLUE}10. Testing AI DM response - Combat...${NC}"
AI_RESPONSE_COMBAT=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I attack the goblin with my warhammer\",
    \"action_type\": \"combat\",
    \"target\": \"goblin\",
    \"dice_result\": {
      \"dice\": \"1d20+5\",
      \"result\": [17],
      \"modifier\": 5,
      \"total\": 22
    }
  }")

echo -e "${BLUE}AI Response:${NC}"
echo $AI_RESPONSE_COMBAT | jq '.'

# Get narrative history
echo -e "\n${GREEN}11. Getting narrative history...${NC}"
NARRATIVE=$(curl -s -X GET "$API_URL/sessions/$SESSION_ID/narrative" \
  -H "Authorization: Bearer $TOKEN")

echo -e "${BLUE}Narrative history:${NC}"
echo $NARRATIVE | jq '.'

# Get events by type
echo -e "\n${GREEN}12. Getting combat events...${NC}"
EVENTS=$(curl -s -X GET "$API_URL/sessions/$SESSION_ID/events/combat" \
  -H "Authorization: Bearer $TOKEN")

echo -e "${BLUE}Combat events:${NC}"
echo $EVENTS | jq '.'

echo -e "\n${GREEN}Test completed!${NC}"
echo -e "Session ID: ${BLUE}$SESSION_ID${NC}"
echo -e "Character ID: ${BLUE}$CHARACTER_ID${NC}"