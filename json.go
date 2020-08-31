package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON represents a nullable json map
type JSON struct {
	Map   map[string]interface{}
	Valid bool
}

// J returns a new nullable JSON with the specified value
func J(value map[string]interface{}) JSON {
	var j JSON
	j.Map = value
	j.Valid = value != nil

	return j
}

// AsJSON parses the input string, and returns a nullable JSON value
func AsJSON(jsonString string) JSON {
	var j JSON
	j.BindStr(jsonString)

	return j
}

// Eq returns true if the nullable json is non-null and is
// equal to other.
func (j JSON) Eq(other JSON) bool {
	if !j.Valid {
		return false
	}

	return j.String() == other.String()
}

// Set sets the value of this nullable.
func (j *JSON) Set(data map[string]interface{}) {
	j.Valid = true
	j.Map = data
}

// IsNull implements Nullable
func (j JSON) IsNull() bool {
	return !j.Valid
}

// BindStr implements Nullable
func (j *JSON) BindStr(val string) {
	result := make(map[string]interface{})

	if err := json.Unmarshal([]byte(val), &result); err != nil {
		j.Valid = false
		return
	}

	j.Valid = true
	j.Map = result
}

// Scan implements the sql.Scanner interface.
func (j *JSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return j.UnmarshalJSON([]byte(v))

	case []byte:
		return j.UnmarshalJSON(v)

	default:
		return fmt.Errorf("Unsupported json scanner value: %v (%t)", value, value)
	}
}

// Value implements the driver.Valuer interface.
func (j JSON) Value() (driver.Value, error) {
	if j.Valid && j.Map != nil {
		b, err := j.MarshalJSON()
		if err != nil {
			return nil, err
		}

		return string(b), nil
	}

	return nil, nil
}

// RemoveFields remove all the keys from the data
func (j *JSON) RemoveFields(fields []string) {
	for _, field := range fields {
		delete(j.Map, field)
	}
}

// UnmarshalJSON converts the value from a json value.
func (j *JSON) UnmarshalJSON(bs []byte) error {
	if len(bs) == 0 {
		j.Valid = false
		j.Map = nil
		return nil
	}

	result := make(map[string]interface{})

	if err := json.Unmarshal(bs, &result); err != nil {
		j.Valid = false
		return err
	}

	j.Valid = true
	j.Map = result
	return nil
}

// MarshalJSON converts the type to a valid json value.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Valid {
		return json.Marshal(j.Map)
	}

	return json.Marshal(nil)
}

// string implements Stringer interface
func (j JSON) String() string {
	bs, err := j.MarshalJSON()
	if err != nil {
		return ""
	}

	return string(bs)
}
