package models

// budget struct represents a budget limit for a specific category
type Budget struct {
	ID       int64   `db:"id"`
	Category string  `db:"category"`
	Limit    float64 `db:"limit"`
}
