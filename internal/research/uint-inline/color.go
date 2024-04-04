package color

type Color uint8

const (
	UndefinedColor Color = iota
	Red
	Green
	Blue
)

type V struct {
	Color Color `json:"color"`
}
