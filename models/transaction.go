package models

type Transaction struct {
	DATE        string  `csv:"date"`
	AMOUNT      float64 `csv:"amount"`
	DESCRIPTION string  `csv:"description"`
	CATEGORY    string  `csv:"category"`
}
