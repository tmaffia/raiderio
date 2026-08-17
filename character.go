package raiderio

// CharacterQuery is a struct that represents the query parameters
// sent for a character profile request
type CharacterQuery struct {
	Region        Region
	Realm         string
	Name          string
	TalentLoadout bool
	Gear          bool
	fields        []string
}

// Character is a struct that represents the response from
// a character profile request
type Character struct {
	Name              string        `json:"name"`
	Race              string        `json:"race"`
	Class             string        `json:"class"`
	ActiveSpec        string        `json:"active_spec_name"`
	ActiveRole        string        `json:"active_spec_role"`
	Gender            string        `json:"gender"`
	Faction           string        `json:"faction"`
	AchievementPoints int64         `json:"achievement_points"`
	HonorableKills    int64         `json:"honorable_kills"`
	ThumbnailUrl      string        `json:"thumbnail_url"`
	Region            string        `json:"region"`
	Realm             string        `json:"realm"`
	LastCrawledAt     string        `json:"last_crawled_at"`
	ProfileUrl        string        `json:"profile_url"`
	ProfileBanner     string        `json:"profile_banner"`
	TalentLoadout     TalentLoadout `json:"talentLoadout"`
	Gear              Gear          `json:"gear"`
}

// Gear is a struct that represents the gear of a character
// in a character profile response
type Gear struct {
	UpdatedAt         string `json:"updated_at"`
	ItemLevelEquipped int    `json:"item_level_equipped"`
	ItemLevelTotal    int    `json:"item_level_total"`
	Items             Items  `json:"items"`
}

// Items is a struct that represents the items of a character
// in a character profile response
type Items struct {
	Head     Item `json:"head"`
	Neck     Item `json:"neck"`
	Shoulder Item `json:"shoulder"`
	Back     Item `json:"back"`
	Chest    Item `json:"chest"`
	Wrist    Item `json:"wrist"`
	Hands    Item `json:"hands"`
	Waist    Item `json:"waist"`
	Legs     Item `json:"legs"`
	Feet     Item `json:"feet"`
	Finger1  Item `json:"finger1"`
	Finger2  Item `json:"finger2"`
	Trinket1 Item `json:"trinket1"`
	Trinket2 Item `json:"trinket2"`
	Mainhand Item `json:"mainhand"`
	Offhand  Item `json:"offhand"`
	Shirt    Item `json:"shirt"`
	Tabard   Item `json:"tabard"`
}

// Item is a struct that represents a single item
type Item struct {
	ID          int    `json:"item_id"`
	ItemLevel   int    `json:"item_level"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	ItemQuality int    `json:"item_quality"`
	IsLegendary bool   `json:"is_legendary"`
	Gems        []int  `json:"gems"`
	Bonuses     []int  `json:"bonuses"`
}

// TalentLoadout is a struct of a talent loadout
// It includes the spec id and talent loadout string
type TalentLoadout struct {
	LoadoutSpecID int    `json:"loadout_spec_id"`
	LoadoutText   string `json:"loadout_text"`
}

// validateCharacterQuery creates and validates a CharacterQuery struct
// It returns an error if any of the required parameters are empty
// or if the fields are invalid
func validateCharacterQuery(cq *CharacterQuery) error {
	if cq.Region == "" {
		return ErrInvalidRegion
	}

	if cq.Realm == "" {
		return ErrInvalidRealm
	}

	if cq.Name == "" {
		return ErrInvalidCharName
	}

	if cq.TalentLoadout {
		cq.fields = append(cq.fields, "talents")
	}

	if cq.Gear {
		cq.fields = append(cq.fields, "gear")
	}

	return nil
}
