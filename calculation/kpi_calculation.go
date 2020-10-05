package calculation

import (
	"strconv"
	"xml-parser/models"
)

func Calculate(mdc models.MDC) {
	accessibilityCalculation(mdc)
}

func accessibilityCalculation(mdc models.MDC) {
	md2761 := populateCountersKeyValue(mdc)
	KPIs := calculateAccessibilityKPIs(md2761)

	for key, value := range KPIs {
		println("-------------", key, "---------------------")
		for k, v := range value {
			println(k, ":", v)
		}
	}
}

func populateCountersKeyValue(mdc models.MDC) map[string]map[string]string {
	md2761 := make(map[string]map[string]string)

	mdList := []uint{27, 61}
	EUtranCellFDD := []string{"LN4277XA0", "LN4277XA10", "LN4277XA2", "LN4277XB0", "LN4277XB1", "LN4277XB2", "LN4277XC0", "LN4277XC1", "LN4277XC2"}
	for moidKey, moid := range EUtranCellFDD {
		md2761[moid] = make(map[string]string)
		for _, md := range mdList {
			mt := mdc.Md[md].Mi.Mt
			mv := mdc.Md[md].Mi.Mv
			for i, el := range mt {
				md2761[moid][el] = mv[moidKey].R[i]
			}
		}
	}

	return md2761
}

func calculateAccessibilityKPIs(md2761 map[string]map[string]string) map[string]map[string]int {
	KPIs := make(map[string]map[string]int)

	for EUtranCellFDDKey, EUtranCellFDDValue := range md2761 {
		counters := convertCounterToInt(EUtranCellFDDValue)

		kpis := calculateKPIs(counters)

		KPIs[EUtranCellFDDKey] = kpis
	}

	return KPIs
}

func convertCounterToInt(EUtranCellFDDValue map[string]string) map[string]int {
	integerCounter := make(map[string]int)
	integerCounter["pmRrcConnEstabSucc"], _ = strconv.Atoi(EUtranCellFDDValue["pmRrcConnEstabSucc"])
	integerCounter["pmRrcConnEstabAtt"], _ = strconv.Atoi(EUtranCellFDDValue["pmRrcConnEstabAtt"])
	return integerCounter

}

func calculateKPIs(counters map[string]int) map[string]int {
	var (
		Acc_RrcConnSetupSuccRate int
		Acc_test                 int
	)
	Acc_RrcConnSetupSuccRate = 100 * counters["pmRrcConnEstabSucc"]
	Acc_test = 100 * counters["pmRrcConnEstabAtt"]
	kpi := make(map[string]int)
	kpi["Acc_RrcConnSetupSuccRate"] = Acc_RrcConnSetupSuccRate
	kpi["Acc_test"] = Acc_test

	return kpi
}
