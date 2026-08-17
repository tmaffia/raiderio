package raiderio

import (
	"net/url"
	"strings"
	"time"
)

// CharacterQuery is a struct that represents the query parameters
// sent for a character profile request
type CharacterQuery struct {
	Region                                   Region
	Realm                                    string
	Name                                     string
	TalentLoadout                            bool
	Gear                                     bool
	Guild                                    bool
	RaidProgression                          bool
	MythicPlusScoresBySeason                 []string
	MythicPlusRanks                          bool
	PreviousMythicPlusRanks                  bool
	MythicPlusRecentRuns                     bool
	MythicPlusBestRuns                       bool
	MythicPlusBestRunsAll                    bool
	MythicPlusHighestLevelRuns               bool
	MythicPlusWeeklyHighestLevelRuns         bool
	MythicPlusPreviousWeeklyHighestLevelRuns bool
	RaidAchievementMeta                      []string
	RaidAchievementCurve                     []string
}

// Character is a struct that represents the response from
// a character profile request
type Character struct {
	Name                                     string                     `json:"name"`
	Race                                     string                     `json:"race"`
	Class                                    string                     `json:"class"`
	ActiveSpec                               string                     `json:"active_spec_name"`
	ActiveRole                               string                     `json:"active_spec_role"`
	Gender                                   string                     `json:"gender"`
	Faction                                  string                     `json:"faction"`
	AchievementPoints                        int64                      `json:"achievement_points"`
	HonorableKills                           int64                      `json:"honorable_kills"`
	ThumbnailUrl                             string                     `json:"thumbnail_url"`
	Region                                   string                     `json:"region"`
	Realm                                    string                     `json:"realm"`
	LastCrawledAt                            time.Time                  `json:"last_crawled_at"`
	ProfileUrl                               string                     `json:"profile_url"`
	ProfileBanner                            string                     `json:"profile_banner"`
	TalentLoadout                            TalentLoadout              `json:"talentLoadout"`
	Gear                                     Gear                       `json:"gear"`
	Guild                                    CharacterGuild             `json:"guild"`
	RaidProgression                          map[string]RaidProgression `json:"raid_progression"`
	MythicPlusScoresBySeason                 []MythicPlusSeasonScores   `json:"mythic_plus_scores_by_season"`
	MythicPlusRanks                          map[string]MythicPlusRank  `json:"mythic_plus_ranks"`
	PreviousMythicPlusRanks                  map[string]MythicPlusRank  `json:"previous_mythic_plus_ranks"`
	MythicPlusRecentRuns                     []MythicPlusRun            `json:"mythic_plus_recent_runs"`
	MythicPlusBestRuns                       []MythicPlusRun            `json:"mythic_plus_best_runs"`
	MythicPlusHighestLevelRuns               []MythicPlusRun            `json:"mythic_plus_highest_level_runs"`
	MythicPlusWeeklyHighestLevelRuns         []MythicPlusRun            `json:"mythic_plus_weekly_highest_level_runs"`
	MythicPlusPreviousWeeklyHighestLevelRuns []MythicPlusRun            `json:"mythic_plus_previous_weekly_highest_level_runs"`
	RaidAchievementMeta                      []RaidAchievementMeta      `json:"raid_achievement_meta"`
	RaidAchievementCurve                     []RaidAchievementCurve     `json:"raid_achievement_curve"`
}

// CharacterGuild is the guild block on a character profile
type CharacterGuild struct {
	Name  string `json:"name"`
	Realm string `json:"realm"`
}

// MythicPlusSeasonScores is one season entry in mythic_plus_scores_by_season
type MythicPlusSeasonScores struct {
	Season   string                       `json:"season"`
	Scores   MythicPlusScores             `json:"scores"`
	Segments map[string]ScoreColorSegment `json:"segments"`
}

// MythicPlusScores is the score breakdown for a season
type MythicPlusScores struct {
	All    float64 `json:"all"`
	DPS    float64 `json:"dps"`
	Healer float64 `json:"healer"`
	Tank   float64 `json:"tank"`
	Spec0  float64 `json:"spec_0"`
	Spec1  float64 `json:"spec_1"`
	Spec2  float64 `json:"spec_2"`
	Spec3  float64 `json:"spec_3"`
}

// ScoreColorSegment is a score plus its display color
type ScoreColorSegment struct {
	Score float64 `json:"score"`
	Color string  `json:"color"`
}

// MythicPlusRank is a world/region/realm rank triple
type MythicPlusRank struct {
	World  int `json:"world"`
	Region int `json:"region"`
	Realm  int `json:"realm"`
}

// MythicPlusRun is a keystone run on a character profile
type MythicPlusRun struct {
	Dungeon             string         `json:"dungeon"`
	ShortName           string         `json:"short_name"`
	MythicLevel         int            `json:"mythic_level"`
	CompletedAt         time.Time      `json:"completed_at"`
	ClearTimeMs         int            `json:"clear_time_ms"`
	KeystoneRunID       int            `json:"keystone_run_id"`
	ParTimeMs           int            `json:"par_time_ms"`
	NumKeystoneUpgrades int            `json:"num_keystone_upgrades"`
	MapChallengeModeID  int            `json:"map_challenge_mode_id"`
	ZoneID              int            `json:"zone_id"`
	ZoneExpansionID     int            `json:"zone_expansion_id"`
	IconURL             string         `json:"icon_url"`
	BackgroundImageURL  string         `json:"background_image_url"`
	Score               float64        `json:"score"`
	Affixes             []Affix        `json:"affixes"`
	URL                 string         `json:"url"`
	Spec                MythicPlusSpec `json:"spec"`
	Role                string         `json:"role"`
}

// Affix is a mythic plus affix
type Affix struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IconURL     string `json:"icon_url"`
	WowheadURL  string `json:"wowhead_url"`
}

// MythicPlusSpec is the spec played on a keystone run
type MythicPlusSpec struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	ClassID int    `json:"class_id"`
	Role    string `json:"role"`
	IsMelee bool   `json:"is_melee"`
	Patch   string `json:"patch"`
	Ordinal int    `json:"ordinal"`
}

// RaidAchievementMeta is one raid-tier meta on a character profile
type RaidAchievementMeta struct {
	Tier                  string            `json:"tier"`
	CompletedCount        int               `json:"completed_count"`
	TotalCount            int               `json:"total_count"`
	MetaAchievement       RaidAchievement   `json:"meta_achievement"`
	CompletedAchievements []RaidAchievement `json:"completed_achievements"`
	RemainingAchievements []RaidAchievement `json:"remaining_achievements"`
}

// RaidAchievement is a raid achievement id/slug pair
type RaidAchievement struct {
	ID        int       `json:"id"`
	Raid      string    `json:"raid"`
	Timestamp time.Time `json:"timestamp"`
}

// RaidAchievementCurve is AOTC / Cutting Edge status for a raid
type RaidAchievementCurve struct {
	Raid        string    `json:"raid"`
	AOTC        time.Time `json:"aotc"`
	CuttingEdge time.Time `json:"cutting_edge"`
}

// Gear is a struct that represents the gear of a character
// in a character profile response
type Gear struct {
	UpdatedAt         time.Time `json:"updated_at"`
	ItemLevelEquipped int       `json:"item_level_equipped"`
	ItemLevelTotal    int       `json:"item_level_total"`
	Items             Items     `json:"items"`
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

// validate checks the required fields of a CharacterQuery.
func (q *CharacterQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}

	if q.Region == "" {
		return ErrInvalidRegion
	}

	if q.Realm == "" {
		return ErrInvalidRealm
	}

	if q.Name == "" {
		return ErrInvalidCharName
	}

	return nil
}

// params builds the URL parameters for a character profile request.
func (q *CharacterQuery) params() url.Values {
	params := url.Values{
		"region": {string(q.Region)},
		"realm":  {q.Realm},
		"name":   {q.Name},
	}
	if fields := q.fields(); len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}
	return params
}

// fields builds the optional request fields from the boolean/slice toggles.
func (q *CharacterQuery) fields() []string {
	var fields []string
	if q.TalentLoadout {
		fields = append(fields, "talents")
	}
	if q.Gear {
		fields = append(fields, "gear")
	}
	if q.Guild {
		fields = append(fields, "guild")
	}
	if q.RaidProgression {
		fields = append(fields, "raid_progression")
	}
	if len(q.MythicPlusScoresBySeason) > 0 {
		fields = append(fields, "mythic_plus_scores_by_season:"+strings.Join(q.MythicPlusScoresBySeason, ":"))
	}
	if q.MythicPlusRanks {
		fields = append(fields, "mythic_plus_ranks")
	}
	if q.PreviousMythicPlusRanks {
		fields = append(fields, "previous_mythic_plus_ranks")
	}
	if q.MythicPlusRecentRuns {
		fields = append(fields, "mythic_plus_recent_runs")
	}
	switch {
	case q.MythicPlusBestRunsAll:
		fields = append(fields, "mythic_plus_best_runs:all")
	case q.MythicPlusBestRuns:
		fields = append(fields, "mythic_plus_best_runs")
	}
	if q.MythicPlusHighestLevelRuns {
		fields = append(fields, "mythic_plus_highest_level_runs")
	}
	if q.MythicPlusWeeklyHighestLevelRuns {
		fields = append(fields, "mythic_plus_weekly_highest_level_runs")
	}
	if q.MythicPlusPreviousWeeklyHighestLevelRuns {
		fields = append(fields, "mythic_plus_previous_weekly_highest_level_runs")
	}
	if len(q.RaidAchievementMeta) > 0 {
		fields = append(fields, "raid_achievement_meta:"+strings.Join(q.RaidAchievementMeta, ":"))
	}
	if len(q.RaidAchievementCurve) > 0 {
		fields = append(fields, "raid_achievement_curve:"+strings.Join(q.RaidAchievementCurve, ":"))
	}

	return fields
}
