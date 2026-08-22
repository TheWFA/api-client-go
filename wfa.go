// Package wfa holds the shared core of the WFA (Wheelchair Football
// Association) Matchday API client: the HTTP transport (Backend), the
// pagination and error types, and the reference types (TeamRef, PersonRef,
// SeasonRef, and so on) shared across every resource.
//
// Most programs won't use this package directly. Instead, use
// github.com/TheWFA/api-client-go/client, which wires a Backend up to every
// resource package (matches, teams, clubs, ...) behind a single Client:
//
//	c, err := client.New(os.Getenv("WFA_API_KEY"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	matches, err := c.Matches.List(ctx, matches.ListQuery{
//		OrderByDateDesc: wfa.Bool(true),
//	})
//
// Constructing resource packages individually against a shared Backend is
// also supported, for programs that only need a subset of the API:
//
//	backend, err := wfa.NewBackend(os.Getenv("WFA_API_KEY"))
//	matchesSvc := matches.New(backend)
//
// The webhooks package (github.com/TheWFA/api-client-go/webhooks) verifies
// and parses inbound webhook deliveries; it doesn't depend on Backend at all.
package wfa

// Int returns a pointer to v. It's a convenience for populating the optional
// *int fields on query types.
func Int(v int) *int { return &v }

// Bool returns a pointer to v. It's a convenience for populating the optional
// *bool fields on query types.
func Bool(v bool) *bool { return &v }

// SnowflakePtr returns a pointer to v. It's a convenience for populating the
// optional *Snowflake fields on query types.
func SnowflakePtr(v Snowflake) *Snowflake { return &v }
