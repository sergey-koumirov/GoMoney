package controllers
import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "net/http"
    "github.com/jinzhu/gorm"
    models "GoMoney/src/models"
    "strconv"
    "time"
    //"fmt"
)

func GetMeters(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var meters []models.Meter
    db.Order("name").Find(&meters)
    r.HTML(200, "meters/index", meters)
}

func GetMeter(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var meter models.Meter
    db.Find(&meter, params["id"])
    r.HTML(200, "meters/edit", meter)
}

func NewMeter(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    meter := models.Meter{}
    r.HTML(200, "meters/new", meter)
}

func CreateMeter(meter models.Meter, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    db.Create(&meter)
    r.Redirect("/meters")
}

func UpdateMeter(meter models.Meter, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    meter.ID, _ = strconv.ParseInt(params["id"], 10, 64)
    db.Save(meter)
    r.Redirect("/meters")
}

func DeleteMeter(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.Meter{})
    r.Redirect("/meters")
}

//Meter Values
func GetMeterValues(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var meters []models.Meter
    db.Order("name").Find(&meters)

    data := models.MeterValuesOnDates(db)

    r.HTML(200, "meter_values/index", models.MeterValuesIndex{D: data, Meters: meters, Prev: data[1], Current: data[0]})
}

func MeterValuesPrint(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var meters []models.Meter
    db.Order("name").Find(&meters)
    data := models.MeterValuesOnDates(db)
    
    layout := "2006-01-02"
    str := data[0].Date
    t, _ := time.Parse(layout, str)
    
    println( t.Format("Januar 2006") )
    month := map[int]string{
        1: "Январь",
        2: "Февраль",
        3: "Март",
        4: "Апрель",
        5: "Май",
        6: "Июнь",
        7: "Июль",
        8: "Август",
        9: "Сентябрь",
        10: "Октябрь",
        11: "Ноябрь",
        12: "Декабрь",
    }
    
    r.HTML(200, "meter_values/print_version",
        models.MeterValuesIndex{D: data, Meters: meters, Prev: data[1], Current: data[0], MonthStr: strconv.Itoa(t.Year())+" "+month[int(t.Month())] },
        render.HTMLOptions{Layout: "print"},
    )
}

func GetMeterValue(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    var meter_value models.MeterValue
    db.Find(&meter_value, params["id"])

    var meters []models.Meter
    db.Find(&meters)

    r.HTML(200, "meter_values/edit", models.MeterValueForm{V: meter_value, Meters: meters})
}

func NewMeterValue(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    meterValue := models.MeterValue{}
    if(req.URL.Query().Get("date") != ""){
        meterValue.Date = req.URL.Query().Get("date")
    }else{
        meterValue.Date = time.Now().Format("2006-01-02")
    }
    if(req.URL.Query().Get("meter_id") != ""){
        meterValue.MeterID, _ = strconv.ParseInt(req.URL.Query().Get("meter_id"), 10, 64)
    }
    var meters []models.Meter
    db.Find(&meters)
    r.HTML(200, "meter_values/new",  models.MeterValueForm{V: meterValue, Meters: meters})
}

func CreateMeterValue(meter_value models.MeterValue, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    db.Create(&meter_value)
    r.Redirect("/meter_values")
}

func UpdateMeterValue(meter_value models.MeterValue, db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    meter_value.ID, _ = strconv.ParseInt(params["id"], 10, 64)

    db.Save(meter_value)
    r.Redirect("/meter_values")
}

func DeleteMeterValue(db *gorm.DB, params martini.Params, req *http.Request, r render.Render){
    id, _ := strconv.ParseInt(params["id"], 10, 64)
    db.Where("id = ?", id).Delete(models.MeterValue{})
    r.Redirect("/meter_values")
}
