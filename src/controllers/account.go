package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
//    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    models "GoMoney/src/models"
    "strconv"
)

func GetAccounts(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var accounts []models.Account
    db.Preload("Currency").Find(&accounts)
    r.HTML(200, "accounts/index", accounts)
}

func GetAccount(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var account models.Account
    db.Find(&account, params["id"])

    var currencyList []models.Currency
    db.Order("num").Find(&currencyList)

    r.HTML(200, "accounts/edit", models.AccountForm{A: account, CurrencyList: currencyList})
}

func NewAccount(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    account := models.Account{}

    var currencyList []models.Currency
    db.Order("num").Find(&currencyList)

    r.HTML(200, "accounts/new", models.AccountForm{A: account, CurrencyList: currencyList})
}

func CreateAccount(account models.Account, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    db.Create(&account)
    r.Redirect("/accounts")
}

func UpdateAccount(account models.Account, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    account.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    db.Save(account)
    r.Redirect("/accounts")
}

func DeleteAccount(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Account{})
    r.Redirect("/accounts")
}