#!/bin/bash

# Simple test for enhanced AI parameters
echo "=== Testing Enhanced AI with Longer Responses ==="
echo ""

# Start MongoDB if needed
docker-compose -f docker-compose.dev.yml up -d

# Wait for MongoDB
sleep 3

# Source environment
export GEMINI_API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
export JWT_SECRET="test-secret-key"
export MONGODB_URI="mongodb://dnduser:dndpass@localhost:27017/dnd_simulator?authSource=dnd_simulator"

# Start server
echo "Starting server..."
go run main.go &
SERVER_PID=$!

# Wait for server
sleep 5

# Quick test
echo "1. Testing login..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "dm@example.com", "password": "password123"}' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "Login failed, creating user..."
  curl -s -X POST http://localhost:8080/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{"username": "dm", "email": "dm@example.com", "password": "password123"}'
  
  TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email": "dm@example.com", "password": "password123"}' | jq -r '.token')
fi

echo "✓ Token obtained"

# Get campaign
CAMPAIGN_ID=$(curl -s -X GET http://localhost:8080/api/campaigns \
  -H "Authorization: Bearer $TOKEN" | jq -r '.[0].id')

if [ -z "$CAMPAIGN_ID" ] || [ "$CAMPAIGN_ID" = "null" ]; then
  echo "Creating campaign..."
  CAMPAIGN_ID=$(curl -s -X POST http://localhost:8080/api/campaigns \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name": "Test Campaign", "description": "Testing AI"}' | jq -r '.id')
fi

echo "✓ Campaign ID: $CAMPAIGN_ID"

# Get or create session
SESSION_ID=$(curl -s -X GET http://localhost:8080/api/campaigns/$CAMPAIGN_ID \
  -H "Authorization: Bearer $TOKEN" | jq -r '.sessions[0].id')

if [ -z "$SESSION_ID" ] || [ "$SESSION_ID" = "null" ]; then
  echo "Creating session..."
  SESSION_ID=$(curl -s -X POST http://localhost:8080/api/sessions \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"campaign_id\": \"$CAMPAIGN_ID\", \"name\": \"Test Session\"}" | jq -r '.id')
fi

echo "✓ Session ID: $SESSION_ID"

# Test AI with a detailed prompt
echo ""
echo "2. Testing AI response length..."
echo ""

AI_RESPONSE=$(curl -s -X POST http://localhost:8080/api/sessions/$SESSION_ID/ai-action \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "I enter the ancient library and look around. Describe in detail what I see - the architecture, any books or scrolls, magical items, creatures, and any interesting features. I want to get a full sense of this mysterious place.",
    "action_type": "exploration"
  }')

# Extract and analyze response
DESCRIPTION=$(echo "$AI_RESPONSE" | jq -r '.description')
LENGTH=$(echo -n "$DESCRIPTION" | wc -c)

echo "Response preview:"
echo "$DESCRIPTION" | head -n 10
echo "..."
echo ""
echo "Response length: $LENGTH characters"

if [ $LENGTH -gt 2000 ]; then
  echo "✓ SUCCESS: Enhanced parameters working! Much longer response generated."
else
  echo "⚠ Response may still be limited"
fi

# Kill server
kill $SERVER_PID 2>/dev/null

echo ""
echo "=== Test Complete ==="