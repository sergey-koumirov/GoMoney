package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
//    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
    "strconv"
)

func GetAccounts(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var accounts []models.Account
    db.Find(&accounts)
    r.HTML(200, "accounts/index", accounts)
}

func GetAccount(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var account models.Account
    db.Find(&account, params["id"])
    r.HTML(200, "accounts/edit", account)
}

func NewAccount(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    account := models.Account{}
    r.HTML(200, "accounts/new", account)
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