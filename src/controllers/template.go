package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sergey-koumirov/GoMoney/src/db"
	"github.com/sergey-koumirov/GoMoney/src/models"
)

func GetTemplates(c *gin.Context) {

	var templates []models.Template
	db.DBI.Preload("AccountFrom").Preload("AccountTo").Order("id desc").Find(&templates)

	c.HTML(200, "templates/index", templates)
}

func GetTemplate(c *gin.Context) {
	var template models.Template
	db.DBI.Find(&template, c.Param("id"))

	var accountList []models.Account
	db.DBI.Order("type, name").Find(&accountList)

	formData := models.TemplateForm{T: template, AccountList: accountList}

	c.HTML(200, "templates/edit", formData)
}

func NewTemplate(c *gin.Context) {
	template := models.Template{}

	var accountList []models.Account
	db.DBI.Order("type, name").Find(&accountList)

	formData := models.TemplateForm{T: template, AccountList: accountList}

	c.HTML(200, "templates/new", formData)
}

func CreateTemplate(c *gin.Context) {
	var template models.Template

	if c.ShouldBind(&template) == nil {
		template.ParseMoney()
		db.DBI.Create(&template)
	}

	c.Redirect(302, "/templates")
}

func UpdateTemplate(c *gin.Context) {
	var template models.Template
	template.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Find(&template)

	if c.ShouldBind(&template) == nil {
		template.ParseMoney()
		db.DBI.Save(&template)
	}

	c.Redirect(302, "/templates")
}

func DeleteTemplate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	db.DBI.Where("id = ?", id).Delete(models.Template{})
	c.Redirect(302, "/templates")
}
