#!/bin/bash

# Test AI integration with Google Gemini

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api"

echo "Testing AI DM Integration with Google Gemini..."
echo "=============================================="

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

# Get session ID (assuming we have one from previous tests)
SESSION_ID="68328ef9a2090fcc38604f94"

# Test AI action
echo -e "\n${GREEN}2. Testing AI DM response...${NC}"
AI_RESPONSE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "character_id": "68328eb4a2090fcc38604f93",
    "action": "I search the room for hidden doors",
    "action_type": "exploration"
  }')

echo "AI Response:"
echo $AI_RESPONSE | jq '.'

# Test AI action with dice roll
echo -e "\n${GREEN}3. Testing AI DM response with dice roll...${NC}"
AI_RESPONSE_WITH_DICE=$(curl -s -X POST "$API_URL/sessions/$SESSION_ID/action" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "character_id": "68328eb4a2090fcc38604f93",
    "action": "I attack the goblin with my sword",
    "action_type": "combat",
    "dice_result": {
      "dice": "1d20+5",
      "result": [18],
      "modifier": 5,
      "total": 23
    }
  }')

echo "AI Response with dice:"
echo $AI_RESPONSE_WITH_DICE | jq '.'

# Get narrative history
echo -e "\n${GREEN}4. Getting narrative history...${NC}"
NARRATIVE=$(curl -s -X GET "$API_URL/sessions/$SESSION_ID/narrative" \
  -H "Authorization: Bearer $TOKEN")

echo "Narrative history:"
echo $NARRATIVE | jq '.'

echo -e "\n${GREEN}Test completed!${NC}"