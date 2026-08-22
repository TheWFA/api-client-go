package wfa

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var (
	isoDateRe     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	isoDateTimeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}(?::?\d{2})?)?$`)
	// fullOffsetRe matches a complete "+HH:MM"/"-HH:MM" offset at the end of
	// a datetime string.
	fullOffsetRe = regexp.MustCompile(`[+-]\d{2}:\d{2}$`)
	// bareOffsetRe matches a minutes-less "+HH"/"-HH" offset at the end of a
	// datetime string — e.g. the API renders UTC timestamps as "...+00"
	// rather than "...Z" or "...+00:00", even when fractional seconds are
	// present, so this can't be anchored to the seconds field.
	bareOffsetRe = regexp.MustCompile(`[+-]\d{2}$`)
)

// Time is an ISO 8601 date or datetime, as returned by the API. It unmarshals
// from either a bare date (e.g. "2026-08-22", assumed UTC midnight) or a
// datetime with an optional fractional-seconds component and an optional
// timezone offset. A datetime with no timezone offset is assumed to be UTC.
type Time struct {
	time.Time
}

// parseAPITime parses s using the same rules as Time.UnmarshalJSON.
func parseAPITime(s string) (time.Time, error) {
	if isoDateRe.MatchString(s) {
		return time.Parse("2006-01-02", s)
	}

	if !isoDateTimeRe.MatchString(s) {
		return time.Time{}, fmt.Errorf("wfa: %q is not an ISO date or datetime string", s)
	}

	normalized := strings.Replace(s, " ", "T", 1)

	switch {
	case strings.HasSuffix(normalized, "Z"), fullOffsetRe.MatchString(normalized):
		// Already fully-specified.
	case bareOffsetRe.MatchString(normalized):
		normalized = bareOffsetRe.ReplaceAllString(normalized, "$0:00")
	default:
		// No timezone info at all; assume UTC.
		normalized += "Z"
	}

	return time.Parse(time.RFC3339Nano, normalized)
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := parseAPITime(s)
	if err != nil {
		return err
	}

	t.Time = parsed

	return nil
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339Nano))
}
