package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/models"
)

func GetAccounts(c *gin.Context) {

	D := time.Now().Day()
	M := time.Now().Month()
	Y := time.Now().Year()
	yearAgo := now.New(time.Date(Y-1, M, D, 0, 0, 0, 0, time.UTC))

	c.HTML(200, "accounts/index", models.AccountDescriptions(db.DBI, yearAgo.Format("2006-01-02")))
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

	c.Redirect(302, "/accounts")
}

func UpdateAccount(c *gin.Context) {
	var account models.Account
	account.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Find(&account)

	if c.ShouldBind(&account) == nil {
		db.DBI.Save(&account)
	}

	c.Redirect(302, "/accounts")
}

func DeleteAccount(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Where("id = ?", id).Delete(models.Account{})
	c.Redirect(302, "/accounts")
}
