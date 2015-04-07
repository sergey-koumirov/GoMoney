package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
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

func NewAccount(params martini.Params){
    fmt.Println(params)
}

func UpdateAccount(params martini.Params){
    fmt.Println(params)
}

func DeleteAccount(params martini.Params){
    fmt.Println(params)
}