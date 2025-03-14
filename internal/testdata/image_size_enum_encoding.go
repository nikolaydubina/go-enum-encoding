// Code generated by go-enum-encoding; DO NOT EDIT.

package image

import "errors"

var ErrUnknownImageSize = errors.New("unknown ImageSize")

func (s *ImageSize) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*s = UndefinedSize
	case "small":
		*s = Small
	case "large":
		*s = Large
	case "xlarge":
		*s = XLarge
	default:
		return ErrUnknownImageSize
	}
	return nil
}

var seq_bytes_ImageSize = [...][]byte{[]byte(""), []byte("small"), []byte("large"), []byte("xlarge")}

func (s ImageSize) MarshalText() ([]byte, error) {
	switch s {
	case UndefinedSize:
		return seq_bytes_ImageSize[0], nil
	case Small:
		return seq_bytes_ImageSize[1], nil
	case Large:
		return seq_bytes_ImageSize[2], nil
	case XLarge:
		return seq_bytes_ImageSize[3], nil
	default:
		return nil, ErrUnknownImageSize
	}
}

var seq_string_ImageSize = [...]string{"", "small", "large", "xlarge"}

func (s ImageSize) String() string {
	switch s {
	case UndefinedSize:
		return seq_string_ImageSize[0]
	case Small:
		return seq_string_ImageSize[1]
	case Large:
		return seq_string_ImageSize[2]
	case XLarge:
		return seq_string_ImageSize[3]
	default:
		return ""
	}
}
