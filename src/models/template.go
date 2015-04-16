package models
import (
	"strconv"
)

type Template struct {
	ID int64             `form:"ID"`
	Title string         `form:"Title"`
	AccountFrom   Account
	AccountFromID int64  `form:"AccountFromID"`
	AccountTo     Account
	AccountToID   int64  `form:"AccountToID"`
	AmountFrom    int64
	AmountFromStr string `form:"AmountFrom" sql:"-"`
	AmountTo      int64
	AmountToStr   string `form:"AmountTo" sql:"-"`
    FocusOn       string `form:"FocusOn"`
}


type TemplateForm struct {
	T Template
	AccountList []Account
}

func (t *Template) ParseMoney() {
	fAmountFrom, _ := strconv.ParseFloat(t.AmountFromStr, 64)
	t.AmountFrom = int64( fAmountFrom * 100 )

	fAmountTo, _ := strconv.ParseFloat(t.AmountToStr, 64)
	if fAmountTo == 0 {
		fAmountTo = fAmountFrom
	}
	t.AmountTo = int64( fAmountTo * 100 )
}


