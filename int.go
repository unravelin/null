package null

import (
	"database/sql"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"strconv"
)

// Int is an nullable int64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Int struct {
	sql.NullInt64
}

// NewInt creates a new Int
func NewInt[T ~int64 | ~int](i T, valid bool) Int {
	return Int{
		NullInt64: sql.NullInt64{
			Int64: int64(i),
			Valid: valid,
		},
	}
}

// I creates a new Int that will always be valid.
func I[T ~int64 | ~int](i T) Int {
	return IntFrom(i)
}

// IntFrom creates a new Int that will always be valid.
func IntFrom[T ~int64 | ~int](i T) Int {
	return NewInt(i, true)
}

// IntFromPtr creates a new Int that be null if i is nil.
func IntFromPtr[T ~int64 | ~int](i *T) Int {
	if i == nil {
		return NewInt(0, false)
	}
	return NewInt(*i, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (i Int) ValueOrZero() int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

func (i Int) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !i.Valid {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteToken(jsontext.Int(i.Int64))
}

func (i *Int) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch kind := dec.PeekKind(); kind {
	case jsontext.KindNull:
		if err := dec.SkipValue(); err != nil {
			return fmt.Errorf("reading null for null.Int: %w", err)
		}
		*i = Int{}
		return nil

	case jsontext.KindNumber:
		val, err := dec.ReadValue()
		if err != nil {
			return fmt.Errorf("reading number for null.Int: %w", err)
		}
		// If we use the inbuilt stuff it will silently convert floats to ints,
		// which is not what we want.
		ival, err := strconv.ParseInt(val.String(), 10, 64)
		if err != nil {
			return fmt.Errorf("parsing number for null.Int: %w", err)
		}
		*i = I(ival)
		return nil

	case jsontext.KindString:
		// We want to be able to read stringified integers.
		tok, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading string for int")
		}
		return json.Unmarshal([]byte(tok.String()), i, dec.Options())

	case jsontext.KindBeginObject:
		// We want to try two different things here, so we extract the object
		// value
		val, err := dec.ReadValue()
		if err != nil {
			return err
		}
		var sq sql.NullInt64
		if err := json.Unmarshal(val, &sq, dec.Options()); err == nil {
			i.NullInt64 = sq
			return nil
		}
		// Try a string version
		var si struct {
			Int64 int64 `json:",string"`
			Valid bool
		}
		if err := json.Unmarshal(val, &si, dec.Options()); err != nil {
			return err
		}
		i.NullInt64 = sql.NullInt64(si)
		return nil

	default:
		return &json.SemanticError{
			Err: fmt.Errorf("unexpected token unmarshalling null.Int: %s", kind),
		}
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	i.Int64, err = strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	i.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Int is null.
func (i Int) MarshalText() ([]byte, error) {
	if !i.Valid {
		return []byte{}, nil
	}
	return strconv.AppendInt(nil, i.Int64, 10), nil
}

// SetValid changes this Int's value and also sets it to be non-null.
func (i *Int) SetValid(n int64) {
	i.Int64 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int's value, or a nil pointer if this Int is null.
func (i Int) Ptr() *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

// IsZero returns true for invalid Ints, for future omitempty support (Go 1.4?)
// A non-null Int with a 0 value will not be considered zero.
func (i Int) IsZero() bool {
	return !i.Valid
}

func (i Int) IsDefined() bool {
	return !i.IsZero()
}

// Equal returns true if both ints have the same value or are both null.
func (i Int) Equal(other Int) bool {
	return i.Valid == other.Valid && (!i.Valid || i.Int64 == other.Int64)
}
