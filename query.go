package wfa

import (
	"fmt"
	"net/url"
	"strconv"
)

// ListParams holds the pagination and free-text query parameters shared by
// (almost) every list endpoint. Embed it in a resource-specific query type.
type ListParams struct {
	// Page is the 1-indexed page number to retrieve.
	Page *int
	// ItemsPerPage is the number of items to return per page.
	ItemsPerPage *int
	// Query is a free-text filter, when the endpoint supports one.
	Query string
}

// Apply writes p's parameters into v.
func (p ListParams) Apply(v url.Values) {
	SetInt(v, "page", p.Page)
	SetInt(v, "itemsPerPage", p.ItemsPerPage)
	SetString(v, "query", p.Query)
}

// SetString sets key to val in v, unless val is empty.
func SetString(v url.Values, key, val string) {
	if val != "" {
		v.Set(key, val)
	}
}

// SetInt sets key to val in v, unless val is nil.
func SetInt(v url.Values, key string, val *int) {
	if val != nil {
		v.Set(key, strconv.Itoa(*val))
	}
}

// SetBool sets key to val in v, unless val is nil.
func SetBool(v url.Values, key string, val *bool) {
	if val != nil {
		v.Set(key, strconv.FormatBool(*val))
	}
}

// SetInts encodes vals as repeated indexed parameters (e.g. teamId[0]=1&teamId[1]=2),
// matching the wire format produced by the reference JavaScript client's `qs.stringify`.
func SetInts(v url.Values, key string, vals []int) {
	for i, n := range vals {
		v.Set(fmt.Sprintf("%s[%d]", key, i), strconv.Itoa(n))
	}
}

// SetSnowflake sets key to val in v, unless val is nil.
func SetSnowflake(v url.Values, key string, val *Snowflake) {
	if val != nil {
		v.Set(key, val.String())
	}
}

// SetSnowflakes encodes vals as repeated indexed parameters, matching SetInts.
func SetSnowflakes(v url.Values, key string, vals []Snowflake) {
	for i, n := range vals {
		v.Set(fmt.Sprintf("%s[%d]", key, i), n.String())
	}
}

// SetStrings encodes vals as repeated indexed parameters, matching SetInts.
func SetStrings(v url.Values, key string, vals []string) {
	for i, s := range vals {
		v.Set(fmt.Sprintf("%s[%d]", key, i), s)
	}
}

// SetEnums encodes vals as repeated indexed parameters, matching SetInts, for
// any named string type.
func SetEnums[T ~string](v url.Values, key string, vals []T) {
	for i, s := range vals {
		v.Set(fmt.Sprintf("%s[%d]", key, i), string(s))
	}
}
