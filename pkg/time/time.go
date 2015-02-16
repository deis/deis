package time

import "time"

// DeisDatetimeFormat is the standard date/time representation used in Deis.
const DeisDatetimeFormat string = "2006-01-02T15:04:05MST"

// Time represents the standard datetime format used across the Deis Platform.
type Time struct {
	time.Time
}

func (t *Time) format() string {
	return t.Time.Format(DeisDatetimeFormat)
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in Deis' datetime format.
func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Format(`"` + DeisDatetimeFormat + `"`)), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be in Deis' datetime format.
func (t *Time) UnmarshalText(data []byte) (err error) {
	tt, err := time.Parse(DeisDatetimeFormat, string(data))
	*t = Time{tt}
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in Deis' datetime format.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(`"`+DeisDatetimeFormat+`"`, string(data))
	*t = Time{tt}
	return
}
