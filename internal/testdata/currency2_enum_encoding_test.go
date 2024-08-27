// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestJSON_Currency2(t *testing.T) {
	type V struct {
		Values []Currency2 `json:"values"`
	}

	values := []Currency2{UndefCurrency2, SGD2, USD2}

	var v V
	s := `{"values":["","SGD","USD"]}`
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
		if !errors.Is(err, ErrUnknownCurrency2) {
			t.Errorf("wrong error: %s", err)
		}
	})
}

func BenchmarkMarshalText_Currency2(b *testing.B) {
	var v []byte
	var err error
	for i := 0; i < b.N; i++ {
		for _, c := range []Currency2{UndefCurrency2, SGD2, USD2} {
			if v, err = c.MarshalText(); err != nil {
				b.Fatal("empty")
			}
		}
	}
	if len(v) > 1000 {
		b.Fatal("noop")
	}
}

func BenchmarkUnmarshalText_Currency2(b *testing.B) {
	var x Currency2
	for i := 0; i < b.N; i++ {
		for _, c := range []string{"", "SGD", "USD"} {
			if err := x.UnmarshalText([]byte(c)); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
