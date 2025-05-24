package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// D&D 5e Race definitions
type Race struct {
	Name              string            `json:"name" bson:"name"`
	Size              string            `json:"size" bson:"size"`
	Speed             int               `json:"speed" bson:"speed"`
	AbilityIncrease   map[string]int    `json:"ability_increase" bson:"ability_increase"`
	Traits            []string          `json:"traits" bson:"traits"`
	Languages         []string          `json:"languages" bson:"languages"`
	Proficiencies     []string          `json:"proficiencies" bson:"proficiencies"`
}

// D&D 5e Class definitions
type Class struct {
	Name           string   `json:"name" bson:"name"`
	HitDie         int      `json:"hit_die" bson:"hit_die"`
	PrimaryAbility []string `json:"primary_ability" bson:"primary_ability"`
	SavingThrows   []string `json:"saving_throws" bson:"saving_throws"`
	SkillChoices   int      `json:"skill_choices" bson:"skill_choices"`
	Skills         []string `json:"skills" bson:"skills"`
	Equipment      []string `json:"equipment" bson:"equipment"`
	Spellcaster    bool     `json:"spellcaster" bson:"spellcaster"`
}

// D&D 5e Background definitions
type Background struct {
	Name         string   `json:"name" bson:"name"`
	SkillProfs   []string `json:"skill_proficiencies" bson:"skill_proficiencies"`
	Languages    []string `json:"languages" bson:"languages"`
	Equipment    []string `json:"equipment" bson:"equipment"`
	Feature      string   `json:"feature" bson:"feature"`
	Personality  []string `json:"personality_traits" bson:"personality_traits"`
	Ideals       []string `json:"ideals" bson:"ideals"`
	Bonds        []string `json:"bonds" bson:"bonds"`
	Flaws        []string `json:"flaws" bson:"flaws"`
}

// Enhanced Character model
type Character struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID   `bson:"user_id" json:"user_id"`
	CampaignID        primitive.ObjectID   `bson:"campaign_id,omitempty" json:"campaign_id,omitempty"`
	Name              string               `bson:"name" json:"name" binding:"required,min=2,max=50"`
	Race              string               `bson:"race" json:"race" binding:"required"`
	Class             string               `bson:"class" json:"class" binding:"required"`
	Background        string               `bson:"background" json:"background" binding:"required"`
	Level             int                  `bson:"level" json:"level"`
	ExperiencePoints  int                  `bson:"experience_points" json:"experience_points"`
	
	// Core Abilities
	Abilities         AbilityScores        `bson:"abilities" json:"abilities"`
	
	// Combat Stats
	HitPoints         int                  `bson:"hit_points" json:"hit_points"`
	MaxHitPoints      int                  `bson:"max_hit_points" json:"max_hit_points"`
	ArmorClass        int                  `bson:"armor_class" json:"armor_class"`
	Initiative        int                  `bson:"initiative" json:"initiative"`
	Speed             int                  `bson:"speed" json:"speed"`
	
	// Proficiencies
	ProficiencyBonus  int                  `bson:"proficiency_bonus" json:"proficiency_bonus"`
	SavingThrows      map[string]int       `bson:"saving_throws" json:"saving_throws"`
	Skills            map[string]int       `bson:"skills" json:"skills"`
	
	// Equipment & Inventory
	Equipment         []Equipment          `bson:"equipment" json:"equipment"`
	Weapons           []Weapon             `bson:"weapons" json:"weapons"`
	Armor             *Armor               `bson:"armor,omitempty" json:"armor,omitempty"`
	
	// Spellcasting
	SpellcastingClass string               `bson:"spellcasting_class,omitempty" json:"spellcasting_class,omitempty"`
	SpellSlots        map[string]int       `bson:"spell_slots,omitempty" json:"spell_slots,omitempty"`
	SpellsKnown       []Spell              `bson:"spells_known,omitempty" json:"spells_known,omitempty"`
	CantripsKnown     []Spell              `bson:"cantrips_known,omitempty" json:"cantrips_known,omitempty"`
	
	// Character Details
	Alignment         string               `bson:"alignment" json:"alignment"`
	PersonalityTraits []string             `bson:"personality_traits" json:"personality_traits"`
	Ideals            []string             `bson:"ideals" json:"ideals"`
	Bonds             []string             `bson:"bonds" json:"bonds"`
	Flaws             []string             `bson:"flaws" json:"flaws"`
	
	// Metadata
	CreatedAt         time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time            `bson:"updated_at" json:"updated_at"`
}

type AbilityScores struct {
	Strength     int `bson:"strength" json:"strength" binding:"min=1,max=30"`
	Dexterity    int `bson:"dexterity" json:"dexterity" binding:"min=1,max=30"`
	Constitution int `bson:"constitution" json:"constitution" binding:"min=1,max=30"`
	Intelligence int `bson:"intelligence" json:"intelligence" binding:"min=1,max=30"`
	Wisdom       int `bson:"wisdom" json:"wisdom" binding:"min=1,max=30"`
	Charisma     int `bson:"charisma" json:"charisma" binding:"min=1,max=30"`
}

type Equipment struct {
	Name        string `bson:"name" json:"name"`
	Quantity    int    `bson:"quantity" json:"quantity"`
	Weight      float64 `bson:"weight" json:"weight"`
	Value       int    `bson:"value" json:"value"` // in copper pieces
	Description string `bson:"description" json:"description"`
}

type Weapon struct {
	Name       string   `bson:"name" json:"name"`
	Damage     string   `bson:"damage" json:"damage"`     // e.g., "1d8"
	DamageType string   `bson:"damage_type" json:"damage_type"` // e.g., "slashing"
	Properties []string `bson:"properties" json:"properties"`
	Range      string   `bson:"range,omitempty" json:"range,omitempty"`
	Weight     float64  `bson:"weight" json:"weight"`
	Value      int      `bson:"value" json:"value"`
}

type Armor struct {
	Name     string  `bson:"name" json:"name"`
	Type     string  `bson:"type" json:"type"` // light, medium, heavy, shield
	AC       int     `bson:"ac" json:"ac"`
	DexMod   bool    `bson:"dex_mod" json:"dex_mod"`
	MaxDex   int     `bson:"max_dex,omitempty" json:"max_dex,omitempty"`
	MinStr   int     `bson:"min_str,omitempty" json:"min_str,omitempty"`
	Stealth  bool    `bson:"stealth_disadvantage,omitempty" json:"stealth_disadvantage,omitempty"`
	Weight   float64 `bson:"weight" json:"weight"`
	Value    int     `bson:"value" json:"value"`
}

type Spell struct {
	Name        string   `bson:"name" json:"name"`
	Level       int      `bson:"level" json:"level"`
	School      string   `bson:"school" json:"school"`
	CastingTime string   `bson:"casting_time" json:"casting_time"`
	Range       string   `bson:"range" json:"range"`
	Components  []string `bson:"components" json:"components"`
	Duration    string   `bson:"duration" json:"duration"`
	Description string   `bson:"description" json:"description"`
}