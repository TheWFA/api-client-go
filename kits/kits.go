// Package kits gives access to the WFA Matchday API's /kit-types and
// /teams/{id}/kits endpoints.
package kits

import (
	"context"
	"fmt"

	"github.com/TheWFA/api-client-go"
)

// Type is a kind of kit a team can wear (home, away, alternative).
type Type struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// TeamKit is a single kit belonging to a team.
type TeamKit struct {
	ID           int     `json:"id"`
	TeamID       int     `json:"teamId"`
	KitType      *Type   `json:"kitType"`
	IsGoalkeeper bool    `json:"isGoalkeeper"`
	ImageURL     string  `json:"imageUrl"`
	TextColour   *string `json:"textColour"`
}

// Service gives access to the /kit-types and /teams/{id}/kits endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// Types retrieves the available kit types (home, away, alternative).
func (s *Service) Types(ctx context.Context) (wfa.UnpaginatedListResponse[Type], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[Type]](ctx, s.backend, "/kit-types", nil)
}

// ForTeam retrieves the kits for a specific team.
func (s *Service) ForTeam(ctx context.Context, teamID int) (wfa.UnpaginatedListResponse[TeamKit], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[TeamKit]](ctx, s.backend, fmt.Sprintf("/teams/%d/kits", teamID), nil)
}
