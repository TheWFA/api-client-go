// Package history gives access to the WFA Matchday API's /history
// endpoints.
package history

import (
	"context"
	"fmt"

	"github.com/TheWFA/api-client-go"
)

// Service gives access to the /history endpoints.
type Service struct {
	backend *wfa.Backend
}

// New constructs a Service backed by backend.
func New(backend *wfa.Backend) *Service {
	return &Service{backend: backend}
}

// List retrieves the superseded identities of an entity, newest window
// first.
//
// Each entry is valid over [ValidFrom, ValidTo); an entry with a nil ValidTo
// is the open-ended current window.
func (s *Service) List(ctx context.Context, entity wfa.HistoryEntity, entityID int) (wfa.UnpaginatedListResponse[wfa.HistoryEntry], error) {
	return wfa.Get[wfa.UnpaginatedListResponse[wfa.HistoryEntry]](ctx, s.backend, fmt.Sprintf("/history/%s/%d", entity, entityID), nil)
}

// Get retrieves a single history entry by ID.
func (s *Service) Get(ctx context.Context, entity wfa.HistoryEntity, entityID, historyID int) (wfa.HistoryEntry, error) {
	return wfa.Get[wfa.HistoryEntry](ctx, s.backend, fmt.Sprintf("/history/%s/%d/%d", entity, entityID, historyID), nil)
}
