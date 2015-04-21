package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
)

func GetReportDateRange(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){

    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Order("date desc, id desc").Find(&transactions)

    var accountList []models.Account
    db.Where("hidden<>1").Order("type, name").Find(&accountList)

    selectedIDs := []int64{}

    r.HTML(
      200, "reports/index",
      models.DateRangeReport{
          T: transactions,
          Total: 0,
          BeginDate: "",
          EndDate: "",
          AccountList: accountList,
          AccountIDs: selectedIDs,
      },
    );
}

