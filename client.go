package raiderio

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tmaffia/raiderio/expansions"
)

// Base URL for the Raider.IO API
const baseUrl string = "https://raider.io/api"

// Client is the main struct for interacting with the Raider.IO API
type Client struct {
	ApiUrl     string
	AccessKey  string
	HttpClient *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithAccessKey returns a ClientOption that sets the AccessKey
func WithAccessKey(key string) ClientOption {
	return func(c *Client) {
		c.AccessKey = key
	}
}

// NewClient creates a new Client struct
func NewClient(opts ...ClientOption) *Client {
	var c Client
	c.ApiUrl = baseUrl + "/v1"
	c.HttpClient = &http.Client{}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}

// GetCharacter retrieves a character profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the CharacterProfile struct
func (c *Client) GetCharacter(ctx context.Context, cq *CharacterQuery) (*Character, error) {
	if err := validateCharacterQuery(cq); err != nil {
		return nil, err
	}

	params := url.Values{
		"region": {cq.Region.Slug},
		"realm":  {cq.Realm},
		"name":   {cq.Name},
	}
	if len(cq.fields) > 0 {
		params.Set("fields", strings.Join(cq.fields, ","))
	}

	return getJSON[Character](c, ctx, "/characters/profile", params)
}

// GetGuild retrieves a guild profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the GuildProfile struct
func (c *Client) GetGuild(ctx context.Context, gq *GuildQuery) (*Guild, error) {
	if err := createGuildQuery(gq); err != nil {
		return nil, err
	}

	params := url.Values{
		"region": {gq.Region.Slug},
		"realm":  {gq.Realm},
		"name":   {gq.Name},
	}
	if len(gq.fields) > 0 {
		params.Set("fields", strings.Join(gq.fields, ","))
	}

	profile, err := getJSON[Guild](c, ctx, "/guilds/profile", params)
	if err != nil {
		return nil, err
	}
	for k, entry := range profile.RaidRankings {
		entry.RaidSlug = k
		profile.RaidRankings[k] = entry
	}
	return profile, nil
}

// GetRaids retrieves a list of raids from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Raids struct
// Takes an Expansion enum as a parameter, in addition to context.Context
func (c *Client) GetRaids(ctx context.Context, e expansions.Expansion) (*Raids, error) {
	params := url.Values{"expansion_id": {strconv.Itoa(int(e))}}
	return getJSON[Raids](c, ctx, "/raiding/static-data", params)
}

// GetRaidRankings retrieves a list of raid rankings from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the RaidRankings struct
// Takes a RaidQuery struct as a parameter, in addition to context.Context
func (c *Client) GetRaidRankings(ctx context.Context, rq *RaidQuery) (*RaidRankings, error) {
	if err := validateRaidRankingsQuery(rq); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {rq.Slug},
		"difficulty": {string(rq.Difficulty)},
		"region":     {rq.Region.Slug},
	}
	if rq.Realm != "" {
		params.Set("realm", rq.Realm)
	}
	if rq.Limit != 0 {
		params.Set("limit", strconv.Itoa(rq.Limit))
	}
	if rq.Page != 0 {
		params.Set("page", strconv.Itoa(rq.Page))
	}

	return getJSON[RaidRankings](c, ctx, "/raiding/raid-rankings", params)
}

// GetGuildBossKill returns a guild's first kill of a given boss
// Takes a context.Context object to facilitate timeout, and a GuildBossKillQuery
// GuildBossKillQuery has only required fields for this request
// returns a BossKill object
func (c *Client) GetGuildBossKill(ctx context.Context, q *GuildBossKillQuery) (*BossKill, error) {
	if err := validateGuildBossKillQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.RaidSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {q.Region.Slug},
		"realm":      {q.Realm},
		"guild":      {q.GuildName},
		"boss":       {q.BossSlug},
	}

	body, err := c.getAPIResponse(ctx, "/guilds/boss-kill", params)
	if err != nil {
		return nil, err
	}
	return unmarshalGuildBossKill(body)
}

// GetBossRankings retrieves the boss rankings for a given raid and boss
func (c *Client) GetBossRankings(ctx context.Context, q *BossRankingsQuery) (*BossRankings, error) {
	if err := validateBossRankingsQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.RaidSlug},
		"boss":       {q.BossSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {q.Region.Slug},
	}
	if q.Realm != "" {
		params.Set("realm", q.Realm)
	}

	return getJSON[BossRankings](c, ctx, "/raiding/boss-rankings", params)
}

// GetHallOfFame retrieves the hall of fame for a given raid
func (c *Client) GetHallOfFame(ctx context.Context, q *HallOfFameQuery) (*HallOfFame, error) {
	if err := validateHallOfFameQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.RaidSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {q.Region.Slug},
	}

	return getJSON[HallOfFame](c, ctx, "/raiding/hall-of-fame", params)
}

// GetRaidProgression retrieves the raid progression for a given raid
func (c *Client) GetRaidProgression(ctx context.Context, q *RaidProgressionQuery) (*RaidProgressionResponse, error) {
	if err := validateRaidProgressionQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.RaidSlug},
		"difficulty": {string(q.Difficulty)},
		"region":     {q.Region.Slug},
	}

	return getJSON[RaidProgressionResponse](c, ctx, "/raiding/progression", params)
}
