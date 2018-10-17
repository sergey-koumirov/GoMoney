package controllers

import (

	//    "fmt"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/models"
)

func GetAccounts(c *gin.Context) {
	var accounts []models.Account
	db.DBI.Preload("Currency").Find(&accounts)
	c.HTML(200, "accounts/index", accounts)
}

func GetAccount(c *gin.Context) {
	var account models.Account
	db.DBI.Find(&account, c.Param("id"))

	var currencyList []models.Currency
	db.DBI.Order("num").Find(&currencyList)

	c.HTML(200, "accounts/edit", models.AccountForm{A: account, CurrencyList: currencyList})
}

func NewAccount(c *gin.Context) {
	account := models.Account{}

	var currencyList []models.Currency
	db.DBI.Order("num").Find(&currencyList)

	c.HTML(200, "accounts/new", models.AccountForm{A: account, CurrencyList: currencyList})
}

func CreateAccount(c *gin.Context) {
	var account models.Account

	if c.ShouldBind(&account) == nil {
		db.DBI.Create(&account)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/accounts")
}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	account.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Find(&account)

	if c.ShouldBind(&account) == nil {
		db.DBI.Save(&account)
	}

	c.Redirect(http.StatusTemporaryRedirect, "/accounts")
}

func DeleteAccount(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Where("id = ?", id).Delete(models.Account{})
	c.Redirect(http.StatusTemporaryRedirect, "/accounts")
}
