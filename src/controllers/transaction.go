package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
//    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
    "strconv"
    "time"
)

func GetTransactions(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Order("id desc").Find(&transactions)
    r.HTML(200, "transactions/index", models.TransactionsIndex{T: transactions, Rests: models.BalanceRest(db)})
}

func GetTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var transaction models.Transaction
    db.Find(&transaction, params["id"])

    var accountList []models.Account
    db.Order("name").Find(&accountList)

    formData := models.TransactionForm{ T: transaction, AccountList: accountList }

    r.HTML(200, "transactions/edit", formData)
}

func NewTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction := models.Transaction{}
    transaction.Date = time.Now().Format("2006-01-02")

    var accountList []models.Account
    db.Order("name").Find(&accountList)

    formData := models.TransactionForm{ T: transaction, AccountList: accountList }

    r.HTML(200, "transactions/new", formData)
}

func CreateTransaction(transaction models.Transaction, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction.ParseMoney()
    db.Create(&transaction)
    r.Redirect("/transactions")
}

func UpdateTransaction(transaction models.Transaction, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    transaction.ParseMoney()
    db.Save(transaction)
    r.Redirect("/transactions")
}

func DeleteTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Transaction{})
    r.Redirect("/transactions")
}