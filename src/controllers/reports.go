package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
    "github.com/jinzhu/now"
    "strings"
    "utils"
    "strconv"
)

func GetReportDateRange(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){

    reportData := models.DateRangeReport{}

    if req.URL.Query().Get("BeginDate") != "" {
        reportData.BeginDate = req.URL.Query().Get("BeginDate")
    }else{
        reportData.BeginDate = now.BeginningOfMonth().Format("2006-01-02")
    }

    if req.URL.Query().Get("EndDate") != "" {
        reportData.EndDate = req.URL.Query().Get("EndDate")
    }else{
        reportData.EndDate = now.EndOfMonth().Format("2006-01-02")
    }

    var accountList []models.Account
    db.Where("hidden<>1").Order("type, name").Find(&accountList)

    idFilter := ""
    if len(req.URL.Query()["IDs"]) > 0 {
        idFilter = " and account_from_id in ("+strings.Join(req.URL.Query()["IDs"],",")+")" //todo add sanitize
        for index, _ := range accountList {
            if utils.Contains(req.URL.Query()["IDs"], strconv.FormatInt(accountList[index].ID,10)) {
                accountList[index].Selected = true
            }
        }
    }

    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Where("date >= ? and date <= ? "+idFilter,reportData.BeginDate,reportData.EndDate).Order("date desc, id desc").Find(&transactions)



    reportData.Sections = models.FillAccountTypeSectionsInfo(transactions)
    reportData.AccountList = accountList

    r.HTML( 200, "reports/index", reportData );
}

