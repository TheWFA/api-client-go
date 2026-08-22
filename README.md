# WFA API Client (Powerchair Football) - Go

A Go client for accessing **Wheelchair Football Association (WFA)** Matchday data.

This is the Go counterpart to [`@thewfa/api-client`](https://github.com/TheWFA/api-client-js), covering the same endpoints.

---

## Features

- Typed request and response models for every endpoint
- API key auth (`x-api-key` header)
- Full coverage of matches, teams, clubs, competitions, organisations, seasons, persons, accreditations, suspensions, ties, kits, search and history
- Webhook signature verification and payload parsing

---

## Package layout

In Go, a directory is a package boundary, so this client is split into one package per resource rather than one flat package:

| Package                                       | Contents                                                                 |
| ---------------------------------------------- | ------------------------------------------------------------------------ |
| `github.com/TheWFA/api-client-go`             | Core: `Backend` (HTTP transport), pagination/error types, and reference types (`TeamRef`, `PersonRef`, `SeasonRef`, ...) shared across every resource |
| `github.com/TheWFA/api-client-go/client`      | `Client`, aggregating every resource behind one struct — most programs only need this plus the resource packages they call into |
| `github.com/TheWFA/api-client-go/matches`     | Matches: list/get, lineups, events, penalties                            |
| `github.com/TheWFA/api-client-go/teams`       | Teams: list/get, rosters, staff, registrations, seasons played, stats    |
| `github.com/TheWFA/api-client-go/clubs`       | Clubs and their teams                                                    |
| `github.com/TheWFA/api-client-go/competitions`| Competitions, tables, match groups, stats                                |
| `github.com/TheWFA/api-client-go/organisations`| Organisations and the competitions they run                             |
| `github.com/TheWFA/api-client-go/seasons`     | Seasons                                                                  |
| `github.com/TheWFA/api-client-go/persons`     | People: registrations, appearances, stats, suspensions                   |
| `github.com/TheWFA/api-client-go/accreditations`| Accreditations and their facets                                        |
| `github.com/TheWFA/api-client-go/suspensions` | Suspensions (global list, split by origin/served-in match)               |
| `github.com/TheWFA/api-client-go/ties`        | Two-legged ties and aggregate scores                                     |
| `github.com/TheWFA/api-client-go/kits`        | Kit types and per-team kits                                              |
| `github.com/TheWFA/api-client-go/locations`   | Venues and courts                                                        |
| `github.com/TheWFA/api-client-go/history`     | Superseded identities of teams, clubs, competitions and organisations    |
| `github.com/TheWFA/api-client-go/search`      | Fuzzy search across persons, teams, clubs, competitions and matches      |
| `github.com/TheWFA/api-client-go/health`      | API health status                                                        |
| `github.com/TheWFA/api-client-go/webhooks`    | Webhook signature verification and payload parsing (no HTTP dependency)  |

A resource package's types drop the resource name where the package name already says it — e.g. `matches.ListQuery`, not `matches.MatchListQuery`.

---

## Installation

```bash
go get github.com/TheWFA/api-client-go
```

Then import the `client` package plus whichever resource packages you call into:

```go
import (
	"github.com/TheWFA/api-client-go/client"
	"github.com/TheWFA/api-client-go/matches"
)
```

---

## Authentication

Pass your API key to `client.New`. It's sent as:

```
x-api-key: <your_api_key>
```

> **Never commit API keys** to source control. Prefer environment variables.

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/client"
	"github.com/TheWFA/api-client-go/matches"
)

func main() {
	c, err := client.New(os.Getenv("WFA_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List the latest matches.
	page, err := c.Matches.List(ctx, matches.ListQuery{
		OrderByDateDesc: wfa.Bool(true),
		ListParams:      wfa.ListParams{ItemsPerPage: wfa.Int(20)},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch a single match by ID.
	match, err := c.Matches.Get(ctx, page.Items[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(match.HomeTeam.Name, "vs", match.AwayTeam.Name)
}
```

---

## Options

Options configure the shared transport (`wfa.Backend`) and are passed to `client.New`:

```go
c, err := client.New(
	apiKey,
	wfa.WithBaseURL("https://api.thewfa.org.uk/v1"), // defaults to the official base
	wfa.WithHTTPClient(customHTTPClient),
	wfa.WithHeader("X-Extra", "value"),
	wfa.WithUserAgent("my-app/1.0"),
)
```

If a program only needs a handful of resources, it can skip the `client` package and construct a `wfa.Backend` plus only the resource packages it needs:

```go
backend, err := wfa.NewBackend(apiKey)
matchesSvc := matches.New(backend)
```

---

## Pagination

Most list endpoints return a `wfa.ListResponse[T]`:

```go
type ListResponse[T any] struct {
	Items        []T
	TotalItems   int
	Page         int
	ItemsPerPage int
}
```

A few endpoints return every item in one response instead, as `wfa.UnpaginatedListResponse[T]`.

Filters that accept multiple values (e.g. `matches.ListQuery.TeamID`) take a Go slice — pass a single-element slice for one value.

---

## Errors

Non-2xx responses are returned as a `*wfa.APIError`:

```go
match, err := c.Matches.Get(ctx, 999999)
if err != nil {
	if wfa.IsNotFound(err) {
		// not found
	}
	var apiErr *wfa.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode, apiErr.Code, apiErr.Message)
	}
}
```

Other predicates: `wfa.IsBadRequest`, `wfa.IsUnauthorized`, `wfa.IsForbidden`, `wfa.IsRateLimited`.

---

## Webhooks

The `webhooks` package verifies and parses inbound webhook deliveries. It doesn't depend on `wfa.Backend` at all, so it can be imported on its own:

```go
import "github.com/TheWFA/api-client-go/webhooks"

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	event, err := webhooks.ConstructEvent(r, os.Getenv("WFA_WEBHOOK_PUBLIC_KEY"), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch e := event.(type) {
	case webhooks.GoalScoredEvent:
		fmt.Println("goal scored", e.MatchTime)
	case webhooks.MatchStatusChangedEvent:
		fmt.Println("status changed", e.PreviousStatus, "->", e.NewStatus)
	}

	w.WriteHeader(http.StatusOK)
}
```

If the raw body and headers are already available separately (e.g. from a framework that has already read the body), use `webhooks.ConstructEventFromHeaders` instead.

---

## Development

```bash
go build ./...
go vet ./...
gofmt -l .
go test ./...
```

---

## License

MIT © WFA / Contributors
