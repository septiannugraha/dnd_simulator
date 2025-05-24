package services

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"dnd-simulator/internal/models"
)

type DiceService struct {
	rng *rand.Rand
}

func NewDiceService() *DiceService {
	return &DiceService{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ParseAndRoll parses a dice string (e.g., "1d20+5", "3d6", "2d8-1") and returns the result
func (ds *DiceService) ParseAndRoll(diceString string, purpose string) (*models.DiceRoll, error) {
	// Clean the input
	diceString = strings.TrimSpace(strings.ToLower(diceString))
	
	// Regex to parse dice notation: XdY+Z or XdY-Z
	re := regexp.MustCompile(`^(\d+)d(\d+)([+-]\d+)?$`)
	matches := re.FindStringSubmatch(diceString)
	
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid dice format: %s", diceString)
	}
	
	// Parse number of dice
	numDice, err := strconv.Atoi(matches[1])
	if err != nil || numDice <= 0 || numDice > 20 { // Limit to reasonable number
		return nil, fmt.Errorf("invalid number of dice: %s", matches[1])
	}
	
	// Parse die type
	dieType, err := strconv.Atoi(matches[2])
	if err != nil || !ds.isValidDieType(dieType) {
		return nil, fmt.Errorf("invalid die type: %s", matches[2])
	}
	
	// Parse modifier
	var modifier int
	if len(matches) > 3 && matches[3] != "" {
		modifier, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid modifier: %s", matches[3])
		}
	}
	
	// Roll the dice
	results := make([]int, numDice)
	total := 0
	
	for i := 0; i < numDice; i++ {
		roll := ds.rng.Intn(dieType) + 1
		results[i] = roll
		total += roll
	}
	
	// Apply modifier
	total += modifier
	
	return &models.DiceRoll{
		Dice:     diceString,
		Result:   results,
		Total:    total,
		Modifier: modifier,
		Purpose:  purpose,
	}, nil
}

// RollStandardDice provides common D&D dice rolls
func (ds *DiceService) RollStandardDice(diceType string, purpose string) (*models.DiceRoll, error) {
	var diceString string
	
	switch diceType {
	case "d20":
		diceString = "1d20"
	case "d12":
		diceString = "1d12"
	case "d10":
		diceString = "1d10"
	case "d8":
		diceString = "1d8"
	case "d6":
		diceString = "1d6"
	case "d4":
		diceString = "1d4"
	case "d100", "percentile":
		diceString = "1d100"
	case "advantage":
		return ds.rollAdvantage(purpose)
	case "disadvantage":
		return ds.rollDisadvantage(purpose)
	default:
		return nil, fmt.Errorf("unknown dice type: %s", diceType)
	}
	
	return ds.ParseAndRoll(diceString, purpose)
}

// RollAbilityScores rolls 4d6 drop lowest for ability score generation
func (ds *DiceService) RollAbilityScores() map[string]int {
	abilities := []string{"strength", "dexterity", "constitution", "intelligence", "wisdom", "charisma"}
	scores := make(map[string]int)
	
	for _, ability := range abilities {
		// Roll 4d6, drop lowest
		rolls := make([]int, 4)
		for i := 0; i < 4; i++ {
			rolls[i] = ds.rng.Intn(6) + 1
		}
		
		// Find and remove the lowest roll
		lowest := 0
		for i := 1; i < 4; i++ {
			if rolls[i] < rolls[lowest] {
				lowest = i
			}
		}
		
		total := 0
		for i := 0; i < 4; i++ {
			if i != lowest {
				total += rolls[i]
			}
		}
		
		scores[ability] = total
	}
	
	return scores
}

// rollAdvantage rolls 2d20 and takes the higher result
func (ds *DiceService) rollAdvantage(purpose string) (*models.DiceRoll, error) {
	roll1 := ds.rng.Intn(20) + 1
	roll2 := ds.rng.Intn(20) + 1
	
	higher := roll1
	if roll2 > roll1 {
		higher = roll2
	}
	
	return &models.DiceRoll{
		Dice:     "2d20 (advantage)",
		Result:   []int{roll1, roll2},
		Total:    higher,
		Modifier: 0,
		Purpose:  fmt.Sprintf("%s (advantage)", purpose),
	}, nil
}

// rollDisadvantage rolls 2d20 and takes the lower result
func (ds *DiceService) rollDisadvantage(purpose string) (*models.DiceRoll, error) {
	roll1 := ds.rng.Intn(20) + 1
	roll2 := ds.rng.Intn(20) + 1
	
	lower := roll1
	if roll2 < roll1 {
		lower = roll2
	}
	
	return &models.DiceRoll{
		Dice:     "2d20 (disadvantage)",
		Result:   []int{roll1, roll2},
		Total:    lower,
		Modifier: 0,
		Purpose:  fmt.Sprintf("%s (disadvantage)", purpose),
	}, nil
}

// isValidDieType checks if the die type is a standard D&D die
func (ds *DiceService) isValidDieType(dieType int) bool {
	validDice := []int{4, 6, 8, 10, 12, 20, 100}
	for _, valid := range validDice {
		if dieType == valid {
			return true
		}
	}
	return false
}

// RollInitiative rolls initiative for a character
func (ds *DiceService) RollInitiative(dexModifier int) *models.DiceRoll {
	roll := ds.rng.Intn(20) + 1
	total := roll + dexModifier
	
	return &models.DiceRoll{
		Dice:     "1d20",
		Result:   []int{roll},
		Total:    total,
		Modifier: dexModifier,
		Purpose:  "Initiative",
	}
}

// RollAttack rolls an attack with modifiers
func (ds *DiceService) RollAttack(attackBonus int, purpose string) *models.DiceRoll {
	roll := ds.rng.Intn(20) + 1
	total := roll + attackBonus
	
	return &models.DiceRoll{
		Dice:     "1d20",
		Result:   []int{roll},
		Total:    total,
		Modifier: attackBonus,
		Purpose:  fmt.Sprintf("Attack: %s", purpose),
	}
}

// RollSavingThrow rolls a saving throw
func (ds *DiceService) RollSavingThrow(saveModifier int, saveName string) *models.DiceRoll {
	roll := ds.rng.Intn(20) + 1
	total := roll + saveModifier
	
	return &models.DiceRoll{
		Dice:     "1d20",
		Result:   []int{roll},
		Total:    total,
		Modifier: saveModifier,
		Purpose:  fmt.Sprintf("%s Saving Throw", saveName),
	}
}

// RollSkillCheck rolls a skill check
func (ds *DiceService) RollSkillCheck(skillModifier int, skillName string) *models.DiceRoll {
	roll := ds.rng.Intn(20) + 1
	total := roll + skillModifier
	
	return &models.DiceRoll{
		Dice:     "1d20",
		Result:   []int{roll},
		Total:    total,
		Modifier: skillModifier,
		Purpose:  fmt.Sprintf("%s Check", skillName),
	}
}

// GetCriticalHitMessage returns a fun message for natural 20s
func (ds *DiceService) GetCriticalHitMessage() string {
	messages := []string{
		"🎯 CRITICAL HIT! Natural 20!",
		"💥 CRIT! The dice gods smile upon you!",
		"⚡ NATURAL 20! Epic success!",
		"🔥 CRITICAL SUCCESS! Maximum awesome!",
		"🎲 NAT 20! Legendary moment!",
	}
	return messages[ds.rng.Intn(len(messages))]
}

// GetCriticalFailMessage returns a fun message for natural 1s
func (ds *DiceService) GetCriticalFailMessage() string {
	messages := []string{
		"💀 CRITICAL FAIL! Natural 1...",
		"😅 FUMBLE! The dice betray you!",
		"🤦 NAT 1! Things go horribly wrong!",
		"💥 EPIC FAIL! Spectacular failure!",
		"🎲 NATURAL 1! Time for chaos!",
	}
	return messages[ds.rng.Intn(len(messages))]
}