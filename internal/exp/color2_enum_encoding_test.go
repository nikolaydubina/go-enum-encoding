// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleColor2_MarshalText() {
	for _, v := range []Color2{UndefinedColor2, Red2} {
		b, _ := v.MarshalText()
		fmt.Printf("%s ", string(b))
	}
	// Output:  red
}

func ExampleColor2_UnmarshalText() {
	for _, s := range []string{"", "red"} {
		var v Color2
		if err := (&v).UnmarshalText([]byte(s)); err != nil {
			fmt.Println(err)
		}
	}
}

func TestColor2_MarshalText_UnmarshalText(t *testing.T) {
	for _, v := range []Color2{UndefinedColor2, Red2} {
		b, err := v.MarshalText()
		if err != nil {
			t.Errorf("cannot encode: %s", err)
		}

		var d Color2
		if err := (&d).UnmarshalText(b); err != nil {
			t.Errorf("cannot decode: %s", err)
		}

		if d != v {
			t.Errorf("exp(%v) != got(%v)", v, d)
		}
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `something`
		var v Color2
		err := (&v).UnmarshalText([]byte(s))
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownColor2) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkColor2_MarshalText(b *testing.B) {
	for b.Loop() {
		for _, c := range []Color2{UndefinedColor2, Red2} {
			if _, err := c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
}

func BenchmarkColor2_UnmarshalText(b *testing.B) {
	var x Color2
	for b.Loop() {
		for _, c := range []string{"", "red"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
