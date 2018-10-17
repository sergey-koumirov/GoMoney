package main

import (
	"fmt"

	"github.com/sergey-koumirov/GoMoney/src/models"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//*** DB INIT ***
	db, error := gorm.Open("sqlite3", "money_0.prod.db")
	if error != nil {
		fmt.Println(error)
	}
	defer db.Close()
	db.DB()

	db.AutoMigrate(&models.Account{})
	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.Currency{})
	db.AutoMigrate(&models.Template{})

}
