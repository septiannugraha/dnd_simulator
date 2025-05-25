#!/bin/bash

# Direct test of Gemini API with enhanced parameters
echo "=== Testing Gemini API with Enhanced Parameters ==="
echo ""

API_KEY="AIzaSyAzZx60x2vUUF42YUODHNAQ1KH6_dvvA-Y"
MODEL="gemini-2.0-flash-exp"

# Test with enhanced parameters
echo "Testing with maxOutputTokens: 8192..."
echo ""

curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:streamGenerateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [{
      "role": "user",
      "parts": [{
        "text": "You are a Dungeon Master for D&D 5e. Describe in rich detail an ancient magical library that the party has just discovered. Include descriptions of the architecture, mysterious tomes, magical artifacts, guardian creatures, hidden sections, and any other interesting features. Make this description at least 500 words long to fully immerse the players in this wondrous location."
      }]
    }],
    "generationConfig": {
      "temperature": 0.8,
      "maxOutputTokens": 8192,
      "topK": 40,
      "topP": 0.95
    }
  }' > /tmp/ai_response.txt

# Parse the streaming response
echo "Parsing response..."
FULL_TEXT=""
while IFS= read -r line; do
  if [[ $line == data:* ]]; then
    # Extract JSON from the data line
    json_data="${line#data: }"
    if [[ -n "$json_data" && "$json_data" != "[DONE]" ]]; then
      # Extract text from the JSON
      text=$(echo "$json_data" | jq -r '.candidates[0].content.parts[0].text // empty' 2>/dev/null)
      if [[ -n "$text" ]]; then
        FULL_TEXT+="$text"
      fi
    fi
  fi
done < /tmp/ai_response.txt

# Display results
echo "Response preview:"
echo "$FULL_TEXT" | head -n 20
echo "..."
echo ""
echo "Response length: $(echo -n "$FULL_TEXT" | wc -c) characters"
echo ""

# Compare with limited tokens
echo "Testing with maxOutputTokens: 500 for comparison..."
curl -s -X POST \
  "https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "contents": [{
      "role": "user",
      "parts": [{
        "text": "Describe a magical library briefly."
      }]
    }],
    "generationConfig": {
      "temperature": 0.8,
      "maxOutputTokens": 500
    }
  }' | jq -r '.candidates[0].content.parts[0].text' > /tmp/short_response.txt

echo ""
echo "Short response length: $(wc -c < /tmp/short_response.txt) characters"
echo ""
echo "=== Comparison ==="
echo "Enhanced (8192 tokens): $(echo -n "$FULL_TEXT" | wc -c) characters"
echo "Limited (500 tokens): $(wc -c < /tmp/short_response.txt) characters"

# Cleanup
rm -f /tmp/ai_response.txt /tmp/short_response.txt