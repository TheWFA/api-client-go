package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/client"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *client.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	c, err := client.New("test-key", wfa.WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}

	return c
}

func TestNewRequiresAPIKey(t *testing.T) {
	if _, err := client.New(""); err == nil {
		t.Fatal("expected an error for an empty API key")
	}
}

func TestCompetitionsGetOptionalSeasonID(t *testing.T) {
	var gotQuery string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1,"name":"C","type":"league","badgeUrl":null,"sortOrder":0,"hidden":false,"organisation":null,"season":null,"matchGroups":[]}`))
	})

	if _, err := c.Competitions.Get(context.Background(), 1, nil); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotQuery != "" {
		t.Errorf("expected no query params when seasonID is nil, got %q", gotQuery)
	}

	if _, err := c.Competitions.Get(context.Background(), 1, wfa.Int(2026)); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if gotQuery != "seasonId=2026" {
		t.Errorf("gotQuery = %q, want %q", gotQuery, "seasonId=2026")
	}
}

func TestHistoryPaths(t *testing.T) {
	var gotPath string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"items":[],"totalItems":0}`))
	})

	if _, err := c.History.List(context.Background(), wfa.HistoryEntityTeam, 42); err != nil {
		t.Fatalf("List: %v", err)
	}

	if want := "/history/team/42"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
}

func TestAccreditationsGetStringID(t *testing.T) {
	var gotPath string

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"abc-123","name":"Media","description":"","category":"media","issuingBody":"WFA","validityPeriod":null,"createdAt":"2026-01-01","holderCount":3}`))
	})

	acc, err := c.Accreditations.Get(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if want := "/accreditations/abc-123"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}

	if acc.HolderCount != 3 {
		t.Errorf("HolderCount = %d, want 3", acc.HolderCount)
	}
}
