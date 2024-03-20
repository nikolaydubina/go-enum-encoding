# go-enum-encoding

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-enum-encoding)](https://goreportcard.com/report/github.com/nikolaydubina/go-enum-encoding)
[![codecov](https://codecov.io/gh/nikolaydubina/go-enum-encoding/graph/badge.svg?token=asZfIddrLV)](https://codecov.io/gh/nikolaydubina/go-enum-encoding)
[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)

```bash
go install github.com/nikolaydubina/go-enum-encoding@latest
```

* 150 LOC
* simple, fast[^1], and strict[^1]
* no dependencies, no template
* works with different enum implementations (`iota`, `struct`, `string`)
* generates tests

given
```go

type Color struct{ c uint }

//go:generate go-enum-encoding -type=Color
var (
	Undefined = Color{}  // json:"-"
	Red       = Color{1} // json:"red"
	Green     = Color{2} // json:"green"
	Blue      = Color{3} // json:"blue"
)
```

generates encoding and decoding routines and tests 
```go
package main

import "errors"

var ErrUnknownColor = errors.New("unknown color")

var colors = map[Color]string{
	Red:   "red",
	Green: "green",
	Blue:  "blue",
}

var colors_inv = map[string]Color{
	"red":   Red,
	"green": Green,
	"blue":  Blue,
}

func (s *Color) UnmarshalText(text []byte) error {
	*s = colors_inv[string(text)]
	if *s == Undefined {
		return ErrUnknownColor
	}
	return nil
}

func (c Color) MarshalText() ([]byte, error) { return []byte(c.String()), nil }

func (c Color) String() string { return colors[c] }
```

## Related Work and References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, has slightly different notation for tests

[^1]: Comparison to other enums methods: http://github.com/nikolaydubina/go-enum-example
