package models
import (
	"github.com/jinzhu/gorm"
)

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

type AccountRecord struct {
	AccountName string
	Amount int64
}

type AccountsInfo struct {
    Records []AccountRecord
	Total int64
}

func BalanceRest(db *gorm.DB) AccountsInfo {
	rows, _ := db.Raw("select a.name as AccountName, ifnull((select sum(t2.amount_to) from transactions t2 where t2.account_to_id = a.id),0)-ifnull((select sum(t1.amount_from) from transactions t1 where t1.account_from_id = a.id),0) as Amount from accounts a where a.type = \"B\" order by name").Rows()
	defer rows.Close()
	var result []AccountRecord
	total := int64(0)
	for rows.Next() {
		item := AccountRecord{}
		rows.Scan(&item.AccountName, &item.Amount)
        total = total + item.Amount
		result = append(result, item)
	}
	return AccountsInfo{Records: result, Total: total}
}

//func IncomeForPeriod(fromDate string, toDate string) AccountsInfo{
//
//}
//
//func ExpenseForPeriod(fromDate string, toDate string) AccountsInfo{
//
//}

