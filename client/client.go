// Package client aggregates every WFA Matchday API resource package behind
// a single Client.
package client

import (
	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/accreditations"
	"github.com/TheWFA/api-client-go/clubs"
	"github.com/TheWFA/api-client-go/competitions"
	"github.com/TheWFA/api-client-go/health"
	"github.com/TheWFA/api-client-go/history"
	"github.com/TheWFA/api-client-go/kits"
	"github.com/TheWFA/api-client-go/locations"
	"github.com/TheWFA/api-client-go/matches"
	"github.com/TheWFA/api-client-go/organisations"
	"github.com/TheWFA/api-client-go/persons"
	"github.com/TheWFA/api-client-go/search"
	"github.com/TheWFA/api-client-go/seasons"
	"github.com/TheWFA/api-client-go/suspensions"
	"github.com/TheWFA/api-client-go/teams"
	"github.com/TheWFA/api-client-go/ties"
)

// Client gives access to every WFA Matchday API resource. Construct one with
// New.
type Client struct {
	Health         *health.Service
	Matches        *matches.Service
	Locations      *locations.Service
	Teams          *teams.Service
	Clubs          *clubs.Service
	Competitions   *competitions.Service
	Organisations  *organisations.Service
	Seasons        *seasons.Service
	Accreditations *accreditations.Service
	Persons        *persons.Service
	Search         *search.Service
	History        *history.Service
	Suspensions    *suspensions.Service
	Ties           *ties.Service
	Kits           *kits.Service
}

// New constructs a Client authenticated with apiKey, sent on every request
// as the x-api-key header. Options configure the shared transport — see
// wfa.WithBaseURL, wfa.WithHTTPClient, wfa.WithHeader and wfa.WithUserAgent.
func New(apiKey string, opts ...wfa.Option) (*Client, error) {
	backend, err := wfa.NewBackend(apiKey, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Health:         health.New(backend),
		Matches:        matches.New(backend),
		Locations:      locations.New(backend),
		Teams:          teams.New(backend),
		Clubs:          clubs.New(backend),
		Competitions:   competitions.New(backend),
		Organisations:  organisations.New(backend),
		Seasons:        seasons.New(backend),
		Accreditations: accreditations.New(backend),
		Persons:        persons.New(backend),
		Search:         search.New(backend),
		History:        history.New(backend),
		Suspensions:    suspensions.New(backend),
		Ties:           ties.New(backend),
		Kits:           kits.New(backend),
	}, nil
}
