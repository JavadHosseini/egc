package models

type CountersJsonModel struct {
	SiteName   string
	SectorName string
	CellName   string
	Counters   CountersKeyValue
}

type CountersKeyValue map[string]string

type KPIJsonModel struct {
	SiteName   string
	SectorName string
	CellName   string
	KPIs       KPIKeyValue
}

type KPIKeyValue map[string]int
