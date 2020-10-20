package null

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/nicksnyder/basen"
)

// UUID represents a UUID v4 value
type UUID struct {
	UUID  uuid.UUID
	Valid bool
}

// NewID returns a new, valid ID
func NewID() UUID {
	return UUID{
		UUID:  uuid.Must(uuid.NewRandom()),
		Valid: true,
	}
}

// ParseID parses a base62 value into an UUID.
func ParseID(base62 string) (UUID, error) {
	b, err := basen.Base62Encoding.DecodeString(base62)
	if err != nil {
		return UUID{}, err
	}

	var array [16]byte
	copy(array[:], b)

	uid := uuid.UUID(array)
	if uid.String() == "" {
		return UUID{}, errors.New("invalid uuid")
	}

	return UUID{
		UUID:  uid,
		Valid: true,
	}, nil
}

// ParseUUID parses a uuid string value into an UUID.
func ParseUUID(strUUID string) (UUID, error) {
	v, err := uuid.Parse(strUUID)
	if err != nil {
		return UUID{}, errors.New("invalid parse uuid")
	}

	return UUID{
		UUID:  v,
		Valid: true,
	}, nil
}

// ID parses the provided base62 string and returns the
// corresponding UUID value. Errors are ignored.
func ID(base62 string) UUID {
	res, _ := ParseID(base62)
	return res
}

// Eq returns true if the nullable uuid is non-null and is
// equal to other.
func (uid UUID) Eq(other uuid.UUID) bool {
	if !uid.Valid {
		return false
	}

	return uid.UUID == other
}

// IsZero returns true if the value is zero.
func (uid UUID) IsZero() bool {
	return uid.Valid && uid.UUID == uuid.Nil
}

// IsNull implements Nullable
func (uid UUID) IsNull() bool {
	return !uid.Valid
}

// BindStr implements Nullable
func (uid *UUID) BindStr(val string) {
	v := ID(val)

	uid.Valid = v.Valid
	uid.UUID = v.UUID
}

// String returns the UUID value as a string
func (uid UUID) String() string {
	return uid.UUID.String()
}

// Base62 returns the base62 encoding of the UUID value
func (uid UUID) Base62() string {
	b, err := uid.UUID.MarshalBinary()
	if err != nil {
		return err.Error()
	}

	return basen.Base62Encoding.EncodeToString(b)
}

// Hex returns the hex encoding of the UUID value
func (uid UUID) Hex() string {
	b, err := uid.UUID.MarshalBinary()
	if err != nil {
		return err.Error()
	}

	return hex.EncodeToString(b)
}

// Scan implements the sql.Scanner interface.
func (uid *UUID) Scan(value interface{}) error {
	var x uuid.UUID
	if err := x.Scan(value); err != nil {
		return err
	}

	uid.UUID = x
	uid.Valid = x != uuid.Nil
	return nil
}

// Value implements the driver.Valuer interface.
func (uid UUID) Value() (driver.Value, error) {
	if !uid.Valid {
		return nil, nil
	}

	return uid.UUID.Value()
}

// UnmarshalJSON converts the value from a json value.
func (uid *UUID) UnmarshalJSON(b []byte) error {

	if err := json.Unmarshal(b, &uid.UUID); err != nil {
		return err
	}

	uid.Valid = true
	return nil
}

// MarshalJSON converts the type to a valid json value.
func (uid UUID) MarshalJSON() ([]byte, error) {
	if uid.Valid {
		return json.Marshal(uid.String())
	}

	return json.Marshal(nil)
}
