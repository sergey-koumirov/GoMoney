package models

type TransactionTypeInfo struct{
    AccountName string
    Total int64
}

type TransactionsSectionInfo struct{
    D map[int64]*TransactionTypeInfo
    AccountName string
    Total int64
}

type AccountTypeSectionInfo struct{
    AccountInfos map[int64]*TransactionsSectionInfo
    Total int64
}

type DateRangeReport struct {
    Sections map[string]*AccountTypeSectionInfo
    BeginDate string
    EndDate string
    AccountList []Account
}

func FillAccountTypeSectionsInfo(ts []Transaction) map[string]*AccountTypeSectionInfo{
    result := map[string]*AccountTypeSectionInfo{}
    for _, t := range ts {
        code := t.AccountFrom.Type + t.AccountTo.Type
        section, ok := result[code]
        if !ok {
            section = &AccountTypeSectionInfo{ AccountInfos: map[int64]*TransactionsSectionInfo{} }
            result[code] = section
        }

        accSection, ok := section.AccountInfos[t.AccountFromID]
        if !ok {
            accSection = &TransactionsSectionInfo{D: map[int64]*TransactionTypeInfo{}, AccountName: t.AccountFrom.Name}
            section.AccountInfos[t.AccountFromID] = accSection
        }

        trInfo, ok := accSection.D[t.AccountToID]
        if !ok {
            trInfo = &TransactionTypeInfo{AccountName: t.AccountTo.Name}
            accSection.D[t.AccountToID] = trInfo
        }

        trInfo.Total += t.AmountFrom
        accSection.Total += t.AmountFrom
        section.Total = section.Total + t.AmountFrom
    }

    return result
}
