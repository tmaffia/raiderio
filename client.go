package raiderio

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Base URL for the Raider.IO API
const baseUrl string = "https://raider.io/api"

// Client is the main struct for interacting with the Raider.IO API
type Client struct {
	ApiUrl     string
	AccessKey  string
	HttpClient *http.Client
}

// NewClient creates a Client. Pass an access key for authenticated requests.
func NewClient(accessKey ...string) *Client {
	c := &Client{
		ApiUrl:     baseUrl + "/v1",
		HttpClient: &http.Client{},
	}
	if len(accessKey) > 0 {
		c.AccessKey = accessKey[0]
	}
	return c
}

// GetCharacter retrieves a character profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Character struct
func (c *Client) GetCharacter(ctx context.Context, cq *CharacterQuery) (*Character, error) {
	fields, err := validateCharacterQuery(cq)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"region": {string(cq.Region)},
		"realm":  {cq.Realm},
		"name":   {cq.Name},
	}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}

	return getJSON[Character](c, ctx, "/characters/profile", params)
}

// GetGuild retrieves a guild profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Guild struct
func (c *Client) GetGuild(ctx context.Context, gq *GuildQuery) (*Guild, error) {
	fields, err := createGuildQuery(gq)
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"region": {string(gq.Region)},
		"realm":  {gq.Realm},
		"name":   {gq.Name},
	}
	if len(fields) > 0 {
		params.Set("fields", strings.Join(fields, ","))
	}

	return getJSON[Guild](c, ctx, "/guilds/profile", params)
}

// GetRaids retrieves a list of raids from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Raids struct
// Takes an Expansion enum as a parameter, in addition to context.Context
func (c *Client) GetRaids(ctx context.Context, e Expansion) (*Raids, error) {
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
		"region":     {string(rq.Region)},
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
		"region":     {string(q.Region)},
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
		"region":     {string(q.Region)},
	}
	if q.Realm != "" {
		params.Set("realm", q.Realm)
	}

	return getJSON[BossRankings](c, ctx, "/raiding/boss-rankings", params)
}

// GetHallOfFame retrieves the hall of fame for a given raid
func (c *Client) GetHallOfFame(ctx context.Context, q *RaidQuery) (*HallOfFame, error) {
	if err := validateRaidQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.Slug},
		"difficulty": {string(q.Difficulty)},
		"region":     {string(q.Region)},
	}

	return getJSON[HallOfFame](c, ctx, "/raiding/hall-of-fame", params)
}

// GetRaidProgression retrieves the raid progression for a given raid
func (c *Client) GetRaidProgression(ctx context.Context, q *RaidQuery) (*RaidProgressionResponse, error) {
	if err := validateRaidQuery(q); err != nil {
		return nil, err
	}

	params := url.Values{
		"raid":       {q.Slug},
		"difficulty": {string(q.Difficulty)},
		"region":     {string(q.Region)},
	}

	return getJSON[RaidProgressionResponse](c, ctx, "/raiding/progression", params)
}
