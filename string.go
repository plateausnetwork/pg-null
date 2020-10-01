package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"

	"golang.org/x/crypto/sha3"
)

// String represents a nullable string value.
type String sql.NullString

// S converts the provided string into a nullable string
func S(s string) String {
	return String{String: s, Valid: true}
}

//CheckNull verifies if the string is equivalent to null and, if so, sets the valid attribute to false
func IsNullStr(nullstr string, s String) String {
	if s.String == nullstr {
		return String{String: s.String, Valid: false}
	}
	return s
}

// Eq returns true if the nullable string is non-null and is
// equal to other.
func (s String) Eq(other string) bool {
	if !s.Valid {
		return false
	}

	return s.String == other
}

// Set sets the value of this nullable.
func (s *String) Set(value string) {
	s.String = value
	s.Valid = true
}

// IsNull implements Nullable
func (s String) IsNull() bool {
	return !s.Valid
}

// BindStr implements StrBind
func (s *String) BindStr(val string) {
	s.String = val
	s.Valid = s.String != ""
}

// Scan implements the sql.Scanner interface.
func (s *String) Scan(value interface{}) error {
	var x sql.NullString
	if err := x.Scan(value); err != nil {
		return err
	}

	s.String = x.String
	s.Valid = x.Valid
	return nil
}

// Value implements the driver.Valuer interface.
func (s String) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}

	return s.String, nil
}

// UnmarshalJSON converts the value from a json value.
func (s *String) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		s.Valid = false
		s.String = ""
		return nil
	}

	if err := json.Unmarshal(b, &s.String); err != nil {
		return err
	}

	s.Valid = true

	return nil
}

// MarshalJSON converts the type to a valid json value.
func (s String) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}

	return json.Marshal(nil)
}

// HashSum256 generates a sha3-256 hash from a String
func (s String) HashSum256() String {

	hash := sha3.Sum256([]byte(s.String))
	pwd := hex.EncodeToString(hash[:])

	return S(pwd)
}
