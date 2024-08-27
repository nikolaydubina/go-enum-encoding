// Code generated by go-enum-encoding DO NOT EDIT
package color

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestJSON_ColorString(t *testing.T) {
	type V struct {
		Values []ColorString `json:"values"`
	}

	values := []ColorString{RedS, GreenS, BlueS}

	var v V
	s := `{"values":["red","green","blue"]}`
	json.Unmarshal([]byte(s), &v)

	if len(v.Values) != len(values) {
		t.Errorf("cannot decode: %d", len(v.Values))
	}
	if !slices.Equal(v.Values, values) {
		t.Errorf("wrong decoded: %v", v.Values)
	}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("cannot encode: %s", err)
	}
	if string(b) != s {
		t.Errorf("wrong encoded: %s != %s", string(b), s)
	}

	t.Run("when unknown value, then error", func(t *testing.T) {
		s := `{"values":["something"]}`
		var v V
		err := json.Unmarshal([]byte(s), &v)
		if err == nil {
			t.Errorf("must be error")
		}
		if !errors.Is(err, ErrUnknownColorString) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkMarshalText_ColorString(b *testing.B) {
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

func BenchmarkUnmarshalText_ColorString(b *testing.B) {
	var x ColorString
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"red", "green", "blue"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}

func TestColorString_String(t *testing.T) {
	values := []ColorString{RedS, GreenS, BlueS}
	tags := []string{"red", "green", "blue"}

	for i := range values {
		if values[i].String() != tags[i] {
			t.Errorf("got(%s) != exp(%s)", values[i].String(), tags[i])
		}
	}
}