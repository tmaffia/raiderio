package raiderio

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// AffixesQuery is the query for a current mythic plus affixes request
type AffixesQuery struct {
	Region Region
	Locale string
}

// MythicPlusAffixes is the response from a mythic plus affixes request
type MythicPlusAffixes struct {
	Region         string  `json:"region"`
	Title          string  `json:"title"`
	LeaderboardURL string  `json:"leaderboard_url"`
	AffixDetails   []Affix `json:"affix_details"`
}

// MythicPlusStaticData is the response from a mythic plus static data request
type MythicPlusStaticData struct {
	Seasons []MythicPlusSeason `json:"seasons"`
}

// MythicPlusSeason is a single season entry in a static data response
type MythicPlusSeason struct {
	Slug             string              `json:"slug"`
	Name             string              `json:"name"`
	BlizzardSeasonID int                 `json:"blizzard_season_id"`
	IsMainSeason     bool                `json:"is_main_season"`
	ShortName        string              `json:"short_name"`
	SeasonalAffix    *string             `json:"seasonal_affix"`
	Starts           RegionTimestamps    `json:"starts"`
	Ends             RegionTimestamps    `json:"ends"`
	Dungeons         []MythicPlusDungeon `json:"dungeons"`
}

// MythicPlusDungeon is the lightweight dungeon shape used in static data
type MythicPlusDungeon struct {
	ID                   int    `json:"id"`
	ChallengeModeID      int    `json:"challenge_mode_id"`
	Slug                 string `json:"slug"`
	Name                 string `json:"name"`
	ShortName            string `json:"short_name"`
	KeystoneTimerSeconds int    `json:"keystone_timer_seconds"`
	IconURL              string `json:"icon_url"`
	BackgroundImageURL   string `json:"background_image_url"`
}

// MythicPlusRunsQuery is the query for a mythic plus runs leaderboard request.
// All fields are optional; the API applies its own defaults when omitted
type MythicPlusRunsQuery struct {
	Region  Region
	Season  string
	Dungeon string
	Affixes string
	Page    int
}

// MythicPlusRuns is the response from a mythic plus runs leaderboard request
type MythicPlusRuns struct {
	Rankings []MythicPlusRanking `json:"rankings"`
}

// MythicPlusRanking is a single ranked entry in a runs leaderboard response
type MythicPlusRanking struct {
	Rank  int                `json:"rank"`
	Score float64            `json:"score"`
	Run   MythicPlusRunEntry `json:"run"`
}

// RunDetailsQuery is the query for a single mythic plus run
type RunDetailsQuery struct {
	Season string
	ID     int
}

// MythicPlusRunEntry represents a single keystone run, shared by the
// runs leaderboard and run-details responses
type MythicPlusRunEntry struct {
	Season          string            `json:"season"`
	Status          string            `json:"status"`
	Dungeon         RunDungeon        `json:"dungeon"`
	KeystoneRunID   int               `json:"keystone_run_id"`
	MythicLevel     int               `json:"mythic_level"`
	ClearTimeMs     int               `json:"clear_time_ms"`
	KeystoneTimeMs  int               `json:"keystone_time_ms"`
	CompletedAt     time.Time         `json:"completed_at"`
	NumChests       int               `json:"num_chests"`
	TimeRemainingMs int               `json:"time_remaining_ms"`
	Faction         string            `json:"faction"`
	Score           float64           `json:"score"`
	WeeklyModifiers []Affix           `json:"weekly_modifiers"`
	Roster          []RunRosterMember `json:"roster"`
}

// RunDungeon is the detailed dungeon shape embedded in run responses.
// This differs from MythicPlusDungeon, the lighter shape used in static data
type RunDungeon struct {
	Type                   string `json:"type"`
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	ShortName              string `json:"short_name"`
	Slug                   string `json:"slug"`
	ExpansionID            int    `json:"expansion_id"`
	IconURL                string `json:"icon_url"`
	Patch                  string `json:"patch"`
	WowInstanceID          int    `json:"wowInstanceId"`
	MapChallengeModeID     int    `json:"map_challenge_mode_id"`
	KeystoneTimerMs        int    `json:"keystone_timer_ms"`
	NumBosses              int    `json:"num_bosses"`
	GroupFinderActivityIDs []int  `json:"group_finder_activity_ids"`
}

// RunRosterMember is a single character in a keystone run roster
type RunRosterMember struct {
	Character RunCharacter `json:"character"`
	Role      string       `json:"role"`
	Loadout   string       `json:"loadout"`
}

// RunCharacter is the character shape embedded in a run roster.
// Note: stream/recruitment/talent-tree fields present on the wire are
// intentionally omitted; add them if a consumer needs more than identity,
// class, spec, and location.
type RunCharacter struct {
	ID      int            `json:"id"`
	Name    string         `json:"name"`
	Class   ClassInfo      `json:"class"`
	Race    RaceInfo       `json:"race"`
	Faction string         `json:"faction"`
	Spec    MythicPlusSpec `json:"spec"`
	Realm   Realm          `json:"realm"`
	Region  RegionInfo     `json:"region"`
}

// ClassInfo is a class id/name/slug triple, used in run character responses
type ClassInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// RaceInfo is a race id/name/slug/faction quad, used in run character responses
type RaceInfo struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Faction string `json:"faction"`
}

// ScoreTiersQuery is the query for a mythic plus score-color tiers request
type ScoreTiersQuery struct {
	Season string
}

// ScoreTier is a single score/color breakpoint
type ScoreTier struct {
	Score      int        `json:"score"`
	RgbHex     string     `json:"rgbHex"`
	RgbDecimal int64      `json:"rgbDecimal"`
	RgbFloat   [3]float64 `json:"rgbFloat"`
	RgbInteger [3]int     `json:"rgbInteger"`
}

// SeasonCutoffsQuery is the query for a mythic plus season cutoffs request
type SeasonCutoffsQuery struct {
	Season string
	Region Region
}

// SeasonCutoffs is the response from a season cutoffs request.
// Note: the season-specific keystoneMyth/allTimedNN/graphData keys present on
// the wire are intentionally dropped; add typed fields if a consumer needs them.
type SeasonCutoffs struct {
	Cutoffs SeasonCutoffData `json:"cutoffs"`
}

// SeasonCutoffData holds the region and percentile breakpoints for a season
type SeasonCutoffData struct {
	UpdatedAt string            `json:"updatedAt"`
	Region    RegionInfo        `json:"region"`
	P999      *SeasonCutoffTier `json:"p999,omitempty"`
	P990      *SeasonCutoffTier `json:"p990,omitempty"`
	P900      *SeasonCutoffTier `json:"p900,omitempty"`
	P750      *SeasonCutoffTier `json:"p750,omitempty"`
	P600      *SeasonCutoffTier `json:"p600,omitempty"`
}

// SeasonCutoffTier is one percentile breakpoint, broken down by faction
type SeasonCutoffTier struct {
	Horde         SeasonCutoffStat `json:"horde"`
	HordeColor    string           `json:"hordeColor"`
	Alliance      SeasonCutoffStat `json:"alliance"`
	AllianceColor string           `json:"allianceColor"`
	All           SeasonCutoffStat `json:"all"`
	AllColor      string           `json:"allColor"`
}

// SeasonCutoffStat is the population/score data for one faction at one tier
type SeasonCutoffStat struct {
	Quantile                   float64 `json:"quantile"`
	QuantileMinValue           float64 `json:"quantileMinValue"`
	QuantilePopulationCount    int     `json:"quantilePopulationCount"`
	QuantilePopulationFraction float64 `json:"quantilePopulationFraction"`
	TotalPopulationCount       int     `json:"totalPopulationCount"`
}

// LeaderboardCapacityQuery is the query for a mythic plus leaderboard capacity request
type LeaderboardCapacityQuery struct {
	Region Region
	Realm  string
	Scope  string
}

// LeaderboardCapacity is the response from a leaderboard capacity request
type LeaderboardCapacity struct {
	RealmListing LeaderboardCapacityListing `json:"realmListing"`
}

// LeaderboardCapacityListing is the region, affixes, and realm breakdown
// in a leaderboard capacity response
type LeaderboardCapacityListing struct {
	Region  RegionInfo                 `json:"region"`
	Affixes []LeaderboardCapacityAffix `json:"affixes"`
	Realms  []LeaderboardCapacityRealm `json:"realms"`
}

// LeaderboardCapacityAffix is an affix with localized name/description maps
type LeaderboardCapacityAffix struct {
	ID          int               `json:"id"`
	Icon        string            `json:"icon"`
	Slug        string            `json:"slug"`
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
}

// LeaderboardCapacityRealm is a connected-realm group and its per-dungeon capacity
type LeaderboardCapacityRealm struct {
	ID              int                          `json:"id"`
	ConnectedRealms []Realm                      `json:"connectedRealms"`
	Dungeons        []LeaderboardCapacityDungeon `json:"dungeons"`
}

// LeaderboardCapacityDungeon pairs a dungeon with its lowest qualifying run.
// Note: "lowest" is undocumented and always null in observed responses; it is
// kept as raw JSON until a non-null sample is available for a typed field.
type LeaderboardCapacityDungeon struct {
	Dungeon RunDungeon      `json:"dungeon"`
	Lowest  json.RawMessage `json:"lowest"`
}

// Periods is the response from a periods request
type Periods struct {
	Periods []RegionPeriods `json:"periods"`
}

// RegionPeriods is the previous/current/next weekly period for one region
type RegionPeriods struct {
	Region   string       `json:"region"`
	Previous PeriodWindow `json:"previous"`
	Current  PeriodWindow `json:"current"`
	Next     PeriodWindow `json:"next"`
}

// PeriodWindow is a single weekly period's id and start/end time
type PeriodWindow struct {
	Period int       `json:"period"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}

func (q *AffixesQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.Region == "" {
		return ErrInvalidRegion
	}
	return nil
}

func (q *AffixesQuery) params() url.Values {
	params := url.Values{"region": {string(q.Region)}}
	if q.Locale != "" {
		params.Set("locale", q.Locale)
	}
	return params
}

// validate only rejects a nil query; every field is optional.
func (q *MythicPlusRunsQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	return nil
}

func (q *MythicPlusRunsQuery) params() url.Values {
	params := url.Values{}
	if q.Region != "" {
		params.Set("region", string(q.Region))
	}
	if q.Season != "" {
		params.Set("season", q.Season)
	}
	if q.Dungeon != "" {
		params.Set("dungeon", q.Dungeon)
	}
	if q.Affixes != "" {
		params.Set("affixes", q.Affixes)
	}
	if q.Page != 0 {
		params.Set("page", strconv.Itoa(q.Page))
	}
	return params
}

func (q *RunDetailsQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.Season == "" {
		return ErrInvalidSeason
	}
	if q.ID <= 0 {
		return ErrInvalidRunID
	}
	return nil
}

func (q *RunDetailsQuery) params() url.Values {
	return url.Values{
		"season": {q.Season},
		"id":     {strconv.Itoa(q.ID)},
	}
}

// validate only rejects a nil query; season is optional.
func (q *ScoreTiersQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	return nil
}

func (q *ScoreTiersQuery) params() url.Values {
	params := url.Values{}
	if q.Season != "" {
		params.Set("season", q.Season)
	}
	return params
}

func (q *SeasonCutoffsQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.Season == "" {
		return ErrInvalidSeason
	}
	if q.Region == "" {
		return ErrInvalidRegion
	}
	return nil
}

func (q *SeasonCutoffsQuery) params() url.Values {
	return url.Values{
		"season": {q.Season},
		"region": {string(q.Region)},
	}
}

func (q *LeaderboardCapacityQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.Region == "" {
		return ErrInvalidRegion
	}
	return nil
}

func (q *LeaderboardCapacityQuery) params() url.Values {
	params := url.Values{"region": {string(q.Region)}}
	if q.Realm != "" {
		params.Set("realm", q.Realm)
	}
	if q.Scope != "" {
		params.Set("scope", q.Scope)
	}
	return params
}
