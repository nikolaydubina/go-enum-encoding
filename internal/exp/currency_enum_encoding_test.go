// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleCurrency_MarshalText() {
	for _, v := range []Currency{UndefinedCurrency, SGD, USD, GBP, KRW, HKD, JPY, MYR, BHT, THC, CBD, XYZ} {
		b, _ := v.MarshalText()
		fmt.Printf("%s ", string(b))
	}
	// Output:  SGD USD GBP KRW HKD JPY MYR BHT THC CBD XYZ
}

func ExampleCurrency_UnmarshalText() {
	for _, s := range []string{"", "SGD", "USD", "GBP", "KRW", "HKD", "JPY", "MYR", "BHT", "THC", "CBD", "XYZ"} {
		var v Currency
		if err := (&v).UnmarshalText([]byte(s)); err != nil {
			fmt.Println(err)
		}
	}
}

func TestCurrency_MarshalText_UnmarshalText(t *testing.T) {
	for _, v := range []Currency{UndefinedCurrency, SGD, USD, GBP, KRW, HKD, JPY, MYR, BHT, THC, CBD, XYZ} {
		b, err := v.MarshalText()
		if err != nil {
			t.Errorf("cannot encode: %s", err)
		}

		var d Currency
		if err := (&d).UnmarshalText(b); err != nil {
			t.Errorf("cannot decode: %s", err)
		}

		if d != v {
			t.Errorf("exp(%v) != got(%v)", v, d)
		}
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `something`
		var v Currency
		err := (&v).UnmarshalText([]byte(s))
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownCurrency) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkCurrency_MarshalText(b *testing.B) {
	for b.Loop() {
		for _, c := range []Currency{UndefinedCurrency, SGD, USD, GBP, KRW, HKD, JPY, MYR, BHT, THC, CBD, XYZ} {
			if _, err := c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
}

func BenchmarkCurrency_UnmarshalText(b *testing.B) {
	var x Currency
	for b.Loop() {
		for _, c := range []string{"", "SGD", "USD", "GBP", "KRW", "HKD", "JPY", "MYR", "BHT", "THC", "CBD", "XYZ"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
