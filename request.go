package raiderio

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type apiErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Err        string `json:"error"`
	Message    string `json:"message"`
}

func (c *Client) getAPIResponse(ctx context.Context, path string, params url.Values) ([]byte, error) {
	if c.AccessKey != "" {
		params.Set("access_key", c.AccessKey)
	}

	reqUrl := c.ApiUrl + path
	if q := params.Encode(); q != "" {
		reqUrl += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, errors.New("error creating HTTP request")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, wrapHttpError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	if resp.StatusCode != 200 {
		var responseBody apiErrorResponse
		_ = json.Unmarshal(body, &responseBody)
		return nil, wrapApiError(&responseBody)
	}

	return body, nil
}

func getJSON[T any](c *Client, ctx context.Context, path string, params url.Values) (*T, error) {
	body, err := c.getAPIResponse(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var v T
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, errors.New("error unmarshalling response")
	}
	return &v, nil
}
