package models

type Account struct {
    ID int64    `form:"ID"`
    Name string `form:"Name"`
	Type string `form:"Type"`

	Currency Currency
	CurrencyID int64 `form:"CurrencyID"`
}

type AccountForm struct {
	A Account
	CurrencyList []Currency
}

