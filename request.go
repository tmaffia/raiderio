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

// getAPIResponse is a helper function that makes a GET request to the Raider.IO API
// It returns an error if the API returns a non-200 status code, or if the
// response body cannot be read
// Returns the error message from the api back to the client method that calls it,
// so in cases where the realm or the character name cannot be found, developer is presented
// with that error state.
func (c *Client) getAPIResponse(ctx context.Context, reqUrl string) ([]byte, error) {
	if c.AccessKey != "" {
		u, err := url.Parse(reqUrl)
		if err != nil {
			return nil, errors.New("error parsing request URL")
		}
		q := u.Query()
		q.Set("access_key", c.AccessKey)
		u.RawQuery = q.Encode()
		reqUrl = u.String()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, errors.New("error creating HTTP request")
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, wrapHttpError(err)
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body")
	}

	// If not 200, api is returning an error state
	if resp.StatusCode != 200 {
		var responseBody apiErrorResponse
		err = json.Unmarshal(body, &responseBody)
		// unmarshal error implies response is in an incorrect format
		// instead of api message, return http status
		if err != nil {
			return nil, wrapApiError(&responseBody)
		}

		// return error with message directly from the api
		return nil, wrapApiError(&responseBody)
	}

	return body, nil
}
