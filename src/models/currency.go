package models

type Currency struct {
    ID int64    `form:"ID"`
    Num string  `form:"Num"`
    Code string `form:"Code"`
}

