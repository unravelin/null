package null

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

// Bool is a nullable bool.
// It does not consider false values to be null.
// It will decode to null, not false, if null.
type Bool struct {
	sql.NullBool
}

// NewBool creates a new Bool
func NewBool(b bool, valid bool) Bool {
	return Bool{
		NullBool: sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

// BoolFrom creates a new Bool that will always be valid.
func BoolFrom(b bool) Bool {
	return NewBool(b, true)
}

// BoolFromPtr creates a new Bool that will be null if f is nil.
func BoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}
	return NewBool(*b, true)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will not be considered a null Bool.
// It also supports unmarshalling a sql.NullBool.
func (b *Bool) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		b.Bool = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &b.NullBool)
	case nil:
		b.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Bool", reflect.TypeOf(v).Name())
	}
	b.Valid = err == nil
	return err
}

// UnmarshalEasyJSON is an easy-JSON specific decoder, that should be more efficient than the standard one.
// We expect the value to be either `null` or `true`, but we also unmarshal if we receive
// `{"Valid":true,"Bool":false}`
func (b *Bool) UnmarshalEasyJSON(w *jlexer.Lexer) {
	if w.IsNull() {
		w.Skip()
		b.Valid = false
		return
	}
	if w.IsDelim('{') {
		w.Skip()
		for !w.IsDelim('}') {
			key := w.UnsafeString()
			w.WantColon()
			if w.IsNull() {
				w.Skip()
				w.WantComma()
				continue
			}
			switch key {
			case "bool", "Bool":
				b.Bool = w.Bool()
			case "valid", "Valid":
				b.Valid = w.Bool()
			}
			w.WantComma()
		}
		return
	}
	b.Bool = w.Bool()
	b.Valid = (w.Error() == nil)
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Bool if the input is a blank or not an integer.
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
		b.Valid = false
		return errors.New("invalid input:" + str)
	}
	b.Valid = true
	return nil
}

var (
	nullLiteral  = []byte("null")
	falseLiteral = []byte("false")
	trueLiteral  = []byte("true")
)

// MarshalJSON implements json.Marshaler.
// It will encode null if this Bool is null.
func (b Bool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return nullLiteral, nil
	}
	if !b.Bool {
		return falseLiteral, nil
	}
	return trueLiteral, nil
}

func (b Bool) MarshalEasyJSON(w *jwriter.Writer) {
	if !b.Valid {
		w.RawString("null")
		return
	}
	w.Bool(b.Bool)
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

func (b Bool) IsDefined() bool {
	return b.Valid
}
