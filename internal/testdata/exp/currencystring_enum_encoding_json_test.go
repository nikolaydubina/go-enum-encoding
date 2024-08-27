// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
)

func TestCurrencyString_JSON(t *testing.T) {
	type V struct {
		Values []CurrencyString `json:"values"`
	}

	values := []CurrencyString{UndefinedCurrencyS, SGDS, USDS, GBPS, KRWS, HKDS, JPYS, MYRS, BHTS, THCS, CBDS, XYZS}

	var v V
	s := `{"values":["","SGD","USD","GBP","KRW","HKD","JPY","MYR","BHT","THC","CBD","XYZ"]}`
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
		if !errors.Is(err, ErrUnknownCurrencyString) {
			t.Errorf("wrong error: %s", err)
		}
	})
}
