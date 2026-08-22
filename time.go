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
	shortTZOffset = regexp.MustCompile(`(:\d{2})([+-]\d{2})$`)
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
	normalized = shortTZOffset.ReplaceAllString(normalized, "$1$2:00")

	hasOffset := strings.HasSuffix(normalized, "Z") ||
		regexp.MustCompile(`[+-]\d{2}:\d{2}$`).MatchString(normalized)
	if !hasOffset {
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
	return json.Marshal(t.Time.Format(time.RFC3339Nano))
}
