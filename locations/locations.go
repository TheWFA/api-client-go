// Package locations gives access to the WFA Matchday API's /locations
// endpoints.
package locations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Location is a venue where matches are played.
type Location = wfa.LocationRef

// CourtRef is a lightweight reference to a court at a location.
type CourtRef struct {
	ID         int    `json:"id"`
	LocationID int    `json:"locationId"`
	Name       string `json:"name"`
}

// LocationWithCourts is a Location with its associated courts.
type LocationWithCourts struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	AddressFirstLine  string     `json:"addressFirstLine"`
	AddressSecondLine *string    `json:"addressSecondLine"`
	Postcode          string     `json:"postcode"`
	County            string     `json:"county"`
	Country           string     `json:"country"`
	Courts            []CourtRef `json:"courts"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)

	return v
}

// Service gives access to the /locations endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of locations.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Location], error) {
	return wfa.Get[wfa.ListResponse[Location]](ctx, s.backend, "/locations", query.encode())
}

// Get retrieves a single location and its associated courts by ID.
func (s *Service) Get(ctx context.Context, id int) (LocationWithCourts, error) {
	return wfa.Get[LocationWithCourts](ctx, s.backend, fmt.Sprintf("/locations/%d", id), nil)
}
