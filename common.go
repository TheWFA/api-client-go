package wfa

// TeamRef is a lightweight reference to a team, embedded in most other
// resources.
type TeamRef struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	BadgeURL string `json:"badgeUrl"`
}

// TeamMiniRef is the smaller team reference used inside suspension servedIn
// entries.
type TeamMiniRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ClubRef is a lightweight reference to a club.
type ClubRef struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	ClubLogo *string `json:"clubLogo"`
}

// OrganisationRef is a lightweight reference to an organisation.
type OrganisationRef struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ShortName *string `json:"shortName"`
	BadgeURL  *string `json:"badgeUrl"`
}

// CompetitionRef is a lightweight reference to a competition.
type CompetitionRef struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	BadgeURL *string `json:"badgeUrl"`
}

// CompetitionRefWithOrganisation is a CompetitionRef that also carries its
// parent organisation.
type CompetitionRefWithOrganisation struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	BadgeURL     *string          `json:"badgeUrl"`
	Organisation *OrganisationRef `json:"organisation"`
}

// CompetitionMiniRef is the smaller competition reference used inside
// suspension and tie entries.
type CompetitionMiniRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SeasonRef is a lightweight reference to a season.
type SeasonRef struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate Time   `json:"startDate"`
	EndDate   Time   `json:"endDate"`
}

// SeasonMiniRef is the smaller season reference used inside suspension and
// tie entries.
type SeasonMiniRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SeasonWithActive is a SeasonRef with its active flag.
type SeasonWithActive struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartDate Time   `json:"startDate"`
	EndDate   Time   `json:"endDate"`
	Active    bool   `json:"active"`
}

// SeasonZoneKind categorizes a SeasonZone.
type SeasonZoneKind string

const (
	SeasonZoneKindPromotion  SeasonZoneKind = "promotion"
	SeasonZoneKindPlayoff    SeasonZoneKind = "playoff"
	SeasonZoneKindRelegation SeasonZoneKind = "relegation"
)

// SeasonZone is a highlighted band of table positions, e.g. a promotion or
// relegation zone.
type SeasonZone struct {
	ID           string         `json:"id"`
	Label        string         `json:"label"`
	Kind         SeasonZoneKind `json:"kind"`
	FromPosition int            `json:"fromPosition"`
	ToPosition   int            `json:"toPosition"`
	Color        string         `json:"color"`
}

// SeasonPoints is the points awarded per match outcome for a season.
type SeasonPoints struct {
	Win  int `json:"win"`
	Draw int `json:"draw"`
	Loss int `json:"loss"`
}

// SeasonSetup describes a season's rules: whether classification checks
// apply, the yellow card suspension threshold, points-per-outcome, and any
// table zones.
type SeasonSetup struct {
	ClassCheck      bool          `json:"classCheck"`
	YellowCardLimit int           `json:"yellowCardLimit"`
	Points          *SeasonPoints `json:"points,omitempty"`
	Zones           []SeasonZone  `json:"zones,omitempty"`
}

// SeasonFull is a SeasonWithActive with its full rules setup.
type SeasonFull struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	StartDate   Time        `json:"startDate"`
	EndDate     Time        `json:"endDate"`
	Active      bool        `json:"active"`
	SeasonSetup SeasonSetup `json:"seasonSetup"`
}

// PersonRef is a lightweight reference to a person.
type PersonRef struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
}

// LocationRef is a lightweight reference to a location.
type LocationRef struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	AddressFirstLine  string  `json:"addressFirstLine"`
	AddressSecondLine *string `json:"addressSecondLine"`
	Postcode          string  `json:"postcode"`
	County            string  `json:"county"`
	Country           string  `json:"country"`
}

// HistoryEntity identifies the kind of entity a HistoryEntry describes.
type HistoryEntity string

const (
	HistoryEntityTeam         HistoryEntity = "team"
	HistoryEntityClub         HistoryEntity = "club"
	HistoryEntityCompetition  HistoryEntity = "competition"
	HistoryEntityOrganisation HistoryEntity = "organisation"
)

// HistoryEntry is one superseded-identity window for an entity.
type HistoryEntry struct {
	ID        int                    `json:"id"`
	Entity    HistoryEntity          `json:"entity"`
	EntityID  int                    `json:"entityId"`
	Values    map[string]interface{} `json:"values"`
	ValidFrom *Time                  `json:"validFrom"`
	ValidTo   *Time                  `json:"validTo"`
	CreatedAt Time                   `json:"createdAt"`
}
