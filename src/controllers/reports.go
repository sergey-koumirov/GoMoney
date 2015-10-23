package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "net/http"
    "github.com/jinzhu/gorm"
    "GoMoney/src/models"
    "github.com/jinzhu/now"
    "strings"
    "GoMoney/src/utils"
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

    var fromIdsFilter, toIdsFilter string
    reportData.AccountFromList, fromIdsFilter = prepareAccountList(db, "account_from_id", req.URL.Query()["fIDs"])
    reportData.AccountToList, toIdsFilter = prepareAccountList(db, "account_to_id", req.URL.Query()["tIDs"])

    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Where("date >= ? and date <= ? "+fromIdsFilter+toIdsFilter,reportData.BeginDate,reportData.EndDate).Order("date desc, id desc").Find(&transactions)

    reportData.Sections = models.FillAccountTypeSectionsInfo(transactions)

    r.HTML( 200, "reports/index", reportData );
}

func prepareAccountList(db *gorm.DB, fieldName string, ids []string) ([]models.Account, string) {
    var accountList []models.Account
    db.Where("hidden<>1").Order("type, name").Find(&accountList)

    idFilter := ""
    if len(ids) > 0 {
        idFilter = " and "+fieldName+" in ("+strings.Join(ids,",")+")" //todo add sanitize
        for index, _ := range accountList {
            if utils.Contains(ids, strconv.FormatInt(accountList[index].ID,10)) {
                accountList[index].Selected = true
            }
        }
    }
    return accountList, idFilter
}
