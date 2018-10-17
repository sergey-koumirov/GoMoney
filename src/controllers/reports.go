package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sergey-koumirov/GoMoney/src/models"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/now"

	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/utils"
)

func GetReportDateRange(c *gin.Context) {

	reportData := models.DateRangeReport{}

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

	var fromIdsFilter, toIdsFilter string
	reportData.AccountFromList, fromIdsFilter = prepareAccountList(db.DBI, "account_from_id", c.Request.URL.Query()["fIDs"])
	reportData.AccountToList, toIdsFilter = prepareAccountList(db.DBI, "account_to_id", c.Request.URL.Query()["tIDs"])

	var transactions []models.Transaction
	db.DBI.Preload("AccountFrom").Preload("AccountTo").Where("date >= ? and date <= ? "+fromIdsFilter+toIdsFilter, reportData.BeginDate, reportData.EndDate).Order("date desc, id desc").Find(&transactions)

	reportData.Sections = models.FillAccountTypeSectionsInfo(transactions)

	c.HTML(http.StatusOK, "reports/index", reportData)
}

func prepareAccountList(db *gorm.DB, fieldName string, ids []string) ([]models.Account, string) {
	var accountList []models.Account
	db.Where("hidden<>1").Order("type, name").Find(&accountList)

	idFilter := ""
	if len(ids) > 0 {
		idFilter = " and " + fieldName + " in (" + strings.Join(ids, ",") + ")" //todo add sanitize
		for index, _ := range accountList {
			if utils.Contains(ids, strconv.FormatInt(accountList[index].ID, 10)) {
				accountList[index].Selected = true
			}
		}
	}
	return accountList, idFilter
}
