package color

type CurrencyString struct{ c uint8 }

var (
	UndefinedCurrencyS = CurrencyString{}   // json:""
	SGDS               = CurrencyString{1}  // json:"SGD"
	USDS               = CurrencyString{2}  // json:"USD"
	GBPS               = CurrencyString{3}  // json:"GBP"
	KRWS               = CurrencyString{4}  // json:"KRW"
	HKDS               = CurrencyString{5}  // json:"HKD"
	JPYS               = CurrencyString{6}  // json:"JPY"
	MYRS               = CurrencyString{7}  // json:"MYR"
	BHTS               = CurrencyString{8}  // json:"BHT"
	THCS               = CurrencyString{9}  // json:"THC"
	CBDS               = CurrencyString{10} // json:"CBD"
	XYZS               = CurrencyString{11} // json:"XYZ"
)
