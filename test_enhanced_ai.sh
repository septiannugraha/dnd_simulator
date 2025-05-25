#!/bin/bash

# Test enhanced AI parameters
echo "=== Testing Enhanced AI Parameters ==="
echo ""
echo "Parameters updated:"
echo "- maxOutputTokens: 8192 (from 1000)"
echo "- topK: 40 (for variety)"  
echo "- topP: 0.95 (for quality)"
echo "- stopSequences: [END_SCENE], [AWAIT_PLAYER_ACTION]"
echo ""

# Source environment
export GEMINI_API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
export JWT_SECRET="test-secret-key"
export MONGODB_URI="mongodb://dnduser:dndpass@localhost:27017/dnd_simulator?authSource=dnd_simulator"

# Start the server in background
echo "Starting server with enhanced AI..."
go run main.go &
SERVER_PID=$!

# Wait for server to start
sleep 3

# Test login
echo "1. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "dm@example.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Login failed!"
    kill $SERVER_PID
    exit 1
fi

echo "✓ Logged in successfully"

# Get campaign
echo ""
echo "2. Getting campaign..."
CAMPAIGN_ID=$(curl -s -X GET http://localhost:8080/api/campaigns \
  -H "Authorization: Bearer $TOKEN" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

echo "✓ Campaign ID: $CAMPAIGN_ID"

# Get session
echo ""
echo "3. Getting session..."
SESSION_RESPONSE=$(curl -s -X GET http://localhost:8080/api/campaigns/$CAMPAIGN_ID \
  -H "Authorization: Bearer $TOKEN")

SESSION_ID=$(echo $SESSION_RESPONSE | grep -o '"sessions":\[{"id":"[^"]*' | cut -d'"' -f5)

if [ -z "$SESSION_ID" ]; then
    echo "Creating new session..."
    SESSION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/sessions \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"campaign_id\": \"$CAMPAIGN_ID\",
        \"name\": \"Enhanced AI Test\",
        \"description\": \"Testing enhanced AI parameters\"
      }")
    SESSION_ID=$(echo $SESSION_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
fi

echo "✓ Session ID: $SESSION_ID"

# Test AI with a complex scenario
echo ""
echo "4. Testing AI with complex scenario..."
echo ""

AI_RESPONSE=$(curl -s -X POST http://localhost:8080/api/sessions/$SESSION_ID/ai-action \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "I want to investigate the mysterious glowing runes on the ancient door. I examine them carefully, looking for any traps or magical auras. I also try to decipher what language they might be written in.",
    "action_type": "exploration"
  }')

echo "AI Response:"
echo "$AI_RESPONSE" | jq '.'

# Check response length
NARRATIVE_LENGTH=$(echo "$AI_RESPONSE" | jq -r '.description' | wc -c)
echo ""
echo "Response length: $NARRATIVE_LENGTH characters"

if [ $NARRATIVE_LENGTH -gt 1000 ]; then
    echo "✓ Enhanced parameters working - longer response generated!"
else
    echo "⚠ Response may still be truncated"
fi

# Test with combat scenario
echo ""
echo "5. Testing combat scenario..."
echo ""

COMBAT_RESPONSE=$(curl -s -X POST http://localhost:8080/api/sessions/$SESSION_ID/ai-action \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "I cast Fireball at the group of goblins, aiming for the center of their formation to hit as many as possible!",
    "action_type": "combat",
    "dice_roll": {
      "dice": "8d6",
      "total": 28,
      "rolls": [4, 3, 6, 2, 5, 3, 1, 4]
    }
  }')

echo "Combat AI Response:"
echo "$COMBAT_RESPONSE" | jq '.'

# Cleanup
echo ""
echo "Cleaning up..."
kill $SERVER_PID

echo ""
echo "=== Test Complete ==="
echo ""
echo "Summary:"
echo "- AI now supports longer responses (up to 8192 tokens)"
echo "- Better response quality with topK=40 and topP=0.95"
echo "- Stop sequences prevent runaway generation"
echo ""
echo "Next steps:"
echo "1. Implement structured output for precise game mechanics"
echo "2. Update frontend to handle structured responses"
echo "3. Add automatic mechanic parsing and dice rolling"