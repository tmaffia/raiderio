package raiderio

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type apiErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Err        string `json:"error"`
	Message    string `json:"message"`
}

// query is implemented by every endpoint query type. It centralizes validation
// and URL-parameter construction so client methods stay thin and consistent.
type query interface {
	validate() error
	params() url.Values
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
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, wrapHttpError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, wrapApiError(resp.StatusCode, body)
	}

	return body, nil
}

// wrapApiError turns a non-200 response into a sentinel error wrapped with the
// API's message (or the raw body when no message is present). Unknown errors
// also carry the HTTP status code.
func wrapApiError(statusCode int, body []byte) error {
	var responseBody apiErrorResponse
	_ = json.Unmarshal(body, &responseBody)

	sentinel := mapApiError(&responseBody)
	if errors.Is(sentinel, ErrUnexpected) {
		return fmt.Errorf("%w: %d %s", sentinel, statusCode, apiMessage(&responseBody, body))
	}
	return fmt.Errorf("%w: %s", sentinel, apiMessage(&responseBody, body))
}

// apiMessage returns the API's human-readable message, falling back to the raw
// body when the message field is absent or empty.
func apiMessage(responseBody *apiErrorResponse, body []byte) string {
	if responseBody.Message != "" {
		return responseBody.Message
	}
	return string(body)
}

func getJSON[T any](c *Client, ctx context.Context, path string, params url.Values) (*T, error) {
	body, err := c.getAPIResponse(ctx, path, params)
	if err != nil {
		return nil, err
	}
	var v T
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return &v, nil
}

// getJSONFor validates a query and fetches its JSON response in one step.
func getJSONFor[T any](c *Client, ctx context.Context, path string, q query) (*T, error) {
	if err := q.validate(); err != nil {
		return nil, err
	}
	return getJSON[T](c, ctx, path, q.params())
}
