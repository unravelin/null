package bench

import (
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/philpearl/plenc"
	plnull "github.com/philpearl/plenc/null"
	"github.com/unravelin/null"
)

var fuzzFuncs = []any{
	func(a *null.Bool, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			a.Bool = c.RandBool()
		}
	},
	func(a *null.Float, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Float64)
		}
	},
	func(a *null.Int, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Int64)
		}
	},
	func(a *null.String, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.String)
		}
	},
	func(a *null.Time, c fuzz.Continue) {
		a.Valid = c.RandBool()
		if a.Valid {
			c.Fuzz(&a.Time)
		}
	},
}

func BenchmarkSerialisation(b *testing.B) {
	f := fuzz.New().Funcs(fuzzFuncs...)

	var in pltest
	f.Fuzz(&in)

	b.Run("plenc", func(b *testing.B) {
		plnull.RegisterCodecs()
		b.ReportAllocs()

		var data []byte
		for b.Loop() {
			var err error
			data, err = plenc.Marshal(data[:0], &in)
			if err != nil {
				b.Fatal(err)
			}
			var out pltest
			if err := plenc.Unmarshal(data, &out); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("json", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			data, err := json.Marshal(&in)
			if err != nil {
				b.Fatal(err)
			}

			var out pltest
			if err := json.Unmarshal(data, &out); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("jsonv2", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			data, err := jsonv2.Marshal(&in)
			if err != nil {
				b.Fatal(err)
			}
			var out pltest
			if err := jsonv2.Unmarshal(data, &out); err != nil {
				b.Fatal(err)
			}
		}
	})
}
