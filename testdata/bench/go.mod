module github.com/unravelin/null/testdata/bench

go 1.27

replace github.com/unravelin/null => ../../

require (
	github.com/google/gofuzz v1.2.0
	github.com/philpearl/plenc v0.0.24
	github.com/unravelin/null v2.1.4+incompatible
)
