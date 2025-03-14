# go-enum-encoding

[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/go-enum-encoding)](https://goreportcard.com/report/github.com/nikolaydubina/go-enum-encoding)
[![Go Reference](https://pkg.go.dev/badge/github.com/nikolaydubina/go-enum-encoding.svg)](https://pkg.go.dev/github.com/nikolaydubina/go-enum-encoding)
[![codecov](https://codecov.io/gh/nikolaydubina/go-enum-encoding/graph/badge.svg?token=asZfIddrLV)](https://codecov.io/gh/nikolaydubina/go-enum-encoding)
[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/go-enum-encoding/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/go-enum-encoding)

```bash
go install github.com/nikolaydubina/go-enum-encoding@latest
```

* 200 LOC
* simple, fast, strict
* generate encoding/decoding, tests, benchmarks

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

`iota` is ok too

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

generated benchmarks

```bash
$ go test -bench=. -benchmem .                  
goos: darwin
goarch: arm64
pkg: github.com/nikolaydubina/go-enum-encoding/internal/testdata
cpu: Apple M3 Max
BenchmarkColor_UnmarshalText-16         752573839                1.374 ns/op           0 B/op          0 allocs/op
BenchmarkColor_AppendText-16            450123993                2.676 ns/op           0 B/op          0 allocs/op
BenchmarkColor_MarshalText-16           80059376                13.68 ns/op            8 B/op          1 allocs/op
BenchmarkImageSize_UnmarshalText-16     751743885                1.601 ns/op           0 B/op          0 allocs/op
BenchmarkImageSize_AppendText-16        500286883                2.402 ns/op           0 B/op          0 allocs/op
BenchmarkImageSize_MarshalText-16       81467318                16.46 ns/op            8 B/op          1 allocs/op
BenchmarkImageSize_String-16            856463289                1.330 ns/op           0 B/op          0 allocs/op
PASS
ok      github.com/nikolaydubina/go-enum-encoding/internal/testdata     8.561s
```

## References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, uses similar notation to trigger go:generate but with different comment directives (non-json field tags)
- http://github.com/nikolaydubina/go-enum-example
