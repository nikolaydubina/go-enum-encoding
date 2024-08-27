// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"errors"
	"testing"
)

func TestColor_MarshalText_UnmarshalText(t *testing.T) {
	for _, v := range []Color{UndefinedColor, Red, Green, Blue} {
		b, err := v.MarshalText()
		if err != nil {
			t.Errorf("cannot encode: %s", err)
		}

		var d Color
		if err := (&d).UnmarshalText(b); err != nil {
			t.Errorf("cannot decode: %s", err)
		}

		if d != v {
			t.Errorf("exp(%v) != got(%v)", v, d)
		}
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `something`
		var v Color
		err := (&v).UnmarshalText([]byte(s))
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownColor) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkColor_MarshalText(b *testing.B) {
	var v []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, c := range []Color{UndefinedColor, Red, Green, Blue} {
			if v, err = c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
	if len(v) > 1000 {
		b.Fatal("noop")
	}
}

func BenchmarkColor_UnmarshalText(b *testing.B) {
	var x Color
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"", "red", "green", "blue"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
