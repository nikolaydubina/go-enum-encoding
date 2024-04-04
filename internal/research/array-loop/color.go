package color

type Color struct{ c uint8 }

//go:generate go-enum-encoding -type=Color
var (
	UndefinedColor = Color{}            // json:""
	Red            = Color{1}           // json:"red"
	Green          = Color{2}           // json:"green"
	Blue           = Color{3}           // json:"blue"
	Purple, Orange = Color{4}, Color{5} // json:"blue"
)

type V struct {
	Color Color `json:"color"`
}
