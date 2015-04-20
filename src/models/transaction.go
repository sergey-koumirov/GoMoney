package models
import (
    "strconv"
)

type Transaction struct {
    ID int64             `form:"ID"`
    AccountFrom   Account
    AccountFromID int64  `form:"AccountFromID"`
    AccountTo     Account
    AccountToID   int64  `form:"AccountToID"`
    Date          string `form:"Date"`
    AmountFrom    int64
    AmountFromStr string `form:"AmountFrom" sql:"-"`
    AmountTo      int64
    AmountToStr   string `form:"AmountTo" sql:"-"`
    Comment       string `form:"Comment"`
}


type TransactionForm struct {
    T Transaction
    AccountFromList []Account
    AccountToList []Account
    FocusOn string
}

type TransactionsIndex struct {
    T []Transaction
    Rests AccountsInfo
    CurrentIncome AccountsInfo
    CurrentExpense AccountsInfo
    PreviousIncome AccountsInfo
    PreviousExpense AccountsInfo
    CurrentDate string
    CurrentMonth string
    PreviousMonth string
    Page int64
    TotalPages []byte
    Templates []Template
}

func (t *Transaction) ParseMoney() {
    fAmountFrom, _ := strconv.ParseFloat(t.AmountFromStr, 64)
    t.AmountFrom = int64( fAmountFrom * 100 )

    fAmountTo, _ := strconv.ParseFloat(t.AmountToStr, 64)
    if fAmountTo == 0 {
        fAmountTo = fAmountFrom
    }
    t.AmountTo = int64( fAmountTo * 100 )
}
