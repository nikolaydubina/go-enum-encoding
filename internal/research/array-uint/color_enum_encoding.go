package color

import "errors"

var ErrUnknownColor = errors.New("unknown Color")

var (
	json_Color = [...]string{"", "red", "green", "blue"}
	vals_Color = [...]Color{UndefinedColor, Red, Green, Blue}
)

func (s *Color) UnmarshalText(text []byte) error {
	for i, v := range json_Color {
		if v == string(text) {
			*s = vals_Color[i]
			return nil
		}
	}
	return ErrUnknownColor
}

func (s Color) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

func (s Color) String() string { return json_Color[s] }
