package color

import "errors"

var ErrUnknownColor = errors.New("unknown Color")

func (s *Color) UnmarshalText(text []byte) error {
	switch string(text) {
	case "red":
		*s = Red
	case "green":
		*s = Green
	case "blue":
		*s = Blue
	case "":
		*s = UndefinedColor
	default:
		return ErrUnknownColor
	}
	return nil
}

func (s Color) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

func (s Color) String() string {
	switch s {
	case UndefinedColor:
		return ""
	case Red:
		return "red"
	case Green:
		return "green"
	case Blue:
		return "blue"
	default:
		return ""
	}
}
