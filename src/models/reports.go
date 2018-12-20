package models

import (
	"github.com/jinzhu/gorm"
)

type DateRangeReport struct {
	BeginDate       string
	EndDate         string
	AccountId       int64
	AccountList     []Account
	InTransactions  []Transaction
	OutTransactions []Transaction
	InMove          []CurrencyMoveInfo
	OutMove         []CurrencyMoveInfo
}

type CurrencyMoveInfo struct {
	CurrencyId   int64
	CurrencyCode string
	Amount       int64
	Rate         float64
}

func GroupByCurrency(db *gorm.DB, from string, to string, accountId int64) ([]CurrencyMoveInfo, []CurrencyMoveInfo) {
	resultTo := make([]CurrencyMoveInfo, 0)
	sqlTo := `select c.id, c.code, sum(t.amount_from) as amount, sum(t.amount_from)*1.0/sum(t.amount_to) as rate
	from transactions t, 
		 accounts a, 
		 currencies c
	where t.account_to_id=?
	  and t.date >= ? and t.date <= ?
	  and a.id = t.account_from_id
	  and c.id = a.currency_id
	group by c.id, c.code
	order by c.id, c.code`

	rowsTo, _ := db.Raw(sqlTo, accountId, from, to).Rows()

	for rowsTo.Next() {
		item := CurrencyMoveInfo{}
		rowsTo.Scan(
			&item.CurrencyId,
			&item.CurrencyCode,
			&item.Amount,
			&item.Rate,
		)
		resultTo = append(resultTo, item)
	}

	resultFrom := make([]CurrencyMoveInfo, 0)

	sqlFrom := `select c.id, c.code, sum(t.amount_to) as amount, sum(t.amount_to)*1.0/sum(t.amount_from) as rate
	from transactions t, 
		 accounts a, 
		 currencies c
	where t.account_from_id=?
	  and t.date >= ? and t.date <= ?
	  and a.id = t.account_to_id
	  and c.id = a.currency_id
	group by c.id, c.code
	order by c.id, c.code`

	rowsFrom, _ := db.Raw(sqlFrom, accountId, from, to).Rows()

	for rowsFrom.Next() {
		item := CurrencyMoveInfo{}
		rowsFrom.Scan(
			&item.CurrencyId,
			&item.CurrencyCode,
			&item.Amount,
			&item.Rate,
		)
		resultFrom = append(resultFrom, item)
	}

	return resultTo, resultFrom
}
