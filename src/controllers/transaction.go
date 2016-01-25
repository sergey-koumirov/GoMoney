package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "net/http"
    "github.com/jinzhu/gorm"
    "GoMoney/src/models"
    "strconv"
    "time"
    "github.com/jinzhu/now"
    "math"
    "fmt"
)

const PER_PAGE = 200

func GetTransactions(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){

    currentPage, _ := strconv.ParseInt( req.URL.Query().Get("page"), 10, 64)
    var totalRecords int64
    db.Model(models.Transaction{}).Count(&totalRecords)

    var transactions []models.Transaction
    db.Preload("AccountFrom").Preload("AccountTo").Order("date desc, id desc").Offset(currentPage*PER_PAGE).Limit(PER_PAGE).Find(&transactions)


    prevD := 1
    prevM := time.Now().Month()
    prevY := time.Now().Year()

    if(prevM == 1){
        prevM = 12
        prevY = prevY - 1
    }else{
        prevM = prevM - 1
    }

    pt := now.New( time.Date(prevY,prevM,prevD,0,0,0,0,time.UTC) )

    totalPages := int64( math.Ceil( float64(totalRecords) / float64(PER_PAGE) ) )

    var templates []models.Template
    db.Order("id desc").Find(&templates)


    var alarms []string

    alarmDate := time.Date(time.Now().Year(), time.Now().Month(), 20, 0, 0, 0, 0, time.UTC)
    var count int
    db.Model(models.MeterValue{}).Where("date > ?", alarmDate.Format("2006-01-02")).Count(&count)
    if( alarmDate.Before(time.Now()) && count < 5 ){
        alarms = append(alarms,  fmt.Sprintf("CHECK METER VALUES"))
    }

    var prevMonth time.Month
    if time.Now().Month()==1 {
        prevMonth = time.December
    }else{
        prevMonth = time.Now().Month()-1
    }

    r.HTML(
      200, "transactions/index",
      models.TransactionsIndex{
          T: transactions,
          Rests: models.BalanceRest(db),
          CurrentIncome: models.IncomeForPeriod(db, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02") ),
          CurrentExpense: models.ExpenseForPeriod(db, now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02") ),

          PreviousIncome: models.IncomeForPeriod(db, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02") ),
          PreviousExpense: models.ExpenseForPeriod(db, pt.BeginningOfMonth().Format("2006-01-02"), pt.EndOfMonth().Format("2006-01-02") ),

          CurrentDate: time.Now().Format("2006-01-02"),
          CurrentMonth: time.Now().Month().String(),
          PreviousMonth: prevMonth.String(),

          Page: currentPage,
          TotalPages: make([]byte, totalPages),

          Templates: templates,
          Alarms: alarms,
      },
    );
}

func GetTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var transaction models.Transaction
    db.Find(&transaction, params["id"])

    var accountFromList []models.Account
    db.Where("type in (\"I\",\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)

    var accountToList []models.Account
    db.Where("type in (\"E\",\"B\") and hidden<>1").Order("type, name").Find(&accountToList)

    formData := models.TransactionForm{ T: transaction, AccountFromList: accountFromList, AccountToList: accountToList }

    r.HTML(200, "transactions/edit", formData)
}

func NewTransaction(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    transaction := models.Transaction{}
    transaction.Date = time.Now().Format("2006-01-02")

    template := models.Template{}



    if(req.URL.Query().Get("template_id") != ""){
        template_id, _ := strconv.ParseInt(req.URL.Query().Get("template_id"), 10, 64)
        db.Find(&template, template_id)
        copyFromTemplate(db, &transaction, template)
    }


    var accountFromList []models.Account
    var accountToList []models.Account

    if(req.URL.Query().Get("type") == "E"){
        db.Where("type in (\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)
        db.Where("type in (\"E\") and hidden<>1").Order("type, name").Find(&accountToList)
    }else if(req.URL.Query().Get("type") == "I"){
        db.Where("type in (\"I\") and hidden<>1").Order("type, name").Find(&accountFromList)
        db.Where("type in (\"B\") and hidden<>1").Order("type, name").Find(&accountToList)
    }else{
        db.Where("type in (\"I\",\"B\") and hidden<>1").Order("type, name").Find(&accountFromList)
        db.Where("type in (\"E\",\"B\") and hidden<>1").Order("type, name").Find(&accountToList)
    }


    formData := models.TransactionForm{ T: transaction, AccountFromList: accountFromList, AccountToList: accountToList, FocusOn: template.FocusOn}

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

func copyFromTemplate(db *gorm.DB, t *models.Transaction, template models.Template){
    t.AccountFromID = template.AccountFromID
    t.AccountToID = template.AccountToID
    t.AmountFrom = template.AmountFrom
    t.AmountTo = template.AmountTo
}