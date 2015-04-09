package models
import (
    "strconv"
    "utils"
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
    AccountList []Account
}

type TransactionsIndex struct {
    T []Transaction
    Rests []AccountRest
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

func (t Transaction) AmountFromAsFloat() float64{
    return float64(t.AmountFrom) / 100.0
}

func (t Transaction) AmountToAsFloat() float64{
    return float64(t.AmountTo) / 100.0
}

func (t Transaction) AmountFromAsMoney() string{
    return utils.RenderFloat( "# ###.##", t.AmountFromAsFloat() )
}

func (t Transaction) AmountToAsMoney() string{
    return utils.RenderFloat( "# ###.##", t.AmountToAsFloat() )
}


