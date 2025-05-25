#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api"

echo -e "${BLUE}Testing D&D Simulator - New Features${NC}"
echo "======================================="

# Step 1: Register and login as DM
echo -e "\n${BLUE}1. Registering DM user...${NC}"
DM_REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dm_test",
    "email": "dm@test.com",
    "password": "password123"
  }')
echo $DM_REGISTER_RESPONSE

echo -e "\n${BLUE}2. Logging in as DM...${NC}"
DM_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "dm_test",
    "password": "password123"
  }')
DM_TOKEN=$(echo $DM_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
echo "DM Token: ${DM_TOKEN:0:20}..."

# Step 2: Create a campaign
echo -e "\n${BLUE}3. Creating campaign...${NC}"
CAMPAIGN_RESPONSE=$(curl -s -X POST "$BASE_URL/campaigns" \
  -H "Authorization: Bearer $DM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "The Lost Mines",
    "description": "Adventure in Phandalin",
    "setting": "fantasy",
    "world_info": "A classic D&D setting with goblins, dragons, and ancient ruins.",
    "is_public": true
  }')
CAMPAIGN_ID=$(echo $CAMPAIGN_RESPONSE | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo "Campaign ID: $CAMPAIGN_ID"

# Step 3: Create a game session
echo -e "\n${BLUE}4. Creating game session...${NC}"
SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions" \
  -H "Authorization: Bearer $DM_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"campaign_id\": \"$CAMPAIGN_ID\",
    \"name\": \"Session 1: Goblin Ambush\",
    \"description\": \"The party travels to Phandalin\"
  }")
SESSION_ID=$(echo $SESSION_RESPONSE | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo "Session ID: $SESSION_ID"

# Step 4: Register and login as player
echo -e "\n${BLUE}5. Registering player...${NC}"
PLAYER_REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@test.com",
    "password": "password123"
  }')

PLAYER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "password": "password123"
  }')
PLAYER_TOKEN=$(echo $PLAYER_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
echo "Player Token: ${PLAYER_TOKEN:0:20}..."

# Step 5: Create a character
echo -e "\n${BLUE}6. Creating character...${NC}"
CHARACTER_RESPONSE=$(curl -s -X POST "$BASE_URL/characters" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Thorin Ironforge",
    "race": "Dwarf",
    "class": "Fighter",
    "background": "Soldier",
    "abilities": {
      "strength": 16,
      "dexterity": 12,
      "constitution": 15,
      "intelligence": 10,
      "wisdom": 13,
      "charisma": 8
    },
    "alignment": "Lawful Good"
  }')
CHARACTER_ID=$(echo $CHARACTER_RESPONSE | grep -o '"id":"[^"]*' | sed 's/"id":"//')
echo "Character ID: $CHARACTER_ID"

# Step 6: Join campaign
echo -e "\n${BLUE}7. Joining campaign...${NC}"
curl -s -X POST "$BASE_URL/campaigns/$CAMPAIGN_ID/join" \
  -H "Authorization: Bearer $PLAYER_TOKEN"

# Step 7: Assign character to campaign
echo -e "\n${BLUE}8. Assigning character to campaign...${NC}"
curl -s -X POST "$BASE_URL/characters/$CHARACTER_ID/assign" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"campaign_id\": \"$CAMPAIGN_ID\"}"

# Step 8: Join session
echo -e "\n${BLUE}9. Joining session...${NC}"
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/join" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"character_id\": \"$CHARACTER_ID\"}"

# Step 9: Start session (as DM)
echo -e "\n${BLUE}10. Starting session...${NC}"
curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/start" \
  -H "Authorization: Bearer $DM_TOKEN"

# Step 10: Update scene (as DM)
echo -e "\n${BLUE}11. Setting the scene...${NC}"
SCENE_RESPONSE=$(curl -s -X PUT "$BASE_URL/sessions/$SESSION_ID/scene" \
  -H "Authorization: Bearer $DM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "scene": "You are traveling along the Triboar Trail when suddenly goblins leap from the bushes!",
    "notes": "Goblin ambush encounter - 4 goblins"
  }')
echo $SCENE_RESPONSE

# Step 11: Test AI DM - Player action (Note: This will fail without OpenAI key)
echo -e "\n${BLUE}12. Testing AI DM - Player action...${NC}"
AI_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $PLAYER_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"character_id\": \"$CHARACTER_ID\",
    \"action\": \"I draw my sword and charge at the nearest goblin!\",
    \"action_type\": \"combat\"
  }")
echo $AI_RESPONSE

# Step 12: Get session status
echo -e "\n${BLUE}13. Getting session status...${NC}"
STATUS_RESPONSE=$(curl -s -X GET "$BASE_URL/sessions/$SESSION_ID/status" \
  -H "Authorization: Bearer $PLAYER_TOKEN")
echo $STATUS_RESPONSE | jq '.' 2>/dev/null || echo $STATUS_RESPONSE

# Step 13: Get narrative history
echo -e "\n${BLUE}14. Getting narrative history...${NC}"
NARRATIVE_RESPONSE=$(curl -s -X GET "$BASE_URL/sessions/$SESSION_ID/narrative" \
  -H "Authorization: Bearer $PLAYER_TOKEN")
echo $NARRATIVE_RESPONSE | jq '.' 2>/dev/null || echo $NARRATIVE_RESPONSE

# Step 14: Test session management - Pause
echo -e "\n${BLUE}15. Pausing session...${NC}"
PAUSE_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/pause" \
  -H "Authorization: Bearer $DM_TOKEN")
echo $PAUSE_RESPONSE

# Step 15: Resume session
echo -e "\n${BLUE}16. Resuming session...${NC}"
RESUME_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/resume" \
  -H "Authorization: Bearer $DM_TOKEN")
echo $RESUME_RESPONSE

# Step 16: End session
echo -e "\n${BLUE}17. Ending session...${NC}"
END_RESPONSE=$(curl -s -X POST "$BASE_URL/sessions/$SESSION_ID/end" \
  -H "Authorization: Bearer $DM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }')
echo $END_RESPONSE

echo -e "\n${GREEN}Testing complete!${NC}"
echo "Campaign ID: $CAMPAIGN_ID"
echo "Session ID: $SESSION_ID"
echo "Character ID: $CHARACTER_ID"