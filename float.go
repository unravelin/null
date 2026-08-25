package null

import (
	"database/sql"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"math"
	"strconv"
)

// Float is a nullable float64.
// It does not consider zero values to be null.
// It will decode to null, not zero, if null.
type Float struct {
	sql.NullFloat64
}

// NewFloat creates a new Float
func NewFloat[T ~float64](f T, valid bool) Float {
	return Float{
		NullFloat64: sql.NullFloat64{
			Float64: float64(f),
			Valid:   valid,
		},
	}
}

// F creates a new Float that will always be valid.
func F[T ~float64](f T) Float {
	return FloatFrom(f)
}

// FloatFrom creates a new Float that will always be valid.
func FloatFrom[T ~float64](f T) Float {
	return NewFloat(f, true)
}

// FloatFromPtr creates a new Float that be null if f is nil.
func FloatFromPtr[T ~float64](f *T) Float {
	if f == nil {
		return NewFloat(float64(0), false)
	}
	return NewFloat(*f, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (f Float) ValueOrZero() float64 {
	if !f.Valid {
		return 0
	}
	return f.Float64
}

func (f Float) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !f.Valid {
		return enc.WriteToken(jsontext.Null)
	}
	if math.IsInf(f.Float64, 0) || math.IsNaN(f.Float64) {
		return &json.SemanticError{
			Err: fmt.Errorf("cannot marshal inf or nan"),
		}
	}

	return enc.WriteToken(jsontext.Float(f.Float64))
}

func (f *Float) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch kind := dec.PeekKind(); kind {
	case jsontext.KindNull:
		_, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading null for null.Float: %w", err)
		}
		*f = Float{}
		return nil

	case jsontext.KindNumber:
		tok, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading number for null.Float: %w", err)
		}
		fval, err := tok.Float()
		if err != nil {
			return fmt.Errorf("parsing number as float: %w", err)
		}
		*f = F(fval)
		return nil

	case jsontext.KindString:
		tok, err := dec.ReadToken()
		if err != nil {
			return fmt.Errorf("reading string for float")
		}
		return json.Unmarshal([]byte(tok.String()), f, dec.Options())

	case jsontext.KindBeginObject:
		// We want to try two different things here, so we extract the object
		// value
		val, err := dec.ReadValue()
		if err != nil {
			return err
		}
		var sq sql.NullFloat64
		if err := json.Unmarshal(val, &sq, dec.Options()); err == nil {
			f.NullFloat64 = sq
			return nil
		}
		// Try a string version
		var sf struct {
			Float64 float64 `json:",string"`
			Valid   bool
		}
		if err := json.Unmarshal(val, &sf, dec.Options()); err != nil {
			return err
		}
		f.NullFloat64 = sql.NullFloat64(sf)
		return nil

	default:
		return &json.SemanticError{
			Err: fmt.Errorf("unexpected token unmarshalling null.Float: %s", kind),
		}
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Float if the input is blank.
// It will return an error if the input is not an integer, blank, or "null".
func (f *Float) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		f.Valid = false
		return nil
	}
	var err error
	f.Float64, err = strconv.ParseFloat(string(text), 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	f.Valid = true
	return err
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string if this Float is null.
func (f Float) MarshalText() ([]byte, error) {
	if !f.Valid {
		return []byte{}, nil
	}
	return strconv.AppendFloat(nil, f.Float64, 'f', -1, 64), nil
}

// SetValid changes this Float's value and also sets it to be non-null.
func (f *Float) SetValid(n float64) {
	f.Float64 = n
	f.Valid = true
}

// Ptr returns a pointer to this Float's value, or a nil pointer if this Float is null.
func (f Float) Ptr() *float64 {
	if !f.Valid {
		return nil
	}
	return &f.Float64
}

// IsZero returns true for invalid Floats, for future omitempty support (Go 1.4?)
// A non-null Float with a 0 value will not be considered zero.
func (f Float) IsZero() bool {
	return !f.Valid
}

func (f Float) IsDefined() bool {
	return !f.IsZero()
}

// Equal returns true if both floats have the same value or are both null.
// Warning: calculations using floating point numbers can result in different ways
// the numbers are stored in memory. Therefore, this function is not suitable to
// compare the result of a calculation. Use this method only to check if the value
// has changed in comparison to some previous value.
func (f Float) Equal(other Float) bool {
	return f.Valid == other.Valid && (!f.Valid || f.Float64 == other.Float64)
}
