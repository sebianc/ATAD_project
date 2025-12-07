package models

type Budget struct {
	ID       int64   `db:"id"`
	Category string  `db:"category"`
	Limit    float64 `db:"limit"`
}
