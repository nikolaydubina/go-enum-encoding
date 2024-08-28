package color

type Color2 struct{ c uint8 }

//go:generate go-enum-encoding -type=Color2
var (
	UndefinedColor2  = Color2{}             // json:""
	Red2             = Color2{1}            // json:"red"
	Purple2, Orange2 = Color2{2}, Color2{3} // json:"blue"
)

type Currency2 uint8

//go:generate go-enum-encoding -type=Currency2
const (
	UndefCurrency2 Currency2 = iota // json:""
	SGD2                            // json:"SGD"
	USD2                            // json:"USD"
)
