// Package search gives access to the WFA Matchday API's /search endpoint.
package search

import (
	"context"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// ItemType categorizes an Item.
type ItemType string

const (
	ItemTypePerson       ItemType = "person"
	ItemTypeTeam         ItemType = "team"
	ItemTypeClub         ItemType = "club"
	ItemTypeCompetition  ItemType = "competition"
	ItemTypeOrganisation ItemType = "organisation"
	ItemTypeMatch        ItemType = "match"
)

// Item is a single fuzzy-search result.
type Item struct {
	Type        ItemType `json:"type"`
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Description *string  `json:"description"`
	ImageURL    *string  `json:"imageUrl"`
	Score       float64  `json:"score"`
}

// Query holds the parameters accepted by Service.List.
type Query struct {
	Page         *int
	ItemsPerPage *int
	// Query is the search term. Required; a blank query returns no results.
	Query string
	Type  []ItemType
}

func (q Query) encode() url.Values {
	v := url.Values{}
	wfa.SetInt(v, "page", q.Page)
	wfa.SetInt(v, "itemsPerPage", q.ItemsPerPage)
	wfa.SetString(v, "query", q.Query)
	wfa.SetEnums(v, "type", q.Type)

	return v
}

// Service gives access to the /search endpoint.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List performs a fuzzy search across persons, teams, clubs, competitions,
// organisations and matches, ranked by trigram similarity. A blank query
// returns no results.
func (s *Service) List(ctx context.Context, query Query) (wfa.ListResponse[Item], error) {
	return wfa.Get[wfa.ListResponse[Item]](ctx, s.backend, "/search", query.encode())
}
