package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
//    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    "models"
    "strconv"
    "time"
    "github.com/jinzhu/now"
    "math"
)

const PER_PAGE = 5

func GetTransactions(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){

    currentPage, _ := strconv.ParseInt( req.URL.Query().Get("page"), 10, 64)
    var totalRecords int64
    db.Model(models.Transaction{}).Count(&totalRecords)

    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Order("date desc, id desc").Offset(currentPage*PER_PAGE).Limit(PER_PAGE).Find(&transactions)

    pt := now.New( time.Now().AddDate(0,-1,0) )

    totalPages := int64( math.Ceil( float64(totalRecords) / float64(PER_PAGE) ) )

    r.HTML(
      200, "transactions/index",
      models.TransactionsIndex{
          T: transactions,
          Rests: models.BalanceRest(db),
          CurrentIncome: models.IncomeForPeriod(db, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02") ),
          CurrentExpense: models.ExpenseForPeriod(db, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02") ),

          PreviousIncome: models.IncomeForPeriod(db, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02") ),
          PreviousExpense: models.ExpenseForPeriod(db, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02") ),

          CurrentMonth: time.Now().Month().String(),
          PreviousMonth: (time.Now().Month()-1).String(),

          Page: currentPage,
          TotalPages: make([]byte, totalPages),
      },
    );
}

func GetTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var transaction models.Transaction
    db.Find(&transaction, params["id"])

    var accountList []models.Account
    db.Order("name").Find(&accountList)

    formData := models.TransactionForm{ T: transaction, AccountList: accountList }

    r.HTML(200, "transactions/edit", formData)
}

func NewTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction := models.Transaction{}
    transaction.Date = time.Now().Format("2006-01-02")

    var accountList []models.Account
    db.Order("name").Find(&accountList)

    formData := models.TransactionForm{ T: transaction, AccountList: accountList }

    r.HTML(200, "transactions/new", formData)
}

func CreateTransaction(transaction models.Transaction, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction.ParseMoney()
    db.Create(&transaction)
    r.Redirect("/transactions")
}

func UpdateTransaction(transaction models.Transaction, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    transaction.ParseMoney()
    db.Save(transaction)
    r.Redirect("/transactions")
}

func DeleteTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Transaction{})
    r.Redirect("/transactions")
}