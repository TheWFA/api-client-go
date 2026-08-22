# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

## [1.2.0] - 2026-08-22

### Added

- `matches.ListQuery.ID`, filtering the matches list by a list of match IDs.
- `matches.GoalEvent.Assister`, the assisting player on a goal event (`*wfa.PersonRef`, `nil` for unassisted goals).

## [1.1.0] - 2026-08-22

### Added

- `wfa.Snowflake`, an unsigned 64-bit ID type used for every entity ID and foreign-key reference across the library. Marshals as a decimal string (so large IDs survive round-tripping through consumers that decode JSON numbers as `float64`) and unmarshals from either a string or a bare number, matching the API's own wire format. `wfa.GetSnowflake` converts an arbitrary caller-supplied value (string, or any integer/float type); `wfa.SnowflakePtr` builds a `*Snowflake` for optional query filters.
- A `.githooks/pre-commit` hook running gofmt, vet, golangci-lint (if installed locally) and the test suite before every commit. Enable with `git config core.hooksPath .githooks`.

### Changed

- **Breaking:** every entity ID field and ID-typed filter (e.g. `Match.ID`, `TeamRef.ID`, `MatchListQuery.TeamID`, `CompetitionsService.Get`'s `seasonID` parameter) is now `wfa.Snowflake` or `*wfa.Snowflake`/`[]wfa.Snowflake`, replacing plain `int`/`string`.

## [1.0.0] - 2026-08-22

Initial release.

### Added

- Full coverage of the WFA Matchday API: matches, teams, clubs, competitions, organisations, seasons, persons, accreditations, suspensions, ties, kits, locations, search, and history, plus health checks — one package per resource, aggregated behind a `client.Client`.
- Webhook signature verification and event parsing (`webhooks` package), with no dependency on the HTTP transport.
- `*wfa.APIError` with `Is*` status-code predicates (`IsNotFound`, `IsForbidden`, `IsRateLimited`, `IsBadRequest`, `IsUnauthorized`).
- A live-API integration test suite (`client/integration_test.go`), gated behind the `integration` build tag and skipped unless `WFA_API_KEY` is set.
- GitHub Actions CI: build and test across two Go versions, golangci-lint, and a live-API integration job.

### Fixed

- `Time` parsing rejected the API's actual timestamp format for UTC values with fractional seconds (a bare `+00` offset with no minutes or colon) — found via the live integration suite.
- `Accreditation.ID` was typed as `string`, following the reference JS client; the API actually returns an integer — also found via the live integration suite.

[Unreleased]: https://github.com/TheWFA/api-client-go/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/TheWFA/api-client-go/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/TheWFA/api-client-go/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/TheWFA/api-client-go/releases/tag/v1.0.0
