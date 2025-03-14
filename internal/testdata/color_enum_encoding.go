// Code generated by go-enum-encoding; DO NOT EDIT.

package image

import "errors"

var ErrUnknownColor = errors.New("unknown Color")

func (s *Color) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*s = UndefinedColor
	case "red":
		*s = Red
	case "green":
		*s = Green
	case "blue":
		*s = Blue
	default:
		return ErrUnknownColor
	}
	return nil
}

var seq_bytes_Color = [...][]byte{[]byte(""), []byte("red"), []byte("green"), []byte("blue")}

func (s Color) MarshalText() ([]byte, error) { return s.AppendText(nil) }

func (s Color) AppendText(b []byte) ([]byte, error) {
	switch s {
	case UndefinedColor:
		return append(b, seq_bytes_Color[0]...), nil
	case Red:
		return append(b, seq_bytes_Color[1]...), nil
	case Green:
		return append(b, seq_bytes_Color[2]...), nil
	case Blue:
		return append(b, seq_bytes_Color[3]...), nil
	default:
		return nil, ErrUnknownColor
	}
}
