package raiderio

import (
	"context"
	"errors"
	"strings"
)

// Errors that the api produces
var (
	ErrInvalidRegion     = errors.New("invalid region")
	ErrInvalidRealm      = errors.New("invalid realm")
	ErrInvalidCharName   = errors.New("invalid character name")
	ErrInvalidGuildName  = errors.New("invalid guild name")
	ErrInvalidRaidName   = errors.New("invalid raid name")
	ErrInvalidRaidDiff   = errors.New("invalid raid difficulty")
	ErrInvalidRaid       = errors.New("invalid raid")
	ErrCharacterNotFound = errors.New("character not found")
	ErrGuildNotFound     = errors.New("guild not found")
	ErrUnsupportedExpac  = errors.New("unsupported expansion")
	ErrLimitOutOfBounds  = errors.New("limit must be a positive int")
	ErrPageOutOfBounds   = errors.New("page must be a positive int")
	ErrInvalidBoss       = errors.New("invalid boss")
	ErrInvalidQuery      = errors.New("invalid query")
	ErrApiTimeout        = errors.New("raiderio api request timeout")
	ErrUnexpected        = errors.New("unexpected error")
)

// Turns api errors into standardized go errors with
// consistent error messages
func wrapApiError(responseBody *apiErrorResponse) error {
	if strings.Contains(responseBody.Message, "Failed to find region") {
		return ErrInvalidRegion
	}

	if strings.Contains(responseBody.Message, "Failed to find realm") {
		return ErrInvalidRealm
	}

	if strings.Contains(responseBody.Message, "Failed to find raid") {
		return ErrInvalidRaid
	}

	if strings.Contains(responseBody.Message, "Failed to find boss") {
		return ErrInvalidBoss
	}

	if strings.Contains(responseBody.Message, "Could not find requested character") {
		return ErrCharacterNotFound
	}

	if strings.Contains(responseBody.Message, "Could not find requested guild") {
		return ErrGuildNotFound
	}

	if strings.Contains(responseBody.Message, "Requested unsupported expansion_id") {
		return ErrUnsupportedExpac
	}

	if strings.Contains(responseBody.Message, "Could not find requested raid") {
		return ErrInvalidRaid
	}

	if strings.Contains(responseBody.Message, "Invalid request query input") {
		return ErrInvalidQuery
	}

	return ErrUnexpected
}

func wrapHttpError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return ErrApiTimeout
	}
	return ErrUnexpected
}
