package models

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// we have a points system, given to users
// according to their behavior on the site
const (
	PointsCreatesAccount = 30
	PointsLogsIn         = 1
	PointsPosts          = 5
	PointsPerDayAway     = -1
	PointsTextStarred    = 1
	PointsTextFlagged    = -10
)

// User is the struct for our users
// we need to use nulls.String rather than string on some fields
// when we also have a unique index on said field(s)
// see actions/shared.go ToNullString func
type User struct {
	ID                uuid.UUID    `json:"id" db:"id"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at" db:"updated_at"`
	AvatarURL         nulls.String `json:"avatar_url" db:"avatar_url"`
	Bio               nulls.String `json:"bio" db:"bio"`
	Email             nulls.String `json:"email" db:"email"`
	InvitationToken   string       `json:"invitation_token" db:"invitation_token"`
	InvitedAt         time.Time    `json:"invited_at" db:"invited_at"`
	IsAdmin           bool         `json:"is_admin" db:"is_admin"`
	LastLoggedAt      time.Time    `json:"last_logged_at" db:"last_logged_at"`
	LastPostedAt      time.Time    `json:"last_posted_at" db:"last_posted_at"`
	Name              nulls.String `json:"name" db:"name"`
	Nickname          nulls.String `json:"nickname" db:"nickname"`
	Provider          nulls.String `json:"provider" db:"provider"`
	ProviderID        nulls.String `json:"provider_id" db:"provider_id"`
	Score             int          `json:"score" db:"score"`
	SignedUpAt        time.Time    `json:"signedup_at" db:"signedup_at"`
	SponsorshipsCount int          `json:"sponsorships_count" db:"sponsorships_count"`
	SponsorID         uuid.UUID    `json:"sponsor_id" db:"sponsor_id"`
	Sponsoring        Users        `has_many:"users"`
	Texts             Texts        `has_many:"texts" order_by:"created_at desc"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
/*
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Nickname, Name: "Nickname"},
		&validators.URLIsPresent{Field: u.AvatarURL, Name: "AvatarURL", Message: "Doesn't look like a valid url..."},
		&validators.StringIsPresent{Field: u.Provider, Name: "Provider"},
		&validators.StringIsPresent{Field: u.ProviderID, Name: "ProviderID"},
		&validators.StringLengthInRange{Name: "Bio", Field: u.Bio, Min: 10, Max: 255, Message: "Too long, too short, not a proper Bio if you ask me..."},
	), nil
}*/

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email.String, Name: "Email"},
	), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Name.String, Name: "Name"},
		&validators.StringIsPresent{Field: u.Nickname.String, Name: "Nickname"},
		&validators.URLIsPresent{Field: u.AvatarURL.String, Name: "AvatarURL", Message: "Doesn't look like a valid url..."},
		&validators.StringIsPresent{Field: u.Provider.String, Name: "Provider"},
		&validators.StringIsPresent{Field: u.ProviderID.String, Name: "ProviderID"},
	), nil
}

// NickValidate checks if a nickname already exists
func NickValidate(nick string, tx *pop.Connection) string {
	q := tx.Where("nickname = ?", nick)
	qUser := User{}
	err := q.First(&qUser)
	if err == nil {
		// found a user with same username, create random nickname
		suffix := nickGenerate()
		return nick + suffix
	}
	return nick
}

func nickGenerate() string {
	var adjectives = []string{"Black", "White", "Gray", "Brown", "Red", "Pink", "Crimson", "Carnelian", "Orange", "Yellow", "Ivory", "Cream", "Green", "Viridian", "Aquamarine", "Cyan", "Blue", "Cerulean", "Azure", "Indigo", "Navy", "Violet", "Purple", "Lavender", "Magenta", "Rainbow", "Iridescent", "Spectrum", "Prism", "Bold", "Vivid", "Pale", "Clear", "Glass", "Translucent", "Misty", "Dark", "Light", "Gold", "Silver", "Copper", "Bronze", "Steel", "Iron", "Brass", "Mercury", "Zinc", "Chrome", "Platinum", "Titanium", "Nickel", "Lead", "Pewter", "Rust", "Metal", "Stone", "Quartz", "Granite", "Marble", "Alabaster", "Agate", "Jasper", "Pebble", "Pyrite", "Crystal", "Geode", "Obsidian", "Mica", "Flint", "Sand", "Gravel", "Boulder", "Basalt", "Ruby", "Beryl", "Scarlet", "Citrine", "Sulpher", "Topaz", "Amber", "Emerald", "Malachite", "Jade", "Abalone", "Lapis", "Sapphire", "Diamond", "Peridot", "Gem", "Jewel", "Bevel", "Coral", "Jet", "Ebony", "Wood", "Tree", "Cherry", "Maple", "Cedar", "Branch", "Bramble", "Rowan", "Ash", "Fir", "Pine", "Cactus", "Alder", "Grove", "Forest", "Jungle", "Palm", "Bush", "Mulberry", "Juniper", "Vine", "Ivy", "Rose", "Lily", "Tulip", "Daffodil", "Honeysuckle", "Fuschia", "Hazel", "Walnut", "Almond", "Lime", "Lemon", "Apple", "Blossom", "Bloom", "Crocus", "Rose", "Buttercup", "Dandelion", "Iris", "Carnation", "Fern", "Root", "Branch", "Leaf", "Seed", "Flower", "Petal", "Pollen", "Orchid", "Mangrove", "Cypress", "Sequoia", "Sage", "Heather", "Snapdragon", "Daisy", "Mountain", "Hill", "Alpine", "Chestnut", "Valley", "Glacier", "Forest", "Grove", "Glen", "Tree", "Thorn", "Stump", "Desert", "Canyon", "Dune", "Oasis", "Mirage", "Well", "Spring", "Meadow", "Field", "Prairie", "Grass", "Tundra", "Island", "Shore", "Sand", "Shell", "Surf", "Wave", "Foam", "Tide", "Lake", "River", "Brook", "Stream", "Pool", "Pond", "Sun", "Sprinkle", "Shade", "Shadow", "Rain", "Cloud", "Storm", "Hail", "Snow", "Sleet", "Thunder", "Lightning", "Wind", "Hurricane", "Typhoon", "Dawn", "Sunrise", "Morning", "Noon", "Twilight", "Evening", "Sunset", "Midnight", "Night", "Sky", "Star", "Stellar", "Comet", "Nebula", "Quasar", "Solar", "Lunar", "Planet", "Meteor", "Sprout", "Pear", "Plum", "Kiwi", "Berry", "Apricot", "Peach", "Mango", "Pineapple", "Coconut", "Olive", "Ginger", "Root", "Plain", "Fancy", "Stripe", "Spot", "Speckle", "Spangle", "Ring", "Band", "Blaze", "Paint", "Pinto", "Shade", "Tabby", "Brindle", "Patch", "Calico", "Checker", "Dot", "Pattern", "Glitter", "Glimmer", "Shimmer", "Dull", "Dust", "Dirt", "Glaze", "Scratch", "Quick", "Swift", "Fast", "Slow", "Clever", "Fire", "Flicker", "Flash", "Spark", "Ember", "Coal", "Flame", "Chocolate", "Vanilla", "Sugar", "Spice", "Cake", "Pie", "Cookie", "Candy", "Caramel", "Spiral", "Round", "Jelly", "Square", "Narrow", "Long", "Short", "Small", "Tiny", "Big", "Giant", "Great", "Atom", "Peppermint", "Mint", "Butter", "Fringe", "Rag", "Quilt", "Truth", "Lie", "Holy", "Curse", "Noble", "Sly", "Brave", "Shy", "Lava", "Foul", "Leather", "Fantasy", "Keen", "Luminous", "Feather", "Sticky", "Gossamer", "Cotton", "Rattle", "Silk", "Satin", "Cord", "Denim", "Flannel", "Plaid", "Wool", "Linen", "Silent", "Flax", "Weak", "Valiant", "Fierce", "Gentle", "Rhinestone", "Splash", "North", "South", "East", "West", "Summer", "Winter", "Autumn", "Spring", "Season", "Equinox", "Solstice", "Paper", "Motley", "Torch", "Ballistic", "Rampant", "Shag", "Freckle", "Wild", "Free", "Chain", "Sheer", "Crazy", "Mad", "Candle", "Ribbon", "Lace", "Notch", "Wax", "Shine", "Shallow", "Deep", "Bubble", "Harvest", "Fluff", "Venom", "Boom", "Slash", "Rune", "Cold", "Quill", "Love", "Hate", "Garnet", "Zircon", "Power", "Bone", "Void", "Horn", "Glory", "Cyber", "Nova", "Hot", "Helix", "Cosmic", "Quark", "Quiver", "Holly", "Clover", "Polar", "Regal", "Ripple", "Ebony", "Wheat", "Phantom", "Dew", "Chisel", "Crack", "Chatter", "Laser", "Foil", "Tin", "Clever", "Treasure", "Maze", "Twisty", "Curly", "Fortune", "Fate", "Destiny", "Cute", "Slime", "Ink", "Disco", "Plume", "Time", "Psychadelic", "Relic", "Fossil", "Water", "Savage", "Ancient", "Rapid", "Road", "Trail", "Stitch", "Button", "Bow", "Nimble", "Zest", "Sour", "Bitter", "Phase", "Fan", "Frill", "Plump", "Pickle", "Mud", "Puddle", "Pond", "River", "Spring", "Stream", "Battle", "Arrow", "Plume", "Roan", "Pitch", "Tar", "Cat", "Dog", "Horse", "Lizard", "Bird", "Fish", "Saber", "Scythe", "Sharp", "Soft", "Razor", "Neon", "Dandy", "Weed", "Swamp", "Marsh", "Bog", "Peat", "Moor", "Muck", "Mire", "Grave", "Fair", "Just", "Brick", "Puzzle", "Skitter", "Prong", "Fork", "Dent", "Dour", "Warp", "Luck", "Coffee", "Split", "Chip", "Hollow", "Heavy", "Legend", "Hickory", "Mesquite", "Nettle", "Rogue", "Charm", "Prickle", "Bead", "Sponge", "Whip", "Bald", "Frost", "Fog", "Oil", "Veil", "Cliff", "Volcano", "Rift", "Maze", "Proud", "Dew", "Mirror", "Shard", "Salt", "Pepper", "Honey", "Thread", "Bristle", "Ripple", "Glow", "Zenith"}
	var nouns = []string{"Head", "Crest", "Crown", "Tooth", "Fang", "Horn", "Frill", "Skull", "Bone", "Tongue", "Throat", "Voice", "Nose", "Snout", "Chin", "Eye", "Sight", "Seer", "Speaker", "Singer", "Song", "Chanter", "Howler", "Chatter", "Shrieker", "Shriek", "Jaw", "Bite", "Biter", "Neck", "Shoulder", "Fin", "Wing", "Arm", "Lifter", "Grasp", "Grabber", "Hand", "Paw", "Foot", "Finger", "Toe", "Thumb", "Talon", "Palm", "Touch", "Racer", "Runner", "Hoof", "Fly", "Flier", "Swoop", "Roar", "Hiss", "Hisser", "Snarl", "Dive", "Diver", "Rib", "Chest", "Back", "Ridge", "Leg", "Legs", "Tail", "Beak", "Walker", "Lasher", "Swisher", "Carver", "Kicker", "Roarer", "Crusher", "Spike", "Shaker", "Charger", "Hunter", "Weaver", "Crafter", "Binder", "Scribe", "Muse", "Snap", "Snapper", "Slayer", "Stalker", "Track", "Tracker", "Scar", "Scarer", "Fright", "Killer", "Death", "Doom", "Healer", "Saver", "Friend", "Foe", "Guardian", "Thunder", "Lightning", "Cloud", "Storm", "Forger", "Scale", "Hair", "Braid", "Nape", "Belly", "Thief", "Stealer", "Reaper", "Giver", "Taker", "Dancer", "Player", "Gambler", "Twister", "Turner", "Painter", "Dart", "Drifter", "Sting", "Stinger", "Venom", "Spur", "Ripper", "Swallow", "Devourer", "Knight", "Lady", "Lord", "Queen", "King", "Master", "Mistress", "Prince", "Princess", "Duke", "Dutchess", "Samurai", "Ninja", "Knave", "Slave", "Servant", "Sage", "Wizard", "Witch", "Warlock", "Warrior", "Jester", "Paladin", "Bard", "Trader", "Sword", "Shield", "Knife", "Dagger", "Arrow", "Bow", "Fighter", "Bane", "Follower", "Leader", "Scourge", "Watcher", "Cat", "Panther", "Tiger", "Cougar", "Puma", "Jaguar", "Ocelot", "Lynx", "Lion", "Leopard", "Ferret", "Weasel", "Wolverine", "Bear", "Raccoon", "Dog", "Wolf", "Kitten", "Puppy", "Cub", "Fox", "Hound", "Terrier", "Coyote", "Hyena", "Jackal", "Pig", "Horse", "Donkey", "Stallion", "Mare", "Zebra", "Antelope", "Gazelle", "Deer", "Buffalo", "Bison", "Boar", "Elk", "Whale", "Dolphin", "Shark", "Fish", "Minnow", "Salmon", "Ray", "Fisher", "Otter", "Gull", "Duck", "Goose", "Crow", "Raven", "Bird", "Eagle", "Raptor", "Hawk", "Falcon", "Moose", "Heron", "Owl", "Stork", "Crane", "Sparrow", "Robin", "Parrot", "Cockatoo", "Carp", "Lizard", "Gecko", "Iguana", "Snake", "Python", "Viper", "Boa", "Condor", "Vulture", "Spider", "Fly", "Scorpion", "Heron", "Oriole", "Toucan", "Bee", "Wasp", "Hornet", "Rabbit", "Bunny", "Hare", "Brow", "Mustang", "Ox", "Piper", "Soarer", "Flasher", "Moth", "Mask", "Hide", "Hero", "Antler", "Chill", "Chiller", "Gem", "Ogre", "Myth", "Elf", "Fairy", "Pixie", "Dragon", "Griffin", "Unicorn", "Pegasus", "Sprite", "Fancier", "Chopper", "Slicer", "Skinner", "Butterfly", "Legend", "Wanderer", "Rover", "Raver", "Loon", "Lancer", "Glass", "Glazer", "Flame", "Crystal", "Lantern", "Lighter", "Cloak", "Bell", "Ringer", "Keeper", "Centaur", "Bolt", "Catcher", "Whimsey", "Quester", "Rat", "Mouse", "Serpent", "Wyrm", "Gargoyle", "Thorn", "Whip", "Rider", "Spirit", "Sentry", "Bat", "Beetle", "Burn", "Cowl", "Stone", "Gem", "Collar", "Mark", "Grin", "Scowl", "Spear", "Razor", "Edge", "Seeker", "Jay", "Ape", "Monkey", "Gorilla", "Koala", "Kangaroo", "Yak", "Sloth", "Ant", "Roach", "Weed", "Seed", "Eater", "Razor", "Shirt", "Face", "Goat", "Mind", "Shift", "Rider", "Face", "Mole", "Vole", "Pirate", "Llama", "Stag", "Bug", "Cap", "Boot", "Drop", "Hugger", "Sargent", "Snagglefoot", "Carpet", "Curtain"}
	rand.Seed(time.Now().UnixNano())
	suffix := adjectives[rand.Intn(len(adjectives)-1)] + nouns[rand.Intn(len(nouns)-1)]
	// FIXME: we're not double checking the nickname still doesn't exist,
	// which would be really unlucky, but still...
	return suffix
}
