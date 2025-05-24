package services

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"dnd-simulator/internal/database"
	"dnd-simulator/internal/data"
	"dnd-simulator/internal/models"
)

type CharacterService struct {
	db *database.DB
}

func NewCharacterService(db *database.DB) *CharacterService {
	return &CharacterService{db: db}
}

type CreateCharacterRequest struct {
	Name       string                   `json:"name" binding:"required,min=2,max=50"`
	Race       string                   `json:"race" binding:"required"`
	Class      string                   `json:"class" binding:"required"`
	Background string                   `json:"background" binding:"required"`
	Abilities  models.AbilityScores     `json:"abilities" binding:"required"`
	Alignment  string                   `json:"alignment" binding:"required"`
}

func (s *CharacterService) CreateCharacter(userID primitive.ObjectID, req CreateCharacterRequest) (*models.Character, error) {
	collection := s.db.GetCollection("characters")

	// Validate race, class, and background exist
	race, exists := data.Races[req.Race]
	if !exists {
		return nil, errors.New("invalid race")
	}

	class, exists := data.Classes[req.Class]
	if !exists {
		return nil, errors.New("invalid class")
	}

	background, exists := data.Backgrounds[req.Background]
	if !exists {
		return nil, errors.New("invalid background")
	}

	// Apply racial ability score increases
	finalAbilities := req.Abilities
	for ability, increase := range race.AbilityIncrease {
		switch ability {
		case "strength":
			finalAbilities.Strength += increase
		case "dexterity":
			finalAbilities.Dexterity += increase
		case "constitution":
			finalAbilities.Constitution += increase
		case "intelligence":
			finalAbilities.Intelligence += increase
		case "wisdom":
			finalAbilities.Wisdom += increase
		case "charisma":
			finalAbilities.Charisma += increase
		}
	}

	// Calculate derived stats
	level := 1
	proficiencyBonus := s.calculateProficiencyBonus(level)
	hitPoints := s.calculateHitPoints(class, finalAbilities.Constitution, level)
	armorClass := s.calculateArmorClass(finalAbilities.Dexterity, nil)
	initiative := s.calculateModifier(finalAbilities.Dexterity)
	savingThrows := s.calculateSavingThrows(finalAbilities, class.SavingThrows, proficiencyBonus)
	skills := s.calculateSkills(finalAbilities, background.SkillProfs, proficiencyBonus)

	character := &models.Character{
		UserID:            userID,
		Name:              req.Name,
		Race:              req.Race,
		Class:             req.Class,
		Background:        req.Background,
		Level:             level,
		ExperiencePoints:  0,
		Abilities:         finalAbilities,
		HitPoints:         hitPoints,
		MaxHitPoints:      hitPoints,
		ArmorClass:        armorClass,
		Initiative:        initiative,
		Speed:             race.Speed,
		ProficiencyBonus:  proficiencyBonus,
		SavingThrows:      savingThrows,
		Skills:            skills,
		Equipment:         []models.Equipment{},
		Weapons:           []models.Weapon{},
		Alignment:         req.Alignment,
		PersonalityTraits: []string{},
		Ideals:            []string{},
		Bonds:             []string{},
		Flaws:             []string{},
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Set spellcasting info if applicable
	if class.Spellcaster {
		character.SpellcastingClass = req.Class
		character.SpellSlots = s.calculateSpellSlots(req.Class, level)
		character.SpellsKnown = []models.Spell{}
		character.CantripsKnown = []models.Spell{}
	}

	result, err := collection.InsertOne(context.Background(), character)
	if err != nil {
		return nil, err
	}

	character.ID = result.InsertedID.(primitive.ObjectID)
	return character, nil
}

func (s *CharacterService) GetCharacterByID(characterID primitive.ObjectID) (*models.Character, error) {
	collection := s.db.GetCollection("characters")

	var character models.Character
	err := collection.FindOne(context.Background(), bson.M{"_id": characterID}).Decode(&character)
	if err != nil {
		return nil, err
	}

	return &character, nil
}

func (s *CharacterService) GetUserCharacters(userID primitive.ObjectID) ([]models.Character, error) {
	collection := s.db.GetCollection("characters")

	cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var characters []models.Character
	if err = cursor.All(context.Background(), &characters); err != nil {
		return nil, err
	}

	return characters, nil
}

func (s *CharacterService) UpdateCharacter(characterID, userID primitive.ObjectID, updates map[string]interface{}) (*models.Character, error) {
	collection := s.db.GetCollection("characters")

	// Check ownership
	var character models.Character
	err := collection.FindOne(context.Background(), bson.M{"_id": characterID}).Decode(&character)
	if err != nil {
		return nil, err
	}

	if character.UserID != userID {
		return nil, errors.New("you can only update your own characters")
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": characterID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, err
	}

	return s.GetCharacterByID(characterID)
}

func (s *CharacterService) DeleteCharacter(characterID, userID primitive.ObjectID) error {
	collection := s.db.GetCollection("characters")

	// Check ownership
	var character models.Character
	err := collection.FindOne(context.Background(), bson.M{"_id": characterID}).Decode(&character)
	if err != nil {
		return err
	}

	if character.UserID != userID {
		return errors.New("you can only delete your own characters")
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": characterID})
	return err
}

func (s *CharacterService) AssignToCampaign(characterID, campaignID, userID primitive.ObjectID) error {
	collection := s.db.GetCollection("characters")

	// Check ownership
	var character models.Character
	err := collection.FindOne(context.Background(), bson.M{"_id": characterID}).Decode(&character)
	if err != nil {
		return err
	}

	if character.UserID != userID {
		return errors.New("you can only assign your own characters")
	}

	_, err = collection.UpdateOne(
		context.Background(),
		bson.M{"_id": characterID},
		bson.M{"$set": bson.M{"campaign_id": campaignID, "updated_at": time.Now()}},
	)
	return err
}

// D&D 5e Calculation Functions

func (s *CharacterService) calculateModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

func (s *CharacterService) calculateProficiencyBonus(level int) int {
	return int(math.Ceil(float64(level)/4)) + 1
}

func (s *CharacterService) calculateHitPoints(class models.Class, constitutionScore, level int) int {
	conMod := s.calculateModifier(constitutionScore)
	baseHP := class.HitDie + conMod
	
	// For level 1, just base HP
	if level == 1 {
		return baseHP
	}
	
	// For higher levels, add average hit die rolls + con mod per level
	averageHitDie := float64(class.HitDie)/2 + 1
	additionalHP := int(averageHitDie) * (level - 1) + conMod*(level-1)
	
	return baseHP + additionalHP
}

func (s *CharacterService) calculateArmorClass(dexterity int, armor *models.Armor) int {
	dexMod := s.calculateModifier(dexterity)
	
	if armor == nil {
		// No armor: 10 + Dex mod
		return 10 + dexMod
	}
	
	ac := armor.AC
	if armor.DexMod {
		if armor.MaxDex > 0 && dexMod > armor.MaxDex {
			ac += armor.MaxDex
		} else {
			ac += dexMod
		}
	}
	
	return ac
}

func (s *CharacterService) calculateSavingThrows(abilities models.AbilityScores, proficientSaves []string, proficiencyBonus int) map[string]int {
	saves := map[string]int{
		"strength":     s.calculateModifier(abilities.Strength),
		"dexterity":    s.calculateModifier(abilities.Dexterity),
		"constitution": s.calculateModifier(abilities.Constitution),
		"intelligence": s.calculateModifier(abilities.Intelligence),
		"wisdom":       s.calculateModifier(abilities.Wisdom),
		"charisma":     s.calculateModifier(abilities.Charisma),
	}
	
	// Add proficiency bonus to proficient saves
	for _, save := range proficientSaves {
		if _, exists := saves[save]; exists {
			saves[save] += proficiencyBonus
		}
	}
	
	return saves
}

func (s *CharacterService) calculateSkills(abilities models.AbilityScores, proficientSkills []string, proficiencyBonus int) map[string]int {
	// D&D 5e skill to ability mapping
	skillAbilities := map[string]string{
		"Acrobatics":      "dexterity",
		"Animal Handling": "wisdom",
		"Arcana":          "intelligence",
		"Athletics":       "strength",
		"Deception":       "charisma",
		"History":         "intelligence",
		"Insight":         "wisdom",
		"Intimidation":    "charisma",
		"Investigation":   "intelligence",
		"Medicine":        "wisdom",
		"Nature":          "intelligence",
		"Perception":      "wisdom",
		"Performance":     "charisma",
		"Persuasion":      "charisma",
		"Religion":        "intelligence",
		"Sleight of Hand": "dexterity",
		"Stealth":         "dexterity",
		"Survival":        "wisdom",
	}
	
	abilityValues := map[string]int{
		"strength":     abilities.Strength,
		"dexterity":    abilities.Dexterity,
		"constitution": abilities.Constitution,
		"intelligence": abilities.Intelligence,
		"wisdom":       abilities.Wisdom,
		"charisma":     abilities.Charisma,
	}
	
	skills := make(map[string]int)
	
	for skill, ability := range skillAbilities {
		modifier := s.calculateModifier(abilityValues[ability])
		
		// Add proficiency bonus if proficient
		for _, profSkill := range proficientSkills {
			if skill == profSkill {
				modifier += proficiencyBonus
				break
			}
		}
		
		skills[skill] = modifier
	}
	
	return skills
}

func (s *CharacterService) calculateSpellSlots(class string, level int) map[string]int {
	// Simplified spell slot calculation for level 1
	slots := make(map[string]int)
	
	switch class {
	case "wizard", "sorcerer", "cleric", "druid":
		slots["1st"] = 2
	case "bard", "ranger", "paladin", "warlock":
		if level >= 2 { // These classes get spells at level 2
			slots["1st"] = 2
		}
	}
	
	return slots
}