package zero

import (
	"encoding/json"
	"errors"
	"testing"
)

var (
	boolJSON     = []byte(`true`)
	falseJSON    = []byte(`false`)
	nullBoolJSON = []byte(`{"Bool":true,"Valid":true}`)
)

func TestBoolFrom(t *testing.T) {
	b := BoolFrom(true)
	assertBool(t, b, "BoolFrom()")

	zero := BoolFrom(false)
	if zero.Valid {
		t.Error("BoolFrom(false)", "is valid, but should be invalid")
	}
}

func TestBoolFromPtr(t *testing.T) {
	v := true
	bptr := &v
	b := BoolFromPtr(bptr)
	assertBool(t, b, "BoolFromPtr()")

	null := BoolFromPtr(nil)
	assertNullBool(t, null, "BoolFromPtr(nil)")
}

func TestUnmarshalBool(t *testing.T) {
	var b Bool
	err := json.Unmarshal(boolJSON, &b)
	maybePanic(err)
	assertBool(t, b, "float json")

	var nb Bool
	err = json.Unmarshal(nullBoolJSON, &nb)
	if err == nil {
		panic("expected error")
	}

	var zero Bool
	err = json.Unmarshal(falseJSON, &zero)
	maybePanic(err)
	assertNullBool(t, zero, "zero json")

	var null Bool
	err = json.Unmarshal(nullJSON, &null)
	maybePanic(err)
	assertNullBool(t, null, "null json")

	var invalid Bool
	err = invalid.UnmarshalJSON(invalidJSON)
	var syntaxError *json.SyntaxError
	if !errors.As(err, &syntaxError) {
		t.Errorf("expected wrapped json.SyntaxError, not %T", err)
	}
	assertNullBool(t, invalid, "invalid json")

	var badType Bool
	err = json.Unmarshal(intJSON, &badType)
	if err == nil {
		panic("err should not be nil")
	}
	assertNullBool(t, badType, "wrong type json")
}

func TestTextUnmarshalBool(t *testing.T) {
	var b Bool
	err := b.UnmarshalText(boolJSON)
	maybePanic(err)
	assertBool(t, b, "UnmarshalText() bool")

	var zero Bool
	err = zero.UnmarshalText(falseJSON)
	maybePanic(err)
	assertNullBool(t, zero, "UnmarshalText() zero bool")

	var blank Bool
	err = blank.UnmarshalText([]byte(""))
	maybePanic(err)
	assertNullBool(t, blank, "UnmarshalText() empty bool")

	var null Bool
	err = null.UnmarshalText(nullJSON)
	maybePanic(err)
	assertNullBool(t, null, `UnmarshalText() "null"`)

	var invalid Bool
	err = invalid.UnmarshalText(invalidJSON)
	if err == nil {
		panic("err should not be nil")
	}
}

func TestMarshalBool(t *testing.T) {
	b := BoolFrom(true)
	data, err := json.Marshal(b)
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty json marshal")

	// invalid values should be encoded as false
	null := NewBool(false, false)
	data, err = json.Marshal(null)
	maybePanic(err)
	assertJSONEquals(t, data, "false", "null json marshal")
}

func TestMarshalBoolText(t *testing.T) {
	b := BoolFrom(true)
	data, err := b.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "true", "non-empty text marshal")

	// invalid values should be encoded as zero
	null := NewBool(false, false)
	data, err = null.MarshalText()
	maybePanic(err)
	assertJSONEquals(t, data, "false", "null text marshal")
}

func TestBoolPointer(t *testing.T) {
	b := BoolFrom(true)
	ptr := b.Ptr()
	if *ptr != true {
		t.Errorf("bad %s bool: %#v ≠ %v\n", "pointer", ptr, true)
	}

	null := NewBool(false, false)
	ptr = null.Ptr()
	if ptr != nil {
		t.Errorf("bad %s bool: %#v ≠ %s\n", "nil pointer", ptr, "nil")
	}
}

func TestBoolIsZero(t *testing.T) {
	b := BoolFrom(true)
	if b.IsZero() {
		t.Errorf("IsZero() should be false")
	}

	null := NewBool(false, false)
	if !null.IsZero() {
		t.Errorf("IsZero() should be true")
	}

	zero := NewBool(false, true)
	if !zero.IsZero() {
		t.Errorf("IsZero() should be true")
	}
}

func TestBoolSetValid(t *testing.T) {
	change := NewBool(false, false)
	assertNullBool(t, change, "SetValid()")
	change.SetValid(true)
	assertBool(t, change, "SetValid()")
}

func TestBoolScan(t *testing.T) {
	var b Bool
	err := b.Scan(true)
	maybePanic(err)
	assertBool(t, b, "scanned bool")

	var null Bool
	err = null.Scan(nil)
	maybePanic(err)
	assertNullBool(t, null, "scanned null")
}

func TestBoolValueOrZero(t *testing.T) {
	valid := NewBool(true, true)
	if valid.ValueOrZero() != true {
		t.Error("unexpected ValueOrZero", valid.ValueOrZero())
	}

	invalid := NewBool(true, false)
	if invalid.ValueOrZero() != false {
		t.Error("unexpected ValueOrZero", invalid.ValueOrZero())
	}
}

func TestBoolEqual(t *testing.T) {
	b1 := NewBool(true, false)
	b2 := NewBool(true, false)
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = NewBool(true, false)
	b2 = NewBool(false, false)
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = NewBool(true, true)
	b2 = NewBool(true, true)
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = NewBool(true, false)
	b2 = NewBool(false, true)
	assertBoolEqualIsTrue(t, b1, b2)

	b1 = NewBool(true, true)
	b2 = NewBool(true, false)
	assertBoolEqualIsFalse(t, b1, b2)

	b1 = NewBool(true, false)
	b2 = NewBool(true, true)
	assertBoolEqualIsFalse(t, b1, b2)

	b1 = NewBool(true, true)
	b2 = NewBool(false, true)
	assertBoolEqualIsFalse(t, b1, b2)
}

func assertBool(t *testing.T, b Bool, from string) {
	t.Helper()
	if b.Bool != true {
		t.Errorf("bad %s bool: %v ≠ %v\n", from, b.Bool, true)
	}
	if !b.Valid {
		t.Error(from, "is invalid, but should be valid")
	}
}

func assertNullBool(t *testing.T, b Bool, from string) {
	if b.Valid {
		t.Error(from, "is valid, but should be invalid")
	}
}

func assertBoolEqualIsTrue(t *testing.T, a, b Bool) {
	t.Helper()
	if !a.Equal(b) {
		t.Errorf("Equal() of Bool{%t, Valid:%t} and Bool{%t, Valid:%t} should return true", a.Bool, a.Valid, b.Bool, b.Valid)
	}
}

func assertBoolEqualIsFalse(t *testing.T, a, b Bool) {
	t.Helper()
	if a.Equal(b) {
		t.Errorf("Equal() of Bool{%t, Valid:%t} and Bool{%t, Valid:%t} should return false", a.Bool, a.Valid, b.Bool, b.Valid)
	}
}
