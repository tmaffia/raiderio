package raiderio

import (
	"context"
	"errors"
	"fmt"
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
	ErrInvalidSeason     = errors.New("invalid season")
	ErrInvalidRunID      = errors.New("invalid mythic plus run id")
	ErrRunNotFound       = errors.New("mythic plus run not found")
	ErrApiTimeout        = errors.New("raiderio api request timeout")
	ErrRequestCanceled   = errors.New("raiderio api request canceled")
	// ErrUnexpected wraps an unclassified API or transport failure. It carries
	// the underlying cause via %w (when one exists) so callers can errors.Is
	// against it and still inspect the wrapped error for the real reason.
	ErrUnexpected = errors.New("unexpected error")
)

// mapApiError maps an API error body to a standardized sentinel error.
// Matching is done case-insensitively against the message because the API's
// messages are human-readable prose rather than a stable machine contract.
func mapApiError(responseBody *apiErrorResponse) error {
	msg := strings.ToLower(responseBody.Message)

	switch {
	case strings.Contains(msg, "failed to find region"),
		strings.Contains(msg, `"region" must be one of`):
		return ErrInvalidRegion
	case strings.Contains(msg, "failed to find realm"):
		return ErrInvalidRealm
	case strings.Contains(msg, "failed to find raid"),
		strings.Contains(msg, "could not find requested raid"):
		return ErrInvalidRaid
	case strings.Contains(msg, "failed to find boss"):
		return ErrInvalidBoss
	case strings.Contains(msg, "could not find requested character"):
		return ErrCharacterNotFound
	case strings.Contains(msg, "could not find requested guild"):
		return ErrGuildNotFound
	case strings.Contains(msg, "requested unsupported expansion_id"):
		return ErrUnsupportedExpac
	case strings.Contains(msg, "invalid request query input"):
		return ErrInvalidQuery
	case strings.Contains(msg, "could not find data for season"):
		return ErrInvalidSeason
	case strings.Contains(msg, "could not find keystone run"):
		return ErrRunNotFound
	default:
		return ErrUnexpected
	}
}

// wrapHttpError maps a transport error to a sentinel while preserving the
// original cause so callers can still inspect it via errors.Is/As.
func wrapHttpError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %w", ErrApiTimeout, err)
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %w", ErrRequestCanceled, err)
	default:
		return fmt.Errorf("%w: %w", ErrUnexpected, err)
	}
}
