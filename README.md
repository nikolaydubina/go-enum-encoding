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
* generates encoding/decoding, tests and benchmarks

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

## Related Work and References

- http://github.com/zarldev/goenums - does much more advanced struct generation, generates all enum utilities besides encoding, does not generate tests, uses similar notation to trigger go:generate but with different comment directives (non-json field tags)


<details><summary>Appendix: Performance</summary>

@mishak87 [proposed](https://github.com/nikolaydubina/go-enum-encoding/issues/19) to use array instead of map for performance.
Similarly, @nikolaydubina faced degradation in performance for loop based array for large enum sets (256 values) while working on fpmoney[^2] and iso4217[^3].

```bash
$ go test -bench=Benchmark -benchmem ./internal/research/map >  map.bench
$ go test -bench=Benchmark -benchmem ./internal/research/inline > inline.bench
$ go test -bench=Benchmark -benchmem ./internal/research/array-loop > array-loop.bench 
$ go test -bench=Benchmark -benchmem ./internal/research/array-index > array-index.bench
$ go test -bench=Benchmark -benchmem ./internal/research/uint-array > uint-array.bench
$ go test -bench=Benchmark -benchmem ./internal/research/uint-inline > uint-inline.bench
$ benchstat -split="XYZ" map.bench inline.bench array-loop.bench array-index.bench uint-array.bench uint-inline.bench
name \ time/op          map.bench    inline.bench  array-loop.bench  array-index.bench  uint-array.bench  uint-inline.bench
MarshalText_Color-16    22.3ns ± 0%    5.3ns ± 0%        7.5ns ± 0%         2.2ns ± 0%        1.9ns ± 0%         5.0ns ± 0%
UnmarshalText_Color-16  11.9ns ± 0%    5.7ns ± 0%       14.5ns ± 0%        11.8ns ± 0%       14.3ns ± 0%         5.7ns ± 0%

name \ alloc/op         map.bench    inline.bench  array-loop.bench  array-index.bench  uint-array.bench  uint-inline.bench
MarshalText_Color-16     0.00B         0.00B             0.00B              0.00B             0.00B              0.00B     
UnmarshalText_Color-16   0.00B         0.00B             0.00B              0.00B             0.00B              0.00B     

name \ allocs/op        map.bench    inline.bench  array-loop.bench  array-index.bench  uint-array.bench  uint-inline.bench
MarshalText_Color-16      0.00          0.00              0.00               0.00              0.00               0.00     
UnmarshalText_Color-16    0.00          0.00              0.00               0.00              0.00               0.00 
```

</details>

[^1]: Comparison to other enums methods: http://github.com/nikolaydubina/go-enum-example
[^2]: iso4217 enums performance loop vs map: https://github.com/ferdypruis/iso4217/issues/4
[^3]: fpmoney: https://github.com/nikolaydubina/fpmoney?tab=readme-ov-file#appendix-a-jsonunmarshal-optimizations
