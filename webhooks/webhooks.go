// Package webhooks verifies and parses inbound WFA Matchday API webhook
// deliveries. It has no dependency on wfa.Backend — signature verification
// and payload parsing are pure functions over a delivery's raw body and
// headers.
package webhooks

import (
	"encoding/json"
	"fmt"

	"github.com/TheWFA/api-client-go"
	"github.com/TheWFA/api-client-go/matches"
)

// EventType discriminates the concrete type of an Event.
type EventType string

const (
	EventTypeMatchStatusChanged     EventType = "MatchStatusChanged"
	EventTypeGoalScored             EventType = "GoalScored"
	EventTypeCardIssued             EventType = "CardIssued"
	EventTypeSubstitutionMade       EventType = "SubstitutionMade"
	EventTypePenaltyShootoutAttempt EventType = "PenaltyShootoutAttempt"
	EventTypeMatchScoreCorrected    EventType = "MatchScoreCorrected"
	EventTypePing                   EventType = "WebhookPing"
)

// CardType is the subset of matches.EventType a CardIssuedEvent can carry:
// matches.EventTypeYellowCard or matches.EventTypeRedCard.
type CardType = matches.EventType

// ResolvedTeamRef holds a team reference that has either been resolved to
// its full wfa.TeamRef, or has fallen back to its bare ID — e.g. the record
// was deleted, or raced the write that triggered the event. Check Resolved
// to tell which.
type ResolvedTeamRef struct {
	Team *wfa.TeamRef
	ID   wfa.Snowflake
}

// Resolved reports whether the reference was resolved to a full wfa.TeamRef.
func (r ResolvedTeamRef) Resolved() bool { return r.Team != nil }

// UnmarshalJSON implements json.Unmarshaler.
func (r *ResolvedTeamRef) UnmarshalJSON(data []byte) error {
	var id wfa.Snowflake
	if err := json.Unmarshal(data, &id); err == nil {
		r.ID = id
		return nil
	}

	var team wfa.TeamRef
	if err := json.Unmarshal(data, &team); err != nil {
		return err
	}

	r.Team = &team

	return nil
}

// ResolvedPersonRef holds a person reference that has either been resolved
// to its full wfa.PersonRef, or has fallen back to its bare ID. Check
// Resolved to tell which.
type ResolvedPersonRef struct {
	Person *wfa.PersonRef
	ID     wfa.Snowflake
}

// Resolved reports whether the reference was resolved to a full
// wfa.PersonRef.
func (r ResolvedPersonRef) Resolved() bool { return r.Person != nil }

// UnmarshalJSON implements json.Unmarshaler.
func (r *ResolvedPersonRef) UnmarshalJSON(data []byte) error {
	var id wfa.Snowflake
	if err := json.Unmarshal(data, &id); err == nil {
		r.ID = id
		return nil
	}

	var person wfa.PersonRef
	if err := json.Unmarshal(data, &person); err != nil {
		return err
	}

	r.Person = &person

	return nil
}

// MatchScore is a match's score, as embedded in Match.
type MatchScore struct {
	Home        int `json:"home"`
	Away        int `json:"away"`
	HomePenalty int `json:"homePenalty"`
	AwayPenalty int `json:"awayPenalty"`
}

// Match is a resolved match reference, as embedded in webhook events.
type Match struct {
	ID           wfa.Snowflake          `json:"id"`
	Status       matches.Status         `json:"status"`
	ScheduledFor *wfa.Time              `json:"scheduledFor"`
	Competition  wfa.CompetitionMiniRef `json:"competition"`
	HomeTeam     wfa.TeamRef            `json:"homeTeam"`
	AwayTeam     wfa.TeamRef            `json:"awayTeam"`
	// Score reflects the state after this event — e.g. a goal's payload shows
	// the tally including that goal.
	Score MatchScore `json:"score"`
}

// ResolvedMatchRef holds a match reference that has either been resolved to
// its full Match, or has fallen back to its bare ID. Check Resolved to tell
// which.
type ResolvedMatchRef struct {
	Match *Match
	ID    wfa.Snowflake
}

// Resolved reports whether the reference was resolved to a full Match.
func (r ResolvedMatchRef) Resolved() bool { return r.Match != nil }

// UnmarshalJSON implements json.Unmarshaler.
func (r *ResolvedMatchRef) UnmarshalJSON(data []byte) error {
	var id wfa.Snowflake
	if err := json.Unmarshal(data, &id); err == nil {
		r.ID = id
		return nil
	}

	var match Match
	if err := json.Unmarshal(data, &match); err != nil {
		return err
	}

	r.Match = &match

	return nil
}

// Event is implemented by every concrete webhook event type. Use a type
// switch on the concrete type, or compare EventType() against the EventType
// constants.
type Event interface {
	EventType() EventType
}

// MatchStatusChangedEvent fires when a match transitions between statuses
// (e.g. kickoff, half-time, full-time).
type MatchStatusChangedEvent struct {
	Match          ResolvedMatchRef `json:"match"`
	PreviousStatus matches.Status   `json:"previousStatus"`
	NewStatus      matches.Status   `json:"newStatus"`
	OccurredAt     wfa.Time         `json:"occurredAt"`
}

// EventType implements Event.
func (MatchStatusChangedEvent) EventType() EventType { return EventTypeMatchStatusChanged }

// GoalScoredEvent fires when a goal is scored.
type GoalScoredEvent struct {
	Match       ResolvedMatchRef     `json:"match"`
	Team        ResolvedTeamRef      `json:"team"`
	Scorer      *ResolvedPersonRef   `json:"scorer"`
	Assister    *ResolvedPersonRef   `json:"assister"`
	GoalType    matches.GoalType     `json:"goalType"`
	IsPenalty   bool                 `json:"isPenalty"`
	MatchPeriod *matches.EventPeriod `json:"matchPeriod"`
	MatchTime   *int                 `json:"matchTime"`
	OccurredAt  wfa.Time             `json:"occurredAt"`
}

// EventType implements Event.
func (GoalScoredEvent) EventType() EventType { return EventTypeGoalScored }

// CardIssuedEvent fires when a yellow or red card is issued.
type CardIssuedEvent struct {
	Match       ResolvedMatchRef     `json:"match"`
	Team        ResolvedTeamRef      `json:"team"`
	Player      ResolvedPersonRef    `json:"player"`
	CardType    CardType             `json:"cardType"`
	MatchPeriod *matches.EventPeriod `json:"matchPeriod"`
	MatchTime   *int                 `json:"matchTime"`
	OccurredAt  wfa.Time             `json:"occurredAt"`
}

// EventType implements Event.
func (CardIssuedEvent) EventType() EventType { return EventTypeCardIssued }

// SubstitutionMadeEvent fires when a substitution is made.
type SubstitutionMadeEvent struct {
	Match       ResolvedMatchRef     `json:"match"`
	Team        ResolvedTeamRef      `json:"team"`
	PlayerOn    ResolvedPersonRef    `json:"playerOn"`
	PlayerOff   ResolvedPersonRef    `json:"playerOff"`
	MatchPeriod *matches.EventPeriod `json:"matchPeriod"`
	MatchTime   *int                 `json:"matchTime"`
	OccurredAt  wfa.Time             `json:"occurredAt"`
}

// EventType implements Event.
func (SubstitutionMadeEvent) EventType() EventType { return EventTypeSubstitutionMade }

// PenaltyShootoutAttemptEvent fires for each attempt in a penalty shootout.
type PenaltyShootoutAttemptEvent struct {
	Match      ResolvedMatchRef  `json:"match"`
	Team       ResolvedTeamRef   `json:"team"`
	Player     ResolvedPersonRef `json:"player"`
	Sequence   int               `json:"sequence"`
	Scored     *bool             `json:"scored"`
	OccurredAt wfa.Time          `json:"occurredAt"`
}

// EventType implements Event.
func (PenaltyShootoutAttemptEvent) EventType() EventType { return EventTypePenaltyShootoutAttempt }

// MatchScoreCorrectedEvent fires only when editing or deleting an existing
// goal or shootout attempt changes the computed score — never alongside the
// event that produced the score in the first place.
type MatchScoreCorrectedEvent struct {
	Match      ResolvedMatchRef `json:"match"`
	OccurredAt wfa.Time         `json:"occurredAt"`
}

// EventType implements Event.
func (MatchScoreCorrectedEvent) EventType() EventType { return EventTypeMatchScoreCorrected }

// PingEvent is a synthetic verification event — never published from the
// real match feed, but delivered through the same signing/HTTP path as real
// events.
type PingEvent struct {
	OccurredAt wfa.Time `json:"occurredAt"`
}

// EventType implements Event.
func (PingEvent) EventType() EventType { return EventTypePing }

// ParsePayload parses a webhook delivery body into a typed Event, narrowed
// by its "detailType" field. It does not verify the signature — see
// VerifySignature.
func ParsePayload(body []byte) (Event, error) {
	var probe struct {
		DetailType EventType `json:"detailType"`
	}

	if err := json.Unmarshal(body, &probe); err != nil {
		return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: webhook payload is not a valid JSON object: %v", err)}
	}

	switch probe.DetailType {
	case EventTypeMatchStatusChanged:
		var e MatchStatusChangedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse MatchStatusChanged payload: %v", err)}
		}

		return e, nil
	case EventTypeGoalScored:
		var e GoalScoredEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse GoalScored payload: %v", err)}
		}

		return e, nil
	case EventTypeCardIssued:
		var e CardIssuedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse CardIssued payload: %v", err)}
		}

		return e, nil
	case EventTypeSubstitutionMade:
		var e SubstitutionMadeEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse SubstitutionMade payload: %v", err)}
		}

		return e, nil
	case EventTypePenaltyShootoutAttempt:
		var e PenaltyShootoutAttemptEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse PenaltyShootoutAttempt payload: %v", err)}
		}

		return e, nil
	case EventTypeMatchScoreCorrected:
		var e MatchScoreCorrectedEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse MatchScoreCorrected payload: %v", err)}
		}

		return e, nil
	case EventTypePing:
		var e PingEvent
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: could not parse WebhookPing payload: %v", err)}
		}

		return e, nil
	default:
		return nil, &Error{Kind: ErrorKindPayload, Message: fmt.Sprintf("wfa: unknown webhook event type: %q", probe.DetailType)}
	}
}
