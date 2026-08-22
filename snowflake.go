package wfa

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Snowflake is a unique numeric entity identifier.
type Snowflake uint64

// GetSnowflake converts id to a Snowflake. It accepts a decimal string or
// any Go integer/float type, which covers values decoded from arbitrary
// JSON into interface{}.
func GetSnowflake(id any) (Snowflake, error) {
	switch v := id.(type) {
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}

		return Snowflake(n), nil
	case int:
		return Snowflake(v), nil
	case int8:
		return Snowflake(v), nil
	case int16:
		return Snowflake(v), nil
	case int32:
		return Snowflake(v), nil
	case int64:
		return Snowflake(v), nil
	case uint:
		return Snowflake(v), nil
	case uint8:
		return Snowflake(v), nil
	case uint16:
		return Snowflake(v), nil
	case uint32:
		return Snowflake(v), nil
	case uint64:
		return Snowflake(v), nil
	case float32:
		return Snowflake(v), nil
	case float64:
		return Snowflake(v), nil
	default:
		return 0, fmt.Errorf("wfa: cannot convert %T to Snowflake", id)
	}
}

// MarshalJSON implements json.Marshaler, encoding the Snowflake as a decimal
// string so large IDs survive round-tripping through consumers that decode
// JSON numbers as float64.
func (i Snowflake) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatUint(uint64(i), 10))
}

// UnmarshalJSON implements json.Unmarshaler, accepting either a decimal
// string (matching MarshalJSON) or a bare JSON number (the API's own wire
// format).
func (i *Snowflake) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*i = Snowflake(n)

		return nil
	}

	return json.Unmarshal(b, (*uint64)(i))
}

// String implements fmt.Stringer.
func (i Snowflake) String() string {
	return strconv.FormatUint(uint64(i), 10)
}
