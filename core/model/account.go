package model

import "time"

type Account struct {
	ID           int
	ClientID     int
	CurrencyID   int
	CurrencyCode string
	Number       string
	Balance      float64
	Status       string
	CreatedAt    time.Time
}
