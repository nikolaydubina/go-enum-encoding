# go-enum-encoding

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-enum-encoding)](https://goreportcard.com/report/github.com/nikolaydubina/go-enum-encoding)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/go-enum-encoding.svg)](https://pkg.go.dev/github.com/nikolaydubina/go-enum-encoding)
[![codecov](https://codecov.io/gh/nikolaydubina/go-enum-encoding/graph/badge.svg?token=asZfIddrLV)](https://codecov.io/gh/nikolaydubina/go-enum-encoding)
[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/go-enum-encoding/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/go-enum-encoding)

```bash
go install github.com/nikolaydubina/go-enum-encoding@latest
```

* 100 LOC
* simple, fast[^1], strict[^1]
* generates tests

given
```go
type Color struct{ c uint8 }

//go:generate go-enum-encoding -type=Color
var (
	UndefinedColor = Color{}            // json:""
	Red            = Color{1}           // json:"red"
	Green          = Color{2}           // json:"green"
	Blue           = Color{3}           // json:"blue"
)
```

<details><summary>generates encoding and decoding routines</summary>
	
```go
// Code generated by go-enum-encoding DO NOT EDIT
package main

import "errors"

var ErrUnknownColor = errors.New("unknown Color")

var vals_Color = map[Color]string{
	Blue:           "blue",
	Green:          "green",
	Red:            "red",
	UndefinedColor: "",
}

var vals_inv_Color = map[string]Color{
	"blue":  Blue,
	"green": Green,
	"red":   Red,
	"":      UndefinedColor,
}

func (s *Color) UnmarshalText(text []byte) error {
	var ok bool
	if *s, ok = vals_inv_Color[string(text)]; !ok {
		return ErrUnknownColor
	}
	return nil
}

func (s Color) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

func (s Color) String() string { return vals_Color[s] }
```

</details>

<details><summary>... and geneartes tests</summary>
	
```go
// Code generated by go-enum-encoding DO NOT EDIT
package main

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestJSON_Color(t *testing.T) {
	type V struct {
		Values []Color `json:"values"`
	}

	values := []Color{Blue, Green, Red, UndefinedColor}

	var v V
	s := `{"values":["blue","green","red",""]}`
	json.Unmarshal([]byte(s), &v)

	if len(v.Values) != len(values) {
		t.Errorf("cannot decode: %d", len(v.Values))
	}
	if !slices.Equal(v.Values, values) {
		t.Errorf("wrong decoded: %v", v.Values)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("cannot encode: %s", err)
	}
	if string(b) != s {
		t.Errorf("wrong encoded: %s != %s", string(b), s)
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `{"values":["something"]}`
		var v V
		err := json.Unmarshal([]byte(s), &v)
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownColor) {
			t.Errorf("wrong error: %s", err)
		}
	})
}
```

</details>

## Related Work and References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, uses similar notation to trigger go:generate but with different comment directivs (non-json field tags)

[^1]: Comparison to other enums methods: http://github.com/nikolaydubina/go-enum-example

## Appendix: Decoding from Array

@mishak87 [proposed](https://github.com/nikolaydubina/go-enum-encoding/issues/19) to use array instead of map for during decoding phase to improve performance.
However, using arrays will require user enums to be contigious, low-number, starting from zero values.
This will also require to read numeric value from enum var/const definition in AST, that increases code complexity.

```bash
$ go test -bench=Benchmark -benchmem ./internal/research/array-unmarshal-index > array-unmarshal-index.bench 
$ go test -bench=Benchmark -benchmem ./internal/research/array-unmarshal-index-string > array-unmarshal-index-string.bench 
$ go test -bench=Benchmark -benchmem ./internal/research/map-unmarshal >  map-unmarshal.bench
$ benchstat -split="XYZ" map-unmarshal.bench array-unmarshal-index.bench  array-unmarshal-index-string.bench
name \ time/op          map-unmarshal.bench  array-unmarshal-index.bench  array-unmarshal-index-string.bench
MarshalText_Color-16            11.0ns ± 0%                   2.2ns ± 0%                          2.3ns ± 0%
UnmarshalText_Color-16          12.0ns ± 0%                  12.0ns ± 0%                         11.9ns ± 0%

name \ alloc/op         map-unmarshal.bench  array-unmarshal-index.bench  array-unmarshal-index-string.bench
MarshalText_Color-16             0.00B                        0.00B                               0.00B     
UnmarshalText_Color-16           0.00B                        0.00B                               0.00B     

name \ allocs/op        map-unmarshal.bench  array-unmarshal-index.bench  array-unmarshal-index-string.bench
MarshalText_Color-16              0.00                         0.00                                0.00     
UnmarshalText_Color-16            0.00                         0.00                                0.00   
```
