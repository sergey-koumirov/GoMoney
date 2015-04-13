package main

import(
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/martini-contrib/binding"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"

    "fmt"

    "controllers"
    "models"
    "utils"
    "html/template"
)

func main() {
    //*** DB INIT ***
    db, error := gorm.Open("sqlite3", "money_0.db")
    if(error !=nil){
        fmt.Println(error)
    }
    defer db.Close()
    db.DB()

    dbi := &db

    //*** APP INIT ***
    m := martini.Classic()
    m.Map(dbi)
    m.Use(render.Renderer(render.Options{
        Layout: "layout",
        Extensions: []string{".tmpl", ".html"},
        Funcs: []template.FuncMap{{ "money": utils.RenderMoney, "float": utils.MoneyAsFloat }},
    }))

    //*** ROUTES ***
    m.Get("/", controllers.GetTransactions)

    m.Get("/accounts", controllers.GetAccounts)
    m.Group("/accounts", func(r martini.Router) {
        r.Get("/new", controllers.NewAccount)
        r.Post("/create", binding.Bind(models.Account{}), controllers.CreateAccount)
        r.Post("/update/:id", binding.Bind(models.Account{}), controllers.UpdateAccount)
        r.Get("/delete/:id", controllers.DeleteAccount)
        r.Get("/:id", controllers.GetAccount)
    })

    m.Get("/currencies", controllers.GetCurrencies)
    m.Group("/currencies", func(r martini.Router) {
        r.Get("/new", controllers.NewCurrency)
        r.Post("/create", binding.Bind(models.Currency{}), controllers.CreateCurrency)
        r.Post("/update/:id", binding.Bind(models.Currency{}), controllers.UpdateCurrency)
        r.Get("/delete/:id", controllers.DeleteCurrency)
        r.Get("/:id", controllers.GetCurrency)
    })

    m.Get("/transactions", controllers.GetTransactions)
    m.Group("/transactions", func(r martini.Router) {
        r.Get("/new", controllers.NewTransaction)
        r.Post("/create", binding.Bind(models.Transaction{}), controllers.CreateTransaction)
        r.Post("/update/:id", binding.Bind(models.Transaction{}), controllers.UpdateTransaction)
        r.Get("/delete/:id", controllers.DeleteTransaction)
        r.Get("/:id", controllers.GetTransaction)
    })

    m.Get("/meters", controllers.GetMeters)
    m.Group("/meters", func(r martini.Router) {
        r.Get("/new", controllers.NewMeter)
        r.Post("/create", binding.Bind(models.Meter{}), controllers.CreateMeter)
        r.Post("/update/:id", binding.Bind(models.Meter{}), controllers.UpdateMeter)
        r.Get("/delete/:id", controllers.DeleteMeter)
        r.Get("/:id", controllers.GetMeter)
    })

    m.Get("/meter_values", controllers.GetMeterValues)
    m.Group("/meter_values", func(r martini.Router) {
        r.Get("/new", controllers.NewMeterValue)
        r.Post("/create", binding.Bind(models.MeterValue{}), controllers.CreateMeterValue)
        r.Post("/update/:id", binding.Bind(models.MeterValue{}), controllers.UpdateMeterValue)
        r.Get("/delete/:id", controllers.DeleteMeterValue)
        r.Get("/:id", controllers.GetMeterValue)
    })

    m.RunOnAddr(":7000")

}