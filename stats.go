package wfa

import "net/url"

// StatsFilterQuery holds the aggregate-stats filters shared by the team,
// competition and person stats endpoints.
type StatsFilterQuery struct {
	SeasonID      []int
	CompetitionID []int
	MatchGroupID  []int
	TeamID        []int
	// From is an ISO date or datetime string.
	From string
	// To is an ISO date or datetime string.
	To string
}

// Apply writes q's parameters into v.
func (q StatsFilterQuery) Apply(v url.Values) {
	SetInts(v, "seasonId", q.SeasonID)
	SetInts(v, "competitionId", q.CompetitionID)
	SetInts(v, "matchGroupId", q.MatchGroupID)
	SetInts(v, "teamId", q.TeamID)
	SetString(v, "from", q.From)
	SetString(v, "to", q.To)
}

// PlayerStatsOrderBy is a sort key accepted by the per-player stats
// endpoints.
type PlayerStatsOrderBy string

const (
	PlayerStatsOrderByName          PlayerStatsOrderBy = "name"
	PlayerStatsOrderByAppearances   PlayerStatsOrderBy = "appearances"
	PlayerStatsOrderByGoals         PlayerStatsOrderBy = "goals"
	PlayerStatsOrderByAssists       PlayerStatsOrderBy = "assists"
	PlayerStatsOrderByContributions PlayerStatsOrderBy = "contributions"
	PlayerStatsOrderByYellowCards   PlayerStatsOrderBy = "yellowCards"
	PlayerStatsOrderByRedCards      PlayerStatsOrderBy = "redCards"
)

// PlayerStatsQuery holds the filters accepted by the team and competition
// per-player stats endpoints.
type PlayerStatsQuery struct {
	ListParams
	StatsFilterQuery
	OrderBy PlayerStatsOrderBy
}

// Encode renders q as URL query parameters.
func (q PlayerStatsQuery) Encode() url.Values {
	v := url.Values{}
	q.ListParams.Apply(v)
	q.StatsFilterQuery.Apply(v)
	SetString(v, "orderBy", string(q.OrderBy))

	return v
}

// PlayerStatsRow is a per-player stats aggregate row, shared by the team and
// competition player-stats endpoints.
type PlayerStatsRow struct {
	Player        PersonRef `json:"player"`
	Team          *TeamRef  `json:"team"`
	Number        *int      `json:"number"`
	Appearances   int       `json:"appearances"`
	Goals         int       `json:"goals"`
	Assists       int       `json:"assists"`
	Contributions int       `json:"contributions"`
	YellowCards   int       `json:"yellowCards"`
	RedCards      int       `json:"redCards"`
}
