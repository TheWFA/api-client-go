package wfa

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "date only",
			input: `"2026-08-22"`,
			want:  time.Date(2026, 8, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "datetime with Z",
			input: `"2026-08-22T12:00:00Z"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "datetime with milliseconds and Z",
			input: `"2026-08-22T12:00:00.123Z"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 123000000, time.UTC),
		},
		{
			name:  "datetime with space separator",
			input: `"2026-08-22 12:00:00Z"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "datetime with no timezone assumes UTC",
			input: `"2026-08-22T12:00:00"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC),
		},
		{
			name:  "datetime with short offset",
			input: `"2026-08-22T12:00:00+01"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 0, time.FixedZone("", 3600)),
		},
		{
			// The live API renders UTC timestamps this way: a bare, minutes-less
			// offset immediately after fractional seconds, with no "Z".
			name:  "datetime with fractional seconds and short offset",
			input: `"2026-08-22T12:13:08.827805+00"`,
			want:  time.Date(2026, 8, 22, 12, 13, 8, 827805000, time.UTC),
		},
		{
			name:  "datetime with fractional seconds and non-UTC short offset",
			input: `"2026-08-22T12:13:08.827805+05"`,
			want:  time.Date(2026, 8, 22, 12, 13, 8, 827805000, time.FixedZone("", 5*3600)),
		},
		{
			name:  "datetime with full offset",
			input: `"2026-08-22T12:00:00+01:00"`,
			want:  time.Date(2026, 8, 22, 12, 0, 0, 0, time.FixedZone("", 3600)),
		},
		{
			name:    "garbage",
			input:   `"not a date"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Time

			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			if !got.Equal(tt.want) {
				t.Errorf("got %v, want %v", got.Time, tt.want)
			}
		})
	}
}

func TestTimeUnmarshalJSONNull(t *testing.T) {
	var got *Time

	if err := json.Unmarshal([]byte(`null`), &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestTimeRoundTrip(t *testing.T) {
	type wrapper struct {
		At Time `json:"at"`
	}

	original := wrapper{At: Time{Time: time.Date(2026, 8, 22, 12, 30, 0, 0, time.UTC)}}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var roundTripped wrapper
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !roundTripped.At.Equal(original.At.Time) {
		t.Errorf("got %v, want %v", roundTripped.At.Time, original.At.Time)
	}
}
