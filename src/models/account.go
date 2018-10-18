package models

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
)

type Account struct {
	ID         int64  `form:"ID"`
	Name       string `form:"Name"`
	Type       string `form:"Type"`
	Currency   Currency
	CurrencyID int64 `form:"CurrencyID"`
	Hidden     int64 `form:"Hidden"`

	Selected bool `sql:"-"`
}

type AccountDescription struct {
	ID           int64
	Name         string
	CurrencyCode string
	Type         string
	Hidden       int64
	TrCnt        int64
	TrCntYear    int64
}

type AccountForm struct {
	A            Account
	CurrencyList []Currency
}

type AccountRecord struct {
	AccountName  string
	CurrencyCode string
	Amount       int64
	Percent      float64
	SumPercent   float64
}

type AccountsInfo struct {
	Records []AccountRecord
	Total   int64
}

func extract(result *AccountsInfo, rows *sql.Rows) {
	for rows.Next() {
		item := AccountRecord{}
		rows.Scan(&item.AccountName, &item.CurrencyCode, &item.Amount)
		result.Total = result.Total + item.Amount
		result.Records = append(result.Records, item)
	}
}

func BalanceRest(db *gorm.DB) AccountsInfo {
	rows, e1 := db.Raw(
		`select a.name as AccountName,
		        c.code as CurrencyCode,
		        ifnull(
							(select sum(t2.amount_to) from transactions t2 where t2.account_to_id = a.id),
							0
						) -
					  ifnull(
							(select sum(t1.amount_from) from transactions t1 where t1.account_from_id = a.id),
							0
						) as Amount
			from accounts a
			       left join currencies c on c.id = a.currency_id
			where a.type = "B" and a.hidden<>1
			order by c.code, name`,
	).Rows()

	if e1 != nil {
		fmt.Println(e1)
	}

	defer rows.Close()

	result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)
	return result
}

func IncomeForPeriod(db *gorm.DB, fromDate string, toDate string) AccountsInfo {
	sql := "select a.name as AccountName," +
		"       ifnull( (select sum(t.amount_from) from transactions t where t.account_from_id = a.id and t.date >= ? and t.date <= ?),0) as Amount" +
		"  from accounts a" +
		"  where a.type = \"I\" and hidden<> 1" +
		"    and Amount > 0 " +
		"  order by Amount desc"
	rows, error := db.Raw(sql, fromDate, toDate).Rows()
	if error != nil {
		fmt.Println(error)
	}
	defer rows.Close()

	result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)
	return result
}

func ExpenseForPeriod(db *gorm.DB, fromDate string, toDate string) AccountsInfo {
	sql := "select a.name as AccountName," +
		"       ifnull( (select sum(t.amount_to) from transactions t where t.account_to_id = a.id and t.date >= ? and t.date <= ?),0) as Amount" +
		"  from accounts a" +
		"  where a.type = \"E\" and hidden<>1 " +
		"    and Amount > 0 " +
		"  order by Amount desc"
	rows, error := db.Raw(sql, fromDate, toDate).Rows()
	if error != nil {
		fmt.Println(error)
	}
	defer rows.Close()

	result := AccountsInfo{Records: []AccountRecord{}, Total: 0}
	extract(&result, rows)

	agg := float64(0)
	for index, record := range result.Records {
		result.Records[index].Percent = 100.0 * float64(record.Amount) / float64(result.Total)
		agg = agg + result.Records[index].Percent
		result.Records[index].SumPercent = agg
	}

	return result
}

func AccountDescriptions(db *gorm.DB, yearAgo string) []AccountDescription {
	result := make([]AccountDescription, 0)

	sql := `select a.id,
	               a.name,
								 c.code,
								 a.type,
								 a.hidden,
								 (select count(1) from transactions t where t.account_from_id = a.id or t.account_to_id = a.id) as cnt,
								 (select count(1) from transactions t where (t.account_from_id = a.id or t.account_to_id = a.id) and t.date > ?) as cnty
	          from accounts a
						       left join currencies c on c.id = a.currency_id
						order by c.code, a.type, a.name`

	rows, e1 := db.Raw(sql, yearAgo).Rows()
	if e1 != nil {
		fmt.Println(e1)
	}

	for rows.Next() {
		item := AccountDescription{}
		rows.Scan(
			&item.ID,
			&item.Name,
			&item.CurrencyCode,
			&item.Type,
			&item.Hidden,
			&item.TrCnt,
			&item.TrCntYear,
		)
		result = append(result, item)
	}

	return result
}
