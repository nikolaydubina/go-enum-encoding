package color

type NoUndefined struct{ c uint }

//go:generate go-enum-encoding -type=NoUndefined
var (
	RedNoUndefined  = NoUndefined{1} // json:"red"
	BlueNoUndefined = NoUndefined{2} // json:"blue"
)
