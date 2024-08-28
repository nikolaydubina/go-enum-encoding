// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"errors"
	"fmt"
	"testing"
)

func ExampleColorString_MarshalText() {
	for _, v := range []ColorString{RedS, GreenS, BlueS} {
		b, _ := v.MarshalText()
		fmt.Printf("%s ", string(b))
	}
	// Output: red green blue
}

func ExampleColorString_UnmarshalText() {
	for _, s := range []string{"red", "green", "blue"} {
		var v ColorString
		if err := (&v).UnmarshalText([]byte(s)); err != nil {
			fmt.Println(err)
		}
	}
}

func TestColorString_MarshalText_UnmarshalText(t *testing.T) {
	for _, v := range []ColorString{RedS, GreenS, BlueS} {
		b, err := v.MarshalText()
		if err != nil {
			t.Errorf("cannot encode: %s", err)
		}

		var d ColorString
		if err := (&d).UnmarshalText(b); err != nil {
			t.Errorf("cannot decode: %s", err)
		}

		if d != v {
			t.Errorf("exp(%v) != got(%v)", v, d)
		}
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `something`
		var v ColorString
		err := (&v).UnmarshalText([]byte(s))
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownColorString) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkColorString_MarshalText(b *testing.B) {
	var v []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, c := range []ColorString{RedS, GreenS, BlueS} {
			if v, err = c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
	if len(v) > 1000 {
		b.Fatal("noop")
	}
}

func BenchmarkColorString_UnmarshalText(b *testing.B) {
	var x ColorString
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"red", "green", "blue"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}