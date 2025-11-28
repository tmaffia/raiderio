package raiderio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// NewClient creates a new Client struct
func NewClient() *Client {
	var c Client
	c.ApiUrl = baseUrl + "/v1"
	c.HttpClient = &http.Client{}
	return &c
}

// GetCharacter retrieves a character profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the CharacterProfile struct
func (c *Client) GetCharacter(ctx context.Context, cq *CharacterQuery) (*Character, error) {
	err := validateCharacterQuery(cq)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("region", cq.Region.Slug)
	params.Add("realm", cq.Realm)
	params.Add("name", cq.Name)
	if len(cq.fields) > 0 {
		params.Add("fields", strings.Join(cq.fields, ","))
	}

	reqUrl := fmt.Sprintf("%s/characters/profile?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var profile Character
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return nil, errors.New("error unmarshalling character profile")
	}

	return &profile, nil
}

// GetGuild retrieves a guild profile from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the GuildProfile struct
func (c *Client) GetGuild(ctx context.Context, gq *GuildQuery) (*Guild, error) {
	err := createGuildQuery(gq)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("region", gq.Region.Slug)
	params.Add("realm", gq.Realm)
	params.Add("name", gq.Name)
	if len(gq.fields) > 0 {
		params.Add("fields", strings.Join(gq.fields, ","))
	}

	reqUrl := fmt.Sprintf("%s/guilds/profile?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	profile, err := unmarshalGuild(body)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// GetRaids retrieves a list of raids from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the Raids struct
// Takes an Expansion enum as a parameter, in addition to context.Context
func (c *Client) GetRaids(ctx context.Context, e expansions.Expansion) (*Raids, error) {
	params := url.Values{}
	params.Add("expansion_id", fmt.Sprintf("%d", e))
	reqUrl := fmt.Sprintf("%s/raiding/static-data?%s", c.ApiUrl, params.Encode())
	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var raids Raids
	err = json.Unmarshal(body, &raids)
	if err != nil {
		return nil, errors.New("error unmarshalling raids")
	}

	return &raids, nil
}

// GetRaidRankings retrieves a list of raid rankings from the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read or mapped to the RaidRankings struct
// Takes a RaidQuery struct as a parameter, in addition to context.Context
func (c *Client) GetRaidRankings(ctx context.Context, rq *RaidQuery) (*RaidRankings, error) {
	err := validateRaidRankingsQuery(rq)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("raid", rq.Slug)
	params.Add("difficulty", string(rq.Difficulty))
	params.Add("region", rq.Region.Slug)

	if rq.Realm != "" {
		params.Add("realm", rq.Realm)
	}

	if rq.Limit != 0 {
		params.Add("limit", fmt.Sprintf("%d", rq.Limit))
	}

	if rq.Page != 0 {
		params.Add("page", fmt.Sprintf("%d", rq.Page))
	}

	reqUrl := fmt.Sprintf("%s/raiding/raid-rankings?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var rankings RaidRankings
	err = json.Unmarshal(body, &rankings)
	if err != nil {
		return nil, errors.New("error unmarshalling raid rankings")
	}

	return &rankings, nil
}

// GetGuildBossKill returns a guild's first kill of a given boss
// Takes a context.Context object to facilitate timeout, and a GuildBossKillQuery
// GuildBossKillQuery has only required fields for this request
// returns a BossKill object
func (c *Client) GetGuildBossKill(ctx context.Context, q *GuildBossKillQuery) (*BossKill, error) {
	err := validateGuildBossKillQuery(q)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("raid", q.RaidSlug)
	params.Add("difficulty", string(q.Difficulty))
	params.Add("region", q.Region.Slug)
	params.Add("realm", q.Realm)
	params.Add("guild", q.GuildName)
	params.Add("boss", q.BossSlug)

	reqUrl := fmt.Sprintf("%s/guilds/boss-kill?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	k, err := unmarshalGuildBossKill(body)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// GetBossRankings retrieves the boss rankings for a given raid and boss
func (c *Client) GetBossRankings(ctx context.Context, q *BossRankingsQuery) (*BossRankings, error) {
	err := validateBossRankingsQuery(q)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("raid", q.RaidSlug)
	params.Add("boss", q.BossSlug)
	params.Add("difficulty", string(q.Difficulty))
	params.Add("region", q.Region.Slug)
	if q.Realm != "" {
		params.Add("realm", q.Realm)
	}

	reqUrl := fmt.Sprintf("%s/raiding/boss-rankings?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var rankings BossRankings
	err = json.Unmarshal(body, &rankings)
	if err != nil {
		return nil, errors.New("error unmarshalling boss rankings")
	}

	return &rankings, nil
}

// GetHallOfFame retrieves the hall of fame for a given raid
func (c *Client) GetHallOfFame(ctx context.Context, q *HallOfFameQuery) (*HallOfFame, error) {
	err := validateHallOfFameQuery(q)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("raid", q.RaidSlug)
	params.Add("difficulty", string(q.Difficulty))
	params.Add("region", q.Region.Slug)

	reqUrl := fmt.Sprintf("%s/raiding/hall-of-fame?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var hof HallOfFame
	err = json.Unmarshal(body, &hof)
	if err != nil {
		return nil, errors.New("error unmarshalling hall of fame")
	}

	return &hof, nil
}

// GetRaidProgression retrieves the raid progression for a given raid
func (c *Client) GetRaidProgression(ctx context.Context, q *RaidProgressionQuery) (*RaidProgressionResponse, error) {
	err := validateRaidProgressionQuery(q)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("raid", q.RaidSlug)
	params.Add("difficulty", string(q.Difficulty))
	params.Add("region", q.Region.Slug)

	reqUrl := fmt.Sprintf("%s/raiding/progression?%s", c.ApiUrl, params.Encode())

	body, err := c.getAPIResponse(ctx, reqUrl)
	if err != nil {
		return nil, err
	}

	var prog RaidProgressionResponse
	err = json.Unmarshal(body, &prog)
	if err != nil {
		return nil, errors.New("error unmarshalling raid progression")
	}

	return &prog, nil
}
