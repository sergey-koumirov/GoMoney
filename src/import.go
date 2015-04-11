package main

import(
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "os"
    "bufio"
    "strings"
    "models"
    "strconv"
    "utils"
)

func parseRecord(
         db gorm.DB,
         dateStr string,
         moneyStr string,
         accountFromStr string,
         accountFromType string,
         accountToStr string,
         accountToType string,
         comment string,
){
    date_parts := strings.Split(dateStr, ".")

    cleanedMoney := strings.Replace(moneyStr,"\u00a0","",-1)
    cleanedMoney = strings.Replace(cleanedMoney,"р","",-1)
    cleanedMoney = strings.Replace(cleanedMoney,",",".",-1)

    amount, error := strconv.ParseFloat(cleanedMoney, 64)
    if(error !=nil){ fmt.Println(error) }

    var af models.Account
    db.Where("name = ? and type = ? ", strings.Trim(accountFromStr," \u00a0"), accountFromType ).First(&af)
    if(af.ID == 0){
        af.CurrencyID = 1
        af.Name = strings.Trim(accountFromStr," \u00a0")
        af.Type = accountFromType
        db.Create(&af)
    }

    var at models.Account
    db.Where("name = ? and type = ?", strings.Trim(accountToStr," \u00a0"), accountToType ).First(&at)
    if(at.ID == 0){
        at.CurrencyID = 1
        at.Name = strings.Trim(accountToStr," \u00a0")
        at.Type = accountToType
        db.Create(&at)
    }

    //fmt.Println( strings.Trim(accountFromStr," \u00a0") )
    //fmt.Println( af )

    temp := models.Transaction{
        Date: date_parts[2]+"-"+date_parts[1]+"-"+date_parts[0],
        AmountFrom: int64( utils.Round(amount*100, 0.5, 0) ),
        AmountTo: int64( utils.Round(amount*100, 0.5, 0) ),
        AccountFromID: af.ID,
        AccountToID: at.ID,
        Comment: comment,
    }
    db.Create(&temp)
    fmt.Println( temp )
}

func main() {
    //*** DB INIT ***
    db, error := gorm.Open("sqlite3", "money_0.db")
    if(error !=nil){ fmt.Println(error) }
    defer db.Close()
    db.DB()

    db.Where("1=1").Delete(models.Transaction{})
    db.Where("1=1").Delete(models.Account{})

    f0, _ := os.Open("./temp/begin.TXT")
    defer f0.Close()
    scanner := bufio.NewScanner(f0)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), ";")
        parseRecord(db, "01.09.2012", parts[1], "Init", "I", parts[0], "B", "")
    }

    f1, _ := os.Open("./temp/incomes.TXT")
    defer f1.Close()
    scanner = bufio.NewScanner(f1)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), ";")
        parseRecord(db, parts[0], parts[5], parts[2]+" "+parts[3], "I", parts[1], "B", parts[6])
    }

    f2, _ := os.Open("./temp/expenses.TXT")
    defer f2.Close()
    scanner = bufio.NewScanner(f2)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), ";")
        parseRecord(db, parts[0], parts[5], parts[1], "B", parts[2]+" "+parts[3], "E", parts[6])
    }

    f3, _ := os.Open("./temp/transfers.TXT")
    defer f3.Close()
    scanner = bufio.NewScanner(f3)
    for scanner.Scan() {
        parts := strings.Split(scanner.Text(), ";")
        parseRecord(db, parts[0], parts[3], parts[1], "B", parts[2], "B","")
    }


}