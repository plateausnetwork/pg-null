package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	"time"
)

// Time represents a nullable timespan
type Time sql.NullTime

// T converts the provided value into a nullable time
func T(time time.Time) Time {
	return Time{Valid: true, Time: time}
}

// TS converts the provided string value into a nullable time
func TS(time string) Time {
	var t Time
	t.BindStr(time)

	return t
}

// Eq returns true if the nullable time is non-null and is
// equal to other.
func (t Time) Eq(other time.Time) bool {
	if !t.Valid {
		return false
	}

	return t.Time.Equal(other)
}

// Set sets the value of this nullable.
func (t *Time) Set(value time.Time) {
	t.Valid = true
	t.Time = value
}

// IsNull implements Nullable
func (t Time) IsNull() bool {
	return !t.Valid
}

// BindStr implements Nullable
func (t *Time) BindStr(val string) {
	v, err := time.Parse(time.RFC3339, val)
	if err != nil {
		t.Valid = false
		t.Time = time.Time{}
		return
	}

	t.Valid = true
	t.Time = v
}

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		tv, err := time.Parse(time.RFC3339, str)
		if err != nil {
			log.Println("ops", str, err)
			return err
		}

		t.Time = tv
		t.Valid = true
		return nil
	}

	var x sql.NullTime
	if err := x.Scan(value); err != nil {
		return err
	}

	t.Time = x.Time
	t.Valid = x.Valid
	return nil
}

// Value implements the driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if t.Valid {
		return t.Time, nil
	}

	return nil, nil
}

// UnmarshalJSON converts the value from a json value.
func (t *Time) UnmarshalJSON(bs []byte) error {
	if len(bs) == 0 {
		t.Valid = false
		t.Time = time.Time{}
		return nil
	}

	var str string
	if err := json.Unmarshal(bs, &str); err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = parsed
	t.Valid = true
	return nil
}

// MarshalJSON converts the type to a valid json value.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.Valid {
		str := t.Time.Format(time.RFC3339)
		return json.Marshal(str)
	}

	return json.Marshal(nil)
}
