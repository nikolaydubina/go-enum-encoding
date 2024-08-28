// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleCurrencyString_MarshalText() {
	for _, v := range []CurrencyString{UndefinedCurrencyS, SGDS, USDS, GBPS, KRWS, HKDS, JPYS, MYRS, BHTS, THCS, CBDS, XYZS} {
		b, _ := v.MarshalText()
		fmt.Printf("%s ", string(b))
	}
	// Output:  SGD USD GBP KRW HKD JPY MYR BHT THC CBD XYZ
}

func ExampleCurrencyString_UnmarshalText() {
	for _, s := range []string{"", "SGD", "USD", "GBP", "KRW", "HKD", "JPY", "MYR", "BHT", "THC", "CBD", "XYZ"} {
		var v CurrencyString
		if err := (&v).UnmarshalText([]byte(s)); err != nil {
			fmt.Println(err)
		}
	}
}

func TestCurrencyString_MarshalText_UnmarshalText(t *testing.T) {
	for _, v := range []CurrencyString{UndefinedCurrencyS, SGDS, USDS, GBPS, KRWS, HKDS, JPYS, MYRS, BHTS, THCS, CBDS, XYZS} {
		b, err := v.MarshalText()
		if err != nil {
			t.Errorf("cannot encode: %s", err)
		}

		var d CurrencyString
		if err := (&d).UnmarshalText(b); err != nil {
			t.Errorf("cannot decode: %s", err)
		}

		if d != v {
			t.Errorf("exp(%v) != got(%v)", v, d)
		}
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `something`
		var v CurrencyString
		err := (&v).UnmarshalText([]byte(s))
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownCurrencyString) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkCurrencyString_MarshalText(b *testing.B) {
	var v []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, c := range []CurrencyString{UndefinedCurrencyS, SGDS, USDS, GBPS, KRWS, HKDS, JPYS, MYRS, BHTS, THCS, CBDS, XYZS} {
			if v, err = c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
	if len(v) > 1000 {
		b.Fatal("noop")
	}
}

func BenchmarkCurrencyString_UnmarshalText(b *testing.B) {
	var x CurrencyString
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"", "SGD", "USD", "GBP", "KRW", "HKD", "JPY", "MYR", "BHT", "THC", "CBD", "XYZ"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}