package models

type Transaction struct {
	DATE        string  `csv:"date" xml:"DTPOSTED"`
	AMOUNT      float64 `csv:"amount" xml:"TRNAMT"`
	DESCRIPTION string  `csv:"description" xml:"NAME"`
	CATEGORY    string  `csv:"category" xml:"TRNTYPE"`
}
