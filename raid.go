package raiderio

import (
	"net/url"
	"strconv"
	"time"
)

// RaidQuery is the query for raid-scoped endpoints: hall-of-fame and
// progression. Slug, Difficulty, and Region are all required.
type RaidQuery struct {
	Slug       string
	Difficulty RaidDifficulty
	Region     Region
}

// RaidRankingsQuery is the query for the raid-rankings endpoint. Embeds
// RaidQuery for the required Slug/Difficulty/Region; Realm, Limit, and Page
// are additional and only honored by this endpoint.
type RaidRankingsQuery struct {
	RaidQuery
	Realm string
	Limit int
	Page  int
}

// RaidRankings is a struct that represents the response from a
// raid rankings request
type RaidRankings struct {
	RaidRanking []RaidRanking `json:"raidRankings"`
}

// RaidRanking is a struct that represents a raid ranking in a
// raid rankings response from the api
// Unfortunately the "Guild" object differs in structure from the
// guild profile response. This requires a separate struct
type RaidRanking struct {
	Rank               int       `json:"rank"`
	RegionalRank       int       `json:"region_rank"`
	Guild              RaidGuild `json:"guild"`
	EncountersDefeated []struct {
		Slug           string    `json:"slug"`
		LastDefeatedAt time.Time `json:"lastDefeated"`
		FirstDefeated  time.Time `json:"firstDefeated"`
	} `json:"encountersDefeated"`
	EncountersPulled []struct {
		Id             int       `json:"id"`
		Slug           string    `json:"slug"`
		Pulls          int       `json:"numPulls"`
		PullsStartedAt time.Time `json:"pullStartedAt"`
		BestPercent    float32   `json:"bestPercent"`
		IsDefeated     bool      `json:"isDefeated"`
	} `json:"encountersPulled"`
}

// Realm is a struct that represents a realm available in the Raider.IO API
type Realm struct {
	Id               int    `json:"id"`
	ConnectedRealmId int    `json:"connectedRealmId"`
	Name             string `json:"name"`
	AltName          string `json:"altName"`
	Slug             string `json:"slug"`
	AltSlug          string `json:"altSlug"`
	Locale           string `json:"locale"`
	IsConnected      bool   `json:"isConnected"`
}

// RegionInfo is a struct that represents a region as returned in API responses
type RegionInfo struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ShortName string `json:"short_name"`
}

// RegionTimestamps holds a per-region timestamp, used for raid and
// mythic plus season start/end dates
type RegionTimestamps struct {
	Us time.Time `json:"us"`
	Eu time.Time `json:"eu"`
	Tw time.Time `json:"tw"`
	Kr time.Time `json:"kr"`
	Cn time.Time `json:"cn"`
}

// RaidGuild represents a guild in raid-related responses
// This structure is used in RaidRankings, BossRankings, HallOfFame, etc.
type RaidGuild struct {
	Id      int        `json:"id"`
	Name    string     `json:"name"`
	Faction string     `json:"faction"`
	Realm   Realm      `json:"realm"`
	Region  RegionInfo `json:"region"`
	Path    string     `json:"path"`
	Logo    string     `json:"logo"`
	Color   string     `json:"color"`
}

// RaidProgression is a struct that contains raid progression for a tier,
// used in both guild and character profile responses
type RaidProgression struct {
	Summary     string `json:"summary"`
	ExpansionID int    `json:"expansion_id"`
	Bosses      int    `json:"total_bosses"`
	NormalKills int    `json:"normal_bosses_killed"`
	HeroicKills int    `json:"heroic_bosses_killed"`
	MythicKills int    `json:"mythic_bosses_killed"`
}

// GuildRaidRanking is a struct that contains the raid rankings of a guild
// in a guild profile response
// Includes Normal Heroic and Mythic rankings
type GuildRaidRanking struct {
	Normal struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Realm  int `json:"realm"`
	} `json:"normal"`

	Heroic struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Realm  int `json:"realm"`
	} `json:"heroic"`

	Mythic struct {
		World  int `json:"world"`
		Region int `json:"region"`
		Realm  int `json:"realm"`
	} `json:"mythic"`
}

// Raids is a struct that represents the response from a
// raid static data request
type Raids struct {
	Raids []Raid `json:"raids"`
}

// Raid is a struct that represents a raid in a raid static
// data response. Includes raid encounters and other static data
type Raid struct {
	Id        int              `json:"id"`
	Slug      string           `json:"slug"`
	Name      string           `json:"name"`
	ShortName string           `json:"short_name"`
	Icon      string           `json:"icon"`
	Starts    RegionTimestamps `json:"starts"`
	Ends      RegionTimestamps `json:"ends"`

	Encounters []Encounter `json:"encounters"`
}

// Encounter is a struct that represents an encounter in a raid
// in a raid static data response
type Encounter struct {
	Id   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// RaidDifficulty is a string type that represents the difficulty of a raid
// in a raid request
// Options are "normal", "heroic", and "mythic"
type RaidDifficulty string

// Options for different raid difficulties

const (
	NORMAL_RAID RaidDifficulty = "normal"
	HEROIC_RAID RaidDifficulty = "heroic"
	MYTHIC_RAID RaidDifficulty = "mythic"
)

// Includes BossKillData along with the roster of characters
// which were present for the first kill
type BossKill struct {
	Kill   BossKillData
	Roster []RosterCharacter
}

// RosterCharacter is a boss-kill roster member. Profile uses Character;
// the two payloads do not share a spec field.
type RosterCharacter struct {
	Name          string
	Class         string
	Spec          string
	Realm         string
	Region        string
	TalentLoadout TalentLoadout
	Gear          Gear
}

// BossKillData provides metadata for the guilds first boss kill
// Includes timestamps and Item Levels etc...
type BossKillData struct {
	PulledAt             time.Time     `json:"pulledAt"`
	DefeatedAt           time.Time     `json:"defeatedAt"`
	Duration             time.Duration `json:"duration"`
	IsSuccess            bool          `json:"isSuccess"`
	ItemLevelEquippedAvg float32       `json:"itemLevelEquippedAvg"`
	ItemLevelEquippedMax float32       `json:"itemLevelEquippedMax"`
	ItemLevelEquippedMin float32       `json:"itemLevelEquippedMin"`
}

// Unexported wire types for /guilds/boss-kill. Mapped to RosterCharacter.
type bossKillResp struct {
	Kill struct {
		PulledAt             time.Time `json:"pulledAt"`
		DefeatedAt           time.Time `json:"defeatedAt"`
		DurationMs           int       `json:"durationMs"`
		IsSuccess            bool      `json:"isSuccess"`
		ItemLevelEquippedAvg float32   `json:"itemLevelEquippedAvg"`
		ItemLevelEquippedMax float32   `json:"itemLevelEquippedMax"`
		ItemLevelEquippedMin float32   `json:"itemLevelEquippedMin"`
	}
	Roster []bossKillCharacter `json:"roster"`
}
type bossKillCharacter struct {
	Character struct {
		Name  string `json:"name"`
		Class struct {
			Slug string `json:"slug"`
		} `json:"class"`
		Spec struct {
			Slug string `json:"slug"`
		} `json:"spec"`
		TalentLoadout struct {
			LoadoutSpecID int    `json:"loadoutSpecId"`
			LoadoutText   string `json:"loadoutText"`
		} `json:"talentLoadout"`
		Realm struct {
			Slug string `json:"slug"`
		} `json:"realm"`
		Region struct {
			Slug string `json:"slug"`
		} `json:"region"`
		ItemLevelEquipped float32 `json:"itemLevelEquipped"`
	} `json:"character"`
}

// GuildBossKillQuery requires all fields to be valid when sending
// a request to the api. Use GetRaids() to see a list of raids and bosses
type GuildBossKillQuery struct {
	Region     Region
	Realm      string
	GuildName  string
	RaidSlug   string
	BossSlug   string
	Difficulty RaidDifficulty
}

// The /guilds/boss-kill endpoint returns an enormous JSON structure for each
// character in the raid roster; this library maps it to a simplified shape.
func mapBossKill(resp bossKillResp) BossKill {
	kd := BossKillData{
		PulledAt:             resp.Kill.PulledAt,
		DefeatedAt:           resp.Kill.DefeatedAt,
		Duration:             time.Duration(resp.Kill.DurationMs) * time.Millisecond,
		IsSuccess:            resp.Kill.IsSuccess,
		ItemLevelEquippedAvg: resp.Kill.ItemLevelEquippedAvg,
		ItemLevelEquippedMax: resp.Kill.ItemLevelEquippedMax,
		ItemLevelEquippedMin: resp.Kill.ItemLevelEquippedMin,
	}
	return BossKill{
		Kill:   kd,
		Roster: unmarshalBossKillRoster(&resp),
	}
}

func unmarshalBossKillRoster(k *bossKillResp) []RosterCharacter {
	var chars []RosterCharacter
	for _, c := range k.Roster {
		chars = append(chars, RosterCharacter{
			Name:   c.Character.Name,
			Class:  c.Character.Class.Slug,
			Spec:   c.Character.Spec.Slug,
			Realm:  c.Character.Realm.Slug,
			Region: c.Character.Region.Slug,
			TalentLoadout: TalentLoadout{
				LoadoutSpecID: c.Character.TalentLoadout.LoadoutSpecID,
				LoadoutText:   c.Character.TalentLoadout.LoadoutText,
			},
			Gear: Gear{
				ItemLevelEquipped: int(c.Character.ItemLevelEquipped),
			},
		})
	}
	return chars
}

func (q *GuildBossKillQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}

	if q.Region == "" {
		return ErrInvalidRegion
	}

	if q.Realm == "" {
		return ErrInvalidRealm
	}

	if q.GuildName == "" {
		return ErrInvalidGuildName
	}

	if q.RaidSlug == "" {
		return ErrInvalidRaidName
	}

	if q.BossSlug == "" {
		return ErrInvalidBoss
	}

	if !raidDifficultyValid(q.Difficulty) {
		return ErrInvalidRaidDiff
	}

	return nil
}

func (q *GuildBossKillQuery) params() url.Values {
	return url.Values{
		"raid":       {q.RaidSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {string(q.Region)},
		"realm":      {q.Realm},
		"guild":      {q.GuildName},
		"boss":       {q.BossSlug},
	}
}

func raidDifficultyValid(d RaidDifficulty) bool {
	return d == NORMAL_RAID || d == HEROIC_RAID || d == MYTHIC_RAID
}

func (q *RaidRankingsQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if err := q.RaidQuery.validate(); err != nil {
		return err
	}
	if q.Limit < 0 {
		return ErrLimitOutOfBounds
	}
	if q.Page < 0 {
		return ErrPageOutOfBounds
	}
	return nil
}

func (r *Raids) GetRaidBySlug(slug string) (*Raid, error) {
	for _, raid := range r.Raids {
		if raid.Slug == slug {
			return &raid, nil
		}
	}
	return nil, ErrInvalidRaid
}

// BossRankingsQuery represents the query parameters for boss rankings
type BossRankingsQuery struct {
	RaidSlug   string
	BossSlug   string
	Difficulty RaidDifficulty
	Region     Region
	Realm      string
}

// BossRankings represents the response from a boss rankings request
type BossRankings struct {
	BossRankings []BossRanking `json:"bossRankings"`
}

// BossRanking represents a single ranking entry for a boss
type BossRanking struct {
	Rank               int       `json:"rank"`
	RegionRank         int       `json:"regionRank"`
	Guild              RaidGuild `json:"guild"`
	EncountersDefeated []struct {
		Slug          string    `json:"slug"`
		LastDefeated  time.Time `json:"lastDefeated"`
		FirstDefeated time.Time `json:"firstDefeated"`
	} `json:"encountersDefeated"`
}

// HallOfFame represents the response from a hall of fame request
type HallOfFame struct {
	HallOfFame HallOfFameEntry `json:"hallOfFame"`
}

// HallOfFameEntry represents a single entry in the hall of fame
type HallOfFameEntry struct {
	BossKills []HallOfFameBossKill `json:"bossKills"`
}

// HallOfFameBossKill is one boss in a hall-of-fame response
type HallOfFameBossKill struct {
	Boss        string                `json:"boss"`
	BossSummary HallOfFameBossSummary `json:"bossSummary"`
	DefeatedBy  HallOfFameDefeatedBy  `json:"defeatedBy"`
}

// HallOfFameBossSummary is the static summary of a hall-of-fame boss
type HallOfFameBossSummary struct {
	EncounterId int    `json:"encounterId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Ordinal     int    `json:"ordinal"`
	WingId      int    `json:"wingId"`
	IconUrl     string `json:"iconUrl"`
}

// HallOfFameDefeatedBy is the set of guilds that killed a hall-of-fame boss
type HallOfFameDefeatedBy struct {
	TotalCount int               `json:"totalCount"`
	Guilds     []HallOfFameGuild `json:"guilds"`
}

// HallOfFameGuild is a guild that killed a hall-of-fame boss
type HallOfFameGuild struct {
	Guild      RaidGuild `json:"guild"`
	DefeatedAt time.Time `json:"defeatedAt"`
}

// RaidProgressionResponse represents the response from a raid progression request
type RaidProgressionResponse struct {
	Progression []RaidProgressionEntry `json:"progression"`
}

// RaidProgressionEntry represents a single progression entry
type RaidProgressionEntry struct {
	Progress    int                    `json:"progress"`
	TotalGuilds int                    `json:"totalGuilds"`
	Guilds      []RaidProgressionGuild `json:"guilds"`
}

// RaidProgressionGuild is a guild in a raid-progression entry
type RaidProgressionGuild struct {
	DefeatedAt time.Time `json:"defeatedAt"`
	Guild      RaidGuild `json:"guild"`
}

func (q *BossRankingsQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.RaidSlug == "" {
		return ErrInvalidRaidName
	}
	if q.BossSlug == "" {
		return ErrInvalidBoss
	}
	if !raidDifficultyValid(q.Difficulty) {
		return ErrInvalidRaidDiff
	}
	if q.Region == "" {
		return ErrInvalidRegion
	}
	return nil
}

func (q *BossRankingsQuery) params() url.Values {
	params := url.Values{
		"raid":       {q.RaidSlug},
		"boss":       {q.BossSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {string(q.Region)},
	}
	if q.Realm != "" {
		params.Set("realm", q.Realm)
	}
	return params
}

func (q *RaidQuery) validate() error {
	if q == nil {
		return ErrInvalidQuery
	}
	if q.Slug == "" {
		return ErrInvalidRaidName
	}
	if !raidDifficultyValid(q.Difficulty) {
		return ErrInvalidRaidDiff
	}
	if q.Region == "" {
		return ErrInvalidRegion
	}
	return nil
}

// params returns the params shared by all raid-scoped endpoints.
func (q *RaidQuery) params() url.Values {
	return url.Values{
		"raid":       {q.Slug},
		"difficulty": {string(q.Difficulty)},
		"region":     {string(q.Region)},
	}
}

func (q *RaidRankingsQuery) params() url.Values {
	params := q.RaidQuery.params()
	if q.Realm != "" {
		params.Set("realm", q.Realm)
	}
	if q.Limit != 0 {
		params.Set("limit", strconv.Itoa(q.Limit))
	}
	if q.Page != 0 {
		params.Set("page", strconv.Itoa(q.Page))
	}
	return params
}
