#!/bin/bash

# Direct test of Gemini API

GEMINI_API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
MODEL_ID="gemini-2.5-pro-preview-05-06"

echo "Testing Gemini API directly..."

cat << EOF > gemini_test.json
{
    "contents": [
      {
        "role": "user",
        "parts": [
          {
            "text": "You are a Dungeon Master. A dwarf fighter named Thorin approaches a hooded figure in a tavern who offers him a quest. Thorin says: 'I examine the parchment carefully. What does it say?' Describe what the parchment says."
          }
        ]
      }
    ],
    "generationConfig": {
      "temperature": 0.8,
      "maxOutputTokens": 1000,
      "responseMimeType": "text/plain"
    }
}
EOF

echo "Request:"
cat gemini_test.json | jq

echo -e "\nResponse:"
curl -s -X POST \
  -H "Content-Type: application/json" \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL_ID}:streamGenerateContent?key=${GEMINI_API_KEY}" \
  -d @gemini_test.json

rm gemini_test.json