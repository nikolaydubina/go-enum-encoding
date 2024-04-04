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

	values := []Color{Blue, Green, Red, UndefinedColor}

	var v V
	s := `{"values":["blue","green","red",""]}`
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
	for i := 0; i < b.N; i++ {
		for _, c := range []Color{Blue, Green, Red} {
			if v, err := c.MarshalText(); err != nil || len(v) == 0 {
				b.Fatal("empty")
			}
		}
	}
}

func BenchmarkUnmarshalText_Color(b *testing.B) {
	var x Color
	for i := 0; i < b.N; i++ {
		for _, c := range [][]byte{[]byte("blue"), []byte("green"), []byte("red")} {
			if err := x.UnmarshalText(c); err != nil {
				b.Fatal("cannot decode")
			}
		}
	}
}
