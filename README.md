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

Generated benchmarks:

```bash
$ go test -bench=. -benchmem .
goos: darwin
goarch: arm64
pkg: test
cpu: Apple M3 Max
BenchmarkColor_UnmarshalText-16         804431073                1.366 ns/op           0 B/op          0 allocs/op
BenchmarkColor_AppendText-16            81371102                14.01 ns/op           24 B/op          1 allocs/op
BenchmarkColor_MarshalText-16           84193539                13.77 ns/op            8 B/op          1 allocs/op
BenchmarkImageSize_UnmarshalText-16     900864548                1.345 ns/op           0 B/op          0 allocs/op
BenchmarkImageSize_AppendText-16        82080981                14.07 ns/op           24 B/op          1 allocs/op
BenchmarkImageSize_MarshalText-16       498537706                2.429 ns/op           0 B/op          0 allocs/op
BenchmarkImageSize_String-16            1000000000               1.076 ns/op           0 B/op          0 allocs/op
PASS
ok      test    8.221s
```

## References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, uses similar notation to trigger go:generate but with different comment directives (non-json field tags)
- http://github.com/nikolaydubina/go-enum-example
