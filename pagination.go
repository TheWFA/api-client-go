package wfa

// ListResponse is a page of results from a paginated list endpoint.
type ListResponse[T any] struct {
	Items        []T `json:"items"`
	TotalItems   int `json:"totalItems"`
	Page         int `json:"page"`
	ItemsPerPage int `json:"itemsPerPage"`
}

// UnpaginatedListResponse is the response from a list endpoint that returns
// every item in one go, without pagination.
type UnpaginatedListResponse[T any] struct {
	Items      []T `json:"items"`
	TotalItems int `json:"totalItems"`
}
