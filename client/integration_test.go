//go:build integration

// Integration tests that run against the real API.
//
// These are gated behind the "integration" build tag AND skipped at runtime
// unless WFA_API_KEY is set, so a plain `go test ./...` never makes network
// calls. Run them with:
//
//	WFA_API_KEY=... go test -tags=integration ./client/...
//
// Set WFA_API_BASE_URL to point at a non-default environment.
package client_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/client"
	"github.com/TheWFA/api-client-go/matches"
	"github.com/TheWFA/api-client-go/persons"
	"github.com/TheWFA/api-client-go/search"
)

func requireAPIKey(t *testing.T) string {
	t.Helper()

	apiKey := os.Getenv("WFA_API_KEY")
	if apiKey == "" {
		t.Skip("skipping integration test: WFA_API_KEY not set")
	}

	return apiKey
}

func newIntegrationClient(t *testing.T, apiKey string) *client.Client {
	t.Helper()

	var opts []wfa.Option
	if baseURL := os.Getenv("WFA_API_BASE_URL"); baseURL != "" {
		opts = append(opts, wfa.WithBaseURL(baseURL))
	}

	c, err := client.New(apiKey, opts...)
	if err != nil {
		t.Fatalf("client.New: %v", err)
	}

	return c
}

func integrationContext(t *testing.T) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)

	return ctx
}

func TestIntegrationHealth(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	health, err := c.Health.Get(ctx)
	if err != nil {
		t.Fatalf("Health.Get: %v", err)
	}

	if health.Status == "" {
		t.Error("expected a non-empty Status")
	}

	if health.Scope == "" {
		t.Error("expected a non-empty Scope")
	}
}

func TestIntegrationMatches(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Matches.List(ctx, matches.ListQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}})
	if err != nil {
		t.Fatalf("Matches.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no matches available to test against")
	}

	match, err := c.Matches.Get(ctx, page.Items[0].ID)
	if err != nil {
		t.Fatalf("Matches.Get: %v", err)
	}

	if match.ID != page.Items[0].ID {
		t.Errorf("ID = %d, want %d", match.ID, page.Items[0].ID)
	}

	if match.ScheduledFor != nil && match.ScheduledFor.Year() < 2000 {
		t.Errorf("ScheduledFor parsed implausibly: %v", match.ScheduledFor)
	}
}

func TestIntegrationTeams(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Teams.List(ctx, teamsListQuery(5))
	if err != nil {
		t.Fatalf("Teams.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no teams available to test against")
	}

	teamID := page.Items[0].ID

	team, err := c.Teams.Get(ctx, teamID)
	if err != nil {
		t.Fatalf("Teams.Get: %v", err)
	}

	if team.ID != teamID {
		t.Errorf("ID = %d, want %d", team.ID, teamID)
	}

	if _, err := c.Teams.Players(ctx, teamID, playersQuery()); err != nil {
		t.Errorf("Teams.Players: %v", err)
	}

	if _, err := c.Teams.Staff(ctx, teamID, staffQuery()); err != nil {
		t.Errorf("Teams.Staff: %v", err)
	}

	if _, err := c.Teams.Registrations(ctx, teamID, registrationsQuery()); err != nil {
		t.Errorf("Teams.Registrations: %v", err)
	}

	if _, err := c.Teams.Seasons(ctx, teamID); err != nil {
		t.Errorf("Teams.Seasons: %v", err)
	}

	if _, err := c.Teams.Stats.Summary(ctx, teamID, wfa.StatsFilterQuery{}); err != nil {
		t.Errorf("Teams.Stats.Summary: %v", err)
	}

	if _, err := c.Teams.Stats.Players(ctx, teamID, wfa.PlayerStatsQuery{}); err != nil {
		t.Errorf("Teams.Stats.Players: %v", err)
	}
}

func TestIntegrationClubs(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Clubs.List(ctx, clubsListQuery(5))
	if err != nil {
		t.Fatalf("Clubs.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no clubs available to test against")
	}

	club, err := c.Clubs.Get(ctx, page.Items[0].ID)
	if err != nil {
		t.Fatalf("Clubs.Get: %v", err)
	}

	if club.ID != page.Items[0].ID {
		t.Errorf("ID = %d, want %d", club.ID, page.Items[0].ID)
	}
}

func TestIntegrationCompetitions(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Competitions.List(ctx, competitionsListQuery(5))
	if err != nil {
		t.Fatalf("Competitions.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no competitions available to test against")
	}

	id := page.Items[0].ID

	competition, err := c.Competitions.Get(ctx, id, nil)
	if err != nil {
		t.Fatalf("Competitions.Get: %v", err)
	}

	if competition.ID != id {
		t.Errorf("ID = %d, want %d", competition.ID, id)
	}

	if _, err := c.Competitions.Teams(ctx, id, nil); err != nil {
		t.Errorf("Competitions.Teams: %v", err)
	}

	if _, err := c.Competitions.Seasons(ctx, id); err != nil {
		t.Errorf("Competitions.Seasons: %v", err)
	}

	if _, err := c.Competitions.Stats.Summary(ctx, id, wfa.StatsFilterQuery{}); err != nil {
		t.Errorf("Competitions.Stats.Summary: %v", err)
	}
}

func TestIntegrationOrganisations(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Organisations.List(ctx, organisationsListQuery(5))
	if err != nil {
		t.Fatalf("Organisations.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no organisations available to test against")
	}

	org, err := c.Organisations.Get(ctx, page.Items[0].ID)
	if err != nil {
		t.Fatalf("Organisations.Get: %v", err)
	}

	if org.ID != page.Items[0].ID {
		t.Errorf("ID = %d, want %d", org.ID, page.Items[0].ID)
	}
}

func TestIntegrationSeasons(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Seasons.List(ctx, seasonsListQuery(5))
	if err != nil {
		t.Fatalf("Seasons.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no seasons available to test against")
	}

	season, err := c.Seasons.Get(ctx, page.Items[0].ID)
	if err != nil {
		t.Fatalf("Seasons.Get: %v", err)
	}

	if season.ID != page.Items[0].ID {
		t.Errorf("ID = %d, want %d", season.ID, page.Items[0].ID)
	}
}

func TestIntegrationAccreditations(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	if _, err := c.Accreditations.List(ctx, accreditationsListQuery(5)); err != nil {
		t.Fatalf("Accreditations.List: %v", err)
	}

	facets, err := c.Accreditations.Facets(ctx)
	if err != nil {
		t.Fatalf("Accreditations.Facets: %v", err)
	}

	if facets.Categories == nil {
		t.Error("expected a non-nil Categories slice")
	}
}

func TestIntegrationPersons(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Persons.List(ctx, personsListQuery(5))
	if err != nil {
		t.Fatalf("Persons.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no persons available to test against")
	}

	personID := page.Items[0].ID

	person, err := c.Persons.Get(ctx, personID)
	if err != nil {
		t.Fatalf("Persons.Get: %v", err)
	}

	if person.ID != personID {
		t.Errorf("ID = %d, want %d", person.ID, personID)
	}

	if _, err := c.Persons.Registrations(ctx, personID, persons.RegistrationsQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}); err != nil {
		t.Errorf("Persons.Registrations: %v", err)
	}

	if _, err := c.Persons.Appearances(ctx, personID, persons.AppearancesQuery{ListParams: wfa.ListParams{ItemsPerPage: wfa.Int(5)}}); err != nil {
		t.Errorf("Persons.Appearances: %v", err)
	}

	if _, err := c.Persons.Suspensions(ctx, personID, personSuspensionsQuery()); err != nil {
		t.Errorf("Persons.Suspensions: %v", err)
	}

	if _, err := c.Persons.Stats.Summary(ctx, personID, wfa.StatsFilterQuery{}); err != nil {
		t.Errorf("Persons.Stats.Summary: %v", err)
	}

	if _, err := c.Persons.Stats.Goals(ctx, personID, persons.StatsQuery{}); err != nil {
		t.Errorf("Persons.Stats.Goals: %v", err)
	}

	if _, err := c.Persons.Stats.Assists(ctx, personID, persons.StatsQuery{}); err != nil {
		t.Errorf("Persons.Stats.Assists: %v", err)
	}

	if _, err := c.Persons.Stats.Cards(ctx, personID, persons.CardsQuery{}); err != nil {
		t.Errorf("Persons.Stats.Cards: %v", err)
	}
}

func TestIntegrationSearch(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Search.List(ctx, search.Query{Query: "a", ItemsPerPage: wfa.Int(5)})
	if err != nil {
		t.Fatalf("Search.List: %v", err)
	}

	if page.Items == nil {
		t.Error("expected a non-nil Items slice")
	}
}

func TestIntegrationLocations(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	page, err := c.Locations.List(ctx, locationsListQuery(5))
	if err != nil {
		t.Fatalf("Locations.List: %v", err)
	}

	if len(page.Items) == 0 {
		t.Skip("no locations available to test against")
	}

	location, err := c.Locations.Get(ctx, page.Items[0].ID)
	if err != nil {
		t.Fatalf("Locations.Get: %v", err)
	}

	if location.ID != page.Items[0].ID {
		t.Errorf("ID = %d, want %d", location.ID, page.Items[0].ID)
	}
}

func TestIntegrationSuspensions(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	if _, err := c.Suspensions.List(ctx, suspensionsListQuery(5)); err != nil {
		t.Fatalf("Suspensions.List: %v", err)
	}
}

func TestIntegrationTies(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	if _, err := c.Ties.List(ctx, tiesListQuery(5)); err != nil {
		t.Fatalf("Ties.List: %v", err)
	}
}

func TestIntegrationKits(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	types, err := c.Kits.Types(ctx)
	if err != nil {
		t.Fatalf("Kits.Types: %v", err)
	}

	if types.Items == nil {
		t.Error("expected a non-nil Items slice")
	}
}

func TestIntegrationHistory(t *testing.T) {
	c := newIntegrationClient(t, requireAPIKey(t))
	ctx := integrationContext(t)

	teamsPage, err := c.Teams.List(ctx, teamsListQuery(1))
	if err != nil {
		t.Fatalf("Teams.List: %v", err)
	}

	if len(teamsPage.Items) == 0 {
		t.Skip("no teams available to test history against")
	}

	if _, err := c.History.List(ctx, wfa.HistoryEntityTeam, teamsPage.Items[0].ID); err != nil {
		t.Fatalf("History.List: %v", err)
	}
}

func TestIntegrationErrorHandling(t *testing.T) {
	apiKey := requireAPIKey(t)
	c := newIntegrationClient(t, apiKey)
	ctx := integrationContext(t)

	t.Run("not found", func(t *testing.T) {
		_, err := c.Matches.Get(ctx, 100000000)
		if err == nil {
			t.Fatal("expected an error for a non-existent match")
		}

		if !wfa.IsNotFound(err) {
			t.Errorf("expected IsNotFound, got %v", err)
		}

		var apiErr *wfa.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode != 404 {
			t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
		}
	})

	t.Run("invalid api key", func(t *testing.T) {
		invalidClient := newIntegrationClient(t, "invalid-api-key")

		// Rejected by the API gateway's usage plan before the request reaches
		// the API, so this surfaces as a 403 rather than an application-level
		// 401.
		_, err := invalidClient.Matches.List(ctx, matches.ListQuery{})
		if err == nil {
			t.Fatal("expected an error for an invalid API key")
		}

		if !wfa.IsForbidden(err) {
			t.Errorf("expected IsForbidden, got %v", err)
		}

		var apiErr *wfa.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode != 403 {
			t.Errorf("StatusCode = %d, want 403", apiErr.StatusCode)
		}
	})
}
