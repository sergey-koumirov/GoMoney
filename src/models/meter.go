package models

type Meter struct{
    ID int64 `form:"ID"`
    Name string `form:"Name"`
    Number string `form:"Number"`
}

type MeterValue struct{
    ID int64 `form:"ID"`
    MeterID int64 `form:"MeterID"`
    Meter Meter
    Value float64 `form:"Value"`
    Date string `form:"Date"`
}

type MeterValueForm struct{
    V MeterValue
    Meters []Meter
}