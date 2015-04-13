package main

import(
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "models"
)


func main() {
    //*** DB INIT ***
    db, error := gorm.Open("sqlite3", "money_0.db")
    if(error !=nil){ fmt.Println(error) }
    defer db.Close()
    db.DB()

    db.AutoMigrate(&models.Meter{})
    db.AutoMigrate(&models.MeterValue{})
}