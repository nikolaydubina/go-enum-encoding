// Code generated by go-enum-encoding DO NOT EDIT
package color

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestJSON_Color2(t *testing.T) {
	type V struct {
		Values []Color2 `json:"values"`
	}

	values := []Color2{UndefinedColor2, Red2}

	var v V
	s := `{"values":["","red"]}`
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
		if !errors.Is(err, ErrUnknownColor2) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkMarshalText_Color2(b *testing.B) {
	var v []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, c := range []Color2{UndefinedColor2, Red2} {
			if v, err = c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
	if len(v) > 1000 {
		b.Fatal("noop")
	}
}

func BenchmarkUnmarshalText_Color2(b *testing.B) {
	var x Color2
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"", "red"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
