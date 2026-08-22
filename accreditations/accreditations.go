// Package accreditations gives access to the WFA Matchday API's
// /accreditations endpoints.
package accreditations

import (
	"context"
	"fmt"
	"net/url"

	"github.com/TheWFA/api-client-go"
)

// Accreditation is a summary of an accreditation.
type Accreditation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	IssuingBody    string   `json:"issuingBody"`
	ValidityPeriod *int     `json:"validityPeriod"`
	CreatedAt      wfa.Time `json:"createdAt"`
}

// FullAccreditation is an Accreditation with its holder count.
type FullAccreditation struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	IssuingBody    string   `json:"issuingBody"`
	ValidityPeriod *int     `json:"validityPeriod"`
	CreatedAt      wfa.Time `json:"createdAt"`
	HolderCount    int      `json:"holderCount"`
}

// Facets lists the categories and issuing bodies currently in use.
type Facets struct {
	Categories    []string `json:"categories"`
	IssuingBodies []string `json:"issuingBodies"`
}

// ListQuery holds the filters accepted by Service.List.
type ListQuery struct {
	wfa.ListParams
	Category    string
	IssuingBody string
}

func (q ListQuery) encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	wfa.SetString(v, "category", q.Category)
	wfa.SetString(v, "issuingBody", q.IssuingBody)

	return v
}

// Service gives access to the /accreditations endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves a paginated list of accreditations.
func (s *Service) List(ctx context.Context, query ListQuery) (wfa.ListResponse[Accreditation], error) {
	return wfa.Get[wfa.ListResponse[Accreditation]](ctx, s.backend, "/accreditations", query.encode())
}

// Facets retrieves the categories and issuing bodies currently in use.
func (s *Service) Facets(ctx context.Context) (Facets, error) {
	return wfa.Get[Facets](ctx, s.backend, "/accreditations/facets", nil)
}

// Get retrieves detailed information about a specific accreditation.
func (s *Service) Get(ctx context.Context, id string) (FullAccreditation, error) {
	return wfa.Get[FullAccreditation](ctx, s.backend, fmt.Sprintf("/accreditations/%s", id), nil)
}
