package data

import "dnd-simulator/internal/models"

// D&D 5e Background data
var Backgrounds = map[string]models.Background{
	"acolyte": {
		Name:       "Acolyte",
		SkillProfs: []string{"Insight", "Religion"},
		Languages:  []string{"Two of your choice"},
		Equipment:  []string{"Holy symbol", "Prayer book", "Incense", "Vestments", "Common clothes", "Belt pouch with 15 gp"},
		Feature:    "Shelter of the Faithful",
		Personality: []string{
			"I idolize a particular hero of my faith, and constantly refer to that person's deeds and example.",
			"I can find common ground between the fiercest enemies, empathizing with them and always working toward peace.",
		},
		Ideals: []string{
			"Tradition. The ancient traditions of worship and sacrifice must be preserved and upheld.",
			"Charity. I always try to help those in need, no matter what the personal cost.",
		},
		Bonds: []string{
			"I would die to recover an ancient relic of my faith that was lost long ago.",
			"I will someday get revenge on the corrupt temple hierarchy who branded me a heretic.",
		},
		Flaws: []string{
			"I judge others harshly, and myself even more severely.",
			"I put too much trust in those who wield power within my temple's hierarchy.",
		},
	},
	"criminal": {
		Name:       "Criminal",
		SkillProfs: []string{"Deception", "Stealth"},
		Languages:  []string{},
		Equipment:  []string{"Crowbar", "Dark common clothes with hood", "Belt pouch with 15 gp"},
		Feature:    "Criminal Contact",
		Personality: []string{
			"I always have a plan for what to do when things go wrong.",
			"I am always calm, no matter what the situation. I never raise my voice or let my emotions control me.",
		},
		Ideals: []string{
			"Honor. I don't steal from others in the trade.",
			"Freedom. Chains are meant to be broken, as are those who would forge them.",
		},
		Bonds: []string{
			"I'm trying to pay off an old debt I owe to a generous benefactor.",
			"My ill-gotten gains go to support my family.",
		},
		Flaws: []string{
			"When I see something valuable, I can't think about anything but how to steal it.",
			"When faced with a choice between money and my friends, I usually choose the money.",
		},
	},
	"folk-hero": {
		Name:       "Folk Hero",
		SkillProfs: []string{"Animal Handling", "Survival"},
		Languages:  []string{},
		Equipment:  []string{"Artisan's tools", "Shovel", "Iron pot", "Common clothes", "Belt pouch with 10 gp"},
		Feature:    "Rustic Hospitality",
		Personality: []string{
			"I judge people by their actions, not their words.",
			"If someone is in trouble, I'm always ready to lend help.",
		},
		Ideals: []string{
			"Respect. People deserve to be treated with dignity and respect.",
			"Fairness. No one should get preferential treatment before the law, and no one is above the law.",
		},
		Bonds: []string{
			"I have a family, but I have no idea where they are. I hope to see them again one day.",
			"I worked the land, I love the land, and I will protect the land.",
		},
		Flaws: []string{
			"The tyrant who rules my land will stop at nothing to see me killed.",
			"I'm convinced of the significance of my destiny, and blind to my shortcomings and the risk of failure.",
		},
	},
	"noble": {
		Name:       "Noble",
		SkillProfs: []string{"History", "Persuasion"},
		Languages:  []string{"One of your choice"},
		Equipment:  []string{"Fine clothes", "Signet ring", "Scroll of pedigree", "Purse with 25 gp"},
		Feature:    "Position of Privilege",
		Personality: []string{
			"My eloquent flattery makes everyone I talk to feel like the most wonderful and important person in the world.",
			"The common folk love me for my kindness and generosity.",
		},
		Ideals: []string{
			"Respect. Respect is due to me because of my position, but all people regardless of station deserve to be treated with dignity.",
			"Noble Obligation. It is my duty to protect and care for the people beneath me.",
		},
		Bonds: []string{
			"I will face any challenge to win the approval of my family.",
			"My house's alliance with another noble family must be sustained at all costs.",
		},
		Flaws: []string{
			"I secretly believe that everyone is beneath me.",
			"I hide a truly scandalous secret that could ruin my family forever.",
		},
	},
	"sage": {
		Name:       "Sage",
		SkillProfs: []string{"Arcana", "History"},
		Languages:  []string{"Two of your choice"},
		Equipment:  []string{"Bottle of black ink", "Quill", "Small knife", "Letter", "Common clothes", "Belt pouch with 10 gp"},
		Feature:    "Researcher",
		Personality: []string{
			"I use polysyllabic words that convey the exact meaning I intend.",
			"I've read every book in the world's greatest libraries—or I like to boast that I have.",
		},
		Ideals: []string{
			"Knowledge. The path to power and self-improvement is through knowledge.",
			"Beauty. What is beautiful points us beyond itself toward what is true.",
		},
		Bonds: []string{
			"It is my duty to protect my students.",
			"I have an ancient text that holds terrible secrets that must not fall into the wrong hands.",
		},
		Flaws: []string{
			"I am easily distracted by the promise of information.",
			"Most people scream and run when they see a demon. I stop and take notes on its anatomy.",
		},
	},
	"soldier": {
		Name:       "Soldier",
		SkillProfs: []string{"Athletics", "Intimidation"},
		Languages:  []string{},
		Equipment:  []string{"Insignia of rank", "Trophy", "Common clothes", "Belt pouch with 10 gp"},
		Feature:    "Military Rank",
		Personality: []string{
			"I'm always polite and respectful.",
			"I'm haunted by memories of war. I can't get the images of violence out of my mind.",
		},
		Ideals: []string{
			"Greater Good. Our lot is to lay down our lives in defense of others.",
			"Responsibility. I do what I must and obey just authority.",
		},
		Bonds: []string{
			"I would still lay down my life for the people I served with.",
			"Someone saved my life on the battlefield. To this day, I will never leave a friend behind.",
		},
		Flaws: []string{
			"The monstrous enemy we faced in battle still leaves me quivering with fear.",
			"I have little respect for anyone who is not a proven warrior.",
		},
	},
}