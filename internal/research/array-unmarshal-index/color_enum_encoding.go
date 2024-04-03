package color

import "errors"

var ErrUnknownColor = errors.New("unknown Color")

var vals_Color = [][]byte{[]byte(""), []byte("red"), []byte("green"), []byte("blue")}

var vals_inv_Color = map[string]Color{
	"blue":  Blue,
	"green": Green,
	"red":   Red,
	"":      UndefinedColor,
}

func (s *Color) UnmarshalText(text []byte) error {
	var ok bool
	if *s, ok = vals_inv_Color[string(text)]; !ok {
		return ErrUnknownColor
	}
	return nil
}

func (s Color) MarshalText() ([]byte, error) { return []byte(vals_Color[s.c]), nil }
