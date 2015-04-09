package models
import (
	"github.com/jinzhu/gorm"
	"fmt"
	"utils"
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

type AccountRest struct {
    AccountName string
	Amount int64
}

func BalanceRest(db *gorm.DB) []AccountRest {
	rows, _ := db.Raw("select a.name as AccountName, ifnull((select sum(t2.amount_to) from transactions t2 where t2.account_to_id = a.id),0)-ifnull((select sum(t1.amount_from) from transactions t1 where t1.account_from_id = a.id),0) as Amount from accounts a where a.type = \"B\" order by name").Rows()
	defer rows.Close()

    var result []AccountRest

	for rows.Next() {
		item := AccountRest{}
        rows.Scan(&item.AccountName, &item.Amount)
		fmt.Println(item)

		result = append(result, item)
    }

    return result
}

func (ar AccountRest) AmountAsFloat() float64{
	return float64(ar.Amount) / 100.0
}
func (ar AccountRest) AmountAsMoney() string{
	return utils.RenderFloat( "# ###.##", ar.AmountAsFloat() )
}

