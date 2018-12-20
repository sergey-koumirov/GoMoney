package controllers

import (
	"fmt"
	"math"
	"strconv"
	"time"

	// "io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/models"
)

const PER_PAGE = 200

func GetTransactions(c *gin.Context) {

	currentPage, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	var totalRecords int64
	db.DBI.Model(models.Transaction{}).Count(&totalRecords)

	var transactions []models.Transaction
	db.DBI.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Order("date desc, id desc").
		Offset(int(currentPage) * PER_PAGE).
		Limit(PER_PAGE).
		Find(&transactions)

	for i := range transactions {
		transactions[i].HasDiff = transactions[i].AccountFrom.CurrencyID != transactions[i].AccountTo.CurrencyID || transactions[i].AmountFrom != transactions[i].AmountTo
	}

	prevD := 1
	prevM := time.Now().Month()
	prevY := time.Now().Year()

	sYear := now.New(time.Date(prevY-1, prevM, prevD, 0, 0, 0, 0, time.UTC))

	if prevM == 1 {
		prevM = 12
		prevY = prevY - 1
	} else {
		prevM = prevM - 1
	}

	pt := now.New(time.Date(prevY, prevM, prevD, 0, 0, 0, 0, time.UTC))

	totalPages := int64(math.Ceil(float64(totalRecords) / float64(PER_PAGE)))

	var templates []models.Template
	db.DBI.Order("id desc").Find(&templates)

	var alarms []string

	var prevMonth time.Month
	if time.Now().Month() == 1 {
		prevMonth = time.December
	} else {
		prevMonth = time.Now().Month() - 1
	}

	c.HTML(
		200, "transactions/index",
		models.TransactionsIndex{
			T:              transactions,
			Rests:          models.BalanceRest(db.DBI),
			CurrentIncome:  models.IncomeForPeriod(db.DBI, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02")),
			CurrentExpense: models.ExpenseForPeriod(db.DBI, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02")),

			PreviousIncome:  models.IncomeForPeriod(db.DBI, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02")),
			PreviousExpense: models.ExpenseForPeriod(db.DBI, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02")),

			YearIncome:  models.IncomeForPeriod(db.DBI, sYear.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02")),
			YearExpense: models.ExpenseForPeriod(db.DBI, sYear.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02")),

			CurrentDate:   time.Now().Format("2006-01-02"),
			CurrentMonth:  time.Now().Month().String(),
			PreviousMonth: prevMonth.String(),

			Page:       currentPage,
			TotalPages: make([]byte, totalPages),

			Templates: templates,
			Alarms:    alarms,
		},
	)
}

func GetTransaction(c *gin.Context) {
	var transaction models.Transaction
	db.DBI.Find(&transaction, c.Param("id"))

	var accountFromList []models.Account
	db.DBI.Where("type in (\"I\",\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)

	var accountToList []models.Account
	db.DBI.Where("type in (\"E\",\"B\") and hidden<>1").Order("type, name").Find(&accountToList)

	formData := models.TransactionForm{T: transaction, AccountFromList: accountFromList, AccountToList: accountToList}

	c.HTML(200, "transactions/edit", formData)
}

func NewTransaction(c *gin.Context) {
	transaction := models.Transaction{}
	transaction.Date = time.Now().Format("2006-01-02")

	template := models.Template{}

	if c.Query("template_id") != "" {
		template_id, _ := strconv.ParseInt(c.Query("template_id"), 10, 64)
		db.DBI.Find(&template, template_id)
		copyFromTemplate(db.DBI, &transaction, template)
	}

	var accountFromList []models.Account
	var accountToList []models.Account

	if c.Query("type") == "E" {
		db.DBI.Preload("Currency").Where("type in (\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)
		db.DBI.Preload("Currency").Where("type in (\"E\") and hidden<>1").Order("type, name").Find(&accountToList)
	} else if c.Query("type") == "I" {
		db.DBI.Preload("Currency").Where("type in (\"I\") and hidden<>1").Order("type, name").Find(&accountFromList)
		db.DBI.Preload("Currency").Where("type in (\"B\") and hidden<>1").Order("type, name").Find(&accountToList)
	} else {
		db.DBI.Preload("Currency").Where("type in (\"I\",\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)
		db.DBI.Preload("Currency").Where("type in (\"E\",\"B\") and hidden<>1").Order("type, name").Find(&accountToList)
	}

	formData := models.TransactionForm{T: transaction, AccountFromList: accountFromList, AccountToList: accountToList, FocusOn: template.FocusOn}

	c.HTML(200, "transactions/new", formData)
}

func CreateTransaction(c *gin.Context) {
	var transaction models.Transaction

	parse := c.ShouldBind(&transaction)
	if parse == nil {
		transaction.ParseMoney()
		db.DBI.Create(&transaction)
	} else {
		fmt.Println("parse:", parse)
	}

	c.Redirect(302, "/transactions")
}

func UpdateTransaction(c *gin.Context) {
	var transaction models.Transaction
	transaction.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Find(&transaction)

	if c.ShouldBind(&transaction) == nil {
		transaction.ParseMoney()
		db.DBI.Save(&transaction)
	}

	c.Redirect(302, "/transactions")
}

func DeleteTransaction(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Where("id = ?", id).Delete(models.Transaction{})
	c.Redirect(302, "/transactions")
}

func copyFromTemplate(db *gorm.DB, t *models.Transaction, template models.Template) {
	t.AccountFromID = template.AccountFromID
	t.AccountToID = template.AccountToID
	t.AmountFrom = template.AmountFrom
	t.AmountTo = template.AmountTo
}
