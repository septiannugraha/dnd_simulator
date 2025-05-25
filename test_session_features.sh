#!/bin/bash

# Test the new session features

echo "1. Login as DM..."
DM_LOGIN=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "dm_test", "password": "password123"}')

DM_TOKEN=$(echo $DM_LOGIN | jq -r '.token')
echo "DM Token obtained"

echo -e "\n2. Get campaigns..."
curl -s -X GET "http://localhost:8080/api/campaigns" \
  -H "Authorization: Bearer $DM_TOKEN" | jq

echo -e "\n3. Create a new session..."
SESSION_CREATE=$(curl -s -X POST "http://localhost:8080/api/sessions" \
  -H "Authorization: Bearer $DM_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "campaign_id": "683211c5013638aaf31e9159",
    "name": "Test Session",
    "description": "Testing new features"
  }')

echo $SESSION_CREATE | jq

SESSION_ID=$(echo $SESSION_CREATE | jq -r '.id')
echo "Session ID: $SESSION_ID"

if [ "$SESSION_ID" != "null" ]; then
  echo -e "\n4. Start the session..."
  curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/start" \
    -H "Authorization: Bearer $DM_TOKEN" | jq

  echo -e "\n5. Update the scene..."
  curl -s -X PUT "http://localhost:8080/api/sessions/$SESSION_ID/scene" \
    -H "Authorization: Bearer $DM_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "scene": "The adventurers arrive at the goblin cave entrance",
      "notes": "First encounter"
    }' | jq

  echo -e "\n6. Get session status..."
  curl -s -X GET "http://localhost:8080/api/sessions/$SESSION_ID/status" \
    -H "Authorization: Bearer $DM_TOKEN" | jq

  echo -e "\n7. Pause the session..."
  curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/pause" \
    -H "Authorization: Bearer $DM_TOKEN" | jq

  echo -e "\n8. Resume the session..."
  curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/resume" \
    -H "Authorization: Bearer $DM_TOKEN" | jq

  echo -e "\n9. Get narrative history..."
  curl -s -X GET "http://localhost:8080/api/sessions/$SESSION_ID/narrative" \
    -H "Authorization: Bearer $DM_TOKEN" | jq

  echo -e "\n10. End the session..."
  curl -s -X POST "http://localhost:8080/api/sessions/$SESSION_ID/end" \
    -H "Authorization: Bearer $DM_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"status": "completed"}' | jq
fi