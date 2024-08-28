package color

type Currency struct{ c uint8 }

//go:generate go-enum-encoding -type=Currency
var (
	UndefinedCurrency = Currency{}   // json:""
	SGD               = Currency{1}  // json:"SGD"
	USD               = Currency{2}  // json:"USD"
	GBP               = Currency{3}  // json:"GBP"
	KRW               = Currency{4}  // json:"KRW"
	HKD               = Currency{5}  // json:"HKD"
	JPY               = Currency{6}  // json:"JPY"
	MYR               = Currency{7}  // json:"MYR"
	BHT               = Currency{8}  // json:"BHT"
	THC               = Currency{9}  // json:"THC"
	CBD               = Currency{10} // json:"CBD"
	XYZ               = Currency{11} // json:"XYZ"
)
