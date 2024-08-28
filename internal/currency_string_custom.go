package color

type CurrencyStringCustom struct{ c uint8 }

var (
	UndefinedCurrencySC = CurrencyStringCustom{}   // json:""
	SGDSC               = CurrencyStringCustom{1}  // json:"SGD"
	USDSC               = CurrencyStringCustom{2}  // json:"USD"
	GBPSC               = CurrencyStringCustom{3}  // json:"GBP"
	KRWSC               = CurrencyStringCustom{4}  // json:"KRW"
	HKDSC               = CurrencyStringCustom{5}  // json:"HKD"
	JPYSC               = CurrencyStringCustom{6}  // json:"JPY"
	MYRSC               = CurrencyStringCustom{7}  // json:"MYR"
	BHTSC               = CurrencyStringCustom{8}  // json:"BHT"
	THCSC               = CurrencyStringCustom{9}  // json:"THC"
	CBDSC               = CurrencyStringCustom{10} // json:"CBD"
	XYZSC               = CurrencyStringCustom{11} // json:"XYZ"
)
