package models
import (
	"github.com/jinzhu/gorm"
	"database/sql"
	"fmt"
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
func extract(result *AccountsInfo, rows *sql.Rows){
	for rows.Next() {
		item := AccountRecord{}
		rows.Scan(&item.AccountName, &item.Amount)
		result.Total = result.Total + item.Amount
		result.Records = append(result.Records, item)
	}
}

func BalanceRest(db *gorm.DB) AccountsInfo {
	rows, _ := db.Raw("select a.name as AccountName, ifnull((select sum(t2.amount_to) from transactions t2 where t2.account_to_id = a.id),0)-ifnull((select sum(t1.amount_from) from transactions t1 where t1.account_from_id = a.id),0) as Amount from accounts a where a.type = \"B\" order by name").Rows()
	defer rows.Close()

	result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)
	return result
}

func IncomeForPeriod(db *gorm.DB, fromDate string, toDate string) AccountsInfo{
	sql := "select a.name as AccountName,"+
	       "       ifnull( (select sum(t.amount_from) from transactions t where t.account_from_id = a.id and t.date >= ? and t.date <= ?),0) as Amount"+
	       "  from accounts a"+
	       "  where a.type = \"I\" "+
	       "  order by Amount desc";
	rows, error := db.Raw(sql, fromDate, toDate).Rows()
	if(error != nil){
		fmt.Println(error)
	}
	defer rows.Close()

    result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)
	return result
}

func ExpenseForPeriod(db *gorm.DB, fromDate string, toDate string) AccountsInfo{
	sql := "select a.name as AccountName,"+
	"       ifnull( (select sum(t.amount_to) from transactions t where t.account_to_id = a.id and t.date >= ? and t.date <= ?),0) as Amount"+
	"  from accounts a"+
	"  where a.type = \"E\" "+
	"  order by Amount desc";
	rows, error := db.Raw(sql, fromDate, toDate).Rows()
	if(error != nil){
		fmt.Println(error)
	}
	defer rows.Close()

	result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)
	return result
}

