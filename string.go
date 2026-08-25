// Package null contains SQL types that consider zero input and null input as separate values,
// with convenient support for JSON and text marshaling.
// Types in this package will always encode to their null value if null.
// Use the zero subpackage if you want zero values and null to be treated the same.
package null

import (
	"database/sql"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// String is a nullable string. It supports SQL and JSON serialization.
// It will marshal to null if null. Blank string input will be considered null.
type String struct {
	sql.NullString
}

// NewString creates a new String
func NewString[T ~string](s T, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: string(s),
			Valid:  valid,
		},
	}
}

// S creates a new String that will never be blank.
func S[T ~string](s T) String {
	return StringFrom(s)
}

// ZS creates a new String that is valid if s is not zero.
func ZS[T ~string](s T) String {
	return NewString(s, s != "")
}

// StringFrom creates a new String that will never be blank.
func StringFrom[T ~string](s T) String {
	return NewString(s, true)
}

// StringFromPtr creates a new String that be null if s is nil.
func StringFromPtr[T ~string](s *T) String {
	if s == nil {
		return NewString("", false)
	}
	return NewString(*s, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (s String) ValueOrZero() string {
	if !s.Valid {
		return ""
	}
	return s.String
}

func (s String) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !s.Valid {
		return enc.WriteToken(jsontext.Null)
	}
	return enc.WriteToken(jsontext.String(s.String))
}

func (s *String) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if err := dec.SkipValue(); err != nil {
			return fmt.Errorf("reading null for null.String: %w", err)
		}
		*s = String{}
		return nil

	case jsontext.KindString:
		tok, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading string for null.String: %w", err)
		}
		*s = S(tok.String())
		return nil

	case jsontext.KindBeginObject:
		return json.UnmarshalDecode(dec, &s.NullString)

	default:
		return fmt.Errorf("unexpected token unmarshalling null.String: %s", dec.PeekKind())
	}
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// SetValid changes this String's value and also sets it to be non-null.
func (s *String) SetValid(v string) {
	s.String = v
	s.Valid = true
}

// Ptr returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null strings, for potential future omitempty support.
func (s String) IsZero() bool {
	return !s.Valid
}

func (s String) IsDefined() bool {
	return !s.IsZero()
}

// Equal returns true if both strings have the same value or are both null.
func (s String) Equal(other String) bool {
	return s.Valid == other.Valid && (!s.Valid || s.String == other.String)
}
