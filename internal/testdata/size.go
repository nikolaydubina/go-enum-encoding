package color

type Size uint8

//go:generate go-enum-encoding -type=Size
const (
	UndefinedSize Size = iota // json:""
	Small                     // json:"small"
	Large                     // json:"large"
	XLarge                    // json:"xlarge"
)
