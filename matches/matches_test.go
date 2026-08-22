package matches_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/matches"
)

func newTestService(t *testing.T, handler http.HandlerFunc) *matches.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	backend, err := wfa.NewBackend("test-key", wfa.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("NewBackend: %v", err)
	}

	return matches.New(backend)
}

func TestListQueryEncoding(t *testing.T) {
	var gotQuery string

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[],"totalItems":0,"page":1,"itemsPerPage":20}`))
	})

	_, err := svc.List(context.Background(), matches.ListQuery{
		TeamID:          []wfa.Snowflake{1, 2},
		Status:          []matches.Status{matches.StatusScheduled, matches.StatusFullTime},
		OrderByDateDesc: wfa.Bool(true),
		ListParams:      wfa.ListParams{ItemsPerPage: wfa.Int(10)},
	})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	values, err := url.ParseQuery(gotQuery)
	if err != nil {
		t.Fatalf("parsing query: %v", err)
	}

	want := map[string]string{
		"teamId[0]":       "1",
		"teamId[1]":       "2",
		"status[0]":       "scheduled",
		"status[1]":       "full-time",
		"orderByDateDesc": "true",
		"itemsPerPage":    "10",
	}

	for key, wantVal := range want {
		if got := values.Get(key); got != wantVal {
			t.Errorf("query[%q] = %q, want %q", key, got, wantVal)
		}
	}
}

func TestGetDecodesPolymorphicEvents(t *testing.T) {
	body := `{
		"id": 1,
		"status": "full-time",
		"scheduledFor": "2026-08-22T12:00:00Z",
		"times": {
			"firstHalfStartedAt": null,
			"secondHalfStartedAt": null,
			"firstHalfExtraTimeStartedAt": null,
			"secondHalfExtraTimeStartedAt": null
		},
		"hidden": false,
		"homeTeam": {"id": 1, "name": "Home", "nickname": "H", "badgeUrl": "https://x/h.png"},
		"awayTeam": {"id": 2, "name": "Away", "nickname": "A", "badgeUrl": "https://x/a.png"},
		"homeScore": 2, "awayScore": 1, "homeScorePenalty": 0, "awayScorePenalty": 0,
		"competition": {"id": 10, "name": "League", "badgeUrl": null, "organisation": null},
		"season": {"id": 5, "name": "2026", "startDate": "2026-01-01", "endDate": "2026-12-31"},
		"court": null,
		"matchGroup": null,
		"officials": {},
		"streams": [],
		"homeLineups": [],
		"awayLineups": [],
		"penalties": [],
		"events": [
			{"type": "goal", "createdAt": "2026-08-22T12:10:00Z", "time": 10, "matchPeriod": "first-half", "teamId": 1, "player": {"id": 3, "name": "Scorer"}, "penalty": false, "goalType": "goal"},
			{"type": "yellow_card", "createdAt": "2026-08-22T12:20:00Z", "time": 20, "matchPeriod": "first-half", "teamId": 2, "player": {"id": 4, "name": "Booked"}},
			{"type": "substitution", "createdAt": "2026-08-22T12:30:00Z", "time": 30, "matchPeriod": "second-half", "teamId": 1, "playerOn": {"id": 5, "name": "On"}, "playerOff": {"id": 3, "name": "Off"}}
		]
	}`

	svc := newTestService(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	})

	match, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(match.Events) != 3 {
		t.Fatalf("len(Events) = %d, want 3", len(match.Events))
	}

	goal, ok := match.Events[0].(matches.GoalEvent)
	if !ok {
		t.Fatalf("Events[0] is %T, want GoalEvent", match.Events[0])
	}
	if goal.EventType() != matches.EventTypeGoal || goal.Player.Name != "Scorer" {
		t.Errorf("unexpected goal event: %+v", goal)
	}

	card, ok := match.Events[1].(matches.CardEvent)
	if !ok {
		t.Fatalf("Events[1] is %T, want CardEvent", match.Events[1])
	}
	if card.EventType() != matches.EventTypeYellowCard || card.Player.Name != "Booked" {
		t.Errorf("unexpected card event: %+v", card)
	}

	sub, ok := match.Events[2].(matches.SubstitutionEvent)
	if !ok {
		t.Fatalf("Events[2] is %T, want SubstitutionEvent", match.Events[2])
	}
	if sub.PlayerOn.Name != "On" || sub.PlayerOff.Name != "Off" {
		t.Errorf("unexpected substitution event: %+v", sub)
	}
}
