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

func GetCurrencies(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var currencies []models.Currency
    db.Find(&currencies)
    r.HTML(200, "currencies/index", currencies)
}

func GetCurrency(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var currency models.Currency
    db.Find(&currency, params["id"])
    r.HTML(200, "currencies/edit", currency)
}

func NewCurrency(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    currency := models.Currency{}
    r.HTML(200, "currencies/new", currency)
}

func CreateCurrency(currency models.Currency, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    db.Create(&currency)
    r.Redirect("/currencies")
}

func UpdateCurrency(currency models.Currency, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    currency.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    db.Save(currency)
    r.Redirect("/currencies")
}

func DeleteCurrency(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Currency{})
    r.Redirect("/currencies")
}