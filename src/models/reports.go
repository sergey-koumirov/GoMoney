package models

type SectionInfo struct{
    SubSections map[string]*SectionInfo
    Title string
    Total int64
}

type DateRangeReport struct {
    Sections map[string]*SectionInfo
    BeginDate string
    EndDate string
    AccountList []Account
}

func FillAccountTypeSectionsInfo(ts []Transaction) map[string]*SectionInfo{
    result := map[string]*SectionInfo{}
    for _, t := range ts {
        code := t.AccountFrom.Type + t.AccountTo.Type
        section, ok := result[code]
        if !ok {
            section = &SectionInfo{ SubSections: map[string]*SectionInfo{} }
            result[code] = section
        }

        accSection, ok := section.SubSections[string(t.AccountFromID)]
        if !ok {
            accSection = &SectionInfo{SubSections: map[string]*SectionInfo{}, Title: t.AccountFrom.Name}
            section.SubSections[string(t.AccountFromID)] = accSection
        }

        trInfo, ok := accSection.SubSections[string(t.AccountToID)]
        if !ok {
            trInfo = &SectionInfo{Title: t.AccountTo.Name}
            accSection.SubSections[string(t.AccountToID)] = trInfo
        }

        trInfo.Total += t.AmountFrom
        accSection.Total += t.AmountFrom
        section.Total = section.Total + t.AmountFrom
    }

    return result
}
