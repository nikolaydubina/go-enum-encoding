# go-enum-encoding

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-enum-encoding)](https://goreportcard.com/report/github.com/nikolaydubina/go-enum-encoding)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/go-enum-encoding.svg)](https://pkg.go.dev/github.com/nikolaydubina/go-enum-encoding)
[![codecov](https://codecov.io/gh/nikolaydubina/go-enum-encoding/graph/badge.svg?token=asZfIddrLV)](https://codecov.io/gh/nikolaydubina/go-enum-encoding)
[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/go-enum-encoding/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/go-enum-encoding)

```bash
go install github.com/nikolaydubina/go-enum-encoding@latest
```

> [!WARNING]
> `v1.8.2` is the last version that supports go1.23 or bellow

* 200 LOC
* simple, fast, strict
* generates encoding/decoding, tests, benchmarks

```go
type Color struct{ c uint8 }

//go:generate go-enum-encoding -type=Color
var (
	UndefinedColor = Color{} 	// json:""
	Red            = Color{1}	// json:"red"
	Green          = Color{2}	// json:"green"
	Blue           = Color{3}	// json:"blue"
)
```

It also works with raw `iota` enums:

```go
type Size uint8

//go:generate go-enum-encoding -type=Size
const (
	UndefinedSize Size = iota // json:""
	Small                     // json:"small"
	Large                     // json:"large"
	XLarge                    // json:"xlarge"
)
```

## References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, uses similar notation to trigger go:generate but with different comment directives (non-json field tags)
- http://github.com/nikolaydubina/go-enum-example
