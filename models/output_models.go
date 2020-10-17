package models

type CountersJsonModel struct {
	SiteName       string
	SectorName     string
	CellName       string
	StringCounters CountersKeyValue
	Counters       IntegerCountersKeyValue
}

type CountersKeyValue map[string]string
type IntegerCountersKeyValue map[string]int

type KPIJsonModel struct {
	SiteName   string
	SectorName string
	CellName   string
	KPIs       KPIKeyValue
}

type KPIKeyValue map[string]int
