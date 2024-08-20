package color

type Color2 struct{ c uint8 }

//go:generate go-enum-encoding -type=Color2
var (
	UndefinedColor = Color2{}             // json:""
	Red            = Color2{1}            // json:"red"
	Purple, Orange = Color2{2}, Color2{3} // json:"blue"
)

type Currency2 uint8

//go:generate go-enum-encoding -type=Currency2
const (
	UndefCurrency Currency2 = iota // json:""
	SGD                            // json:"SGD"
	USD                            // json:"USD"
)
