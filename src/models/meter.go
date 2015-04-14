package models
import (
    "github.com/jinzhu/gorm"
    "fmt"
    "database/sql"
)

type Meter struct{
    ID int64 `form:"ID"`
    Name string `form:"Name"`
    Number string `form:"Number"`
}

type MeterValue struct{
    ID sql.NullInt64 `form:"ID"`
    MeterID int64 `form:"MeterID"`
    Meter Meter
    Value sql.NullFloat64 `form:"Value"`
    Date string `form:"Date"`
}

type MeterValueForm struct{
    V MeterValue
    Meters []Meter
}

type MeterValuesIndex struct{
    D []MeterValuesOnDate
    Meters []Meter
}

type MeterValuesOnDate struct{
    Date string
    Values []MeterValue
}

func MeterValuesOnDates(db *gorm.DB) []MeterValuesOnDate{
    result := []MeterValuesOnDate{}

    var meters []Meter
    db.Order("id").Find(&meters)

    joins := ""
    slct := ""
    for index, meter := range meters {
        joins = joins + fmt.Sprintf("left join meter_values m%d on m%d.meter_id = %d and m%d.date = x.Date\n", index, index, meter.ID, index)
        slct = slct + fmt.Sprintf(", m%d.value as v%d,  m%d.id as id%d", index, index, index, index)
    }
    dataSql := "select x.Date" + slct + "\n" + "from (select distinct date as Date from meter_values order by date desc) x \n"+joins

    rows, error := db.Raw(dataSql).Rows()
    if(error != nil){ fmt.Println(error) }

    //cols, _ := rows.Columns()


    for rows.Next() {
        item := MeterValuesOnDate{}
        item.Values = make([]MeterValue, len(meters))
        for i, _ := range item.Values {
            item.Values[i] = MeterValue{}
        }

        vals := make([]interface{}, len(meters) * 2 + 1)
        vals[0] = &item.Date
        for i, v := range meters {
            vals[i*2+1] = &item.Values[i].Value
            vals[i*2+2] = &item.Values[i].ID
            item.Values[i].MeterID = v.ID
        }
        error := rows.Scan(vals...)
        if(error != nil){ fmt.Println(error) }

        result = append(result, item)
    }

    return result
}