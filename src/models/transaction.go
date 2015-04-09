package models
import (
    "strconv"
    "fmt"
    "utils"
)

type Transaction struct {
    ID int64             `form:"ID"`
    AccountFrom   Account
    AccountFromID int64  `form:"AccountFromID"`
    AccountTo     Account
    AccountToID   int64  `form:"AccountToID"`
    Date          string `form:"Date"`
    Amount        int64
    AmountStr     string `form:"Amount" sql:"-"`
    Comment       string `form:"Comment"`
}


type TransactionForm struct {
    T Transaction
    AccountList []Account
}


func (t *Transaction) ParseAmount() {
    fAmount, _ := strconv.ParseFloat(t.AmountStr, 64)
    t.Amount = int64( fAmount * 100 )
    fmt.Println(t)
}

func (t Transaction) AmountAsFloat() float64{
    return float64(t.Amount) / 100.0
}

func (t Transaction) AmountAsMoney() string{
    return utils.RenderFloat( "#,###.##", t.AmountAsFloat() )
}


