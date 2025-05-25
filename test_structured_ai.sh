#!/bin/bash

# Test structured AI output - demonstrating the benefits

echo "=== D&D Simulator - Structured AI Output Test ==="
echo ""
echo "This test demonstrates how structured output improves game mechanics handling"
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Example of current text-based response vs structured response
echo -e "${BLUE}Current Text-Based AI Response:${NC}"
echo "The goblin snarls and swings its rusty scimitar at you. Make a Dexterity saving throw, DC 14. If you fail, you take 1d6+2 slashing damage."
echo ""

echo -e "${YELLOW}Problems with text parsing:${NC}"
echo "- Might miss the DC value (14)"
echo "- Could miss damage formula (1d6+2)"
echo "- Doesn't capture damage type (slashing)"
echo "- No structured way to handle success/failure"
echo ""

echo -e "${GREEN}Structured JSON Response:${NC}"
cat << 'EOF'
{
  "narrative": "The goblin snarls and swings its rusty scimitar at you in a vicious arc!",
  "game_mechanics": [
    {
      "type": "saving_throw",
      "description": "Dodge the goblin's scimitar attack",
      "target": "Player Character",
      "requirements": {
        "roll_type": "d20",
        "save_type": "dexterity",
        "dc": 14
      },
      "effects": {
        "on_success": "You dodge the attack completely",
        "on_failure": "The scimitar slashes across your body",
        "damage_dice": "1d6+2",
        "damage_type": "slashing"
      }
    }
  ]
}
EOF

echo ""
echo -e "${GREEN}Benefits of Structured Output:${NC}"
echo "✓ Guaranteed extraction of all game mechanics"
echo "✓ Precise DC values, damage formulas, and types"
echo "✓ Clear success/failure conditions"
echo "✓ Easy to process programmatically"
echo "✓ No parsing errors or missed information"
echo ""

echo -e "${BLUE}Complex Combat Example:${NC}"
cat << 'EOF'
{
  "narrative": "The ancient dragon rears back, flames gathering in its maw. The temperature in the cavern rises dramatically as it prepares to unleash its fury!",
  "game_mechanics": [
    {
      "type": "saving_throw",
      "description": "Avoid the dragon's breath weapon",
      "target": "All creatures in 60-foot cone",
      "requirements": {
        "roll_type": "d20",
        "save_type": "dexterity",
        "dc": 21
      },
      "effects": {
        "on_success": "Take half damage from the flames",
        "on_failure": "Take full damage from the dragon's breath",
        "damage_dice": "16d6",
        "damage_type": "fire"
      }
    },
    {
      "type": "condition",
      "description": "Frightening Presence",
      "target": "All creatures within 120 feet",
      "requirements": {
        "roll_type": "d20",
        "save_type": "wisdom",
        "dc": 19
      },
      "effects": {
        "on_failure": "Frightened for 1 minute",
        "conditions": ["frightened"],
        "duration": "1 minute"
      }
    }
  ],
  "npc_actions": [
    {
      "npc_name": "Ancient Red Dragon",
      "action": "Breath Weapon",
      "initiative": 20
    }
  ],
  "scene_changes": {
    "hazards": ["extreme heat", "falling rocks from ceiling"]
  }
}
EOF

echo ""
echo -e "${YELLOW}Implementation Benefits:${NC}"
echo "1. Frontend can automatically create roll buttons with correct DCs"
echo "2. Damage calculations are precise and automated"
echo "3. Conditions and durations are trackable"
echo "4. NPC actions can be queued in initiative order"
echo "5. Environmental hazards are clearly defined"
echo ""

echo -e "${GREEN}Recommended Next Steps:${NC}"
echo "1. Update AI service to use structured output (ai_enhanced.go)"
echo "2. Update frontend to parse structured responses"
echo "3. Create UI components for each mechanic type"
echo "4. Add automatic dice rolling with DC comparison"
echo "5. Track conditions and durations automatically"