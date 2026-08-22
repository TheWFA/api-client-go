package matches

import (
	"encoding/json"
	"fmt"

	"github.com/TheWFA/api-client-go"
)

// EventType discriminates the concrete type of an Event.
type EventType string

const (
	EventTypeGoal         EventType = "goal"
	EventTypeYellowCard   EventType = "yellow_card"
	EventTypeRedCard      EventType = "red_card"
	EventTypeSubstitution EventType = "substitution"
)

// GoalType distinguishes an ordinary goal from an own goal.
type GoalType string

const (
	GoalTypeGoal    GoalType = "goal"
	GoalTypeOwnGoal GoalType = "own-goal"
)

// EventPeriod is the period of play an Event occurred in.
type EventPeriod string

const (
	EventPeriodFirstHalf           EventPeriod = "first-half"
	EventPeriodSecondHalf          EventPeriod = "second-half"
	EventPeriodExtraTime           EventPeriod = "extra-time"
	EventPeriodExtraTimeFirstHalf  EventPeriod = "extra-time-first-half"
	EventPeriodExtraTimeSecondHalf EventPeriod = "extra-time-second-half"
	EventPeriodPenalties           EventPeriod = "penalties"
)

// Event is implemented by every concrete match event type: GoalEvent,
// CardEvent and SubstitutionEvent. Use a type switch on the concrete type,
// or compare EventType() against the EventType constants.
type Event interface {
	EventType() EventType
}

// baseEvent holds the fields common to every match event.
type baseEvent struct {
	CreatedAt   wfa.Time      `json:"createdAt"`
	Time        *int          `json:"time"`
	MatchPeriod *EventPeriod  `json:"matchPeriod"`
	TeamID      wfa.Snowflake `json:"teamId"`
}

// GoalEvent records a goal, including own goals and penalties.
type GoalEvent struct {
	baseEvent
	Player   wfa.PersonRef `json:"player"`
	Penalty  bool          `json:"penalty"`
	GoalType GoalType      `json:"goalType"`
}

// EventType implements Event.
func (GoalEvent) EventType() EventType { return EventTypeGoal }

// CardEvent records a yellow or red card. Check EventType() to tell which.
type CardEvent struct {
	baseEvent
	Type   EventType     `json:"type"`
	Player wfa.PersonRef `json:"player"`
}

// EventType implements Event.
func (e CardEvent) EventType() EventType { return e.Type }

// SubstitutionEvent records a substitution.
type SubstitutionEvent struct {
	baseEvent
	PlayerOn  wfa.PersonRef `json:"playerOn"`
	PlayerOff wfa.PersonRef `json:"playerOff"`
}

// EventType implements Event.
func (SubstitutionEvent) EventType() EventType { return EventTypeSubstitution }

// Events is a list of match events of mixed concrete type, decoded from
// JSON based on each element's "type" field.
type Events []Event

// UnmarshalJSON implements json.Unmarshaler.
func (events *Events) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	result := make(Events, 0, len(raw))

	for _, r := range raw {
		var probe struct {
			Type EventType `json:"type"`
		}

		if err := json.Unmarshal(r, &probe); err != nil {
			return err
		}

		switch probe.Type {
		case EventTypeGoal:
			var e GoalEvent
			if err := json.Unmarshal(r, &e); err != nil {
				return err
			}

			result = append(result, e)
		case EventTypeYellowCard, EventTypeRedCard:
			var e CardEvent
			if err := json.Unmarshal(r, &e); err != nil {
				return err
			}

			result = append(result, e)
		case EventTypeSubstitution:
			var e SubstitutionEvent
			if err := json.Unmarshal(r, &e); err != nil {
				return err
			}

			result = append(result, e)
		default:
			return fmt.Errorf("wfa: unknown match event type %q", probe.Type)
		}
	}

	*events = result

	return nil
}
