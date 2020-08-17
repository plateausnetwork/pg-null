package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Int64 represents a nullable int64 value.
type Int64 sql.NullInt64

// I64 converts the provided value into a nullable string
func I64(i int64) Int64 {
	return Int64{Int64: i, Valid: true}
}

// Eq returns true if the nullable int64 is non-null and is
// equal to other.
func (i Int64) Eq(other int64) bool {
	if !i.Valid {
		return false
	}

	return i.Int64 == other
}

// Set sets the value of this nullable.
func (i *Int64) Set(value int64) {
	i.Int64 = value
	i.Valid = true
}

// IsNull implements Nullable
func (i Int64) IsNull() bool {
	return !i.Valid
}

// BindStr implements Nullable
func (i *Int64) BindStr(val string) {
	v, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		i.Valid = false
		i.Int64 = 0
		return
	}

	i.Valid = true
	i.Int64 = v
}

// Scan implements the sql.Scanner interface.
func (i *Int64) Scan(value interface{}) error {
	var x sql.NullInt64
	if err := x.Scan(value); err != nil {
		return err
	}

	i.Int64 = x.Int64
	i.Valid = x.Valid
	return nil
}

// Value implements the driver.Valuer interface.
func (i Int64) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	return i.Int64, nil
}

// UnmarshalJSON converts the value from a json value.
func (i *Int64) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		i.Valid = false
		i.Int64 = 0
		return nil
	}

	if err := json.Unmarshal(b, &i.Int64); err != nil {
		return err
	}

	i.Valid = true
	return nil
}

// MarshalJSON converts the type to a valid json value.
func (i Int64) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int64)
	}

	return json.Marshal(nil)
}
