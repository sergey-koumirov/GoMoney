package main

import(
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"

    "fmt"

    "controllers"
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
        r.Get("/:id", controllers.GetAccount)
        r.Post("/new", controllers.NewAccount)
        r.Put("/update/:id", controllers.UpdateAccount)
        r.Delete("/delete/:id", controllers.DeleteAccount)
    })


    m.RunOnAddr(":7000")

}