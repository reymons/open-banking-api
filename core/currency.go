package core

import "errors"

const (
	CurrencyEUR = 1
	CurrencyUSD = 2
	CurrencyRUB = 3
	CurrencyRSD = 4
)

var currencyCodeToID = map[string]int{
	"EUR": 1,
	"USD": 2,
	"RUB": 3,
	"RSD": 4,
}

func GetCurrencyID(code string) (int, error) {
	if id, ok := currencyCodeToID[code]; !ok {
		// TODO: handler error better?
		return -1, errors.New("unknown currency")
	} else {
		return id, nil
	}
}
