package null

import (
	"database/sql"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
)

var (
	nullLiteral  = []byte("null")
	falseLiteral = []byte("false")
	trueLiteral  = []byte("true")
)

// Bool is a nullable bool.
// It does not consider false values to be null.
// It will decode to null, not false, if null.
type Bool struct {
	sql.NullBool
}

func True() Bool { return B(true) }

func False() Bool { return B(false) }

// NewBool creates a new Bool
func NewBool[T ~bool](b T, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  bool(b),
			Valid: valid,
		},
	}
}

// B creates a new Bool that will always be valid.
func B[T ~bool](b T) Bool {
	return BoolFrom(b)
}

// BoolFrom creates a new Bool that will always be valid.
func BoolFrom[T ~bool](b T) Bool {
	return NewBool(b, true)
}

// BoolFromPtr creates a new Bool that will be null if f is nil.
func BoolFromPtr[T ~bool](b *T) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// ValueOrZero returns the inner value if valid, otherwise false.
func (b Bool) ValueOrZero() bool {
	return b.Valid && b.Bool
}

func (b Bool) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !b.Valid {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteToken(jsontext.Bool(b.Bool))
}

func (b *Bool) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if err := dec.SkipValue(); err != nil {
			return fmt.Errorf("reading null for null.Bool: %w", err)
		}
		*b = Bool{}
		return nil

	case jsontext.KindFalse, jsontext.KindTrue:
		tok, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading bool for null.Bool: %w", err)
		}
		*b = B(tok.Bool())
		return nil

	case jsontext.KindBeginObject:
		return json.UnmarshalDecode(dec, &b.NullBool)

	default:
		return fmt.Errorf("unexpected token unmarshalling null.Bool: %s", dec.PeekKind())
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (b *Bool) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "", "null":
		b.Valid = false
		return nil
	case "true":
		b.Bool = true
	case "false":
		b.Bool = false
	default:
		return errors.New("null: invalid input for UnmarshalText:" + str)
	}
	b.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Bool is null.
func (b Bool) MarshalText() ([]byte, error) {
	if !b.Valid {
		return []byte{}, nil
	}
	if !b.Bool {
		return falseLiteral, nil
	}
	return trueLiteral, nil
}

// SetValid changes this Bool's value and also sets it to be non-null.
func (b *Bool) SetValid(v bool) {
	b.Bool = v
	b.Valid = true
}

// Ptr returns a pointer to this Bool's value, or a nil pointer if this Bool is null.
func (b Bool) Ptr() *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

// IsZero returns true for invalid Bools, for future omitempty support (Go 1.4?)
// A non-null Bool with a 0 value will not be considered zero.
func (b Bool) IsZero() bool {
	return !b.Valid
}

// IsDefined implements the easyjson.Optional interface for omitempty-ing.
func (b Bool) IsDefined() bool {
	return !b.IsZero()
}

// Equal returns true if both booleans have the same value or are both null.
func (b Bool) Equal(other Bool) bool {
	return b.Valid == other.Valid && (!b.Valid || b.Bool == other.Bool)
}
