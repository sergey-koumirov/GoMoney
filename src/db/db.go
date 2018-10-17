package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var DBI *gorm.DB

func Connect(workFileName string) {
	//*** DB INIT ***
	db, error := gorm.Open("sqlite3", workFileName)
	if error != nil {
		fmt.Println(error)
	} else {
		db.DB()
		DBI = db
	}
}
