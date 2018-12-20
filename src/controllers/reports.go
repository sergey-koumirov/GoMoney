package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/sergey-koumirov/GoMoney/src/models"

	"github.com/jinzhu/now"

	"github.com/sergey-koumirov/GoMoney/src/db"
)

func GetReportDateRange(c *gin.Context) {

	reportData := models.DateRangeReport{
		InTransactions:  []models.Transaction{},
		OutTransactions: []models.Transaction{},
		AccountList:     []models.Account{},
	}

	if c.Query("BeginDate") != "" {
		reportData.BeginDate = c.Query("BeginDate")
	} else {
		reportData.BeginDate = now.BeginningOfMonth().Format("2006-01-02")
	}

	if c.Query("EndDate") != "" {
		reportData.EndDate = c.Query("EndDate")
	} else {
		reportData.EndDate = now.EndOfMonth().Format("2006-01-02")
	}

	if c.Query("AccountId") != "" {
		reportData.AccountId, _ = strconv.ParseInt(c.Query("AccountId"), 10, 64)
	} else {
		reportData.AccountId = 0
	}

	reportData.InMove, reportData.OutMove = models.GroupByCurrency(db.DBI, reportData.BeginDate, reportData.EndDate, reportData.AccountId)

	db.DBI.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Where("date >= ? and date <= ? and account_to_id = ?", reportData.BeginDate, reportData.EndDate, reportData.AccountId).
		Order("date desc, id desc").
		Find(&(reportData.InTransactions))

	db.DBI.
		Preload("AccountFrom.Currency").
		Preload("AccountTo.Currency").
		Where("date >= ? and date <= ? and account_from_id = ?", reportData.BeginDate, reportData.EndDate, reportData.AccountId).
		Order("date desc, id desc").
		Find(&(reportData.OutTransactions))

	db.DBI.
		Preload("Currency").
		Where("Type='B'").
		Order("currency_id, id desc").
		Find(&(reportData.AccountList))

	c.HTML(http.StatusOK, "reports/index", reportData)
}
