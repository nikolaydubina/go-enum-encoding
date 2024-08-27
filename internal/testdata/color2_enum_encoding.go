// Code generated by go-enum-encoding; DO NOT EDIT.

package color

import "errors"

var ErrUnknownColor2 = errors.New("unknown Color2")

func (s *Color2) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*s = UndefinedColor2
	case "red":
		*s = Red2
	default:
		return ErrUnknownColor2
	}
	return nil
}

var seq_bytes_Color2 = [...][]byte{[]byte(""), []byte("red")}

func (s Color2) MarshalText() ([]byte, error) {
	switch s {
	case UndefinedColor2:
		return seq_bytes_Color2[0], nil
	case Red2:
		return seq_bytes_Color2[1], nil
	default:
		return nil, ErrUnknownColor2
	}
}
