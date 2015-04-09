package models

type Account struct {
    ID int64    `form:"ID"`
    Name string `form:"Name"`
	CurrencyID string `form:"CurrencyID"`
}

type AccountForm struct {
	T Account
	CurrencyList []Currency
}

