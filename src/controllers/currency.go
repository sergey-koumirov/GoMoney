package controllers

import (

	//    "fmt"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/models"
)

func GetCurrencies(c *gin.Context) {
	var currencies []models.Currency
	db.DBI.Find(&currencies)
	c.HTML(200, "currencies/index", currencies)
}

func GetCurrency(c *gin.Context) {
	var currency models.Currency
	db.DBI.Find(&currency, c.Param("id"))
	c.HTML(200, "currencies/edit", currency)
}

func NewCurrency(c *gin.Context) {
	currency := models.Currency{}
	c.HTML(200, "currencies/new", currency)
}

func CreateCurrency(c *gin.Context) {

	var currency models.Currency

	if c.ShouldBind(&currency) == nil {
		db.DBI.Create(&currency)
	}

	c.Redirect(http.StatusSeeOther, "/currencies")
}

func UpdateCurrency(c *gin.Context) {
	var currency models.Currency
	currency.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Find(&currency)

	if c.ShouldBind(&currency) == nil {
		db.DBI.Save(&currency)
	}

	c.Redirect(http.StatusSeeOther, "/currencies")
}

func DeleteCurrency(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	err := db.DBI.Delete(models.Currency{ID: id}).Error

	fmt.Println("123", id, err)

	c.Redirect(302, "/currencies")
}
