package color

type Color struct{ c uint8 }

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

type ColorString int

const (
	RedS   ColorString = iota // json:"red"
	GreenS                    // json:"green"
	BlueS                     // json:"blue"
)