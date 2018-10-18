package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sergey-koumirov/GoMoney/src/controllers"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/utils"
)

func main() {
	workFileName := "money_0.prod.db"

	e1 := utils.CopyFile(workFileName, "F:/GoMoneyBackup/"+workFileName+"."+time.Now().Format("2006-01-02"))
	if e1 != nil {
		fmt.Println(e1)
	}

	db.Connect(workFileName)
	defer db.DBI.Close()

	r := gin.Default()

	r.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "templates",
		Extension: ".html",
		Master:    "layout",
		Partials: []string{
			"currencies/_fields",
			"accounts/_fields",
			"templates/_fields",
			"transactions/_fields",
			"transactions/_income",
			"transactions/_expense",
		},
		Funcs: template.FuncMap{
			"money":    utils.RenderMoney,
			"float":    utils.MoneyAsFloat,
			"format3":  utils.RenderFloat3,
			"format64": utils.RenderFloat64,
			"minus": func(a, b sql.NullFloat64) float64 {
				return a.Float64 - b.Float64
			},
		},
		DisableCache: true,
	})

	r.GET("/", controllers.GetTransactions)

	r.GET("/accounts", controllers.GetAccounts)
	r.GET("/accounts/new", controllers.NewAccount)
	r.POST("/accounts", controllers.CreateAccount)
	r.POST("/account/:id", controllers.UpdateAccount)
	r.GET("/account/:id/delete", controllers.DeleteAccount)
	r.GET("/account/:id", controllers.GetAccount)

	r.GET("/currencies", controllers.GetCurrencies)
	r.GET("/currencies/new", controllers.NewCurrency)
	r.POST("/currencies", controllers.CreateCurrency)
	r.POST("/currency/:id", controllers.UpdateCurrency)
	r.GET("/currency/:id/delete", controllers.DeleteCurrency)
	r.GET("/currency/:id", controllers.GetCurrency)

	r.GET("/transactions", controllers.GetTransactions)
	r.GET("/transactions/new", controllers.NewTransaction)
	r.POST("/transactions", controllers.CreateTransaction)
	r.POST("/transaction/:id", controllers.UpdateTransaction)
	r.GET("/transaction/:id/delete", controllers.DeleteTransaction)
	r.GET("/transaction/:id", controllers.GetTransaction)

	r.GET("/templates", controllers.GetTemplates)
	r.GET("/templates/new", controllers.NewTemplate)
	r.POST("/templates", controllers.CreateTemplate)
	r.POST("/template/:id", controllers.UpdateTemplate)
	r.GET("/template/:id/delete", controllers.DeleteTemplate)
	r.GET("/template/:id", controllers.GetTemplate)

	r.GET("/reports", controllers.GetReportDateRange)

	r.Static("/s/", "./public")

	r.Run(":7000")
}
