package models

type DateRangeReport struct {

    D map[string]map[int64][]Transaction

    T []Transaction
    Total int64
    BeginDate string
    EndDate string
    AccountIDs []int64
    AccountList []Account
}
