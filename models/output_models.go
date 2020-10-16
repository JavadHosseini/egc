package models

type KPIJsonModel struct {
	SiteName   string
	SectorName string
	CellName   string
	KPI        kpiKeyValue
}

type kpiKeyValue map[string]int
