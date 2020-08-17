package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
)

// Bool represents a nullable bool value.
type Bool sql.NullBool

// B converts the provided value into a nullable bool
func B(b bool) Bool {
	return Bool{Valid: true, Bool: b}
}

// Eq returns true if the nullable bool is non-null and is
// equal to other.
func (b Bool) Eq(other bool) bool {
	if !b.Valid {
		return false
	}

	return b.Bool == other
}

// Set sets the value of this nullable.
func (b *Bool) Set(value bool) {
	b.Valid = true
	b.Bool = value
}

// IsNull implements Nullable
func (b Bool) IsNull() bool {
	return !b.Valid
}

// BindStr implements Nullable
func (b *Bool) BindStr(val string) {
	v, err := strconv.ParseBool(val)
	if err != nil {
		b.Valid = false
		b.Bool = false
		return
	}

	b.Valid = true
	b.Bool = v
}

// Scan implements the sql.Scanner interface.
func (b *Bool) Scan(value interface{}) error {
	var x sql.NullBool
	if err := x.Scan(value); err != nil {
		return err
	}

	b.Bool = x.Bool
	b.Valid = x.Valid
	return nil
}

// Value implements the driver.Valuer interface.
func (b Bool) Value() (driver.Value, error) {
	if b.Valid && b.Bool {
		return "1", nil
	}

	return "0", nil
}

// UnmarshalJSON converts the value from a json value.
func (b *Bool) UnmarshalJSON(bs []byte) error {
	if len(bs) == 0 {
		b.Valid = false
		b.Bool = false
		return nil
	}

	if err := json.Unmarshal(bs, &b.Bool); err != nil {
		return err
	}

	b.Valid = true
	return nil
}

// MarshalJSON converts the type to a valid json value.
func (b Bool) MarshalJSON() ([]byte, error) {
	if b.Valid {
		return json.Marshal(b.Bool)
	}

	return json.Marshal(nil)
}
