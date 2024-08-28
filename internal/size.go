package color

type Size uint8

const (
	UndefinedSize Size = iota // json:""
	Small                     // json:"small"
	Large                     // json:"large"
	XLarge                    // json:"xlarge"
)
