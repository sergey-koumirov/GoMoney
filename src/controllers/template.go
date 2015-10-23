package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
//    "fmt"
    "net/http"
    "github.com/jinzhu/gorm"
    "GoMoney/src/models"
    "strconv"
    "fmt"
)

func GetTemplates(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){

    var templates []models.Template
    db.Preload("AccountFrom").Preload("AccountTo").Order("id desc").Find(&templates)

    r.HTML(
      200, "templates/index",
      templates,
    );
}

func GetTemplate(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var template models.Template
    db.Find(&template, params["id"])

    var accountList []models.Account
    db.Order("type, name").Find(&accountList)

    formData := models.TemplateForm{ T: template, AccountList: accountList }

    r.HTML(200, "templates/edit", formData)
}

func NewTemplate(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    template := models.Template{}

    var accountList []models.Account
    db.Order("type, name").Find(&accountList)

    formData := models.TemplateForm{ T: template, AccountList: accountList }

    r.HTML(200, "templates/new", formData)
}

func CreateTemplate(template models.Template, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    template.ParseMoney()
    db.Create(&template)
    r.Redirect("/templates")
}

func UpdateTemplate(template models.Template, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    template.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    template.ParseMoney()

    fmt.Println(template)

    db.Save(template)
    r.Redirect("/templates")
}

func DeleteTemplate(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Template{})
    r.Redirect("/templates")
}