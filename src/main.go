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
    }))

    //*** ROUTES ***
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


    m.RunOnAddr(":7000")

}