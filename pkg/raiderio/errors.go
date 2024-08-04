package raiderio

import (
	"errors"
	"strings"
)

var (
	ErrInvalidRegion     = errors.New("invalid region")
	ErrInvalidRealm      = errors.New("invalid realm")
	ErrInvalidCharName   = errors.New("invalid character name")
	ErrInvalidGuildName  = errors.New("invalid guild name")
	ErrInvalidRaidName   = errors.New("invalid raid name")
	ErrInvalidRaidDiff   = errors.New("invalid raid difficulty")
	ErrInvalidRaid       = errors.New("invalid raid")
	ErrFieldMissing      = errors.New("field missing from api response")
	ErrCharacterNotFound = errors.New("character not found")
	ErrGuildNotFound     = errors.New("guild not found")
	ErrUnsupportedExpac  = errors.New("unsupported expansion")
	ErrLimitOutOfBounds  = errors.New("limit must be a positive int")
	ErrPageOutOfBounds   = errors.New("page must be a positive int")
	ErrUnexpected        = errors.New("unexpected error")
)

// Turns api errors into standardized go errors with
// consistent error messages
func wrapAPIError(responseBody *apiErrorResponse) error {
	if strings.Contains(responseBody.Message, "Failed to find region") {
		return ErrInvalidRegion
	}

	if strings.Contains(responseBody.Message, "Failed to find realm") {
		return ErrInvalidRealm
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

	return ErrUnexpected
}