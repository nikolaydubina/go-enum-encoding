// Code generated by go-enum-encoding DO NOT EDIT
package color

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestJSON_Color(t *testing.T) {
	type V struct {
		Values []Color `json:"values"`
	}

	values := []Color{UndefinedColor, Red, Green, Blue}

	var v V
	s := `{"values":["","red","green","blue"]}`
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
		if !errors.Is(err, ErrUnknownColor) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkMarshalText_Color(b *testing.B) {
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

func BenchmarkUnmarshalText_Color(b *testing.B) {
	var x Color
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"", "red", "green", "blue"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
